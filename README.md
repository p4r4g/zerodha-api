
# How to use




func main() {
	
    var zd zerodhaapi.ZerodhaApi        // create variable

	zd.UserId = ""                      // provide credentials
	zd.Password = ""
	zd.TotpKey = ""
	zd.ApiKey = ""
	zd.ApiSecret = ""

	err := zerodhaapi.New(&zd)          // Initiate connection
	if err != nil {
		fmt.Println(err)
	} else {
		cash, _ := zd.CashBalance()     // on conn, prints cash in account
		fmt.Println("\nCash Balance: ", cash)
	}

    // ------------------ FOR TICKs ------------------
	// ticker settings
	var TicksCh = make(chan kitemodels.Tick, 1000)          // create buffered ch
	zd.TickerSubscribeTokens = []uint32{8972034, 8972290}   // provide instruments
	zd.TickerCh = TicksCh                                   // assign the channel
	zd.StartTicker()                                        // start ticks websocket


	go demoTicksReceiver()                      // start ticks receiver
	time.Sleep(10 * time.Second)                // wait for 10 seconds
	zd.CloseTicker()                            // close ticker & ticks channel
}

// to receive ticks, closes when channel is closed
func demoTicksReceiver() {

	for v := range zd.TickerCh { // read from tick channel

		fmt.Println("\nTime: ", v.Timestamp,
			"Instrument: ", v.InstrumentToken,
			"LastPrice: ", v.LastPrice)
		fmt.Printf("ticksRatePerSec:%d ,", zd.TicksPerSec)
	}

	fmt.Println("ticks channel closed")
}

# Submit package


Commit changes
git tag v0.2.0
git push origin --tags

GOPROXY=proxy.golang.org go list -m github.com/parag-b/zerodha-api@v0.2.0


# Use local package

go mod edit -replace=github.com/parag-b/zerodha-api@v0.0.0-unpublished=/repos/zerodha-api/
go mod tidy
