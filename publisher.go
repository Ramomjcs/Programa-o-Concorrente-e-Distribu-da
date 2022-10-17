package main

import (
	"os"
	"time"

	"github.com/streadway/amqp"
)

func main() {

	con, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic("could not connect to server")
	}
	ch, err := con.Channel()
	if err != nil {
		panic("falha ao criar canal")
	}

	FileQ, err := ch.QueueDeclare(
		"arquivos",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("falha ao criar fila req")
	}

	for i := 0; i < 10000; i++ {
		conteudo, err := os.ReadFile("1024.txt")
		if err != nil {
			panic("Arquivo nao encontrado")
		}
		timeHeader := amqp.Table{
			"publishTime": time.Now().UnixMilli(),
		}
		err = ch.Publish(
			"",         // exchange
			FileQ.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				Headers:     timeHeader,
				ContentType: "text/plain",
				Body:        conteudo,
			})
		if err != nil {
			panic("falhou ao publicar na fila")
		}
	}

}
