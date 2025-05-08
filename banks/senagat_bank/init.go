package senagat_bank

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

type SenagatBank struct {
	username string
	password string
}

func Init(user banks.BankUser) banks.Bank {
	return &SenagatBank{
		username: user.Username,
		password: user.Password,
	}
}

func (h *SenagatBank) CheckStatus(orderID string) (banks.OrderStatus, error) {
	urlParams := util.StructToURLParams(OrderStatusRequest{
		Username: h.username,
		Password: h.password,
		OrderID:  orderID,
	})

	fullURL := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatOrderStatusURL, urlParams)

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

func (h *SenagatBank) OrderRegister(form banks.RegisterForm) (banks.OrderRegistrationResponse, error) {
	if form.OrderNumber == "" {
		form.OrderNumber = util.GenerateOrderNumber(1, 32)
	}

	requestPayload := banks.OrderRegistrationRequest{
		Username:           h.username,
		Password:           h.password,
		Amount:             form.Amount,
		SessionTimeoutSecs: form.SessionTimeout,
		Language:           form.Language,
		Currency:           banks.CurrencyTMT,
		ReturnURL:          "/", // Consider making this configurable
		OrderNumber:        form.OrderNumber,
	}

	urlParams := util.StructToURLParams(requestPayload)
	registerURL := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatRegisterURL, urlParams)

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

func (h *SenagatBank) SubmitCard(form banks.SubmitCard) (string, error) {
	urlParams := util.StructToURLParams(form)
	fullUrl := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatConfirmPaymentURL, urlParams)

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

func (h *SenagatBank) ResendOtpCode(requestID string) error {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("resendButton", "Kody gaýtadan ugratmak")
	formData.Add("passwordEdit", "")
	encodedData := formData.Encode()

	requestUrl := fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatOTPURL)
	if _, err := util.Post(requestUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}

func (h *SenagatBank) ConfirmPayment(form banks.ConfirmPaymentRequest) error {
	paRes, err := h.confirmOtp(form)
	if err != nil {
		return err
	}

	return h.finishPayment(paRes, form.MDORDER)
}

func (h *SenagatBank) Refund(form banks.RefundRequest) error {
	form.Username = h.username
	form.Password = h.password

	urlParams := util.StructToURLParams(form)
	fullUrl := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatRefundURL, urlParams)

	if _, err := util.Get(fullUrl); err != nil {
		return err
	}

	return nil
}
