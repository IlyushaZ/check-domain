FROM golang:1.14

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go mod download

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait /wait
RUN chmod +x /wait

EXPOSE 8084

CMD /wait && go run ./cmd/migrations/main.go up && go run ./cmd/google-domain-checker/main.go