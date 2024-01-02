package zerodhaapi_test

import (
	"testing"
	"time"

	zerodhaapi "github.com/parag-b/zerodha-api"
	"github.com/rs/zerolog/log"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
)

var zd zerodhaapi.ZerodhaApi // create variable

func TestStartTicker(t *testing.T) {
	// read credentials from config file

	zd.UserId = "" // provide credentials
	zd.Password = ""
	zd.TotpKey = ""
	zd.ApiKey = ""
	zd.ApiSecret = ""

	err := zerodhaapi.New(&zd) // Initiate connection
	if err != nil {
		log.Error().Err(err).
			Bool("ticker", zd.IsTickerConnected).
			Msg("authentication failed")
	}

	margins, _ := zd.KiteConn.GetUserMargins()
	log.Info().
		Float64("cash", margins.Equity.Available.Cash).
		Float64("collateral", margins.Equity.Available.Collateral).
		Float64("adHoc margin", margins.Equity.Available.AdHocMargin).
		Float64("intraday payin", margins.Equity.Available.IntradayPayin).
		Float64("live balance", margins.Equity.Available.LiveBalance).
		Float64("opening balance", margins.Equity.Available.OpeningBalance).
		Msg("margins")

	// ticker settings
	var TicksCh = make(chan kitemodels.Tick, 1000)        // create buffered ch
	zd.TickerSubscribeTokens = []uint32{8972034, 8972290} // provide instruments
	zd.TickerCh = TicksCh                                 // assign the channel
	zd.StartTicker()                                      // start ticks websocket

	go demoTicksReceiver()       // start ticks receiver
	time.Sleep(10 * time.Second) // wait for 10 seconds
	zd.CloseTicker()             // close ticker & ticks channel
}

// to receive ticks, closes when channel is closed
func demoTicksReceiver() {
	for v := range zd.TickerCh { // read from tick channel
		log.Info().
			Str("\nTime: ", v.Timestamp.String()). // convert v.Timestamp to string
			Uint32("Instrument: ", v.InstrumentToken).
			Float64("LastPrice: ", v.LastPrice).
			Int64("ticksRatePerSec:%d ,", zd.TicksPerSec).
			Msg("Tick Received")
	}
	log.Info().
		Msg("ticks channel closed")
}
