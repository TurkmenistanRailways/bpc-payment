package rysgal_bank

import (
	"encoding/json"
	"fmt"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

type RysgalBank struct {
	userName string
	password string
}

func Init(user banks.BankUser) banks.Bank {
	return &RysgalBank{
		userName: user.Username,
		password: user.Password,
	}
}

func (h *RysgalBank) CheckStatus(orderID string) (banks.OrderStatus, error) {
	return banks.OrderStatusError, nil
}

func (h *RysgalBank) OrderRegister(form banks.RegisterForm) (banks.OrderRegistrationResponse, error) {
	requestPayload := banks.OrderRegistrationRequest{
		Username:           h.userName,
		Password:           h.password,
		Amount:             form.Amount,
		SessionTimeoutSecs: form.SessionTimeout,
		Language:           form.Language,
		Currency:           banks.CurrencyTMT,
		ReturnURL:          "/", // Consider making this configurable
		OrderNumber:        util.GenerateOrderNumber(1, 32),
	}

	urlParams := util.StructToURLParams(requestPayload)
	registerURL := fmt.Sprintf("%s%s?%s", banks.RysgalBankBaseUrl, banks.RysgalRegisterURL, urlParams)

	responseBody, err := util.Post(registerURL, nil)
	if err != nil {
		return banks.OrderRegistrationResponse{}, fmt.Errorf("failed to register order: %w", err)
	}

	var orderRegistrationResponse OrderRegistrationResponse
	if err = json.Unmarshal(responseBody, &orderRegistrationResponse); err != nil {
		return banks.OrderRegistrationResponse{}, fmt.Errorf("failed to unmarshal orderRegistrationResponse: %w", err)
	}

	return banks.OrderRegistrationResponse{
		OrderId: orderRegistrationResponse.OrderId,
		FormUrl: orderRegistrationResponse.FormUrl,
	}, nil
}

func (h *RysgalBank) SubmitCard(form banks.SubmitCard) (string, error) {
	return "", nil
}

func (h *RysgalBank) ResendOtpCode(requestID string) error {
	return nil
}

func (h *RysgalBank) ConfirmPayment(form banks.ConfirmPaymentRequest) error {
	return nil
}

func (h *RysgalBank) Refund(form banks.RefundRequest) error {

	return nil
}
