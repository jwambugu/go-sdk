package elarian

import (
	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *service) customerNumber(customer *Customer) *hera.CustomerNumber {
	return &hera.CustomerNumber{
		Number:    customer.CustomerNumber.Number,
		Provider:  hera.CustomerNumberProvider(customer.CustomerNumber.Provider),
		Partition: wrapperspb.String(customer.CustomerNumber.Partition),
	}
}

func (s *service) customerNumbers(customerNumbers []*CustomerNumber) []*hera.CustomerNumber {
	var numbers []*hera.CustomerNumber
	for _, number := range customerNumbers {
		numbers = append(numbers, &hera.CustomerNumber{
			Number:    number.Number,
			Provider:  hera.CustomerNumberProvider(number.Provider),
			Partition: wrapperspb.String(number.Partition),
		})
	}
	return numbers
}

func (s *service) secondaryID(customer *Customer) *hera.IndexMapping {
	return &hera.IndexMapping{
		Key:   customer.SecondaryID.Key,
		Value: wrapperspb.String(customer.SecondaryID.Value),
	}
}

func (s *service) reminderNotification(notf *hera.ReminderNotification) *ReminderNotification {
	return &ReminderNotification{
		CustomerID: notf.CustomerId,
		Reminder: &Reminder{
			Key:        notf.Reminder.Key,
			Payload:    notf.Reminder.Payload.Value,
			Expiration: notf.Reminder.Expiration.AsTime(),
			Interval:   notf.Reminder.Interval.AsDuration(),
		},
		Tag: &Tag{
			Key:        notf.Tag.Mapping.Key,
			Value:      notf.Tag.Mapping.Value.Value,
			Expiration: notf.Tag.Expiration.AsTime(),
		},
		WorkID: notf.WorkId.Value,
	}
}
