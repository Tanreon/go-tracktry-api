package tracktry

import (
	"errors"
	"time"

	"go.uber.org/ratelimit"

	HttpRunner "github.com/Tanreon/go-http-runner"
)

const API_SERVER = "https://api.tracktry.com"

const FEDEX_CARRIER = "fedex"
const UPS_CARRIER = "ups"
const USPS_CARRIER = "usps"
const ONTRAC_CARRIER = "ontrac"
const LASERSHIP_CARRIER = "lasership"
const DHL_CARRIER = "dhl"
const SHIPT_CARRIER = "shipt"
const AMAZON_CARRIER = "amazon"

var ErrCarrierUnrecognized = errors.New("CARRIER_UNRECOGNIZED")
var ErrTrackCodeIsNotValid = errors.New("TRACK_CODE_IS_NOT_VALID")
var ErrCarrierDisabled = errors.New("CARRIER_DISABLED")
var ErrApiUnknownError = errors.New("API_UNKNOWN_ERROR")

type Tracktry struct {
	apiLimiter ratelimit.Limiter
	runner     HttpRunner.IHttpRunner
	apiToken   string
	code       string
}

type ITracktry interface {
	Code() string
	IsValid() (isValid bool)
	IsDelivered() (isDelivered bool, err error)
	RecognizeCarrier() (carrier string, err error)
}

func (t *Tracktry) Code() string {
	return t.code
}

func NewTracktry(runner *HttpRunner.IHttpRunner, apiToken string, apiLimit int, code string) ITracktry {
	return &Tracktry{
		apiLimiter: ratelimit.New(apiLimit, ratelimit.Per(time.Second)),
		runner:     *runner,
		apiToken:   apiToken,
		code:       code,
	}
}
