package halk_bank

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

type HalkBank struct {
	username string
	password string
}

// CheckStatus implements banks.Bank.
func (h *HalkBank) CheckStatus(orderID string) (banks.OrderStatus, error) {
	urlParams := util.StructToURLParams(OrderStatusRequest{
		Username: h.username,
		Password: h.password,
		OrderID:  orderID,
	})
	fullURL := fmt.Sprintf("%s%s?%s", banks.HalkBankBaseUrl, banks.HalkBankOrderStatusURL, urlParams)

	res, err := util.Post(fullURL, nil)
	if err != nil {
		return banks.OrderStatusError, err
	}

	var response OrderStatusResponse
	if err = json.Unmarshal(res, &response); err != nil {
		return banks.OrderStatusError, err
	}

	if status, ok := statusCodes[response.ErrorCode]; ok {
		return status, nil
	}

	return banks.OrderStatusError, errors.New("invalid status code")
}

// ConfirmPayment implements banks.Bank.
func (h *HalkBank) ConfirmPayment(form banks.ConfirmPaymentRequest) error {
	paRes, err := h.confirmOtp(form)
	if err != nil {
		return err
	}

	return h.finishPayment(paRes, form.MDORDER)
}

// OrderRegister implements banks.Bank.
func (h *HalkBank) OrderRegister(form banks.RegisterForm) (banks.OrderRegistrationResponse, error) {
	requestPayload := banks.OrderRegistrationRequest{
		Username:           h.username,
		Password:           h.password,
		Amount:             form.Amount,
		SessionTimeoutSecs: form.SessionTimeout,
		Language:           form.Language,
		Currency:           banks.CurrencyTMT,
		ReturnURL:          "/", // Consider making this configurable
		OrderNumber:        util.GenerateOrderNumber(1, 32),
	}

	urlParams := util.StructToURLParams(requestPayload)
	registerUrl := fmt.Sprintf("%s%s?%s", banks.HalkBankBaseUrl, banks.HalkBankRegisterURL, urlParams)

	responseBody, err := util.Post(registerUrl, nil)
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

// Refund implements banks.Bank.
func (h *HalkBank) Refund(form banks.RefundRequest) error {
	form.Username = h.username
	form.Password = h.password

	urlParams := util.StructToURLParams(form)
	fullUrl := fmt.Sprintf("%s%s?%s", banks.HalkBankBaseUrl, banks.HalkBankRefundURL, urlParams)

	if _, err := util.Get(fullUrl); err != nil {
		return err
	}
	return nil
}

// ResendOtpCode implements banks.Bank.
func (h *HalkBank) ResendOtpCode(requestId string) error {
	return errors.New("not exist")
}

// SubmitCard implements banks.Bank.
func (h *HalkBank) SubmitCard(form banks.SubmitCard) (string, error) {
	urlParams := util.StructToURLParams(form)
	fullUrl := fmt.Sprintf("%s%s?%s", banks.HalkBankBaseUrl, banks.HalkBankConfirmPaymentURL, urlParams)

	responseBody, err := util.Post(fullUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to register order: %w", err)
	}

	var response SubmitCardResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", err
	}

	requestID, err := h.getOtpRequestID(form.PAN, response)
	if err != nil {
		return "", err
	}

	return requestID, h.sendOtp(requestID)
}

func Init(user banks.BankUser) banks.Bank {
	return &HalkBank{
		username: user.Username,
		password: user.Password,
	}
}
