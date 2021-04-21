package elarian

import (
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

func (s *elarian) paymentStatusNotf(notf *hera.PaymentStatusNotification) *PaymentStatusNotification {
	return &PaymentStatusNotification{
		TransactionID: notf.TransactionId,
		Status:        PaymentStatus(notf.Status),
	}
}

func (s *elarian) walletPaymentStatusNotf(notf *hera.WalletPaymentStatusNotification) *WalletPaymentStatusNotification {
	return &WalletPaymentStatusNotification{
		Status:        PaymentStatus(notf.Status),
		TransactionID: notf.TransactionId,
		WalletID:      notf.WalletId,
	}
}

func (s *elarian) recievedPaymentNotf(notf *hera.ReceivedPaymentNotification) *ReceivedPaymentNotification {
	notification := &ReceivedPaymentNotification{
		ChannelNumber: &PaymentChannelNumber{
			Number:  notf.ChannelNumber.Number,
			Channel: PaymentChannel(notf.ChannelNumber.Channel),
		},
		PurseID:       notf.PurseId,
		Status:        PaymentStatus(notf.Status),
		TransactionID: notf.TransactionId,
		Value: &Cash{
			CurrencyCode: notf.Value.CurrencyCode,
			Amount:       notf.Value.Amount,
		},
	}
	if notf.CustomerNumber != nil {
		customerNumber := s.customerNumber(notf.CustomerNumber)
		if notf.CustomerNumber.Partition != nil {
			customerNumber.Partition = notf.CustomerNumber.Partition.Value
		}
		notification.CustomerNumber = customerNumber
	}
	return notification
}

func (s *elarian) paymentCounterPartyAsPurse(purse *Purse) *hera.PaymentCounterParty_Purse {
	return &hera.PaymentCounterParty_Purse{
		Purse: &hera.PaymentPurseCounterParty{
			PurseId: purse.PurseID,
		},
	}
}

func (s *elarian) paymentCounterPartyAsCustomer(customer *Customer, channel *PaymentChannelNumber) *hera.PaymentCounterParty_Customer {
	return &hera.PaymentCounterParty_Customer{
		Customer: &hera.PaymentCustomerCounterParty{
			CustomerNumber: s.heraCustomerNumber(customer.CustomerNumber),
			ChannelNumber: &hera.PaymentChannelNumber{
				Channel: hera.PaymentChannel(channel.Channel),
				Number:  channel.Number,
			},
		},
	}
}

func (s *elarian) paymentCounterPartyAsWallet(wallet *Wallet) *hera.PaymentCounterParty_Wallet {
	return &hera.PaymentCounterParty_Wallet{
		Wallet: &hera.PaymentWalletCounterParty{
			CustomerId: wallet.CustomerID,
			WalletId:   wallet.WalletID,
		},
	}
}
