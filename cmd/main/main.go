package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	bpc2 "github.com/TurkmenistanRailways/bpc-payment"
	"github.com/TurkmenistanRailways/bpc-payment/banks"
)

const (
	SenagatBankProfile string = "SenagatBank"
	Amount                    = 100
)

func SubmitForms() {
	bpc := bpc2.New()

	if err := bpc.AddProfile(SenagatBankProfile, bpc2.SenagatBank, banks.BankUser{Username: "demir_yollary", Password: "demir_yollary1"}); err != nil {
		log.Fatal(err)
	}

	senagatBankResponse, err := bpc.OrderRegister(SenagatBankProfile, 100, Amount)
	if err != nil {
		log.Fatal(err)
	}

	formRequest := banks.SubmitCard{
		MDORDER:  senagatBankResponse.OrderId,
		EXPIRY:   "202703",
		PAN:      "9934701839255812",
		TEXT:     "ABDYLOW EZIZ",
		CVC:      "913",
		Language: "ru",
	}

	requestId, err := bpc.Submit(SenagatBankProfile, formRequest)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Write OTP code: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	confirmPaymentRequest := banks.ConfirmPaymentRequest{
		MDORDER:      senagatBankResponse.OrderId,
		RequestID:    requestId,
		PasswordEdit: text,
	}

	if err = bpc.ConfirmPayment(SenagatBankProfile, confirmPaymentRequest); err != nil {
		log.Fatal(err)
	}

	refundPaymentRequest := banks.RefundRequest{
		OrderID: senagatBankResponse.OrderId,
		Amount:  Amount,
	}

	if err = bpc.Refund(SenagatBankProfile, refundPaymentRequest); err != nil {
		log.Fatal(err)
	}
}

func main() {
	SubmitForms()
}
