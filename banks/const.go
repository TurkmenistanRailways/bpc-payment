package banks

const (
	HalkBankBaseUrl    = "https://mpi.gov.tm"
	SenagatBankBaseUrl = "https://epg.senagatbank.com.tm"
	RysgalBankBaseUrl  = "https://epg2.rysgalbank.tm"

	SenagatRegisterURL       = "/epg/rest/register.do"
	SenagatConfirmPaymentURL = "/epg/rest/processform.do"
	SenagatFinishURL         = "/payments/rest/finish3ds.do"
	SenagatOTPURL            = "/acs/api/3ds/form/otp"
	SenagatOrderStatusURL    = "/epg/rest/getOrderStatusExtended.do"
	SenagatRefundURL         = "/epg/rest/refund.do"

	HalkBankRegisterURL       = "/payment/rest/register.do"
	HalkBankConfirmPaymentURL = "/payment/rest/processform.do"
	HalkBankFinishURL         = "/payment/rest/finish3ds.do"
	HalkBankOrderStatusURL    = "/payment/rest/getOrderStatusExtended.do"
	HalkBankOtpUrl            = "acs.gov.tm/acs/pages/enrollment/authentication.jsf"
	HalkBankRefundURL         = "/payment/rest/refund.do"

	RysgalRegisterURL       = "/epg/rest/register.do"
	RysgalConfirmPaymentURL = "/epg/rest/processform.do"
)

const (
	CurrencyTMT = "934"
)
