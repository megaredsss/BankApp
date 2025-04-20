package models

type TransferUser struct {
	SenderUser   UserDb  `binding:"required"`
	ReceiverUser UserDb  `binding:"required"`
	Amount       float64 `binding:"required,number,min=0"`
}
