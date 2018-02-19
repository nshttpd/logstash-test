package main

import (
	"flag"
	"fmt"
	"os"

	"time"

	"github.com/bshuster-repo/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
)

func main() {
	sleep := flag.Int("sleep", 5, "time to sleep in seconds")
	logstash := flag.String("logstash", "", "logstash server")
	flag.Parse()

	q := NewQuotes()

	if *logstash == "" {
		fmt.Println("logstash server to use not provided")
		os.Exit(1)
	}

	hook, err := logrustash.NewHook("tcp", *logstash, "test-quotes")

	if err != nil {
		fmt.Printf("error creating logstash logging hook : %s\n", *logstash)
		fmt.Printf("%v\n", err)
	}

	log.SetFormatter(&log.JSONFormatter{})

	log.AddHook(hook)

	x := 0

	for {
		log.WithFields(log.Fields{
			"sequence": x,
		}).Debug(q.RandomQuote())

		x++

		time.Sleep(time.Duration(*sleep) * time.Second)
	}

}
