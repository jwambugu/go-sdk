package elarian

import (
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
)

// func (s *service) paymentStatusNotf(notf *hera.PaymentStatusNotification) *PaymentStatusNotification {
// 	return &PaymentStatusNotification{
// 		CustomerID:    notf.CustomerId.Value,
// 		TransactionID: notf.TransactionId,
// 		Status:        PaymentStatus(notf.GetStatus()),
// 	}
// }

// func (s *service) walletPaymentStatusNotf(notf *hera.WalletPaymentStatusNotification) *WalletPaymentStatusNotification {
// 	return &WalletPaymentStatusNotification{
// 		CustomerID:    notf.CustomerId,
// 		Status:        PaymentStatus(notf.Status),
// 		TransactionID: notf.TransactionId,
// 		WalletID:      notf.WalletId,
// 	}
// }

// func (s *service) recievedPaymentNotf(notf *hera.ReceivedPaymentNotification) *ReceivedPaymentNotification {
// 	return &ReceivedPaymentNotification{
// 		CustomerID: notf.CustomerId,
// 		CustomerNumber: &CustomerNumber{
// 			Number:    notf.CustomerNumber.Number,
// 			Provider:  NumberProvider(notf.CustomerNumber.Provider),
// 			Partition: notf.CustomerNumber.Partition.Value,
// 		},
// 		ChannelNumber: &PaymentChannelNumber{
// 			Number:  notf.ChannelNumber.Number,
// 			Channel: PaymentChannel(notf.ChannelNumber.Channel),
// 		},
// 		PurseID:       notf.PurseId,
// 		Status:        PaymentStatus(notf.Status),
// 		TransactionID: notf.TransactionId,
// 		Value: &Cash{
// 			CurrencyCode: notf.Value.CurrencyCode,
// 			Amount:       notf.Value.Amount,
// 		},
// 	}
// }

func (s *service) paymentCounterPartyAsPurse(purse *Purse) *hera.PaymentCounterParty_Purse {
	return &hera.PaymentCounterParty_Purse{
		Purse: &hera.PaymentPurseCounterParty{
			PurseId: purse.PurseID,
		},
	}
}

func (s *service) paymentCounterPartyAsCustomer(
	customer *Customer,
	channel *PaymentChannelNumber,
) *hera.PaymentCounterParty_Customer {
	return &hera.PaymentCounterParty_Customer{
		Customer: &hera.PaymentCustomerCounterParty{
			CustomerNumber: s.customerNumber(customer),
			ChannelNumber: &hera.PaymentChannelNumber{
				Channel: hera.PaymentChannel(channel.Channel),
				Number:  channel.Number,
			},
		},
	}
}

func (s *service) paymentCounterPartyAsWallet(wallet *Wallet) *hera.PaymentCounterParty_Wallet {
	return &hera.PaymentCounterParty_Wallet{
		Wallet: &hera.PaymentWalletCounterParty{
			CustomerId: wallet.CustomerID,
			WalletId:   wallet.WalletID,
		},
	}
}
