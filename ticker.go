package zerodhaapi

import (
	"fmt"
	"time"

	"github.com/paulbellamy/ratecounter"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var pZerodhaApi *ZerodhaApi
var ticksRate ratecounter.RateCounter

type Tick struct {
	Timestamp       time.Time
	InstrumentToken uint32
	LastTradedPrice float64
	LastPrice       float64
	BuyDemand       uint32
	SellDemand      uint32
	TradesTillNow   uint32
	OpenInterest    uint32
}

// Triggered when any error is raised
func onError(err error) {
	fmt.Println("\nticker:false ConnError: ", err)
	pZerodhaApi.IsTickerConnected = false
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	fmt.Println("\nticker:false ConnClosed: ", code, reason)
	pZerodhaApi.IsTickerConnected = false
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {

	err := pZerodhaApi.Ticker.Subscribe(pZerodhaApi.TickerSubscribeTokens)
	if err != nil {
		fmt.Println("ticker:false tokens-subscribed:", len(pZerodhaApi.TickerSubscribeTokens), " mode:Default")
		pZerodhaApi.IsTickerConnected = false
	} else {
		fmt.Println("ticker:true tokens-subscribed:", len(pZerodhaApi.TickerSubscribeTokens), " mode:Default")
		pZerodhaApi.IsTickerConnected = true
	}
	err = pZerodhaApi.Ticker.SetMode("full", pZerodhaApi.TickerSubscribeTokens)
	if err != nil {
		fmt.Println("ticker:false tokens-subscribed:", len(pZerodhaApi.TickerSubscribeTokens), " mode:FULL")
		pZerodhaApi.IsTickerConnected = false
	} else {
		fmt.Println("ticker:true tokens-subscribed:", len(pZerodhaApi.TickerSubscribeTokens), " mode:FULL")
		pZerodhaApi.IsTickerConnected = true
	}

	ticksRate = *ratecounter.NewRateCounter(1 * time.Second)
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {

	ticksRate.Incr(1)
	pZerodhaApi.TicksPerSec = ticksRate.Rate()

	pZerodhaApi.TickerCh <- tick

	if pZerodhaApi.TickDebug {
		fmt.Println("\nTime: ", tick.Timestamp.Time,
			"Instrument: ", tick.InstrumentToken,
			"LastPrice: ", tick.LastPrice)
		fmt.Printf("\nticksRatePerSec:%d ,", pZerodhaApi.TicksPerSec)
	}
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	pZerodhaApi.IsTickerConnected = true
	fmt.Printf("ticker:false ConnReconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	pZerodhaApi.IsTickerConnected = false
	fmt.Printf("ticker:false ConnNoReconnect - Maximum no of reconnect attempt reached: %d", attempt)
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	fmt.Println("order update ID:", order.OrderID)
}

func (z *ZerodhaApi) StartTicker() {

	if len(z.TickerSubscribeTokens) > 0 { // tokens provided?
		pZerodhaApi = z

		// Create new Kite ticker instance
		z.Ticker = kiteticker.New(z.ApiKey, z.AccessToken)

		// Assign callbacks
		z.Ticker.OnError(onError)
		z.Ticker.OnClose(onClose)
		z.Ticker.OnConnect(onConnect)
		z.Ticker.OnReconnect(onReconnect)
		z.Ticker.OnNoReconnect(onNoReconnect)
		z.Ticker.OnTick(onTick)
		z.Ticker.OnOrderUpdate(onOrderUpdate)

		go z.Ticker.Serve()
		return
	}
	fmt.Printf("ticker:false TickerSubscribeTokens:<empty>")
}

func (z *ZerodhaApi) CloseTicker() {

	if !z.IsTickerConnected {
		z.Ticker.Stop()
	}
	time.Sleep(time.Second * 3) // delay for ticker to terminte connection before we close channel

	select {
	case <-z.TickerCh:
		// channel is closed, this is executed
	default:
		// channel is still open, the previous case is not executed
		close(z.TickerCh)
	}

	fmt.Println("ticker:false channel:closed")
}
