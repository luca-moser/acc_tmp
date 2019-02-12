package main

import (
	"fmt"
	"github.com/iotaledger/iota.go/account/builder"
	"github.com/iotaledger/iota.go/account/plugins/promoter"
	"github.com/iotaledger/iota.go/account/plugins/transfer/poller"
	"github.com/iotaledger/iota.go/account/store/inmemory"
	"github.com/iotaledger/iota.go/api"
	"time"
)

func main() {

	// init IRI API object
	iotaAPI, err := api.ComposeAPI(api.HTTPClientSettings{})
	must(err)

	// init store for the account
	dataStore := inmemory.NewInMemoryStore()

	// init account builder
	ab := builder.NewBuilder().WithAPI(iotaAPI).WithStore(dataStore)

	// init plugins
	eachThirdySeconds := time.Duration(30) * time.Second
	promoterReattacher := promoter.NewPromoter(ab.Settings(), eachThirdySeconds)
	receiveFilter := poller.NewPerTailReceiveEventFilter(true)
	transferPoller := poller.NewTransferPoller(ab.Settings(), receiveFilter, eachThirdySeconds)

	// build the account
	acc, err := ab.Build(promoterReattacher, transferPoller)
	must(err)

	// start the account
	must(acc.Start())
	// shutdown the account when the function returns
	defer acc.Shutdown()

	availBalance, err := acc.AvailableBalance()
	must(err)
	fmt.Printf("available balance %d\n", availBalance)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
