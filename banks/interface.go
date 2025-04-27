package banks

type Bank interface {
	OrderRegister(form RegisterForm) (OrderRegistrationResponse, error)
	SubmitCard(form SubmitCard) (string, error)
	ResendOtpCode(request string) error
}
