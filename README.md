# Elarian go SDK

> The wrapper provides convenient access to the Elarian APIs.

## Documentation

Take a look at the [API docs here](http://docs.elarian.com).

## Install

```bash

    go get github.com/elarianltd/go-sdk

```

## Usage

```go
import (
    elarian elarian github.com/elarianltd/go-sdk
)

func main(){
    client, err := elarian.Initialize("api_key");
    if err != nil {
        log.Fatal(err)
    }
    var customer elarian.Customer
    customer.Id = "customer_id"

    var request elarian.CustomerStaterequest
    request.AppId = "app_id"

    response, err := client.GetCustomerState(&customer, &request)
    if err != nil {
        log.Fatalf("could not get customer state %v", err)
    }
    log.Printf("customer state %v", response)
}
```

## Methods

- `AuthToken()`: Generate auth token

- `GetCustomerState()`:
- `AdoptCustomerState()`:

- `AddCustomerReminder()`:
- `AddCustomerReminderByTag()`:
- `CancelCustomerReminder()`:
- `CancelCustomerReminderByTag()`:

- `UpdateCustomerTag()`:
- `DeleteCustomerTag()`:

- `UpdateCustomerSecondaryId()`:
- `DeleteCustomerSecondaryId()`:

- `UpdateCustomerMetadata()`:
- `DeleteCustomerMetadata ()`:

- `SendMessage()`: Sending a message to your customer
- `SendMessageByTag()`: Sending a message to a group of customers using tags
- `ReplyToMessage()`: Replying to a message from your customer
- `MessagingConsent()`: Opting a customer in or out of receiving messages from your app

- `SendPayment()`:
- `CheckoutPayment()`:

- `MakeVoiceCall()`:

- `StreamNotifications()`:
- `SendWebhookResponse()`:

## Development

```bash

git clone --recurse-submodules https://github.com/ElarianLtd/go-sdk.git
cd go-sdk
make run

```

## Issues

If you find a bug, please file an issue on [our issue tracker on GitHub](https://github.com/ElarianLtd/go-sdk/issues).
