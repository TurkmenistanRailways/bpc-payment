package senagat_bank

import (
	"bytes"
	"errors"
	"strings"

	"golang.org/x/net/html"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/util"
)

func (h *SenagatBank) getOtpRequestID(orderId string, form SubmitCardResponse) (string, error) {
	body := bytes.NewBufferString(generateGetOtpRequestIDForm(orderId, form.PaReq, form.TermUrl))
	resp, err := util.Post(form.AcsUrl, body)
	if err != nil {
		return "", errors.Join(err, errors.New("error sending RequestID request"))
	}

	root, err := html.Parse(bytes.NewReader(resp))
	if err != nil {
		return "", errors.Join(err, errors.New("error parsing HTML response"))
	}

	return util.FindRequestId(root), nil
}

func (h *SenagatBank) sendOtp(requestID string) error {
	body := bytes.NewBufferString(generateSendOtpForm(requestID))
	if _, err := util.Post(sendOtpUrl, body); err != nil {
		return errors.Join(err, errors.New("error sending OTP"))
	}

	return nil
}

func (h *SenagatBank) confirmOtp(form banks.ConfirmPaymentRequest) (string, error) {
	body := bytes.NewBufferString(generateConfirmOtpForm(form.RequestID, form.PasswordEdit))
	res, err := util.Post(confirmOtpUrl, body)
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

func (h *SenagatBank) finishPayment(paRes, orderID string) error {
	body := bytes.NewBufferString(generateFinishPaymentForm(orderID, paRes))
	if _, err := util.Post(finishPaymentUrl, body); err != nil {
		return errors.Join(err, errors.New("error finishing payment"))
	}

	return nil
}

// isErrorDiv checks if the node is an error div
func (h *SenagatBank) isErrorDiv(n *html.Node) bool {
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

	return id == "codeErrorContainer" &&
		strings.Contains(class, "row") &&
		strings.Contains(class, "error")
}

// Recursively search for the error div
func (h *SenagatBank) hasOTPError(n *html.Node) bool {
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
