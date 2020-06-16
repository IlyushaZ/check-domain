package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type config struct {
	TelegramBotToken string `yaml:"telegram_bot_token"`
	RabbitMQ         struct {
		URL       string `yaml:"url"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
}

const chatIDsFile = "ids.txt"

func main() {
	var dir string
	flag.StringVar(&dir, "confDir", "./config/app.yaml", "path to config file")
	flag.Parse()

	file, err := ioutil.ReadFile(dir)
	handleError(err)

	var config config
	handleError(yaml.Unmarshal(file, &config))

	bot := runBot(config.TelegramBotToken)
	consume(config.RabbitMQ.URL, config.RabbitMQ.QueueName, bot)
}

type message struct {
	Domain  string `json:"domain"`
	Request string `json:"request"`
}

func consume(url, queueName string, bot *tgbotapi.BotAPI) {
	conn, err := amqp.Dial(url)
	handleError(err)
	defer conn.Close()

	ch, err := conn.Channel()
	handleError(err)
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err)

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	handleError(err)

	forever := make(chan bool)

	go func() {
		for m := range msgs {
			var message message
			_ = json.Unmarshal(m.Body, &message)
			fmt.Println(message)

			f, err := os.Open(chatIDsFile)
			handleError(err)

			s := bufio.NewScanner(f)
			s.Split(bufio.ScanLines)

			text := fmt.Sprintf("domain %s is not on first page for request %s", message.Domain, message.Request)

			for s.Scan() {
				chatID, _ := strconv.ParseInt(s.Text(), 10, 64)
				tgMessage := tgbotapi.NewMessage(chatID, text)
				bot.Send(tgMessage)
			}

			f.Close()
		}
	}()

	fmt.Println("started consuming messages")
	<-forever
}

func runBot(token string) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(token)
	handleError(err)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	handleError(err)

	go func() {
		for upd := range updates {
			if upd.Message == nil {
				continue
			}

			if upd.Message.Command() == "start" {
				f, err := os.OpenFile(chatIDsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				handleError(err)
				_, err = f.WriteString(strconv.FormatInt(upd.Message.Chat.ID, 10) + "\n")
				handleError(err)
				f.Close()
			}
		}
	}()

	return bot
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
