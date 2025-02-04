package main

import (
	"context"
	"fmt"
	"os"
	"time"

	kp "github.com/IBM/keyprotect-go-client"
	progressbar "github.com/schollz/progressbar/v3"
)

func main() {

	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) != 3 {
		panic("Usage: count-key-versions [end point, for example: [https://us-south.kms.cloud.ibm.com] [api key] [kp instance id]")
	}

	cc := kp.ClientConfig{
		BaseURL:    argsWithoutProg[0],
		APIKey:     argsWithoutProg[1],
		InstanceID: argsWithoutProg[2],
	}

	// Build a new client from the config
	client, err := kp.New(cc, kp.DefaultTransport())

	if err != nil {
		panic(err)
	}

	keys, err := client.ListKeys(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	numKeys := keys.Metadata.NumberOfKeys
	fmt.Printf("Total number of non deleted keys in instance: %s: %d\n", cc.InstanceID, numKeys)

	tc := true
	opts := kp.ListKeyVersionsOptions{
		TotalCount: &tc,
	}

	fmt.Printf("Starting to count and sum up all key versions for all non deleted keys...\n")

	bar := progressbar.NewOptions64(
		int64(numKeys),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionSetItsString("keys"),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			fmt.Printf("\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)

	//for each key, grab its number of keyversions and accumulate total sum
	countSum := 0
	for _, key := range keys.Keys {

		keyVersion, err := client.ListKeyVersions(context.Background(), key.ID, &opts)

		if keyVersion.Metadata.TotalCount != nil {
			countSum = countSum + int(*keyVersion.Metadata.TotalCount)
		}

		bar.Add(1)

		if err != nil {
			panic(err)
		}
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("The total number of all key versions for all non deleted keys is %d\n", countSum)
}
