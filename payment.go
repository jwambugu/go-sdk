package elarian

import (
	"context"
	"errors"
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
		CurrencyCode string  `json:"currencyCode,omitempty"`
		Amount       float64 `json:"amount,omitempty"`
	}

	// PaymentChannelNumber defines a payment channel number
	PaymentChannelNumber struct {
		Number  string         `json:"number,omitempty"`
		Channel PaymentChannel `json:"channel,omitempty"`
	}

	// Wallet struct
	Wallet struct {
		CustomerID string `json:"customerId,omitempty"`
		WalletID   string `json:"walletId,omitempty"`
	}

	// Purse struct
	Purse struct {
		PurseID string `json:"purseId,omitempty"`
	}

	// PaymentParty struct
	PaymentParty struct {
		Customer *Customer `json:"customer,omitempty"`
		Wallet   *Wallet   `json:"wallet,omitempty"`
		Purse    *Purse    `json:"purse,omitempty"`
	}

	// Paymentrequest defines arguments required to make a payment request
	Paymentrequest struct {
		Cash        Cash                 `json:"cash"`
		Channel     PaymentChannelNumber `json:"channel"`
		CreditParty PaymentParty         `json:"creditparty,omitempty"`
		DebitParty  PaymentParty         `json:"debitparty,omitempty"`
	}
	// InitiatePaymentReply struct
	InitiatePaymentReply struct {
		CreditCustomerID string        `json:"creditCustomerID,omitempty"`
		DebitCustomerID  string        `json:"debitCustomerID,omitempty"`
		Description      string        `json:"description,omitempty"`
		Status           PaymentStatus `json:"status,omitempty"`
		TransactionID    string        `json:"transactionId,omitempty"`
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

func (s *elarian) InitiatePayment(customer *Customer, params *Paymentrequest) (*InitiatePaymentReply, error) {
	if params == nil || reflect.ValueOf(params).IsZero() {
		return nil, errors.New("Initiate payment params required")
	}

	command := &hera.InitiatePaymentCommand{
		Value: &hera.Cash{
			Amount:       params.Cash.Amount,
			CurrencyCode: params.Cash.CurrencyCode,
		},
	}
	if !reflect.ValueOf(params.CreditParty.Customer).IsZero() {
		command.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Purse).IsZero() {
		command.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsPurse(params.CreditParty.Purse),
		}
	}
	if !reflect.ValueOf(params.CreditParty.Wallet).IsZero() {
		command.CreditParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsWallet(params.CreditParty.Wallet),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Customer).IsZero() {
		command.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsCustomer(customer, &params.Channel),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Purse).IsZero() {
		command.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsPurse(params.DebitParty.Purse),
		}
	}
	if !reflect.ValueOf(params.DebitParty.Wallet).IsZero() {
		command.DebitParty = &hera.PaymentCounterParty{
			Party: s.paymentCounterPartyAsWallet(params.DebitParty.Wallet),
		}
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_InitiatePayment{InitiatePayment: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}

	reply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), reply); err != nil {
		return nil, err
	}
	return &InitiatePaymentReply{
		CreditCustomerID: reply.GetInitiatePayment().CreditCustomerId.Value,
		DebitCustomerID:  reply.GetInitiatePayment().DebitCustomerId.Value,
		Description:      reply.GetInitiatePayment().Description,
		Status:           PaymentStatus(reply.GetInitiatePayment().Status),
		TransactionID:    reply.GetInitiatePayment().TransactionId.Value,
	}, nil
}

func (s *elarian) ReceivePayment(customerNumber, transactionID string, channel *PaymentChannelNumber, cash *Cash, paymentStatus PaymentStatus) (*SimulatorToServerCommandReply, error) {
	if channel == nil || reflect.ValueOf(channel).IsZero() {
		return nil, errors.New("paymentChannel is required")
	}
	if cash == nil || reflect.ValueOf(cash).IsZero() {
		return nil, errors.New("cash is required")
	}

	command := &hera.ReceivePaymentSimulatorCommand{
		TransactionId:  transactionID,
		CustomerNumber: customerNumber,
		ChannelNumber: &hera.PaymentChannelNumber{
			Channel: hera.PaymentChannel(channel.Channel),
			Number:  channel.Number,
		},
		Value: &hera.Cash{
			CurrencyCode: cash.CurrencyCode,
			Amount:       cash.Amount,
		},
		Status: hera.PaymentStatus(paymentStatus),
	}
	req := &hera.SimulatorToServerCommand{
		Entry: &hera.SimulatorToServerCommand_ReceivePayment{ReceivePayment: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.SimulatorToServerCommandReply)
	if err = proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return &SimulatorToServerCommandReply{
		Status:      reply.Status,
		Description: reply.Description,
		Message:     s.OutboundMessage(reply.Message),
	}, nil
}

func (s *elarian) UpdatePaymentStatus(transactionID string, paymentStatus PaymentStatus) (*SimulatorToServerCommandReply, error) {
	command := &hera.UpdatePaymentStatusSimulatorCommand{
		TransactionId: transactionID,
		Status:        hera.PaymentStatus(paymentStatus),
	}
	req := &hera.SimulatorToServerCommand{
		Entry: &hera.SimulatorToServerCommand_UpdatePaymentStatus{UpdatePaymentStatus: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	payload, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}
	reply := new(hera.SimulatorToServerCommandReply)
	if err = proto.Unmarshal(payload.Data(), reply); err != nil {
		return nil, err
	}
	return &SimulatorToServerCommandReply{
		Status:      reply.Status,
		Message:     s.OutboundMessage(reply.Message),
		Description: reply.Description,
	}, nil
}
