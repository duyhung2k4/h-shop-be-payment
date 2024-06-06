package request

type CreateBillOnlineRequest struct {
	Amount           int    `json:"amount"`
	BankCode         string `json:"bankCode"`
	OrderDescription string `json:"orderDescription"`
	OrderType        string `json:"orderType"`
}
