package main

import (
	"encoding/json"
	"fmt"
)

type TransferService interface {
	Pay(fromAccountId string, toAccountId string, amount uint64) error
}

type TransferServiceImpl struct {
	Accounts AccountGateway
	Queue    Queue
}

func NewTransferService(accounts AccountGateway, queue Queue) TransferService {
	return &TransferServiceImpl{
		Accounts: accounts,
		Queue:    queue,
	}
}

func (service *TransferServiceImpl) Pay(fromAccountId string, toAccountId string, amount uint64) error {
	if sourceAccount, err := service.Accounts.FindBy(fromAccountId); err == nil {
		if sourceAccount.Balance > amount {
			payload, _ := json.Marshal(struct {
				fromAccountId string `json:from_account_id`
				toAccountId   string `json:to_account_id`
				amount        uint64 `json:amount`
			}{fromAccountId, toAccountId, amount})

			return service.Queue.Publish("transfers", payload)
		}
		return fmt.Errorf("Not enough balance.")
	}
	return fmt.Errorf("Account not found.")
}
