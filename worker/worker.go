package worker

import (
	"encoding/json"
	"fmt"

	nsq "github.com/bitly/go-nsq"
	app "github.com/tehAnswer/zivwi"
)

func Run() {
	fmt.Println("Starting worker...")
	appCtx := app.NewAppCtx()
	consumer, _ := nsq.NewConsumer("transfers", "ch", nsq.NewConfig())
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		var transfer app.Transfer
		json.Unmarshal(message.Body, &transfer)

		fromAccount, _ := appCtx.Accounts.FindBy(transfer.FromAccountId)
		_, notFoundDestErr := appCtx.Accounts.FindBy(transfer.ToAccountId)

		if notFoundDestErr != nil {
			_, err := appCtx.Database.Connection.Query(`UPDATE transfers SET status = $2, explanation = $3 WHERE id = $1`,
				transfer.Id, "failed", "Dest. account does not exist")
			return err
		}

		if fromAccount.Balance < transfer.Amount {
			_, err := appCtx.Database.Connection.Query(`UPDATE transfers SET status = $2, explanation = $3 WHERE id = $1`,
				transfer.Id, "failed", "Not enough funds")
			return err
		}

		return appCtx.TransferService.Resolve(transfer)
	}))

	for {
		select {
		case <-consumer.StopChan:
			return
		}
	}
}
