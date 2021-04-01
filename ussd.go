package elarian

type (
	// UssdChannel type
	UssdChannel int32

	// UssdMenu struct
	UssdMenu struct {
		IsTerminal bool   `json:"isTerminal,omitempty"`
		Text       string `json:"text,omitempty"`
	}

	// UssdChannelNumber struct
	UssdChannelNumber struct {
		Channel UssdChannel `json:"channel,omitempty"`
		Number  string      `json:"number,omitempty"`
	}

	// UssdOptions struct
	UssdOptions struct {
		SessionID string    `json:"sessionId,omitempty"`
		UssdMenu  *UssdMenu `json:"UssdMenu,omitempty"`
	}

	// UssdSessionNotification struct
	UssdSessionNotification struct {
		SessionID      string            `json:"sessionId,omitempty"`
		CustomerID     string            `json:"customerId,omitempty"`
		Input          string            `json:"input,omitempty"`
		CustomerNumber *CustomerNumber   `json:"customerNumber,omitempty"`
		ChannelNumber  UssdChannelNumber `json:"channelNumber,omitempty"`
	}
)
