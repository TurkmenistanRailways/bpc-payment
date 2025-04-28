package bpc

import (
	"slices"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/banks/senagat_bank"
)

type BankType string

const (
	SenagatBank BankType = "SenagatBank"
)

var bankTypes = []BankType{
	SenagatBank,
}

func checkBankType(bankType BankType) error {
	if !slices.Contains(bankTypes, bankType) {
		return ErrBankTypeNotFound
	}

	return nil
}

func (b BankType) Register(user banks.BankUser) banks.Bank {
	return senagat_bank.Init(user)
}
