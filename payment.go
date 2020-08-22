package elariango

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

type (
	// PaymentChannelNumber defines a payment channel number
	PaymentChannelNumber struct {
		Number  string `json:"number"`
		Channel hera.PaymentChannel
	}

	// PaymentRequest defines arguments required to make a payment request
	PaymentRequest struct {
		AppID     string               `json:"appId,omitempty"`
		ProductID string               `json:"productId,omitempty"`
		Cash      Cash                 `json:"cash"`
		Channel   PaymentChannelNumber `json:"channel"`
	}
)

func (e *elarian) SendPayment(customer *Customer, params *PaymentRequest) (*hera.SendPaymentReply, error) {
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
				Provider: customer.PhoneNumber.Provider,
			},
		}
	}
	request.ChannelNumber = &hera.PaymentChannelNumber{
		Channel: params.Channel.Channel,
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
	return e.client.SendPayment(ctx, &request)
}

func (e *elarian) CheckoutPayment(customer *Customer, params *PaymentRequest) (*hera.CheckoutPaymentReply, error) {
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
				Provider: customer.PhoneNumber.Provider,
			},
		}
	}

	request.ChannelNumber = &hera.PaymentChannelNumber{
		Channel: params.Channel.Channel,
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
	return e.client.CheckoutPayment(ctx, &request)
}
