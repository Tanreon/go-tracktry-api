package tracktry

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	HttpRunner "github.com/Tanreon/go-http-runner"
	log "github.com/sirupsen/logrus"
)

func (t *Tracktry) IsDelivered() (isDelivered bool, err error) {
	type RealtimeRequest struct {
		TrackingNumber           string `json:"tracking_number"`
		CarrierCode              string `json:"carrier_code"`
		DestinationCode          string `json:"destination_code"`
		TrackingShipDate         string `json:"tracking_ship_date"`
		TrackingPostalCode       string `json:"tracking_postal_code"`
		SpecialNumberDestination string `json:"specialNumberDestination"`
		Order                    string `json:"order"`
		OrderCreateTime          string `json:"order_create_time"`
		Lang                     string `json:"lang"`
	}

	type MetaRealtimeResponse struct {
		Code    int    `json:"code"`
		Type    string `json:"type"`
		Message string `json:"message"`
	}
	type ItemsData struct {
		Id                        string      `json:"id"`
		TrackingNumber            string      `json:"tracking_number"`
		CarrierCode               string      `json:"carrier_code"`
		OrderCreateTime           string      `json:"order_create_time"`
		DestinationCode           string      `json:"destination_code"`
		Status                    string      `json:"status"`
		TrackUpdate               bool        `json:"track_update"`
		OriginalCountry           string      `json:"original_country"`
		DestinationCountry        string      `json:"destination_country"`
		ItemTimeLength            int         `json:"itemTimeLength"`
		StayTimeLength            int         `json:"stayTimeLength"`
		ServiceCode               interface{} `json:"service_code"`
		PackageStatus             interface{} `json:"packageStatus"`
		Substatus                 interface{} `json:"substatus"`
		LastMileTrackingSupported interface{} `json:"last_mile_tracking_supported"`
		/*OriginInfo                struct {
			ItemReceived       string      `json:"ItemReceived"`
			ItemDispatched     interface{} `json:"ItemDispatched"`
			DepartfromAirport  interface{} `json:"DepartfromAirport"`
			ArrivalfromAbroad  interface{} `json:"ArrivalfromAbroad"`
			CustomsClearance   interface{} `json:"CustomsClearance"`
			DestinationArrived interface{} `json:"DestinationArrived"`
			Weblink            string      `json:"weblink"`
			Phone              interface{} `json:"phone"`
			CarrierCode        string      `json:"carrier_code"`
			Trackinfo          []struct {
				StatusDescription string `json:"StatusDescription"`
				Date              string `json:"Date"`
				Details           string `json:"Details"`
				CheckpointStatus  string `json:"checkpoint_status"`
				Substatus         string `json:"substatus"`
				ItemNode          string `json:"ItemNode,omitempty"`
			} `json:"trackinfo"`
		} `json:"origin_info"`*/
		/*DestinationInfo struct {
			ItemReceived       interface{} `json:"ItemReceived"`
			ItemDispatched     interface{} `json:"ItemDispatched"`
			DepartfromAirport  interface{} `json:"DepartfromAirport"`
			ArrivalfromAbroad  interface{} `json:"ArrivalfromAbroad"`
			CustomsClearance   interface{} `json:"CustomsClearance"`
			DestinationArrived interface{} `json:"DestinationArrived"`
			Weblink            interface{} `json:"weblink"`
			Phone              interface{} `json:"phone"`
			CarrierCode        interface{} `json:"carrier_code"`
			Trackinfo          interface{} `json:"trackinfo"`
		} `json:"destination_info"`*/
		LastEvent      string `json:"lastEvent"`
		LastUpdateTime string `json:"lastUpdateTime"`
	}
	type DataRealtimeResponse struct {
		Items []ItemsData `json:"items"`
	}
	type RealtimeResponse struct {
		Meta MetaRealtimeResponse `json:"meta"`
		Data DataRealtimeResponse `json:"data"`
	}

	//

	t.apiLimiter.Take()

	if !t.IsValid() {
		return isDelivered, ErrTrackCodeIsNotValid
	}

	carrier, err := t.RecognizeCarrier()
	if err != nil {
		return isDelivered, err
	}

	switch carrier {
	case "fedex":
		fallthrough
	case "ups":
		fallthrough
	case "usps":
		fallthrough
	case "ontrac":
		fallthrough
	case "lasership":
		fallthrough
	case "dhl":
		realtimeRequest := RealtimeRequest{
			TrackingNumber: t.code,
			CarrierCode:    carrier,
			Lang:           "en",
		}
		realtimeRequestBytes, err := json.Marshal(realtimeRequest)
		if err != nil {
			return isDelivered, fmt.Errorf("marshaling error: %w", err)
		}

		jsonRequestData := HttpRunner.NewJsonRequestData(API_SERVER + "/v1/trackings/realtime")
		jsonRequestData.SetHeaders(map[string]string{
			"Tracktry-Api-Key": t.apiToken,
		})
		jsonRequestData.SetValue(realtimeRequestBytes)
		jsonRequestData.SetTimeoutOption(time.Second * 120)

		response, err := t.runner.PostJson(jsonRequestData)
		if err != nil {
			return isDelivered, fmt.Errorf("/v1/trackings/realtime response error: %w", err)
		}

		var realtimeResponse RealtimeResponse
		if err := json.Unmarshal(response.Body(), &realtimeResponse); err != nil {
			log.Debugf("response code: %d, response body: %v", response.StatusCode(), string(response.Body()))
			return isDelivered, fmt.Errorf("unmarshaling error: %w", err)
		}
		if realtimeResponse.Meta.Code != 200 {
			return isDelivered, fmt.Errorf("api error: %s", realtimeResponse.Meta.Message)
		}

		for _, item := range realtimeResponse.Data.Items {
			if strings.EqualFold(item.TrackingNumber, t.code) {
				return strings.EqualFold(item.Status, "delivered"), err
			}
		}

		return isDelivered, ErrApiUnknownError
	case "shipt":
		fallthrough
	case "amazon":
		return isDelivered, ErrCarrierDisabled
	}

	return isDelivered, ErrCarrierUnrecognized
}
