package rysgal_bank

import (
	"encoding/json"
	"fmt"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/request"
	"github.com/TurkmenistanRailways/bpc-payment/utils"
)

type RysgalBank struct {
	UserName string
	Password string
}

type OrderRegistrationResponse struct {
	OrderId      string `json:"orderId,omitempty"`
	FormUrl      string `json:"formUrl,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	RecurrenceId string `json:"recurrenceId,omitempty"`
}

func Init(form banks.BankUser) banks.Bank {
	return &RysgalBank{
		UserName: form.Username,
		Password: form.Password,
	}
}

func (h *RysgalBank) OrderRegister(form banks.RegisterForm) (banks.OrderRegistrationResponse, error) {
	requestPayload := banks.OrderRegistrationRequest{
		Username:           h.UserName,
		Password:           h.Password,
		Amount:             form.Amount,
		SessionTimeoutSecs: form.SessionTimeout,
		Language:           form.Language,
		Currency:           banks.CurrencyTMT,
		ReturnURL:          "/", // Consider making this configurable
		OrderNumber:        utils.GenerateOrderNumber(1, 32),
	}

	urlParams := utils.StructToURLParams(requestPayload)
	registerURL := fmt.Sprintf("%s%s?%s", banks.RysgalBankBaseUrl, banks.RysgalRegisterURL, urlParams)

	responseBody, err := request.Post(registerURL, nil)
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
	panic("implement me")
}

func (h *RysgalBank) ResendOtpCode(request string) error {
	panic("implement me")
}
