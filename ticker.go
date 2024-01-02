package zerodhaapi

import (
	"strconv"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/rs/zerolog/log"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var pZerodhaApi *ZerodhaApi
var ticksRate ratecounter.RateCounter

// Triggered when any error is raised
func onError(err error) {
	log.Warn().Err(err).
		Msg("\nticker:false ConnError")
	pZerodhaApi.IsTickerConnected = false
}

// Triggered when websocket connection is closed
func onClose(code int, reason string) {
	log.Info().
		Bool("ticker", pZerodhaApi.IsTickerConnected).
		Int("ConnClosed", code).
		Str("ConnClosed", reason).
		Msg("ConnClosed ")
	pZerodhaApi.IsTickerConnected = false
}

// Triggered when connection is established and ready to send and accept data
func onConnect() {

	var m string

	err := pZerodhaApi.Ticker.Subscribe(pZerodhaApi.TickerSubscribeTokens)
	m = "default"
	if err != nil {
		pZerodhaApi.IsTickerConnected = false
	} else {
		pZerodhaApi.IsTickerConnected = true
	}

	err = pZerodhaApi.Ticker.SetMode("full", pZerodhaApi.TickerSubscribeTokens)
	m = "full"
	if err != nil {
		pZerodhaApi.IsTickerConnected = false
	} else {
		pZerodhaApi.IsTickerConnected = true

	}

	log.Info().
		Bool("ticker", pZerodhaApi.IsTickerConnected).
		Int("tokens-subscribed", len(pZerodhaApi.TickerSubscribeTokens)).
		Str("mode", m).
		Msg("tokens subscribed:" + strconv.Itoa(len(pZerodhaApi.TickerSubscribeTokens)))

	ticksRate = *ratecounter.NewRateCounter(1 * time.Second)
}

// Triggered when tick is recevived
func onTick(tick kitemodels.Tick) {

	ticksRate.Incr(1)
	pZerodhaApi.TicksPerSec = ticksRate.Rate()

	pZerodhaApi.TickerCh <- tick

	if pZerodhaApi.TickDebug {
		log.Info().
			Bool("ticker", pZerodhaApi.IsTickerConnected).
			Time("\nTime: ", tick.Timestamp.Time).
			Uint32("Instrument: ", tick.InstrumentToken).
			Float64("LastPrice: ", tick.LastPrice).
			Int64("ticksRatePerSec ", pZerodhaApi.TicksPerSec).
			Msg("Tick Received")
	}
}

// Triggered when reconnection is attempted which is enabled by default
func onReconnect(attempt int, delay time.Duration) {
	pZerodhaApi.IsTickerConnected = true
	log.Info().
		Bool("ticker", pZerodhaApi.IsTickerConnected).
		Int("attempt", attempt).
		Float64("delay", delay.Seconds()).
		Msg("ticker reconnect attempt")
}

// Triggered when maximum number of reconnect attempt is made and the program is terminated
func onNoReconnect(attempt int) {
	pZerodhaApi.IsTickerConnected = false
	log.Error().
		Bool("ticker", pZerodhaApi.IsTickerConnected).
		Int("attempt", attempt).
		Msg("ticker reconnect failed")
}

// Triggered when order update is received
func onOrderUpdate(order kiteconnect.Order) {
	log.Info().
		Str("ID:", order.OrderID).
		Msg("order updated")
}

// Registers instruments with zerodha for tick data
// Setups call backs
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
	log.Info().
		Bool("ticker", z.IsTickerConnected).
		Msg("ticker : started")
}

// Closes the ticker and channel
func (z *ZerodhaApi) CloseTicker() {

	if z.IsTickerConnected {
		z.Ticker.Close()
		time.Sleep(time.Second * 1)
		z.Ticker.Stop()
		time.Sleep(time.Second * 1)
		z.IsTickerConnected = false
	}

	// check if channel is not nil
	if z.TickerCh != nil {
		select {
		case <-z.TickerCh:
			// channel is closed
		default:
			// channel is still open, then close
			close(z.TickerCh)
		}
	}
	log.Info().
		Bool("ticker", z.IsTickerConnected).
		Msg("ticker channel : closed")
}
