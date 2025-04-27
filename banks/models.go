package banks

type RegisterForm struct {
	Amount         int64
	SessionTimeout int
	Language       string
}

type BankUser struct {
	Username string
	Password string
}

type OrderRegistrationRequest struct {
	ApiClient          string `json:"api_client" binding:"required"`
	Username           string `json:"userName" binding:"required"`
	Password           string `json:"password" binding:"required"`
	OrderNumber        string `json:"orderNumber" minLength:"1" maxLength:"32"`
	Amount             int64  `json:"amount" minLength:"1" maxLength:"19" binding:"required"`
	Currency           string `json:"currency" minLength:"3" maxLength:"3"`
	ReturnURL          string `json:"returnUrl"`
	Language           string `json:"language,omitempty" minLength:"2" maxLength:"2"`
	SessionTimeoutSecs int    `json:"sessionTimeoutSecs,omitempty" maxLength:"9"`
}

type OrderRegistrationResponse struct {
	FormUrl string `json:"formUrl"`
	OrderId string `json:"orderId"`
}

type SubmitCard struct {
	// Another meaning of order_id
	MDORDER string `json:"MDORDER"`

	// should be in form like YYYYMM example 202709
	EXPIRY string `json:"$EXPIRY"`

	// PAN is a number on card 12-digit code
	PAN string `json:"$PAN"`

	// Text is cardholder name and surname
	TEXT string `json:"TEXT"`

	// CVC is secure code on back side of your card
	CVC string `json:"$CVC"`

	// language is language
	Language string `json:"language"`
}
