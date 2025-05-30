package banks

import "errors"

const (
	HalkBankBaseUrl    = "https://mpi.gov.tm"
	SenagatBankBaseUrl = "https://epg.senagatbank.com.tm"
	RysgalBankBaseUrl  = "https://epg.rysgalbank.tm"

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
	HalkBankOtpUrl            = "https://acs.gov.tm/acs/pages/enrollment/authentication.jsf"
	HalkBankRefundURL         = "/payment/rest/refund.do"

	RysgalRegisterURL       = "/epg/rest/register.do"
	RysgalConfirmPaymentURL = "/epg/rest/processform.do"
	RysgalFinishURL         = "/payments/rest/finish3ds.do"
	RysgalBankOtpUrl        = "/acs/pages/enrollment/authentication.jsf"
	RysgalOrderStatusURL    = "/epg/rest/getOrderStatusExtended.do"
	RysgalRefundURL         = "/epg/rest/refund.do"
)

const (
	CurrencyTMT = "934"
	// URLFormat is the format for constructing URLs.
	// It takes three parameters: base URL, endpoint, and query string.
	// The final URL will be constructed as: baseURL + endpoint + "?" + queryString
	URLFormat = "%s%s?%s"
)

var (
	ErrorInvalidCardCredentials = errors.New("invalid card credentials")
)

type OrderStatus string

const (
	OrderStatusNotPaid       OrderStatus = "not_paid"       // Order was registered but not paid
	OrderStatusAuthorized    OrderStatus = "authorized"     // Payment was authorized but not yet captured
	OrderStatusPaid          OrderStatus = "paid"           // Payment was authorized and captured
	OrderStatusAuthCanceled  OrderStatus = "auth_canceled"  // Authorization was canceled
	OrderStatusRefunded      OrderStatus = "refunded"       // Payment was refunded
	OrderStatus3DSecure      OrderStatus = "3d_secure"      // Access control server of issuing bank initiated authorization
	OrderStatusDeclined      OrderStatus = "declined"       // Authorization was declined
	OrderStatusPending       OrderStatus = "pending"        // Payment is pending
	OrderStatusPartiallyPaid OrderStatus = "partially_paid" // Intermediate status for multiple partial completions
	OrderStatusError         OrderStatus = "error"          // Unexpected error occurred during payment
)

var StatusCodes = map[int]OrderStatus{
	0: OrderStatusNotPaid,
	1: OrderStatusAuthorized,
	2: OrderStatusPaid,
	3: OrderStatusAuthCanceled,
	4: OrderStatusRefunded,
	5: OrderStatus3DSecure,
	6: OrderStatusDeclined,
	7: OrderStatusPending,
	8: OrderStatusPartiallyPaid,
}
