package halk_bank

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
)

var (
	finishPaymentUrl = fmt.Sprintf("%s%s", banks.HalkBankBaseUrl, banks.HalkBankFinishURL)
)

func generateConfirmOtpForm(requestId, passwordEdit string) *bytes.Buffer {
	formData := url.Values{}
	formData.Add("authForm", "authForm")
	formData.Add("request_id", requestId)
	formData.Add("pwdInputVisible", passwordEdit)
	formData.Add("submitPasswordButton", "Submit")

	return bytes.NewBufferString(formData.Encode())
}

func generateFinishPaymentForm(orderID, paRes string) *bytes.Buffer {
	formData := url.Values{}
	formData.Add("MD", orderID)
	formData.Add("PaRes", paRes)

	return bytes.NewBufferString(formData.Encode())
}

func generateGetOtpRequestIDForm(orderId, paReq, termUrl string) *bytes.Buffer {
	formData := url.Values{}
	formData.Add("MD", orderId)
	formData.Add("PaReq", paReq)
	formData.Add("TermUrl", termUrl)

	return bytes.NewBufferString(formData.Encode())
}

func generateSendOtpForm(requestID string) *bytes.Buffer {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("authForm", "authForm")
	formData.Add("sendPasswordButton", "Send password")

	return bytes.NewBufferString(formData.Encode())
}
