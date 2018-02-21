package main

import (
	"flag"
	"fmt"
	"os"

	"time"

	"io/ioutil"

	"net"

	"crypto/tls"
	"crypto/x509"

	"github.com/bshuster-repo/logrus-logstash-hook"
	log "github.com/sirupsen/logrus"
)

func main() {
	sleep := flag.Int("sleep", 5, "time to sleep in seconds")
	logstash := flag.String("logstash", "", "logstash server")
	secure := flag.Bool("tls", false, "use tls for communications")
	cert := flag.String("cert", "", "pem cert to user for tls")
	flag.Parse()

	q := NewQuotes()

	if *logstash == "" {
		fmt.Println("logstash server to use not provided")
		os.Exit(1)
	}

	if *secure && *cert == "" {
		fmt.Println("tls wanted, but no cert name supplied")
		os.Exit(1)
	}

	var c []byte

	if *secure && *cert != "" {
		var err error
		if c, err = ioutil.ReadFile(*cert); err != nil {
			fmt.Printf("error reading certificat file : %s\n", *cert)
			fmt.Printf("error : %v\n", err)
			os.Exit(1)
		}
	}

	var conn net.Conn

	if !*secure {
		var err error
		if conn, err = net.Dial("tcp", *logstash); err != nil {
			fmt.Printf("error connecting to logstash : %s\n", *logstash)
			fmt.Printf("error : %v\n", err)
			os.Exit(1)
		}
	} else {
		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(c)
		if !ok {
			fmt.Println("failed to parse certificate")
			os.Exit(1)
		}
		var err error
		conn, err = tls.Dial("tcp", *logstash, &tls.Config{
			RootCAs: roots,
		})
		if err != nil {
			fmt.Printf("failed to connect to logstash : %s\n", *logstash)
			fmt.Printf("error : %v\n", err)
			os.Exit(1)
		}
	}

	defer conn.Close()

	hook, err := logrustash.NewHookWithConn(conn, "test-quotes")

	if err != nil {
		fmt.Printf("error creating logstash logging hook : %s\n", *logstash)
		fmt.Printf("%v\n", err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel) // default is Info

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
