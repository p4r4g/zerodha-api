Golang package to provide fully automated zerodha login/authentication and ticker.
_This package is for demo purpose only. Use at your own will & risk._

#### Installation

```plaintext
go get github.com/parag-b/zerodha-api
```

## How to use

#### To Authenticate - Zerodha API

```plaintext
    var zd zerodhaapi.ZerodhaApi        // create variable

    zd.UserId = ""                      // provide credentials
    zd.Password = ""
    zd.TotpKey = ""
    zd.ApiKey = ""
    zd.ApiSecret = ""    

    err := zerodhaapi.New(&zd)          // Authenticate with zerodha
    if err != nil {
        fmt.Println(err)
    }
```

##### Error responses

- Invalid credentials - Check zerodha credentials.  (`zd.UserId`) and/or (`zd.Password`)
- Invalid ToptKey - Check the seed key (`zd.TotpKey)` while enabling TOTP on zerodha
- Invalid ApiKey -  Check the ApiKey generated at kite.trade (`zd.ApiKey)`
- Invalid ApiSecret -  Check the ApiSecret generated at kite.trade (`zd.ApiSecret`)

---

#### Using Ticker

```plaintext
zd.TickerSubscribeTokens = []uint32{8972034, 8972290}   // provide instruments

var TicksCh = make(chan kitemodels.Tick, 1000)          // create buffered ch
zd.TickerCh = TicksCh                                   // assign the channel
zd.StartTicker()                                        // start ticks websocket
```

#### Zerodha APIs

```plaintext
margins, _ := zd.KiteConn.GetUserMargins()
fmt.Println(margins)
```

---

_Refer ticker_test.go for complete example_
