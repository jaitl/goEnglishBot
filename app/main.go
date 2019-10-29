package main

import (
	"github.com/globalsign/mgo"
	"github.com/jaitl/goEnglishBot/app/action"
	"github.com/jaitl/goEnglishBot/app/action/add"
	"github.com/jaitl/goEnglishBot/app/action/audio"
	"github.com/jaitl/goEnglishBot/app/action/list"
	"github.com/jaitl/goEnglishBot/app/action/voice"
	"github.com/jaitl/goEnglishBot/app/aws"
	"github.com/jaitl/goEnglishBot/app/phrase"
	"github.com/jaitl/goEnglishBot/app/settings"
	"github.com/jaitl/goEnglishBot/app/telegram"
	"github.com/jessevdk/go-flags"
	"log"
)

var opts struct {
	TelegramToken   string `long:"token" env:"TOKEN" required:"true"`
	MongoDbUrl      string `long:"mongo-db-url" env:"MONGO_DB_URL" required:"true"`
	AWSKey          string `long:"aws-key" env:"AWS_KEY" required:"true"`
	AWSSecret       string `long:"aws-secret" env:"AWS_SECRET" required:"true"`
	AWSRegion       string `long:"aws-region" env:"AWS_REGION" required:"true"`
	AWSS3BucketName string `long:"aws-s3-bucket-name" env:"AWS_S3_BUCKET_NAME" required:"true"`
	AWSS3VoicePath  string `long:"aws-s3-voice-path" env:"AWS_S3_VOICE_PATH" required:"true"`
	PathToTmpFolder string `long:"tmp-folder" env:"TMP_FOLDER" required:"true"`
}

func main() {
	log.Println("[INFO] start goEnglishBot")

	if _, err := flags.Parse(&opts); err != nil {
		log.Panic(err)
	}

	commonSettings := &settings.CommonSettings{
		TmpFolder:    opts.PathToTmpFolder,
		AwsRegion:    opts.AWSRegion,
		S3BucketName: opts.AWSS3BucketName,
		S3VoicePath:  opts.AWSS3VoicePath,
	}

	mongoSession, err := mgo.Dial(opts.MongoDbUrl)

	if err != nil {
		log.Panic(err)
	}

	phraseModel := phrase.NewModel(mongoSession, "goEnglishBot")
	actionSession := action.NewSessionMongoModel(mongoSession, "goEnglishBot")

	awsSession, err := aws.New(opts.AWSKey, opts.AWSSecret, commonSettings)

	if err != nil {
		log.Panic(err)
	}

	telegramBot, err := telegram.New(opts.TelegramToken)

	if err != nil {
		log.Panic(err)
	}

	actions := []action.Action{
		&add.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel},
		&list.Action{Bot: telegramBot, PhraseModel: phraseModel},
		&audio.Action{Bot: telegramBot, PhraseModel: phraseModel, AwsSession: awsSession},
		&voice.Action{AwsSession: awsSession, ActionSession: actionSession, Bot: telegramBot, PhraseModel: phraseModel, CommonSettings: commonSettings},
	}

	actionExecutor := action.NewExecutor(actionSession, actions)

	telegramBot.Start(actionExecutor)
}
