package banks

type Bank interface {
	OrderRegister(form RegisterForm) (OrderRegistrationResponse, error)
	ConfirmPayment(form ConfirmPaymentRequest) error
	SubmitCard(form SubmitCard) (string, error)
	ResendOtpCode(request string) error
	Refund(form RefundRequest) error
}
