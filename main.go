package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
)

type AWSPlugin struct {
	provider *AWSProvider
}

type PropertyValuesEqualsRequest struct {
	Name     string `json:"name"`
	Previous string `json:"previous"`
	Current  string `json:"current"`
}

type PropertiesValuesEqualsResponse struct {
	Equals bool `json:"equals"`
}

func (p *AWSPlugin) awsProviderHandler(w http.ResponseWriter, req *http.Request) {
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

func (p *AWSPlugin) propertyValuesEquals(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet { // propertyValuesEquals
		pve := PropertyValuesEqualsRequest{}
		if err := json.NewDecoder(req.Body).Decode(&pve); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		b := p.provider.PropertyValuesEqual(pve.Name, pve.Previous, pve.Current)
		r := PropertiesValuesEqualsResponse{
			Equals: b,
		}
		out, _ := json.Marshal(&r)
		w.Write(out)
	}

}

func (p *AWSPlugin) adjustEndpoints(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet { // propertyValuesEquals
		pve := []*endpoint.Endpoint{}
		if err := json.NewDecoder(req.Body).Decode(&pve); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pve = p.provider.AdjustEndpoints(pve)
		out, _ := json.Marshal(&pve)
		w.Write(out)
	}

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
	m.HandleFunc("/records", p.awsProviderHandler)
	m.HandleFunc("/propertyvaluesequals", p.propertyValuesEquals)
	m.HandleFunc("/adjustendpoints", p.adjustEndpoints)
	if err := http.ListenAndServe(":8888", m); err != nil {
		log.Fatal(err)
	}
}
