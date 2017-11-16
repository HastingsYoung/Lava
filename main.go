package main

import (
	"encoding/json"
	"github.com/Lava/blockchain"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var bc *blockchain.Blockchain

var node_identifier = strings.Replace(uuid.NewV4().String(), "-", "", -1)

func main() {
	bc = blockchain.NewBlockchain()
	Handle()
}

func Handle() {
	router := mux.NewRouter()
	router.HandleFunc("/mine", mine).Methods("GET").Name("Mine")
	router.HandleFunc("/transactions/new", new_transaction).Methods("POST").Name("New Transaction")
	router.HandleFunc("/chain", chain).Methods("GET").Name("Chain")
	// launch server
	serv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	serv.ListenAndServe()
}

func mine(w http.ResponseWriter, r *http.Request) {
	last_block := bc.LastBlock()
	last_proof := last_block.Proof
	proof := bc.ProofOfWork(last_proof)

	bc.NewTranX("0", node_identifier, 1)

	block := bc.NewBlock(proof)

	mal, _ := json.Marshal(struct {
		Message       string              `json:"message"`
		Index         int                 `json:"index"`
		Transactions  []*blockchain.TranX `json:"transactions"`
		Proof         int                 `json:"proof"`
		Previous_Hash string              `json:"previous_hash"`
	}{
		Message:       "New block forged",
		Index:         block.Index,
		Transactions:  block.Transactions,
		Proof:         block.Proof,
		Previous_Hash: block.PreviousHash,
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(mal))
}

func new_transaction(w http.ResponseWriter, r *http.Request) {
	var (
		tranx blockchain.TranX
		bytes []byte
		err   error
	)

	bytes, err = ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bytes, &tranx)

	if err != nil || len(tranx.Recipient)*len(tranx.Sender) == 0 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"message": "Missing values"}`)
		return
	}

	index := bc.NewTranX(tranx.Sender, tranx.Recipient, tranx.Amount)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, `{"message": "Transaction will be added to Block `+strconv.Itoa(index)+`"}`)
}

func chain(w http.ResponseWriter, r *http.Request) {
	mal, _ := json.Marshal(struct {
		Chain  []*blockchain.Block `json:"chain"`
		Length int                 `json:"length"`
	}{
		Chain:  bc.Chain(),
		Length: len(bc.Chain()),
	})
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(mal))
}
