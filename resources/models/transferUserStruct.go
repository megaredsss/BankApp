package models

type TransferUser struct {
	SenderUser   UserDb
	ReceiverUser UserDb
	Amount       int
}
