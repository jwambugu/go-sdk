package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// PaymentChannel type
	PaymentChannel int32

	// PaymentStatus type
	PaymentStatus int32

	// PaymentChannelNumber defines a payment channel number
	PaymentChannelNumber struct {
		Number  string `json:"number"`
		Channel PaymentChannel
	}

	// Wallet struct
	Wallet struct {
		CustomerId string
		WalletId   string
	}

	// Purse struct
	Purse struct {
		PurseId string
	}

	// PaymentParty struct
	PaymentParty struct {
		Customer Customer
		Wallet   Wallet
		Purse    Purse
	}

	// Paymentrequest defines arguments required to make a payment request
	Paymentrequest struct {
		Cash        Cash                 `json:"cash"`
		Channel     PaymentChannelNumber `json:"channel"`
		CreditParty PaymentParty         `json:"creditparty,omitempty"`
		DebitParty  PaymentParty         `json:"debitparty,omitempty"`
	}

	// PaymentStatusNotification struct
	PaymentStatusNotification struct {
		CustomerId    string        `json:"customerId,omitempty"`
		TransactionId string        `json:"transactionId,omitempty"`
		Status        PaymentStatus `json:"status,omitempty"`
	}

	// ReceivedPaymentNotification struct
	ReceivedPaymentNotification struct {
		ChannelNumber  *PaymentChannelNumber `json:"channelNumber,omitempty"`
		CustomerId     string                `json:"customerId,omitempty"`
		CustomerNumber *CustomerNumber       `json:"customerNumber,omitempty"`
		PurseId        string                `json:"purseId,omitempty"`
		Status         PaymentStatus         `json:"status,omitempty"`
		TransactionId  string                `json:"transactionId,omitempty"`
		Value          *Cash                 `json:"value,omitempty"`
	}

	WalletPaymentStatusNotification struct {
		CustomerId    string        `json:"customerId,omitempty"`
		Status        PaymentStatus `json:"status,omitempty"`
		TransactionId string        `json:"transactionId,omitempty"`
		WalletId      string        `json:"walletId,omitempty"`
	}
)

const (
	PAYMENT_CHANNEL_UNSPECIFIED PaymentChannel = iota
	PAYMENT_CHANNEL_TELCO
)

const (
	PAYMENT_STATUS_UNSPECIFIED                PaymentStatus = 0
	PAYMENT_STATUS_QUEUED                     PaymentStatus = 101
	PAYMENT_STATUS_PENDING_CONFIRMATION       PaymentStatus = 102
	PAYMENT_STATUS_PENDING_VALIdATION         PaymentStatus = 103
	PAYMENT_STATUS_VALIdATED                  PaymentStatus = 104
	PAYMENT_STATUS_INVALId_REQUEST            PaymentStatus = 200
	PAYMENT_STATUS_NOT_SUPPORTED              PaymentStatus = 201
	PAYMENT_STATUS_INSUFFICIENT_FUNDS         PaymentStatus = 202
	PAYMENT_STATUS_APPLICATION_ERROR          PaymentStatus = 203
	PAYMENT_STATUS_NOT_ALLOWED                PaymentStatus = 204
	PAYMENT_STATUS_DUPLICATE_REQUEST          PaymentStatus = 205
	PAYMENT_STATUS_INVALId_PURSE              PaymentStatus = 206
	PAYMENT_STATUS_INVALId_WALLET             PaymentStatus = 207
	PAYMENT_STATUS_DECOMMISSIONED_CUSTOMER_Id PaymentStatus = 299
	PAYMENT_STATUS_SUCCESS                    PaymentStatus = 300
	PAYMENT_STATUS_PASS_THROUGH               PaymentStatus = 301
	PAYMENT_STATUS_FAILED                     PaymentStatus = 400
	PAYMENT_STATUS_THROTTLED                  PaymentStatus = 401
	PAYMENT_STATUS_EXPIRED                    PaymentStatus = 402
	PAYMENT_STATUS_REJECTED                   PaymentStatus = 403
	PAYMENT_STATUS_REVERSED                   PaymentStatus = 500
)

func (s *service) InitiatePayment(customer *Customer, params *Paymentrequest) (*hera.InitiatePaymentReply, error) {
	var request hera.InitiatePaymentRequest
	request.AppId = s.appId
	request.OrgId = s.orgId
	request.Value = &hera.Cash{
		Amount:       params.Cash.Amount,
		CurrencyCode: params.Cash.CurrencyCode,
	}

	if !reflect.ValueOf(params.CreditParty.Customer).IsZero() {
		request.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Purse).IsZero() {
		request.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsPurse(&params.CreditParty.Purse),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Wallet).IsZero() {
		request.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsWallet(&params.CreditParty.Wallet),
		}
	}

	if !reflect.ValueOf(params.DebitParty.Customer).IsZero() {
		request.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Purse).IsZero() {
		request.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsPurse(&params.DebitParty.Purse),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Wallet).IsZero() {
		request.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsWallet(&params.DebitParty.Wallet),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return s.client.InitiatePayment(ctx, &request)
}
