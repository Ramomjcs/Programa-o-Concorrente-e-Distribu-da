package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {

	args := os.Args[1:]
	prefetchCount, err := strconv.Atoi(args[0])
	if err != nil {
		panic("Argumento invalido")
	}

	con, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic("could not connect to server")
	}
	ch, err := con.Channel()
	if err != nil {
		panic("falha ao criar canal")
	}

	ch.Qos(prefetchCount, 0, false)

	FilesQ, err := ch.QueueDeclare(
		"arquivos",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("falha ao criar fila de arquivos")
	}

	CanalArquivos, err := ch.Consume(
		FilesQ.Name, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		panic("falha ao obter canal de arquivos")
	}

	tempos := make(chan time.Duration)

	go func() {
		var tt time.Duration
		counter := 0
		for t := range tempos {
			tt += t
			counter++
			fmt.Println("Ultimo tempo:", t)
			fmt.Println("Media de entrega:", tt/time.Duration(counter))
		}
	}()

	for f := range CanalArquivos {
		err := f.Ack(false)
		if err != nil {
			panic("falha ao enviar ack")
		}
		tempoEnvioInt, _ := f.Headers["publishTime"].(int64)
		tempoEnvio := time.UnixMilli(tempoEnvioInt)
		tempoEntrega := time.Since(tempoEnvio)
		tempos <- tempoEntrega
	}

}
