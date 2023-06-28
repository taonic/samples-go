package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/temporalio/samples-go/mtls"
	query "github.com/temporalio/samples-go/query-bench"
	"go.temporal.io/sdk/client"
)

const (
	queryType   = "state"
	concurrency = 20
	requests    = 50
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	clientOptions, wid, err := mtls.ParseClientOptionFlags(os.Args[1:])
	if err != nil {
		log.Fatalf("Invalid arguments: %v", err)
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	then := time.Now()
	test(c, wid)
	fmt.Printf("\nTime elapsed: %v\n", time.Since(then))
	fmt.Printf("Total queries: %v\n", concurrency*requests)
	fmt.Printf("RPS: %v\n", concurrency*requests/time.Since(then).Seconds())
}

func test(c client.Client, wid string) {
	var wg sync.WaitGroup
	for i := 1; i <= concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 1; j <= requests; j++ {
				resp, err := c.QueryWorkflow(context.Background(), query.WorkflowIDPrefix+wid, "", queryType)
				if err != nil {
					log.Fatalln("Unable to query workflow", err)
				}
				var result interface{}
				if err := resp.Get(&result); err != nil {
					log.Fatalln("Unable to decode query result", err)
				}
				fmt.Print(".")
			}
		}()
	}

	wg.Wait()
}
