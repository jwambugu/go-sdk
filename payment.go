package elarian

import (
	"context"
	"reflect"
	"time"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go/payload"
	"google.golang.org/protobuf/proto"
)

type (
	// PaymentChannel type
	PaymentChannel int32

	// PaymentStatus type
	PaymentStatus int32

	// Cash defines a cash object
	Cash struct {
		CurrencyCode string  `json:"currencyCode"`
		Amount       float64 `json:"amount"`
	}

	// PaymentChannelNumber defines a payment channel number
	PaymentChannelNumber struct {
		Number  string `json:"number"`
		Channel PaymentChannel
	}

	// Wallet struct
	Wallet struct {
		CustomerID string
		WalletID   string
	}

	// Purse struct
	Purse struct {
		PurseID string
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
)

// PaymentChannel constants
const (
	PaymentChannelUnspecified PaymentChannel = iota
	PaymentChannelCellular
)

// PaymentStatus constants
const (
	PaymentStatusUnspecified              PaymentStatus = 0
	PaymentStatusQueued                   PaymentStatus = 100
	PaymentStatusPendingConfirmation      PaymentStatus = 101
	PaymentStatusPendingValidation        PaymentStatus = 102
	PaymentStatusValidated                PaymentStatus = 103
	PaymentStatusInvalidRequest           PaymentStatus = 200
	PaymentStatusNotSupported             PaymentStatus = 201
	PaymentStatusInsufficientFunds        PaymentStatus = 202
	PaymentStatusApplicationError         PaymentStatus = 203
	PaymentStatusNotAllowed               PaymentStatus = 204
	PaymentStatusDuplicateRequest         PaymentStatus = 205
	PaymentStatusInvalidPurse             PaymentStatus = 206
	PaymentStatusInvalidWallet            PaymentStatus = 207
	PaymentStatusDecommissionedCustomerID PaymentStatus = 299
	PaymentStatusSuccess                  PaymentStatus = 300
	PaymentStatusPassThrough              PaymentStatus = 301
	PaymentStatusFailed                   PaymentStatus = 400
	PaymentStatusThrottled                PaymentStatus = 401
	PaymentStatusExpired                  PaymentStatus = 402
	PaymentStatusRejected                 PaymentStatus = 403
	PaymentStatusReversed                 PaymentStatus = 500
)

func (s *service) InitiatePayment(customer *Customer, params *Paymentrequest) (*hera.InitiatePaymentReply, error) {
	req := new(hera.AppToServerCommand)
	command := new(hera.AppToServerCommand_InitiatePayment)
	command.InitiatePayment = &hera.InitiatePaymentCommand{}
	req.Entry = command

	command.InitiatePayment.Value = &hera.Cash{
		Amount:       params.Cash.Amount,
		CurrencyCode: params.Cash.CurrencyCode,
	}

	if !reflect.ValueOf(params.CreditParty.Customer).IsZero() {
		command.InitiatePayment.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Purse).IsZero() {
		command.InitiatePayment.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsPurse(&params.CreditParty.Purse),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Wallet).IsZero() {
		command.InitiatePayment.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsWallet(&params.CreditParty.Wallet),
		}
	}

	if !reflect.ValueOf(params.DebitParty.Customer).IsZero() {
		command.InitiatePayment.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Purse).IsZero() {
		command.InitiatePayment.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsPurse(&params.DebitParty.Purse),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Wallet).IsZero() {
		command.InitiatePayment.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsWallet(&params.DebitParty.Wallet),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.InitiatePaymentReply{}, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return &hera.InitiatePaymentReply{}, err
	}

	reply := new(hera.AppToServerCommandReply)
	err = proto.Unmarshal(res.Data(), reply)
	return reply.GetInitiatePayment(), err
}

func (s *service) ReceivePayment(channel *PaymentChannelNumber, customerNumber, transactionID string) (*hera.SimulatorToServerCommandReply, error) {
	req := new(hera.SimulatorToServerCommand)
	command := new(hera.SimulatorToServerCommand_ReceivePayment)
	req.Entry = command
	if !reflect.ValueOf(customerNumber).IsZero() {
		command.ReceivePayment.CustomerNumber = customerNumber
	}
	if !reflect.ValueOf(channel).IsZero() {
		command.ReceivePayment.ChannelNumber = &hera.PaymentChannelNumber{
			Channel: hera.PaymentChannel(channel.Channel),
			Number:  channel.Number,
		}
	}
	command.ReceivePayment.TransactionId = transactionID

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.SimulatorToServerCommandReply{}, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	reply := new(hera.SimulatorToServerCommandReply)
	if err != nil {
		return reply, err
	}
	err = proto.Unmarshal(payload.Data(), reply)
	return reply, err
}

func (s *service) UpdatePaymentStatus(transactionID string, paymentStatus PaymentStatus) (*hera.SimulatorToServerCommandReply, error) {
	req := new(hera.SimulatorToServerCommand)
	command := new(hera.SimulatorToServerCommand_UpdatePaymentStatus)
	req.Entry = command
	command.UpdatePaymentStatus = &hera.UpdatePaymentStatusSimulatorCommand{}
	command.UpdatePaymentStatus.Status = hera.PaymentStatus(paymentStatus)
	command.UpdatePaymentStatus.TransactionId = transactionID

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	data, err := proto.Marshal(req)
	if err != nil {
		return &hera.SimulatorToServerCommandReply{}, err
	}
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	reply := new(hera.SimulatorToServerCommandReply)
	if err != nil {
		return reply, err
	}
	err = proto.Unmarshal(payload.Data(), reply)
	return reply, err
}
