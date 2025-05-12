package senagat_bank

import (
	"fmt"
	"net/url"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
)

var (
	confirmOtpUrl    = fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatOTPURL)
	finishPaymentUrl = fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatFinishURL)
	sendOtpUrl       = fmt.Sprintf("%s%s", banks.SenagatBankBaseUrl, banks.SenagatOTPURL)
)

func generateConfirmOtpForm(requestId, passwordEdit string) string {
	formData := url.Values{}
	formData.Add("request_id", requestId)
	formData.Add("passwordEdit", passwordEdit)
	formData.Add("submitButton", "Tassyklamak")
	return formData.Encode()
}

func generateFinishPaymentForm(orderID, paRes string) string {
	formData := url.Values{}
	formData.Add("MD", orderID)
	formData.Add("PaRes", paRes)

	return formData.Encode()
}

func generateGetOtpRequestIDForm(orderId, paReq, termUrl string) string {
	formData := url.Values{}
	formData.Add("MD", orderId)
	formData.Add("PaReq", paReq)
	formData.Add("TermUrl", termUrl)

	return formData.Encode()
}

func generateSendOtpForm(requestID string) string {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("sendButton", "Ugratmak")

	return formData.Encode()
}
