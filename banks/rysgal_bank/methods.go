package rysgal_bank

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/net/html"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

func (h *RysgalBank) getOtpRequestID(orderId string, form SubmitCardResponse) (string, error) {
	resp, err := util.Post(form.AcsUrl, generateGetOtpRequestIDForm(orderId, form.PaReq, form.TermUrl))
	if err != nil {
		return "", errors.Join(err, errors.New("error sending RequestID request"))
	}

	root, err := html.Parse(bytes.NewReader(resp))
	if err != nil {
		return "", errors.Join(err, errors.New("error parsing HTML response"))
	}

	return util.FindRequestId(root), nil
}

func (h *RysgalBank) sendOtp(requestID string) error {
	if _, err := util.Post(sendOtpUrl, generateSendOtpForm(requestID)); err != nil {
		return errors.Join(err, errors.New("error sending OTP"))
	}

	return nil
}

func (h *RysgalBank) confirmOtp(form banks.ConfirmPaymentRequest) (string, error) {
	res, err := util.Post(confirmOtpUrl, generateConfirmOtpForm(form.RequestID, form.PasswordEdit))
	if err != nil {
		return "", errors.Join(err, errors.New("error confirming OTP"))
	}

	root, err := html.Parse(bytes.NewReader(res))
	if err != nil {
		return "", errors.Join(err, errors.New("error parsing HTML response"))
	}

	if h.hasOTPError(root) {
		return "", errors.New("error: OTP code is incorrect")
	}

	return util.FindPaRes(root), nil
}

func (h *RysgalBank) finishPayment(paRes, orderID string) error {
	if _, err := util.Post(finishPaymentUrl, generateFinishPaymentForm(orderID, paRes)); err != nil {
		return errors.Join(err, errors.New("error finishing payment"))
	}

	return nil
}

// isErrorDiv checks if the node is an error div
func (h *RysgalBank) isErrorDiv(n *html.Node) bool {
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
func (h *RysgalBank) hasOTPError(n *html.Node) bool {
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
