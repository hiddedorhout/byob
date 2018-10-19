package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Server struct {
	NodeID              string        `json:"nodeId"`
	Chain               []*Block      `json:"chain"`
	CurrentTransactions []Transaction `json:"currentTransactions"`
	Nodes               []Node        `json:"nodes"`
	ListURL             string        `json:"listURL"`
	IsLeader            bool          `json:"isLeader"`
}

func NewServer(self Node, nodeListURL string) *Server {
	s := &Server{
		NodeID:              self.Name,
		Chain:               make([]*Block, 0),
		CurrentTransactions: make([]Transaction, 0),
		Nodes:               []Node{self},
		ListURL:             nodeListURL,
		IsLeader:            false,
	}
	return s
}

func (s *Server) seedChain() {
	s.newBlock([]byte{1})
}

func (s *Server) lastBlock() *Block {
	return s.Chain[len(s.Chain)-1]
}

func (s *Server) newTransaction(sender, recepient string, amount int) int {
	s.CurrentTransactions = append(s.CurrentTransactions, Transaction{
		Sender:    sender,
		Recepient: recepient,
		Amount:    amount,
	})
	return s.lastBlock().Index + 1
}

func (s *Server) newBlock(previousHash []byte) {
	b := Block{
		Index:          len(s.Chain) + 1,
		Timestamp:      time.Now(),
		Transactions:   s.CurrentTransactions,
		PreviousDigest: previousHash,
	}
	s.CurrentTransactions = make([]Transaction, 0)

	s.Chain = append(s.Chain, &b)
}

func (s *Server) hash(block *Block) []byte {
	raw, err := asn1.Marshal(*block)
	if err != nil {
		log.Fatal(err)
	}
	h := sha256.New()
	h.Write(raw)
	hashed := h.Sum(nil)
	return hashed
}

func (s *Server) broadcast(toAnnounce interface{}) error {

	switch v := toAnnounce.(type) {
	case Node:
		ta, err := asn1.Marshal(v)
		if err != nil {
			return err
		}
		for _, node := range s.Nodes {
			if node.Name != s.NodeID {
				req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/node-announcement/", node.URL), bytes.NewBuffer(ta))
				if err != nil {
					return err
				}
				if _, err := doReq(req); err != nil {
					log.Println(err)
				}
			}
		}
	case Transaction:
		ta, err := asn1.Marshal(v)
		if err != nil {
			return err
		}
		for _, node := range s.Nodes {
			if node.Name != s.NodeID {
				req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/transaction/", node.URL), bytes.NewBuffer(ta))
				if err != nil {
					return err
				}
				if _, err := doReq(req); err != nil {
					log.Println(err)
				}
			}
		}
	}

	return nil
}

func (s *Server) resolve() {

}

func (s *Server) getNodes() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/status", s.ListURL), nil)
	if err != nil {
		return err
	}
	resp, err := doReq(req)
	if err != nil {
		return err
	}
	rawResponseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var currentStatus Server
	if err := json.Unmarshal(rawResponseBody, &currentStatus); err != nil {
		return err
	}
	s.CurrentTransactions = currentStatus.CurrentTransactions
	s.Chain = currentStatus.Chain
	s.Nodes = currentStatus.Nodes
	return nil
}
