package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *service) customerNumber(number *CustomerNumber) *hera.CustomerNumber {
	return &hera.CustomerNumber{
		Number:    number.Number,
		Provider:  hera.CustomerNumberProvider(number.Provider),
		Partition: wrapperspb.String(number.Partition),
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

func (s *service) customerActivity(activity *hera.CustomerActivityNotification) *CustomerActivityNotification {
	return &CustomerActivityNotification{
		SessionID: activity.SessionId,
		Activity: &CustomerActivity{
			Key:        activity.Activity.Key,
			Properties: activity.Activity.Properties,
			CreatedAt:  activity.Activity.CreatedAt.AsTime(),
		},
		CustomerNumber: &CustomerNumber{
			Number:    activity.CustomerNumber.Number,
			Provider:  NumberProvider(activity.CustomerNumber.Provider),
			Partition: activity.CustomerNumber.Partition.Value,
		},
		ActivityChannel: &ActivityChannelNumber{
			Number:  activity.ChannelNumber.Number,
			Channel: ActivityChannel(activity.ChannelNumber.Channel),
		},
	}
}

func (s *service) secondaryID(customer *Customer) *hera.IndexMapping {
	return &hera.IndexMapping{
		Key:   customer.SecondaryID.Key,
		Value: wrapperspb.String(customer.SecondaryID.Value),
	}
}

func (s *service) reminderNotification(notf *hera.ReminderNotification) *ReminderNotification {
	if reflect.ValueOf(notf).IsZero() {
		return &ReminderNotification{}
	}
	reminderNotf := &ReminderNotification{}
	if !reflect.ValueOf(notf.Reminder).IsZero() {
		reminderNotf.Reminder =
			&Reminder{
				Key:      notf.Reminder.Key,
				Payload:  notf.Reminder.Payload.Value,
				Interval: notf.Reminder.Interval.AsDuration(),
				RemindAt: notf.Reminder.RemindAt.AsTime(),
			}
	}
	if !reflect.ValueOf(notf.Tag).IsZero() {
		reminderNotf.Tag = &Tag{
			Key:        notf.Tag.Mapping.Key,
			Value:      notf.Tag.Mapping.Value.Value,
			Expiration: notf.Tag.ExpiresAt.AsTime(),
		}
	}
	if !reflect.ValueOf(notf.WorkId).IsZero() {
		reminderNotf.WorkID = notf.WorkId.Value
	}
	return reminderNotf
}
