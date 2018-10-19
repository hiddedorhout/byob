package main

import (
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (s *Server) setupRoutes() {
	http.HandleFunc("/api/new-transaction/", s.transactionHandler)
	http.HandleFunc("/api/mine", s.mineHandler)
	http.HandleFunc("/api/chain", s.chainHandler)
	http.HandleFunc("/api/node-announcement/", s.handleNodeAnnouncement)
	http.HandleFunc("/api/transaction/", s.handleTransaction)

	// utils
	http.HandleFunc("/api/list", s.getListHandler)
	http.HandleFunc("/api/status", s.statusHandler)

}

func (s *Server) transactionHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var transaction Transaction
	if err := json.Unmarshal(body, &transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	nextBlockIndex := s.newTransaction(transaction.Sender, transaction.Recepient, transaction.Amount)
	go s.broadcast(transaction)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte(fmt.Sprintf("A transaction has been added to block: %d", nextBlockIndex)))
}

func (s *Server) mineHandler(w http.ResponseWriter, r *http.Request) {
	previousHash := s.hash(s.lastBlock())
	s.newBlock(previousHash)

	response := "A New block will be mined"
	w.Header().Set("Content-Type", "plain/text")
	w.Write([]byte(response))
}

func (s *Server) chainHandler(w http.ResponseWriter, r *http.Request) {
	response := make([]Block, 0)
	for _, block := range s.Chain {
		response = append(response, *block)
	}
	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func (s *Server) handleNodeAnnouncement(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var announcement Node
	if _, err := asn1.Unmarshal(body, &announcement); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.Nodes = append(s.Nodes, announcement)
	w.WriteHeader(http.StatusAccepted)
	return
}

func (s *Server) handleTransaction(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var announcement Transaction
	if _, err := asn1.Unmarshal(body, &announcement); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s.newTransaction(announcement.Sender, announcement.Recepient, announcement.Amount)
	w.WriteHeader(http.StatusAccepted)
	return
}

func (s *Server) getListHandler(w http.ResponseWriter, r *http.Request) {
	res, err := asn1.Marshal(s.Nodes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(err.Error()))
	}
	w.Write(res)
}

func (s *Server) statusHandler(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
