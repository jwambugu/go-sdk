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
		CustomerId    string `json:"customerId,omitempty"`
		Status        int    `json:"status,omitempty"`
		TransactionId string `json:"transactionId,omitempty"`
	}
)

const (
	PAYMENT_CHANNEL_UNSPECIFIED PaymentChannel = iota
	PAYMENT_CHANNEL_TELCO
)

const (
	PAYMENT_STATUS_UNSPECIFIED                = 0
	PAYMENT_STATUS_QUEUED                     = 101
	PAYMENT_STATUS_PENDING_CONFIRMATION       = 102
	PAYMENT_STATUS_PENDING_VALIdATION         = 103
	PAYMENT_STATUS_VALIdATED                  = 104
	PAYMENT_STATUS_INVALId_REQUEST            = 200
	PAYMENT_STATUS_NOT_SUPPORTED              = 201
	PAYMENT_STATUS_INSUFFICIENT_FUNDS         = 202
	PAYMENT_STATUS_APPLICATION_ERROR          = 203
	PAYMENT_STATUS_NOT_ALLOWED                = 204
	PAYMENT_STATUS_DUPLICATE_REQUEST          = 205
	PAYMENT_STATUS_INVALId_PURSE              = 206
	PAYMENT_STATUS_INVALId_WALLET             = 207
	PAYMENT_STATUS_DECOMMISSIONED_CUSTOMER_Id = 299
	PAYMENT_STATUS_SUCCESS                    = 300
	PAYMENT_STATUS_PASS_THROUGH               = 301
	PAYMENT_STATUS_FAILED                     = 400
	PAYMENT_STATUS_THROTTLED                  = 401
	PAYMENT_STATUS_EXPIRED                    = 402
	PAYMENT_STATUS_REJECTED                   = 403
	PAYMENT_STATUS_REVERSED                   = 500
)

func (s *service) setPaymentCounterPartyAsPurse(purse *Purse) *hera.PaymentCounterParty_Purse {
	return &hera.PaymentCounterParty_Purse{
		Purse: &hera.PaymentPurseCounterParty{
			PurseId: purse.PurseId,
		},
	}
}
func (s *service) setPaymentCounterPartyAsCustomer(
	customer *Customer,
	channel *PaymentChannelNumber,
) *hera.PaymentCounterParty_Customer {
	return &hera.PaymentCounterParty_Customer{
		Customer: &hera.PaymentCustomerCounterParty{
			CustomerNumber: s.setCustomerNumber(customer),
			ChannelNumber: &hera.PaymentChannelNumber{
				Channel: hera.PaymentChannel(channel.Channel),
				Number:  channel.Number,
			},
		},
	}
}
func (s *service) setPaymentCounterPartyAsWallet(wallet *Wallet) *hera.PaymentCounterParty_Wallet {
	return &hera.PaymentCounterParty_Wallet{
		Wallet: &hera.PaymentWalletCounterParty{
			CustomerId: wallet.CustomerId,
			WalletId:   wallet.WalletId,
		},
	}
}

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
			Party: s.setPaymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Purse).IsZero() {
		request.CreditParty = &hera.PaymentCounterParty{
			Party: s.setPaymentCounterPartyAsPurse(&params.CreditParty.Purse),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Wallet).IsZero() {
		request.CreditParty = &hera.PaymentCounterParty{
			Party: s.setPaymentCounterPartyAsWallet(&params.CreditParty.Wallet),
		}
	}

	if !reflect.ValueOf(params.DebitParty.Customer).IsZero() {
		request.DebitParty = &hera.PaymentCounterParty{
			Party: s.setPaymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Purse).IsZero() {
		request.DebitParty = &hera.PaymentCounterParty{
			Party: s.setPaymentCounterPartyAsPurse(&params.DebitParty.Purse),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Wallet).IsZero() {
		request.DebitParty = &hera.PaymentCounterParty{
			Party: s.setPaymentCounterPartyAsWallet(&params.DebitParty.Wallet),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.InitiatePayment(ctx, &request)
}
