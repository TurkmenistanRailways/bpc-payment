package rysgal_bank

type OrderRegistrationResponse struct {
	OrderId      string `json:"orderId,omitempty"`
	FormUrl      string `json:"formUrl,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	RecurrenceId string `json:"recurrenceId,omitempty"`
}
