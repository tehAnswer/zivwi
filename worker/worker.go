package worker

import (
	"encoding/json"
	"fmt"

	nsq "github.com/bitly/go-nsq"
	app "github.com/tehAnswer/zivwi"
)

type Worker interface {
	HandleMessage(message *nsq.Message) error
}

type TransferWorker struct {
	Accounts        app.AccountGateway
	Transfers       app.TransferGateway
	TransferService app.TransferService
}

func NewTransferWorker(accounts app.AccountGateway, transfers app.TransferGateway, srv app.TransferService) Worker {
	return &TransferWorker{
		Accounts:        accounts,
		Transfers:       transfers,
		TransferService: srv,
	}
}

func (worker *TransferWorker) HandleMessage(message *nsq.Message) error {
	if len(message.Body) == 0 {
		return fmt.Errorf("Blank body, omitting event.")
	}

	var transfer app.Transfer
	json.Unmarshal(message.Body, &transfer)

	fromAccount, _ := worker.Accounts.FindBy(transfer.FromAccountId)
	_, notFoundDestErr := worker.Accounts.FindBy(transfer.ToAccountId)

	if notFoundDestErr != nil {
		transfer.Error = "to_account_not_found"
		transfer.Status = "cancelled"
		worker.Transfers.Update(transfer)
		return notFoundDestErr
	}

	if fromAccount.Balance < transfer.Amount {
		transfer.Error = "not_enough_funds"
		transfer.Status = "cancelled"
		worker.Transfers.Update(transfer)
		return fmt.Errorf("Not enough balance to complete transfer (Ref: %v)", transfer.Id)
	}

	return worker.TransferService.Resolve(transfer)
}
