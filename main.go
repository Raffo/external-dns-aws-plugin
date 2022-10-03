package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"sigs.k8s.io/external-dns/plan"
)

type AWSPlugin struct {
	provider *AWSProvider
}

func (p *AWSPlugin) awsProviderHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("much request")
	if req.Method == http.MethodGet { // records
		log.Println("get records")
		records, err := p.provider.Records(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(records)
		return
	} else if req.Method == http.MethodPost { // applychanges
		log.Println("post applychanges")
		// extract changes from the request body
		var changes plan.Changes
		if err := json.NewDecoder(req.Body).Decode(&changes); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		p.provider.ApplyChanges(context.Background(), &changes)

		err := p.provider.ApplyChanges(context.Background(), &changes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	log.Println("this should never happen")
}

func main() {
	// instantiate the config
	config := AWSConfig{} // TODO populate the config

	// instantiate the aws provider
	awsProvider, err := NewAWSProvider(config)
	if err != nil {
		panic(err)
	}

	p := AWSPlugin{
		provider: awsProvider,
	}

	m := http.NewServeMux()
	m.HandleFunc("/", p.awsProviderHandler)
	if err := http.ListenAndServe(":8888", m); err != nil {
		log.Fatal(err)
	}
}
