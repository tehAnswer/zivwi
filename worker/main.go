package worker

import (
	"fmt"
	"log"

	nsq "github.com/bitly/go-nsq"
	app "github.com/tehAnswer/zivwi"
)

func Run() {
	fmt.Println("Starting worker...")
	appCtx := app.NewAppCtx()
	worker := NewTransferWorker(appCtx.Accounts, appCtx.Transfers, appCtx.TransferService)
	consumer, _ := nsq.NewConsumer("transfers", "process", nsq.NewConfig())

	consumer.ChangeMaxInFlight(1000)
	consumer.AddConcurrentHandlers(worker, 100)

	nsqlds := []string{"127.0.0.1:4161"}
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
