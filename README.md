# Elarian go SDK

> The wrapper provides convenient access to the Elarian APIs.

## Documentation

Take a look at the [API docs here](http://docs.elarian.com).

## Install

```bash

    go get github.com/elarian/go-sdk

```

## Usage

```go
import elarian github.com/elarian/go-sdk

func main(){
    client, err := elariango.Initialize("api_key", true);
    if err != nil {
        log.Fatal(err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    res, err := client.GetCustomerState(
        ctx,
        &elarian.GetCustomerStateRequest{
            AppId: "app_id",
            Customer: &elarian.GetCustomerStateRequest_CustomerId{
            CustomerId: "customer_id",
        },
    })
    if err != nil {
        log.Fatalf("could not get customer state %v", err)
    }
    log.Printf("customer state %v", res.GetData())
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

If you find a bug, please file an issue on [our issue tracker on GitHub](https://github.com/ElarianLtd/elariango/issues).
