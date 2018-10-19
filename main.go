package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func main() {

	var port string
	defaultPort := "8080"
	nodeListURL := fmt.Sprintf("http://localhost:%s", defaultPort)

	flag.StringVar(&port, "port", defaultPort, "The nodes port")
	flag.Parse()

	node := Node{
		Name: uuid.New().String(),
		URL:  fmt.Sprintf("http://localhost:%s", port),
	}

	s := NewServer(node, nodeListURL)
	s.seedChain()
	s.setupRoutes()
	if port != defaultPort {
		if err := s.getNodes(); err != nil {
			log.Fatal(err)
		}
		if err := s.broadcast(node); err != nil {
			log.Fatal(err)
		}
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

type Transaction struct {
	Sender    string `json:"sender"`
	Recepient string `json:"recepient"`
	Amount    int    `json:"amount"`
}

type Block struct {
	Index          int           `json:"index"`
	Timestamp      time.Time     `json:"timestamp"`
	Transactions   []Transaction `json:"transactions"`
	PreviousDigest []byte        `json:"previousDigest"`
}

type Node struct {
	Name string `json:"name" ans1:"Name"`
	URL  string `json:"url" asn1:"URL"`
}
