package main

import (
	"time"

	"sigs.k8s.io/external-dns/provider"
	"sigs.k8s.io/external-dns/provider/webhook"
)

func main() {

	// instantiate the config
	config := AWSConfig{
		ZoneIDFilter:    provider.NewZoneIDFilter([]string{"external.dns"}),
		BatchChangeSize: 10,
		DryRun:          false,
	}

	// instantiate the aws provider
	awsProvider, err := NewAWSProvider(config)
	if err != nil {
		panic(err)
	}

	startedChan := make(chan struct{})

	go webhook.StartHTTPApi(awsProvider, startedChan, 5*time.Second, 5*time.Second, "127.0.0.1:8888")
	<-startedChan

	time.Sleep(100000 * time.Second)

}
