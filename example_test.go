package techlog_test

import (
	"github.com/k0kubun/pp"
	"log"
	"v8platform/techlog"
)

func ExampleRead_file() {

	file := "./logs/20100521.log"

	events, err := techlog.Read(file)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("readed <%d> events", len(events))
}

func ExampleStreamRead_file() {

	file := "./logs/20100521.log"

	events, err := techlog.StreamRead(file, 10)
	if err != nil {
		log.Fatal(err)
	}

	count := 0

	for event := range events {
		pp.Println(event) // Не ракомендую использовать на больших объемах
		count++
	}

	log.Printf("readed <%d> events", count)
}
