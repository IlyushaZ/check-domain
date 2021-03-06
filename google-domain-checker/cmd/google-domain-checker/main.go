package main

import (
	"flag"
	"fmt"
	internalErr "github.com/IlyushaZ/check-domain/google-domain-checker/internal/error"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/notifier"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/request"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/search"
	"github.com/IlyushaZ/check-domain/google-domain-checker/internal/task"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type config struct {
	Port            int    `yaml:"port"`
	DBString        string `yaml:"db_string"`
	SerpStackAPIKey string `yaml:"serp_stack_api_key"`
	RabbitMQ        struct {
		URL       string `yaml:"url"`
		QueueName string `yaml:"queue_name"`
	} `yaml:"rabbitmq"`
}

func main() {
	var dir string
	flag.StringVar(&dir, "confDir", "./config/app.yaml", "path to config file")
	flag.Parse()

	file, err := ioutil.ReadFile(dir)
	handleError(err)

	var config config
	handleError(yaml.Unmarshal(file, &config))

	db, err := connectToDB(config.DBString)
	handleError(err)
	defer db.Close()

	taskRepository := task.NewRepository(db)
	requestRepository := request.NewRepository(db)
	taskHandler := task.NewHandler(task.NewService(taskRepository, requestRepository))

	amqpConn, err := connectToRabbitMQ(config.RabbitMQ.URL)
	defer amqpConn.Close()
	amqpNotifier, err := notifier.NewRabbitMQ(amqpConn, config.RabbitMQ.QueueName)
	handleError(err)
	taskProcessor := task.NewProcessor(
		taskRepository,
		task.NewGoogleChecker(requestRepository, search.NewGoogleSearcher(config.SerpStackAPIKey), amqpNotifier),
		task.NewMinuteSleeper(1, time.Sleep),
	)

	go func() {
		for {
			taskProcessor.Process()
		}
	}()

	handler := http.NewServeMux()
	handler.HandleFunc("/tasks", internalErr.Handler(taskHandler.CreateTask).RespondError)

	server := http.Server{
		Handler: handler,
		Addr:    fmt.Sprintf(":%d", config.Port),
	}

	handleError(server.ListenAndServe())
}

func connectToDB(DBString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", DBString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute)

	return db, nil
}

func connectToRabbitMQ(urlString string) (*amqp.Connection, error) {
	return amqp.Dial(urlString)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
