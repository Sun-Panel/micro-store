package paddle

import (
	"net/http"
	"sun-panel/lib/paddle/types"
)

type TransactionService service

type TransactionItem struct {
	PriceID  string `json:"price_id"`
	Quantity int    `json:"quantity"`
}

type TransactionBillingTerms struct {
	Interval  string `json:"interval"`
	Frequency int    `json:"frequency"`
}

type TransactionBillingDetails struct {
	PaymentTerms          TransactionBillingTerms `json:"payment_terms"`
	EnableCheckout        bool                    `json:"enable_checkout"`
	PurchaseOrderNumber   string                  `json:"purchase_order_number"`
	AdditionalInformation string                  `json:"additional_information,omitempty"`
}

type TransactionBillingPeriod struct {
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
}

type TransactionCheckout struct {
	URL string `json:"url,omitempty"`
}

type TransactionCreateParamRequest struct {
	Items          []TransactionItem          `json:"items"`
	Status         *string                    `json:"status,omitempty"`
	CustomerID     *string                    `json:"customer_id,omitempty"`
	AddressID      *string                    `json:"address_id,omitempty"`
	BusinessID     *string                    `json:"business_id,omitempty"`
	CustomData     interface{}                `json:"custom_data,omitempty"`
	CurrencyCode   *string                    `json:"currency_code,omitempty"`
	CollectionMode *string                    `json:"collection_mode,omitempty"`
	DiscountID     *string                    `json:"discount_id,omitempty"`
	BillingDetails *TransactionBillingDetails `json:"billing_details,omitempty"`
	BillingPeriod  *TransactionBillingPeriod  `json:"billing_period,omitempty"`
	Checkout       *TransactionCheckout       `json:"checkout,omitempty"`
}

type TransactionCreateParamResponse struct {
	ID      string                  `json:"id"`
	Details *TransactionDataDetails `json:"details,omitempty"`
}

type TransactionDataDetailsTotals struct {
	Subtotal     string `json:"subtotal"`
	Tax          string `json:"tax"`
	Discount     string `json:"discount"`
	Total        string `json:"total"`
	CurrencyCode string `json:"currency_code"`
}

type TransactionDataDetails struct {
	// TaxRatesUsed         []TaxRateUsed `json:"tax_rates_used"`
	Totals TransactionDataDetailsTotals `json:"totals"`
	// AdjustedTotals       Totals        `json:"adjusted_totals"`
	// PayoutTotals         interface{}   `json:"payout_totals"`          // It could be null or an object, so interface{} is used.
	// AdjustedPayoutTotals interface{}   `json:"adjusted_payout_totals"` // It could be null or an object, so interface{} is used.
	// LineItems            []LineItem    `json:"line_items"`
}

type TransactionData struct {
	ID             string                  `json:"id"`
	Status         string                  `json:"status"`
	CustomerID     string                  `json:"customer_id"`
	AddressID      string                  `json:"address_id"`
	BusinessID     *string                 `json:"business_id,omitempty"`
	CustomData     interface{}             `json:"custom_data,omitempty"`
	Origin         string                  `json:"origin"`
	CollectionMode string                  `json:"collection_mode"`
	SubscriptionID *string                 `json:"subscription_id,omitempty"`
	Details        *TransactionDataDetails `json:"details,omitempty"`
	// InvoiceID      string                    `json:"invoice_id"` // 暂时不需要
	InvoiceNumber  *string                   `json:"invoice_number,omitempty"`
	BillingDetails TransactionBillingDetails `json:"billing_details,omitempty"`
	BillingPeriod  TransactionBillingPeriod  `json:"billing_period,omitempty"`
	CurrencyCode   string                    `json:"currency_code"`
	// DiscountID     interface{}               `json:"discount_id"`  // 暂时不需要
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	// BilledAt       interface{}               `json:"billed_at"`  // 暂时不需要
	Items []TransactionItem `json:"items"`
	// Details        Details                   `json:"details"` // 暂时不需要
	// Payments []interface{}       `json:"payments"` // 暂时不需要
	Checkout TransactionCheckout `json:"checkout,omitempty"`
}

// 创建订单
// https://developer.paddle.com/api-reference/transactions/create-transaction
func (t *TransactionService) Create(options TransactionCreateParamRequest) (*TransactionCreateParamResponse, *http.Response, error) {
	url := "/transactions"
	resp := types.SuccessResponse[TransactionCreateParamResponse]{}
	httpResponse, err := t.client.post(url, options, &resp)
	return &resp.Data, httpResponse, err
}

// 获取订单详情
// https://developer.paddle.com/api-reference/transactions/get-transaction
func (t *TransactionService) Get(transactionsId string) (*TransactionData, *http.Response, error) {
	url := "/transactions/" + transactionsId
	resp := types.SuccessResponse[TransactionData]{}
	httpResponse, err := t.client.get(url, &resp)
	return &resp.Data, httpResponse, err
}
