package write

import (
	"errors"
	"fmt"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/command"
	"github.com/jaitl/goEnglishBot/app/exercises"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"strings"
)

type Action struct {
	AwsSession    *aws.Session
	ActionSession *action.SessionModel
	Bot           *telegram.Telegram
	PhraseModel   *phrase.Model
	Audio         *telegram.AudioService
}

const (
	Start           action.Stage = "start"
	WaitWrittenText action.Stage = "waitWrittenText"
)

const (
	Mode        action.SessionKey = "mode"
	Session     action.SessionKey = "writeSession"
	CountErrors action.SessionKey = "countErrors"
)

const (
	AudioMode string = "AudioMode"
	TransMode string = "TransMode"
)

const (
	maxCountErrors int = 2
)

func (a *Action) GetType() action.Type {
	return action.Write
}

func (a *Action) GetStartStage() action.Stage {
	return Start
}

func (a *Action) GetWaitCommands(stage action.Stage) map[command.Type]bool {
	switch stage {
	case Start:
		return map[command.Type]bool{command.WriteAudio: true, command.WriteTrans: true}
	case WaitWrittenText:
		return map[command.Type]bool{command.Text: true}
	}

	return nil
}

func (a *Action) Execute(stage action.Stage, cmd command.Command, session *action.Session) error {
	switch stage {
	case Start:
		return a.startStage(cmd)
	case WaitWrittenText:
		return a.waitWrittenText(cmd, session)
	}

	return fmt.Errorf("stage %s not found in WriteAction", stage)
}

func (a *Action) startStage(cmd command.Command) error {
	var mode string
	var from, to *int

	switch mcmd := cmd.(type) {
	case *command.WriteAudioCommand:
		mode = AudioMode
		from = mcmd.From
		to = mcmd.To
	case *command.WriteTransCommand:
		mode = TransMode
		from = mcmd.From
		to = mcmd.To
	default:
		return errors.New("command does not belong to Start stage in WriteAction")
	}

	phrs, err := a.PhraseModel.SmartFindByRange(cmd.GetUserId(), from, to)

	if err != nil {
		return err
	}

	if len(phrs) == 0 {
		return errors.New("don't correct range")
	}

	write := exercises.NewComposite(phrs, exercises.WriteMode, true)

	err = a.newWrite(write, mode)

	ses := action.CreateSession(cmd.GetUserId(), action.Write, WaitWrittenText)
	ses.AddData(Mode, mode)
	ses.AddData(Session, write)
	ses.AddData(CountErrors, 0)
	a.ActionSession.UpdateSession(ses)

	return err
}

func (a *Action) waitWrittenText(cmd command.Command, session *action.Session) error {
	text, ok := cmd.(*command.TextCommand)

	if !ok {
		return errors.New("command does not belong to WaitWrittenText stage in WriteAction")
	}

	write := session.Data[Session].(*exercises.Composite)
	mode := session.GetStringData(Mode)

	words := strings.Split(exercises.ClearText(text.Text), " ")

	writeRes := write.HandleAnswer(words)

	msg := fmt.Sprintf("Фраза №%d из %d", writeRes.Pos+1, writeRes.CountPhrases)
	msg += fmt.Sprintf("\nФраза: %s", writeRes.Result.AnsweredText)

	if writeRes.Result.IsFinish && writeRes.IsFinish {
		a.ActionSession.ClearSession(cmd.GetUserId())
		msg += "\nФраза успешно завершена!"
		msg += "\nУпражнение успешно завершено!"
		return a.Bot.Send(cmd.GetUserId(), msg)
	}

	if writeRes.Result.IsFinish && !writeRes.IsFinish {
		msg += "\nФраза успешно завершена!"
		err := a.Bot.Send(cmd.GetUserId(), msg)

		if err != nil {
			return err
		}

		return a.newWrite(write, mode)
	}

	countErrors := session.GetIntData(CountErrors)

	msg += fmt.Sprintf("\nОсталось слов: %d", writeRes.Result.WordsLeft)

	if writeRes.Result.IsCorrectAnswer {
		countErrors = 0
	} else {
		msg += "\nНекорректное слово!"
		countErrors += 1
	}

	if countErrors >= maxCountErrors {
		msg += fmt.Sprintf("\nСледующее слово: %s", writeRes.Result.NextAnswer)
		countErrors = 0
	}

	session.AddData(CountErrors, countErrors)
	a.ActionSession.UpdateSession(session)

	return a.Bot.Send(cmd.GetUserId(), msg)
}

func (a *Action) newWrite(puzzle *exercises.Composite, mode string) error {
	puzzleRes := puzzle.Next()

	msg := fmt.Sprintf("Фраза №%d из %d", puzzleRes.Pos+1, puzzleRes.CountPhrases)

	if mode == AudioMode {
		msg += "\nНапишите фразу, которую вы слышите"

		err := a.Audio.SendAudio(puzzleRes.Phrase)
		if err != nil {
			return err
		}
	} else {
		msg += fmt.Sprintf("\nНапишите фразу: %s", puzzleRes.Phrase.RussianText)
	}

	return a.Bot.Send(puzzleRes.Phrase.UserId, msg)
}
