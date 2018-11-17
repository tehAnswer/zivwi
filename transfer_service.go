package main

import (
	"encoding/json"
	"fmt"
)

type TransferService interface {
	Perform(fromAccountId string, toAccountId string, amount uint64, msg string) (*Transfer, error)
}

type TransferServiceImpl struct {
	Accounts  AccountGateway
	Transfers TransferGateway
	Queue     Queue
}

func NewTransferService(accounts AccountGateway, transfers TransferGateway, queue Queue) TransferService {
	return &TransferServiceImpl{
		Accounts:  accounts,
		Transfers: transfers,
		Queue:     queue,
	}
}

func (service *TransferServiceImpl) Perform(fromAccountId string, toAccountId string, amount uint64, msg string) (*Transfer, error) {
	if fromAccountId == toAccountId {
		return nil, fmt.Errorf("You can't send money to the same account you're sending from.")
	}

	if amount <= uint64(0) {
		return nil, fmt.Errorf("You need to send at least one cent.")
	}

	if sourceAccount, err := service.Accounts.FindBy(fromAccountId); err == nil {
		if sourceAccount.Balance >= amount {
			transfer := Transfer{
				FromAccountId: fromAccountId,
				ToAccountId:   toAccountId,
				Amount:        amount,
				Message:       msg,
				Status:        "processing",
			}
			transferFromDb, createErr := service.Transfers.Create(transfer)

			if createErr == nil {
				payload, _ := json.Marshal(&transferFromDb)
				queueErr := service.Queue.Publish("transfers", payload)
				if queueErr == nil {
					return transferFromDb, nil
				}
				return nil, queueErr
			}
		}
		return nil, fmt.Errorf("Not enough balance.")
	}
	return nil, fmt.Errorf("Source account not found.")
}
