package halk_bank

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

func (h *HalkBank) getOtpRequestID(orderId string, form SubmitCardResponse) (string, error) {
	formData := url.Values{}
	formData.Add("MD", orderId)
	formData.Add("PaReq", form.PaReq)
	formData.Add("TermUrl", form.TermUrl)
	encodedData := formData.Encode()

	resp, err := util.Post(form.AcsUrl, bytes.NewBufferString(encodedData))
	if err != nil {
		return "", errors.Join(err, errors.New("error sending RequestID request"))
	}

	root, err := html.Parse(bytes.NewReader(resp))
	if err != nil {
		return "", errors.Join(err, errors.New("error parsing HTML response"))
	}

	return util.FindRequestId(root), nil
}

func (h *HalkBank) sendOtp(requestID string) error {
	formData := url.Values{}
	formData.Add("request_id", requestID)
	formData.Add("authForm", "authForm")
	formData.Add("sendPasswordButton", "Send password")
	encodedData := formData.Encode()

	if _, err := util.Post(banks.HalkBankOtpUrl, bytes.NewBufferString(encodedData)); err != nil {
		return errors.Join(err, errors.New("error sending OTP"))
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

	res, err := util.Post(banks.HalkBankOtpUrl, bytes.NewBufferString(encodedData))
	if err != nil {
		return "", errors.Join(err, errors.New("error confirming OTP"))
	}

	root, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		return "", errors.Join(err, errors.New("error parsing HTML response"))
	}

	if h.hasOTPError(root) {
		return "", errors.New("error confirming OTP")
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
		return errors.Join(err, errors.New("error finishing payment"))
	}

	return nil
}

// isErrorDiv checks if the node is an error div
func (h *HalkBank) isErrorDiv(n *html.Node) bool {
	if n.Type != html.ElementNode || n.Data != "div" {
		return false
	}

	var id, class string
	for _, attr := range n.Attr {
		switch attr.Key {
		case "id":
			id = attr.Val
		case "class":
			class = attr.Val
		}
	}

	return id == "errorContainer" && strings.Contains(class, "error")
}

// Recursively search for the error div
func (h *HalkBank) hasOTPError(n *html.Node) bool {
	if h.isErrorDiv(n) {
		return true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if h.hasOTPError(c) {
			return true
		}
	}
	return false
}
