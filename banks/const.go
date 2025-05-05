package banks

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // Payment is initiated but not completed yet
	OrderStatusPaid      OrderStatus = "paid"      // Payment completed successfully
	OrderStatusNotPaid   OrderStatus = "not_paid"  // Payment was not made (user abandoned, etc.)
	OrderStatusFailed    OrderStatus = "failed"    // Payment attempt failed (e.g., declined by bank)
	OrderStatusCanceled  OrderStatus = "canceled"  // Order was canceled before payment completed
	OrderStatusRefunded  OrderStatus = "refunded"  // Payment was returned to customer
	OrderStatusError     OrderStatus = "error"     // Unexpected error occurred during payment
	OrderStatusUnderpaid OrderStatus = "underpaid" // Payment received, but amount is less than expected
)

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
	HalkBankOtpUrl            = "acs.gov.tm/acs/pages/enrollment/authentication.jsf"
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
)
