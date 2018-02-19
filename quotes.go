package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

const (
	QUOTEFILE = "zippy.json"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Quotes struct {
	quotes []string
}

func NewQuotes() *Quotes {
	q := &Quotes{}

	b, err := ioutil.ReadFile(QUOTEFILE)

	if err != nil {
		log.Fatalf("error reading quote file : %s", QUOTEFILE)
	}

	err = json.Unmarshal(b, &q.quotes)

	if err != nil {
		log.Fatalf("error unmarshalling quotes file : %v", err)
	}

	return q

}

func (q Quotes) RandomQuote() string {
	return q.quotes[rand.Intn(len(q.quotes))]
}
