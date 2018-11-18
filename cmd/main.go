package main

import (
	"os"

	app "github.com/tehAnswer/zivwi"
	worker "github.com/tehAnswer/zivwi/worker"
)

func main() {
	concept := os.Args[1]
	if concept == "worker" {
		worker.Run()
	} else if concept == "seeds" {
		if os.Getenv("APP_ENV") == "production" {
			panic("Seeding the database in PROD would remove valuable data.")
		}

		appCtx := app.NewAppCtx()

		// Clean up
		appCtx.Accounts.DeleteAll()
		appCtx.Users.DeleteAll()
		appCtx.Transfers.DeleteAll()

		// Accounts
		account1, _ := appCtx.Accounts.Create(app.Account{Balance: 1.2e10})
		account2, _ := appCtx.Accounts.Create(app.Account{Balance: 5.6e15})
		account3, _ := appCtx.Accounts.Create(app.Account{Balance: 1.6e10})

		// Users
		appCtx.Users.Create(app.User{
			FirstName:  "Amancio",
			LastName:   "Ortega",
			Email:      "jefe@inditex.es",
			Password:   "dameDiner0",
			AccountIds: []string{account1.Id, account2.Id},
		})

		appCtx.Users.Create(app.User{
			FirstName:  "John",
			LastName:   "Heitinga",
			Email:      "hoofd@ajax.nl",
			Password:   "lekkerlekker",
			AccountIds: []string{account3.Id},
		})

		// Transfers
		appCtx.Transfers.Create(app.Transfer{
			FromAccountId: "",
			ToAccountId:   account1.Id,
			Status:        "completed",
			Message:       "Initial funding.",
			Amount:        account1.Balance,
		})
		appCtx.Transfers.Create(app.Transfer{
			FromAccountId: "",
			ToAccountId:   account2.Id,
			Status:        "completed",
			Message:       "Initial funding.",
			Amount:        account2.Balance,
		})
		appCtx.Transfers.Create(app.Transfer{
			FromAccountId: "",
			ToAccountId:   account3.Id,
			Status:        "completed",
			Message:       "Initial funding.",
			Amount:        account3.Balance,
		})

	} else {
		app.Run()
	}
}
