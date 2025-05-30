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
	fullURL := fmt.Sprintf(banks.URLFormat, banks.HalkBankBaseUrl, banks.HalkBankOrderStatusURL, urlParams)

	res, err := util.Post(fullURL, nil)
	if err != nil {
		return banks.OrderStatusError, errors.Join(err, errors.New("error checking order status"))
	}

	var response OrderStatusResponse
	if err = json.Unmarshal(res, &response); err != nil {
		return banks.OrderStatusError, errors.Join(err, errors.New("error unmarshalling order status response"))
	}

	if response.ActionCode != 0 {
		return banks.OrderStatusError, errors.Join(err, fmt.Errorf("error code: %d, description: %s", response.ActionCode, response.ActionCodeDescription))
	}

	if status, ok := banks.StatusCodes[response.OrderStatus]; ok {
		return status, nil
	}

	return banks.OrderStatusError, errors.New("unknown error code")
}

// ConfirmPayment implements banks.Bank.
func (h *HalkBank) ConfirmPayment(form banks.ConfirmPaymentRequest) error {
	paRes, err := h.confirmOtp(form)
	if err != nil {
		return errors.Join(err, errors.New("error confirming OTP"))
	}

	return h.finishPayment(paRes, form.MDORDER)
}

// OrderRegister implements banks.Bank.
func (h *HalkBank) OrderRegister(form banks.RegisterForm) (banks.OrderRegistrationResponse, error) {
	if form.OrderNumber == "" {
		form.OrderNumber = util.GenerateOrderNumber(1, 32)
	}

	if form.ReturnURL == "" {
		form.ReturnURL = "/"
	}

	requestPayload := banks.OrderRegistrationRequest{
		Username:           h.username,
		Password:           h.password,
		Amount:             form.Amount,
		SessionTimeoutSecs: form.SessionTimeout,
		Language:           form.Language,
		Currency:           banks.CurrencyTMT,
		ReturnURL:          form.ReturnURL,
		OrderNumber:        form.OrderNumber,
	}

	urlParams := util.StructToURLParams(requestPayload)
	registerUrl := fmt.Sprintf(banks.URLFormat, banks.HalkBankBaseUrl, banks.HalkBankRegisterURL, urlParams)

	responseBody, err := util.Post(registerUrl, nil)
	if err != nil {
		return banks.OrderRegistrationResponse{}, errors.Join(err, errors.New("error registering order"))
	}

	var orderRegistrationResponse OrderRegistrationResponse
	if err = json.Unmarshal(responseBody, &orderRegistrationResponse); err != nil {
		return banks.OrderRegistrationResponse{}, errors.Join(err, errors.New("error unmarshalling order registration response"))
	}

	if orderRegistrationResponse.ErrorMessage != "" {
		return banks.OrderRegistrationResponse{}, errors.New(orderRegistrationResponse.ErrorMessage)
	}

	return banks.OrderRegistrationResponse{
		OrderId: orderRegistrationResponse.OrderId,
		FormUrl: orderRegistrationResponse.FormUrl,
	}, nil
}

// Refund implements banks.Bank.
func (h *HalkBank) Refund(form banks.RefundRequest) error {
	urlParams := util.StructToURLParams(banks.Refund{
		Username: h.username,
		Password: h.password,
		Amount:   form.Amount,
		OrderID:  form.OrderID,
	})

	fullUrl := fmt.Sprintf(banks.URLFormat, banks.HalkBankBaseUrl, banks.HalkBankRefundURL, urlParams)
	if _, err := util.Get(fullUrl); err != nil {
		return errors.Join(err, errors.New("error refunding order"))
	}
	return nil
}

// ResendOtpCode implements banks.Bank.
func (h *HalkBank) ResendOtpCode(requestId string) error {
	return errors.New("not exist")
}

// SubmitCard implements banks.Bank.
func (h *HalkBank) SubmitCard(form banks.SubmitCard) (string, error) {
	if ok := util.IsValidExpiry(form.EXPIRY); !ok {
		return "", errors.New("invalid expiry date")
	}

	if ok := util.IsValidPAN(form.PAN); !ok {
		return "", errors.New("invalid PAN")
	}

	urlParams := util.StructToURLParams(form)
	fullUrl := fmt.Sprintf(banks.URLFormat, banks.HalkBankBaseUrl, banks.HalkBankConfirmPaymentURL, urlParams)

	responseBody, err := util.Post(fullUrl, nil)
	if err != nil {
		return "", errors.Join(err, errors.New("error submitting card data"))
	}

	var response SubmitCardResponse
	if err = json.Unmarshal(responseBody, &response); err != nil {
		return "", errors.Join(err, errors.New("error unmarshalling card data response"))
	}

	requestID, err := h.getOtpRequestID(form.PAN, response)
	if err != nil {
		return "", errors.Join(err, errors.New("error getting OTP request ID"))
	}

	return requestID, h.sendOtp(requestID)
}

func Init(user banks.BankUser) banks.Bank {
	return &HalkBank{
		username: user.Username,
		password: user.Password,
	}
}
