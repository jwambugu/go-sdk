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

	// PaymentChannelNumber defines a payment channel number
	PaymentChannelNumber struct {
		Number  string `json:"number"`
		Channel PaymentChannel
	}

	// PaymentRequest defines arguments required to make a payment request
	PaymentRequest struct {
		AppID     string               `json:"appId,omitempty"`
		ProductID string               `json:"productId,omitempty"`
		Cash      Cash                 `json:"cash"`
		Channel   PaymentChannelNumber `json:"channel"`
	}
)

const (
	// Unspecfied type of payment channel
	Unspecfied PaymentChannel = iota
	// Telco type of payment channel represets a telecommunication company such as safaricon
	Telco
)

func (s *service) SendPayment(customer *Customer, params *PaymentRequest) (*hera.SendPaymentReply, error) {
	var request hera.SendPaymentRequest

	if customer.ID != "" {
		request.Customer = &hera.SendPaymentRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.SendPaymentRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}
	request.ChannelNumber = &hera.PaymentChannelNumber{
		Channel: hera.PaymentChannel(params.Channel.Channel),
		Number:  params.Channel.Number,
	}
	request.Value = &hera.Cash{
		Amount:       params.Cash.Amount,
		CurrencyCode: params.Cash.CurrencyCode,
	}

	request.AppId = params.AppID
	request.ProductId = params.ProductID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.SendPayment(ctx, &request)
}

func (s *service) CheckoutPayment(customer *Customer, params *PaymentRequest) (*hera.CheckoutPaymentReply, error) {
	var request hera.CheckoutPaymentRequest

	if customer.ID != "" {
		request.Customer = &hera.CheckoutPaymentRequest_CustomerId{
			CustomerId: customer.ID,
		}
	}
	if !reflect.ValueOf(customer.PhoneNumber).IsZero() {
		request.Customer = &hera.CheckoutPaymentRequest_CustomerNumber{
			CustomerNumber: &hera.CustomerNumber{
				Number:   customer.PhoneNumber.Number,
				Provider: hera.CustomerNumberProvider(customer.PhoneNumber.Provider),
			},
		}
	}

	request.ChannelNumber = &hera.PaymentChannelNumber{
		Channel: hera.PaymentChannel(params.Channel.Channel),
		Number:  params.Channel.Number,
	}
	request.Value = &hera.Cash{
		Amount:       params.Cash.Amount,
		CurrencyCode: params.Cash.CurrencyCode,
	}

	request.AppId = params.AppID
	request.ProductId = params.ProductID

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.client.CheckoutPayment(ctx, &request)
}
