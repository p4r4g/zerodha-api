package zerodhaapi

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/asmcos/requests"
	totp "github.com/pquerna/otp/totp"
	kiteconnect "github.com/zerodha/gokiteconnect/v4"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

const (
	kiteLoginUrl     = "https://kite.zerodha.com/api/login"
	twofactorAuthUrl = "https://kite.zerodha.com/api/twofa"
	reqTokenUrl      = "https://kite.zerodha.com/connect/login?v=3&api_key="
)

// struct has 3 categories of variables
// 1. User information to be used for authentication
// 2. Tikcer config varaibles
// 3. Instance variables, there are updated by this package. For read only by the application
type ZerodhaApi struct {
	// authentication settings (zerodha user)
	UserId    string
	Password  string
	ApiKey    string
	ApiSecret string
	TotpKey   string

	// ticker settings
	TickerSubscribeTokens []uint32
	TickerCh              chan kitemodels.Tick
	TickDebug             bool // to print every tick

	// kite instance data
	Ticker            *kiteticker.Ticker
	KiteConn          *kiteconnect.Client
	TicksPerSec       int64 // number of ticks received per seconds
	IsTickerConnected bool
	IsKiteAuth        bool
	KiteReqId         string
	KiteReqToken      string
	AccessToken       string
}

// Authencticates the user and stores the access token
func New(za *ZerodhaApi) error {

	za.IsKiteAuth = false
	req := requests.Requests()

	if err := za.doCredentialAuth(req); err != nil {
		return err
	}

	if err := za.doTwoFactorAuth(req); err != nil {
		return err
	}

	if err := za.getReqToken(req); err != nil {
		return err
	}

	za.KiteConn = kiteconnect.New(za.ApiKey)
	data, err := za.KiteConn.GenerateSession(za.KiteReqToken, za.ApiSecret)
	if err != nil {
		return errors.New("invalid api secret")
	}

	za.KiteConn.SetAccessToken(data.AccessToken)
	za.AccessToken = data.AccessToken

	za.IsKiteAuth = true

	return nil
}

func (z *ZerodhaApi) doCredentialAuth(req *requests.Request) error {

	data := requests.Datas{
		"user_id":  z.UserId,
		"password": z.Password,
	}

	resp, err := req.Post(kiteLoginUrl, data)

	if (err != nil) || (resp.R.StatusCode != 200) {
		return errors.New("invalid credentials")
	}

	z.KiteReqId = extractValue(resp.Text(), "request_id")

	return nil
}

func (z *ZerodhaApi) doTwoFactorAuth(req *requests.Request) error {

	totpCode, err := totp.GenerateCode(z.TotpKey, time.Now())

	if err != nil {
		return errors.New("invalid two factor authentication")
	}

	data := requests.Datas{
		"user_id":     z.UserId,
		"request_id":  z.KiteReqId,
		"twofa_value": totpCode,
	}

	resp, err := req.Post(twofactorAuthUrl, data)

	if (err != nil) || (resp.R.StatusCode != 200) {
		return errors.New("invalid two factor authentication")
	}
	return nil
}

func (z *ZerodhaApi) getReqToken(req *requests.Request) error {

	req.SetTimeout(5)
	resp, err := req.Get(reqTokenUrl + z.ApiKey)

	if err != nil {
		// Not able to reach http handler, parse the URL for key
		arr := strings.Split(err.Error(), `"`) // split on '&'
		val, key := extractKeyValue(arr[1], "request_token")

		if key {
			z.KiteReqToken = val
			return nil
		}
	} else {
		// Able to reach http handler, parse the response
		if resp.R.StatusCode != 200 {
			return errors.New("invalid api key")
		}

		m, _ := url.ParseQuery(resp.R.Request.URL.RawQuery)
		if _, ok := m["request_token"]; ok {
			z.KiteReqToken = m["request_token"][0]
			return nil
		}

	}
	return errors.New("invalid api key")
}

// Returns the Equity part of margins data
func (z ZerodhaApi) CashBalance() (float64, error) {

	if !z.IsKiteAuth {
		return 0, errors.New("account not authenticated")
	}
	margins, err := z.KiteConn.GetUserMargins()
	if err != nil {
		return 0, err
	}
	return margins.Equity.Net, nil
}

func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}

func extractKeyValue(body string, key string) (string, bool) {
	arr := strings.Split(body, `&`) // split on '&'

	for index := range arr {
		if strings.Contains(arr[index], key) { // if key is found
			arrVal := strings.Split(arr[index], `=`)
			return arrVal[1], true // Extract value
		}
	}
	return "", false
}
