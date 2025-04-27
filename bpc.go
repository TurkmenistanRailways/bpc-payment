package bpc

import (
	"fmt"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
)

type BPC struct {
	banks map[string]banks.Bank
}

func New() *BPC {
	return &BPC{
		banks: make(map[string]banks.Bank),
	}
}

func (bpc *BPC) AddProfile(profileName string, bankType BankType, claims banks.BankUser) error {
	if err := checkBankType(bankType); err != nil {
		return err
	}

	if _, ok := bpc.banks[profileName]; ok {
		return fmt.Errorf("profile %s already exists", profileName)
	}

	bpc.banks[profileName] = bankType.Register(claims)

	return nil
}

func (bpc *BPC) OrderRegister(profileName string, sessionTimeOut int, amount int64) (banks.OrderRegistrationResponse, error) {
	if err := bpc.checkProfile(profileName); err != nil {
		return banks.OrderRegistrationResponse{}, err
	}

	registerForm := banks.RegisterForm{
		Amount:         amount,
		SessionTimeout: sessionTimeOut,
		Language:       "ru",
	}

	response, err := bpc.banks[profileName].OrderRegister(registerForm)
	if err != nil {
		return banks.OrderRegistrationResponse{}, err
	}

	return response, nil
}

func (bpc *BPC) Submit(profileName string, submitCard banks.SubmitCard) (string, error) {
	if err := bpc.checkProfile(profileName); err != nil {
		return "", err
	}

	requestID, err := bpc.banks[profileName].SubmitCard(submitCard)
	if err != nil {
		return "", err
	}

	return requestID, nil
}

func (bpc *BPC) ResendOtpCode(profileName string, requestID string) error {
	if err := bpc.checkProfile(profileName); err != nil {
		return err
	}

	if err := bpc.banks[profileName].ResendOtpCode(requestID); err != nil {
		return err
	}

	return nil
}

func (bpc *BPC) ConfirmPayment(profileName string, form banks.ConfirmPaymentRequest) error {
	if err := bpc.checkProfile(profileName); err != nil {
		return err
	}

	if err := bpc.banks[profileName].ConfirmPayment(form); err != nil {
		return err
	}

	return nil
}

func (bpc *BPC) Refund(profileName string, form banks.RefundRequest) error {
	if err := bpc.checkProfile(profileName); err != nil {
		return err
	}

	if err := bpc.banks[profileName].Refund(form); err != nil {
		return err
	}

	return nil
}

func (bpc *BPC) checkProfile(profileName string) error {
	if _, ok := bpc.banks[profileName]; !ok {
		return fmt.Errorf("profile %s not found", profileName)
	}

	return nil
}
