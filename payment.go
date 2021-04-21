package elarian

import (
	"context"
	"errors"
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"github.com/rsocket/rsocket-go/payload"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type (
	// PaymentChannel type
	PaymentChannel int32

	// PaymentStatus type
	PaymentStatus int32

	// IsPaymentParty is an interface implemented by a CustomerPaymentParty, WalletCounterParty, PurseCounterParty, ChannelCounterParty
	IsPaymentParty interface {
		paymentParty()
	}

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

	// CustomerCounterParty struct
	CustomerPaymentParty struct {
		CustomerNumber *CustomerNumber
		ChannelNumber  *PaymentChannelNumber
	}

	// ChannelCounterParty struct
	ChannelCounterParty struct {
		ChannelNumber *PaymentChannelNumber
		Account       string
		ChannelCode   int32
	}

	// PaymentParty struct
	PaymentParty struct {
		CustomerCounterParty *CustomerPaymentParty `json:"CustomerCounterParty,omitempty"`
		WalletCounterParty   *Wallet               `json:"walletCounterParty,omitempty"`
		PurseCounterParty    *Purse                `json:"purseCounterParty,omitempty"`
		ChannelCounterParty  *ChannelCounterParty  `json:"channelCounterParty,omitempty"`
	}

	// PaymentCounterParty defines arguments required to make a payment request
	PaymentCounterParty struct {
		CreditParty IsPaymentParty `json:"creditparty,omitempty"`
		DebitParty  IsPaymentParty `json:"debitparty,omitempty"`
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

func (*CustomerPaymentParty) paymentParty() {}
func (*Wallet) paymentParty()               {}
func (*Purse) paymentParty()                {}
func (*ChannelCounterParty) paymentParty()  {}

func (s *elarian) InitiatePayment(ctx context.Context, party *PaymentCounterParty, cash *Cash) (*InitiatePaymentReply, error) {
	if party == nil || reflect.ValueOf(party).IsZero() {
		return nil, errors.New("Payment Party required")
	}

	command := &hera.InitiatePaymentCommand{
		Value: &hera.Cash{
			Amount:       cash.Amount,
			CurrencyCode: cash.CurrencyCode,
		},
	}

	if party.CreditParty != nil {
		counterParty := &hera.PaymentCounterParty{}
		if customerCounterParty, ok := party.CreditParty.(*CustomerPaymentParty); ok {
			counterParty.Party = &hera.PaymentCounterParty_Customer{
				Customer: &hera.PaymentCustomerCounterParty{
					CustomerNumber: s.heraCustomerNumber(customerCounterParty.CustomerNumber),
					ChannelNumber: &hera.PaymentChannelNumber{
						Channel: hera.PaymentChannel(customerCounterParty.ChannelNumber.Channel),
						Number:  customerCounterParty.ChannelNumber.Number,
					},
				},
			}
			command.CreditParty = counterParty
		}

		if purseCounterParty, ok := party.CreditParty.(*Purse); ok {
			counterParty.Party = &hera.PaymentCounterParty_Purse{
				Purse: &hera.PaymentPurseCounterParty{
					PurseId: purseCounterParty.PurseID,
				},
			}
			command.CreditParty = counterParty
		}
		if walletCounterParty, ok := party.CreditParty.(*Wallet); ok {
			counterParty.Party = &hera.PaymentCounterParty_Wallet{
				Wallet: &hera.PaymentWalletCounterParty{
					CustomerId: walletCounterParty.CustomerID,
					WalletId:   walletCounterParty.WalletID,
				},
			}
			command.CreditParty = counterParty
		}
		if channelCounterParty, ok := party.CreditParty.(*ChannelCounterParty); ok {
			counterParty.Party = &hera.PaymentCounterParty_Channel{
				Channel: &hera.PaymentChannelCounterParty{
					ChannelNumber: &hera.PaymentChannelNumber{
						Channel: hera.PaymentChannel(channelCounterParty.ChannelNumber.Channel),
						Number:  channelCounterParty.ChannelNumber.Number,
					},
					ChannelCode: channelCounterParty.ChannelCode,
					Account:     wrapperspb.String(channelCounterParty.Account),
				},
			}
			command.CreditParty = counterParty
		}
	}

	if party.DebitParty != nil {
		counterParty := &hera.PaymentCounterParty{}
		if customerCounterParty, ok := party.DebitParty.(*CustomerPaymentParty); ok {
			counterParty.Party = &hera.PaymentCounterParty_Customer{
				Customer: &hera.PaymentCustomerCounterParty{
					CustomerNumber: s.heraCustomerNumber(customerCounterParty.CustomerNumber),
					ChannelNumber: &hera.PaymentChannelNumber{
						Channel: hera.PaymentChannel(customerCounterParty.ChannelNumber.Channel),
						Number:  customerCounterParty.ChannelNumber.Number,
					},
				},
			}
			command.DebitParty = counterParty
		}
		if purseCounterParty, ok := party.DebitParty.(*Purse); ok {
			counterParty.Party = &hera.PaymentCounterParty_Purse{
				Purse: &hera.PaymentPurseCounterParty{
					PurseId: purseCounterParty.PurseID,
				},
			}
			command.DebitParty = counterParty
		}
		if walletCounterParty, ok := party.DebitParty.(*Wallet); ok {
			counterParty.Party = &hera.PaymentCounterParty_Wallet{
				Wallet: &hera.PaymentWalletCounterParty{
					CustomerId: walletCounterParty.CustomerID,
					WalletId:   walletCounterParty.WalletID,
				},
			}
			command.DebitParty = counterParty
		}
		if channelCounterParty, ok := party.DebitParty.(*ChannelCounterParty); ok {
			counterParty.Party = &hera.PaymentCounterParty_Channel{
				Channel: &hera.PaymentChannelCounterParty{
					ChannelNumber: &hera.PaymentChannelNumber{
						Channel: hera.PaymentChannel(channelCounterParty.ChannelNumber.Channel),
						Number:  channelCounterParty.ChannelNumber.Number,
					},
					ChannelCode: 12,
					Account:     wrapperspb.String(""),
				},
			}
			command.DebitParty = counterParty
		}
	}

	req := &hera.AppToServerCommand{
		Entry: &hera.AppToServerCommand_InitiatePayment{InitiatePayment: command},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	res, err := s.client.RequestResponse(payload.New(data, []byte{})).Block(ctx)
	if err != nil {
		return nil, err
	}

	commandReply := &hera.AppToServerCommandReply{}
	if err = proto.Unmarshal(res.Data(), commandReply); err != nil {
		return nil, err
	}
	paymentReply := commandReply.GetInitiatePayment()
	reply := &InitiatePaymentReply{
		Status:        PaymentStatus(paymentReply.Status),
		Description:   paymentReply.Description,
		TransactionID: paymentReply.TransactionId.Value,
	}
	if paymentReply.CreditCustomerId != nil {
		reply.CreditCustomerID = paymentReply.CreditCustomerId.Value
	}
	if paymentReply.DebitCustomerId != nil {
		reply.DebitCustomerID = paymentReply.DebitCustomerId.Value
	}
	return reply, nil
}

func (s *elarian) ReceivePayment(ctx context.Context, customerNumber, transactionID string, channel *PaymentChannelNumber, cash *Cash, paymentStatus PaymentStatus) (*SimulatorToServerCommandReply, error) {
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

func (s *elarian) UpdatePaymentStatus(ctx context.Context, transactionID string, paymentStatus PaymentStatus) (*SimulatorToServerCommandReply, error) {
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
