package halk_bank

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
	"golang.org/x/net/html"
)

func (h *HalkBank) getOtpRequestID(orderId string, form SubmitCardResponse) (string, error) {
	formData := url.Values{}
	formData.Add("MD", orderId)
	formData.Add("PaReq", form.PaReq)
	formData.Add("TermUrl", form.TermUrl)
	encodedData := formData.Encode()

	resp, err := util.Post(form.AcsUrl, bytes.NewBufferString(encodedData))
	if err != nil {
		return "", err
	}

	root, err := html.Parse(bytes.NewReader(resp))
	if err != nil {
		return "", err
	}

	return util.FindRequestId(root), nil
}

func (h *HalkBank) sendOtp(requestID string) error {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("authForm", "authForm")
	formData.Add("sendPasswordButton", "Send password")
	encodedData := formData.Encode()

	requestUrl := fmt.Sprintf("https://%s", banks.HalkBankOtpUrl)
	if _, err := util.Post(requestUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}

func (h *HalkBank) confirmOtp(form banks.ConfirmPaymentRequest) (string, error) {
	formData := url.Values{}
	formData.Add("authForm", "authForm")
	formData.Add("request_id", form.RequestID)
	formData.Add("pwdInputVisible", form.PasswordEdit)
	formData.Add("submitPasswordButton", "Submit")
	encodedData := formData.Encode()

	fullUrl := fmt.Sprintf("https://%s", banks.HalkBankOtpUrl)

	res, err := util.Post(fullUrl, bytes.NewBufferString(encodedData))
	if err != nil {
		return "", err
	}

	root, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		return "", err
	}

	return util.FindPaRes(root), nil
}

func (h *HalkBank) finishPayment(paRes, orderID string) error {
	formData := url.Values{}
	formData.Add("MD", orderID)
	formData.Add("PaRes", paRes)
	encodedData := formData.Encode()
	fullUrl := fmt.Sprintf("%s%s", banks.HalkBankBaseUrl, banks.HalkBankFinishURL)

	if _, err := util.Post(fullUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}
