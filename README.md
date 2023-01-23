Golang package to provide automated zerodha login/authentication and ticker setup. 
Supports TOTP based authentication.
This package is for demo purpose only. Use at your own will & risk.

## Installation

```plaintext
go get github.com/parag-b/zerodha-api
```

## How to use

### To Authenticate - Zerodha API

*   [ ] Create variable

```plaintext
    var zd zerodhaapi.ZerodhaApi        // create variable
```

*   [ ] Provide user details

```plaintext
    zd.UserId = ""                      // provide credentials
    zd.Password = ""
    zd.TotpKey = ""
    zd.ApiKey = ""
    zd.ApiSecret = ""    
```

*   [ ] Authenticate

```plaintext
    err := zerodhaapi.New(&zd)          // Authenticate with zerodha
    if err != nil {
        fmt.Println(err)
    }
```

### Error responses

*   Invalid credentials - Check zerodha credentials.  (`zd.UserId`) and/or (`zd.Password`)
*   Invalid ToptKey - Check the seed key (`zd.TotpKey)` while enabling TOTP on zerodha
*   Invalid ApiKey -  Check the ApiKey generated at kite.trade (`zd.ApiKey)`
*   Invalid ApiSecret -  Check the ApiSecret generated at kite.trade (`zd.ApiSecret`)

---

## Using Ticker

*   [ ] Setup tokens for ticker subscription

```plaintext
zd.TickerSubscribeTokens = []uint32{8972034, 8972290}   // provide instruments
```

*   [ ] Create channel and share reference

```plaintext
var TicksCh = make(chan kitemodels.Tick, 1000)          // create buffered ch
zd.TickerCh = TicksCh                                   // assign the channel
```

*   [ ] Start ticker websocket

```plaintext
zd.StartTicker()                                        // start ticks websocket
```

*   [ ] Demo function to receive the ticks.

> Starts the ticker, waits for 10 sec and closes the tickers.

```plaintext
go demoTicksReceiver()                      // start ticks receiver
time.Sleep(10 * time.Second)                // wait for 10 seconds
zd.CloseTicker()                            // close ticker & ticks channel

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
```

## Examples

### Authentication

```plaintext
var zd zerodhaapi.ZerodhaApi        // create variable

zd.UserId = ""                      // provide credentials
zd.Password = ""
zd.TotpKey = ""
zd.ApiKey = ""
zd.ApiSecret = ""

err := zerodhaapi.New(&zd)          // Initiate connection
if err != nil {
	fmt.Println(err)
} 
```
### Call Zerodha APIs

```plaintext
margins, _ := zd.KiteConn.GetUserMargins()
fmt.Println(margins)
```


### Start Ticker Service

```plaintext

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
```

## Developer Settings

#### Submit package

> Commit changes
> 
> git tag v0.2.0
> 
> git push origin --tags
> 
> GOPROXY=proxy.golang.org go list -m github.com/parag-b/zerodha-api@v0.2.0

#### Use local package

> go mod edit -replace=github.com/parag-b/[zerodha-api@v0.0.0-unpublished](mailto:zerodha-api@v0.0.0-unpublished)\=/repos/zerodha-api/
> 
> go mod tidy
