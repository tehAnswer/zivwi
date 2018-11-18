package worker

import (
	"fmt"
	"log"
	"os"

	nsq "github.com/bitly/go-nsq"
	app "github.com/tehAnswer/zivwi"
)

func Run() {
	fmt.Println("Starting worker...")
	appCtx := app.NewAppCtx()
	config := nsq.NewConfig()
	worker := NewTransferWorker(appCtx.Accounts, appCtx.Transfers, appCtx.TransferService)
	consumer, _ := nsq.NewConsumer("transfers", "process", config)

	consumer.ChangeMaxInFlight(1000)
	consumer.AddConcurrentHandlers(worker, 100)

	nsqlds := []string{os.Getenv("NSQ_URL")}
	if err := consumer.ConnectToNSQLookupds(nsqlds); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-consumer.StopChan:
			return
		}
	}
}
