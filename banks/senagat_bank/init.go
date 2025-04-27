package senagat_bank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"

	"golang.org/x/net/html"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/request"
	"github.com/TurkmenistanRailways/bpc-payment/utils"
)

type SenagatBank struct {
	UserName string
	Password string
}

type OrderRegistrationResponse struct {
	OrderId      string `json:"orderId,omitempty"`
	FormUrl      string `json:"formUrl,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	RecurrenceId string `json:"recurrenceId,omitempty"`
}

type SubmitCardResponse struct {
	Info      string `json:"info"`
	AcsUrl    string `json:"acsUrl"`
	PaReq     string `json:"paReq"`
	TermUrl   string `json:"termUrl"`
	ErrorCode int    `json:"errorCode"`
}

func Init(user banks.BankUser) banks.Bank {
	return &SenagatBank{
		UserName: user.Username,
		Password: user.Password,
	}
}

func (h *SenagatBank) OrderRegister(form banks.RegisterForm) (banks.OrderRegistrationResponse, error) {
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
	registerURL := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatRegisterURL, urlParams)

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

func (h *SenagatBank) SubmitCard(form banks.SubmitCard) (string, error) {
	urlParams := utils.StructToURLParams(form)
	fullUrl := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatConfirmPaymentURL, urlParams)

	responseBody, err := request.Post(fullUrl, nil)
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
	if _, err := request.Post(requestUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}

func (h *SenagatBank) ConfirmPayment(form banks.ConfirmPaymentRequest) error {
	paRes, err := h.confirmOtp(form)
	if err != nil {
		return err
	}

	if err = h.finishPayment(paRes, form.MDORDER); err != nil {
		return err
	}
	return nil
}

func (h *SenagatBank) Refund(form banks.RefundRequest) error {
	form.Username = h.UserName
	form.Password = h.Password

	urlParams := utils.StructToURLParams(form)
	fullUrl := fmt.Sprintf("%s%s?%s", banks.SenagatBankBaseUrl, banks.SenagatRefundURL, urlParams)

	if _, err := request.Get(fullUrl); err != nil {
		return err
	}

	return nil
}

func (h *SenagatBank) getOtpRequestID(orderId string, form SubmitCardResponse) (string, error) {
	formData := url.Values{}
	formData.Add("MD", orderId)
	formData.Add("PaReq", form.PaReq)
	formData.Add("TermUrl", form.TermUrl)
	encodedData := formData.Encode()

	resp, err := request.Post(form.AcsUrl, bytes.NewBufferString(encodedData))
	if err != nil {
		return "", err
	}

	root, err := html.Parse(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	return utils.FindRequestId(root), nil
}

func (h *SenagatBank) sendOtp(requestID string) error {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("sendButton", "Ugratmak")
	encodedData := formData.Encode()

	requestUrl := fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatOTPURL)
	if _, err := request.Post(requestUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}

func (h *SenagatBank) confirmOtp(form banks.ConfirmPaymentRequest) (string, error) {
	formData := url.Values{}
	formData.Add("request_id", form.RequestID)
	formData.Add("passwordEdit", form.PasswordEdit)
	formData.Add("submitButton", "Tassyklamak")
	encodedData := formData.Encode()

	fullUrl := fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatOTPURL)

	res, err := request.Post(fullUrl, bytes.NewBufferString(encodedData))
	if err != nil {
		return "", err
	}

	root, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		return "", err
	}

	return utils.FindPaRes(root), nil
}

func (h *SenagatBank) finishPayment(paRes, orderID string) error {
	formData := url.Values{}
	formData.Add("MD", orderID)
	formData.Add("PaRes", paRes)
	encodedData := formData.Encode()
	fullUrl := fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatFinishURL)

	if _, err := request.Post(fullUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}
