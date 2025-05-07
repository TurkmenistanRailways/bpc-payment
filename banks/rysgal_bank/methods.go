package rysgal_bank

import (
	"bytes"
	"fmt"
	"net/url"

	"golang.org/x/net/html"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

func (h *RysgalBank) getOtpRequestID(orderId string, form SubmitCardResponse) (string, error) {
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

func (h *RysgalBank) sendOtp(requestID string) error {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("authForm", "authForm")
	formData.Add("sendPasswordButton", "Send password")
	encodedData := formData.Encode()

	requestUrl := fmt.Sprintf("%s%s", banks.RysgalBankBaseUrl, banks.RysgalBankOtpUrl)
	if _, err := util.Post(requestUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}

func (h *RysgalBank) confirmOtp(form banks.ConfirmPaymentRequest) (string, error) {
	formData := url.Values{}
	formData.Add("authForm", "authForm")
	formData.Add("request_id", form.RequestID)
	formData.Add("pwdInputVisible", form.PasswordEdit)
	formData.Add("submitPasswordButton", "Submit")
	encodedData := formData.Encode()

	fullUrl := fmt.Sprintf("%s%s", banks.RysgalBankBaseUrl, banks.RysgalBankOtpUrl)

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

func (h *RysgalBank) finishPayment(paRes, orderID string) error {
	formData := url.Values{}
	formData.Add("MD", orderID)
	formData.Add("PaRes", paRes)
	encodedData := formData.Encode()
	fullUrl := fmt.Sprintf("%s%s", banks.RysgalBankBaseUrl, banks.SenagatFinishURL)

	if _, err := util.Post(fullUrl, bytes.NewBufferString(encodedData)); err != nil {
		return err
	}

	return nil
}
