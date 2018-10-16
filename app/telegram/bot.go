package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/telegram/command"
	"log"
)

type Telegram struct {
	bot          *tgbotapi.BotAPI
	updateChanel tgbotapi.UpdatesChannel
}

func New(token string) (*Telegram, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	return &Telegram{bot: bot, updateChanel: updates}, nil
}

func (t *Telegram) Start(executor *action.Executor) {
	for update := range t.updateChanel {
		log.Printf("[DEBUG] new telegram message: %v", update)
		go handleMessage(update, executor)
	}
}

func (t *Telegram) SendWithKeyboard(chatId int, message string, keyboard map[ButtonValue]ButtonName) error {
	msg := tgbotapi.NewMessage(int64(chatId), message)

	keys := CreateKeyboard(keyboard)
	msg.ReplyMarkup = keys

	_, err := t.bot.Send(msg)

	return err
}

func handleMessage(update tgbotapi.Update, executor *action.Executor) {
	cmd, err := command.Parse(update)

	if err != nil {
		log.Printf("[ERROR] error during parse: %v", err)
		return
	}

	err = executor.Execute(cmd)

	if err != nil {
		log.Printf("[ERROR] error during execution cmd: %v", err)
		return
	}
}
