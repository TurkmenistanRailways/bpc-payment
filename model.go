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
	SenagatBank BankType = "SenagatBank"
	HalkBank    BankType = "HalkBank"
	RysgalBank  BankType = "RysgalBank"
)

var bankTypes = []BankType{
	SenagatBank,
	HalkBank,
	RysgalBank,
}

func checkBankType(bankType BankType) error {
	if !slices.Contains(bankTypes, bankType) {
		return ErrBankTypeNotFound
	}

	return nil
}

func (b BankType) Register(user banks.BankUser) banks.Bank {
	switch b {
	case SenagatBank:
		return senagat_bank.Init(user)
	case HalkBank:
		return halk_bank.Init(user)
	case RysgalBank:
		return rysgal_bank.Init(user)
	default:
		return nil
	}
}
