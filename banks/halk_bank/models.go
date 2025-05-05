package halk_bank

import "github.com/TurkmenistanRailways/bpc-payment/banks"

type OrderRegistrationResponse struct {
	OrderId      string `json:"orderId,omitempty"`
	FormUrl      string `json:"formUrl,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	RecurrenceId string `json:"recurrenceId,omitempty"`
}

type SubmitCardResponse struct {
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Error        string `json:"error,omitempty"`
	Info         string `json:"info,omitempty"`
	AcsUrl       string `json:"acsUrl,omitempty"`
	PaReq        string `json:"paReq,omitempty"`
	TermUrl      string `json:"termUrl,omitempty"`
}

type OrderStatusRequest struct {
	// bank merchant username
	Username string `json:"userName"`
	// bank merchant password
	Password string `json:"password"`
	// Order id the same as order number
	OrderID string `json:"orderId"`
}

type OrderStatusResponse struct {
	ErrorCode             string               `json:"errorCode,omitempty"`
	ErrorMessage          string               `json:"errorMessage,omitempty"`
	OrderNumber           string               `json:"orderNumber,omitempty"`
	OrderStatus           int                  `json:"orderStatus,omitempty"`
	ActionCode            int                  `json:"actionCode,omitempty"`
	ActionCodeDescription string               `json:"actionCodeDescription,omitempty"`
	Amount                float64              `json:"amount,omitempty"`
	Currency              string               `json:"currency,omitempty"`
	Date                  int64                `json:"date,omitempty"`
	OrderDescription      string               `json:"orderDescription,omitempty"`
	Ip                    string               `json:"ip,omitempty"`
	MerchantOrderParams   []MerchantOrderParam `json:"merchantOrderParams,omitempty"`
	Attributes            []Attribute          `json:"attributes,omitempty"`
	CardAuthInfo          CardAuthInfo         `json:"cardAuthInfo,omitempty"`
	FraudLevel            int                  `json:"fraudLevel,omitempty"`
}

type MerchantOrderParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Attribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CardAuthInfo struct {
	Expiration              string `json:"expiration"`
	CardholderName          string `json:"cardholderName"`
	AuthorizationResponseId string `json:"authorizationResponseId"`
	Pan                     string `json:"pan"`
}

var statusCodes = map[string]banks.OrderStatus{
	"0": banks.OrderStatusNotPaid,
	"2": banks.OrderStatusPaid,
	"6": banks.OrderStatusUnderpaid,
}
