package bpc

import (
	"slices"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/banks/rysgal_bank"
	"github.com/TurkmenistanRailways/bpc-payment/banks/senagat_bank"
)

type BankType string

const (
	SenagatBank BankType = "SenagatBank"
	RysgalBank  BankType = "RysgalBank"
)

var bankTypes = []BankType{
	SenagatBank,
	RysgalBank,
}

func checkBankType(bankType BankType) error {
	if !slices.Contains(bankTypes, bankType) {
		return ErrBankTypeNotFound
	}

	return nil
}

func (b BankType) Register(user banks.BankUser) banks.Bank {
	if b == SenagatBank {
		return senagat_bank.Init(user)
	} else if b == RysgalBank {
		return rysgal_bank.Init(user)
	}
	return nil
}
