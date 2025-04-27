package bpc

import (
	"slices"

	"github.com/TurkmenistanRailways/bpc-payment/banks"
	"github.com/TurkmenistanRailways/bpc-payment/banks/halk_bank"
	"github.com/TurkmenistanRailways/bpc-payment/banks/rysgal_bank"
	"github.com/TurkmenistanRailways/bpc-payment/banks/senagat_bank"
)

type BankType string

const (
	RysgalBank  BankType = "RysgalBank"
	SenagatBank BankType = "SenagatBank"
	HalkBank    BankType = "HalkBank"
)

var bankTypes = []BankType{
	RysgalBank, SenagatBank, HalkBank,
}

func checkBankType(bankType BankType) error {
	if !slices.Contains(bankTypes, bankType) {
		return ErrBankTypeNotFound
	}

	return nil
}

func (b BankType) Register(user banks.BankUser) banks.Bank {
	if b == RysgalBank {
		return rysgal_bank.Init(user)
	}

	if b == SenagatBank {
		return senagat_bank.Init(user)
	}

	return halk_bank.Init(user)
}
