package elarian

import (
	"reflect"

	hera "github.com/elarianltd/go-sdk/com_elarian_hera_proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *elarian) heraCustomerNumber(number *CustomerNumber) *hera.CustomerNumber {
	return &hera.CustomerNumber{
		Number:    number.Number,
		Provider:  hera.CustomerNumberProvider(number.Provider),
		Partition: wrapperspb.String(number.Partition),
	}
}

func (s *elarian) heraCustomerNumbers(customerNumbers []*CustomerNumber) []*hera.CustomerNumber {
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

func (s *elarian) customerActivity(activity *hera.CustomerActivityNotification) *CustomerActivityNotification {
	activityNotification := &CustomerActivityNotification{SessionID: activity.SessionId}
	if activity.ChannelNumber != nil {
		activityNotification.Activity = &CustomerActivity{
			Key:        activity.Activity.Key,
			Properties: activity.Activity.Properties,
			CreatedAt:  activity.Activity.CreatedAt.AsTime(),
		}
	}
	if activity.ChannelNumber != nil {
		activityNotification.ActivityChannel = &ActivityChannelNumber{
			Number:  activity.ChannelNumber.Number,
			Channel: ActivityChannel(activity.ChannelNumber.Channel),
		}
	}
	if activity.CustomerNumber != nil {
		activityNotification.CustomerNumber = s.customerNumber(activity.CustomerNumber)
	}
	return activityNotification
}

func (s *elarian) heraSecondaryID(secondaryID *SecondaryID) *hera.IndexMapping {
	return &hera.IndexMapping{
		Key:   secondaryID.Key,
		Value: wrapperspb.String(secondaryID.Value),
	}
}

func (s *elarian) reminderNotification(notf *hera.ReminderNotification) *ReminderNotification {
	if notf == nil || reflect.ValueOf(notf).IsZero() {
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

func (s *elarian) customerNumber(heraCustomerNumber *hera.CustomerNumber) *CustomerNumber {
	custNumber := &CustomerNumber{
		Number:   heraCustomerNumber.Number,
		Provider: NumberProvider(heraCustomerNumber.Provider),
	}
	if heraCustomerNumber.Partition != nil {
		custNumber.Partition = heraCustomerNumber.Partition.Value
	}
	return custNumber
}
