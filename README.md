# Elarian

> A convenient way to interact with the Elarian APIs.
> ***Project Status: Still under ACTIVE DEVELOPMENT, APIs are unstable and may change at any time until release of v1.0.0.***

## Install

You can install the package from [github](https://www.github.com/elarianltd/go-sdk) by running:

```bash
go get github.com/elarianltd/go-sdk
```

## Usage

```go
package main

import (
 "log"
 "os"
 "os/signal"
 "syscall"

 elarian "github.com/elarianltd/go-sdk"
)

const (
 AppID  string = "appID"
 OrgID  string = "orgID"
 APIKey string = "apiKey"
)

func main() {
 var (
  custNumber *elarian.CustomerNumber
  channel    *elarian.MessagingChannelNumber
  opts       *elarian.Options
 )

 opts = &elarian.Options{
  APIKey:             aPIKey,
  OrgID:              orgID,
  AppID:              appID,
  AllowNotifications: true,
  Log:                true,
 }
 service, err := elarian.Connect(opts, nil)
 if err != nil {
  log.Fatalf("Error Initializing Elarian: %v \n", err)
 }
 defer service.Disconnect()

 sigs := make(chan os.Signal, 1)
 signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
 go func() {
  <-sigs
  log.Println("Disconnecting From Elarian")
  service.Disconnect()
  os.Exit(0)
 }()

 custNumber = &elarian.CustomerNumber{Number: "+254708752502", Provider: elarian.CustomerNumberProviderCellular}
 channel = &elarian.MessagingChannelNumber{Number: "21356", Channel: elarian.MessagingChannelSms}

 ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*30))
 defer cancel()
 response, err := service.SendMessage(ctx, custNumber, channel, elarian.TextMessage("Hello world from the go sdk"))
 if err != nil {
  log.Fatalf("Message not send %v \n", err.Error())
 }
 log.Printf("Status %d Description %s \n customerID %s \n", response.Status, response.Description, response.CustomerID)
}

```

See [example](example/) for a full sample app.

## Documentation

Take a look at the [API docs here](http://developers.elarian.com). For detailed info on this SDK, see the [reference](docs/).

## Development

Run all tests:

```bash
make test
```

See [SDK Spec](https://github.com/ElarianLtd/sdk-spec) for reference.

## Issues

If you find a bug, please file an issue on [our issue tracker on GitHub](https://github.com/ElarianLtd/go-sdk/issues).
