package app

type AppCtx struct {
	Accounts         AccountGateway
	AuthorizeService AuthorizeService
	Transfers        TransferGateway
	TransferService  TransferService
	Users            UserGateway
	Queue            Queue
	Database         *Database
}

func NewAppCtx() AppCtx {
	// Persistance
	database := NewDatabase()
	accounts := NewAccountGateway(database)
	users := NewUserGateway(database)
	transfers := NewTransferGateway(database)

	// Background processing
	queue := NewQueue()

	// Logic
	authService := NewAuthorizeService(users)
	transferService := NewTransferService(accounts, transfers, queue)

	return AppCtx{
		Accounts:         accounts,
		AuthorizeService: authService,
		Transfers:        transfers,
		TransferService:  transferService,
		Users:            users,
		Queue:            queue,
	}
}
