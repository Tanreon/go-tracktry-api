package tracktry

import "regexp"

var FedexPatterns = []*regexp.Regexp{ // https://github.com/jkeen/tracking_number_data/blob/main/couriers/fedex.json
	regexp.MustCompile(`(?i)^\s*(?P<SerialNumber>([0-9]\s*){11})(?P<CheckDigit>[0-9]\s*)$`),                                                                                                                                   // FedEx Express (12)
	regexp.MustCompile(`(?i)^\s*1\s*0\s*0\s*[0-9]\s*[0-9]\s*([0-9]\s*){10}(?P<DestinationZip>([0-9]\s*){5})(?P<SerialNumber>([0-9]\s*){13})(?P<CheckDigit>[0-9]\s*)$`),                                                        // FedEx Express (34)
	regexp.MustCompile(`(?i)^\s*(?P<SerialNumber>([0-9]\s*){14})(?P<CheckDigit>([0-9]\s*))$`),                                                                                                                                 // FedEx Ground
	regexp.MustCompile(`(?i)^\s*(?P<ShippingContainerType>([0-9]\s*){2})(?P<SerialNumber>([0-9]\s*){15})(?P<CheckDigit>[0-9]\s*)$`),                                                                                           // FedEx Ground (SSCC-18)
	regexp.MustCompile(`(?i)^\s*(?P<ApplicationIdentifier>9\s*6\s*)(?P<SCNC>([0-9]\s*){2})(?P<ServiceType>([0-9]\s*){3})(?P<SerialNumber>(?P<ShipperId>([0-9]\s*){7})(?P<PackageId>([0-9]\s*){7}))(?P<CheckDigit>[0-9]\s*)$`), // FedEx Ground 96 (22)
	regexp.MustCompile(`(?i)^\s*(?P<ApplicationIdentifier>9\s*6\s*)(?P<SCNC>([0-9]\s*){2})([0-9]\s*){5}(?P<GSN>([0-9]\s*){10})[0-9]\s*(?P<SerialNumber>([0-9]\s*){13})(?P<CheckDigit>[0-9]\s*)$`),                             // FedEx Ground GSN
}

var UPSPatterns = []*regexp.Regexp{ // https://github.com/jkeen/tracking_number_data/blob/main/couriers/ups.json
	regexp.MustCompile(`(?i)^\s*1\s*Z\s*(?P<SerialNumber>(?P<ShipperId>(?:[A-Z0-9]\s*){6,6})(?P<ServiceType>(?:[A-Z0-9]\s*){2,2})(?P<PackageId>(?:[A-Z0-9]\s*){7,7}))(?P<CheckDigit>[A-Z0-9]\s*)$`),
}

var USPSPatterns = []*regexp.Regexp{ // https://github.com/jkeen/tracking_number_data/blob/main/couriers/usps.json
	regexp.MustCompile(`(?i)^\s*(?P<SerialNumber>(?P<ServiceType>([0-9]\s*){2})(?P<ShipperId>([0-9]\s*){9})(?P<PackageId>([0-9]\s*){8}))(?P<CheckDigit>[0-9]\s*)$`),                                                                                                                                                             // 20 digit USPS numbers
	regexp.MustCompile(`(?i)^\s*(?P<RoutingApplicationId>4\s*2\s*0\s*)(?P<DestinationZip>([0-9]\s*){5})(?P<RoutingNumber>([0-9]\s*){4})(?P<SerialNumber>(?P<ApplicationIdentifier>9\s*[2345]\s*)?(?P<ShipperId>([0-9]\s*){8})(?P<PackageId>([0-9]\s*){11}))(?P<CheckDigit>[0-9]\s*)$`),                                          // variation on 34 digit USPS IMpd numbers
	regexp.MustCompile(`(?i)^\s*(?:(?P<RoutingApplicationId>4\s*2\s*0\s*)(?P<DestinationZip>([0-9]\s*){5}))?(?P<SerialNumber>(?P<ApplicationIdentifier>9\s*[12345]\s*)?(?P<SCNC>([0-9]\s*){2})(?P<ServiceType>([0-9]\s*){2})(?P<ShipperId>([0-9]\s*){8})(?P<PackageId>([0-9]\s*){11}|([0-9]\s*){7}))(?P<CheckDigit>[0-9]\s*)$`), // USPS now calls this the IMpd barcode format
}

var OnTracPatterns = []*regexp.Regexp{ // https://www.trackingmore.com/tracking-status-detail-en-253.html
	regexp.MustCompile(`(?i)^\s*(([DC])\d{14})$`),
}

var LaserShipPatterns = []*regexp.Regexp{ // https://www.trackingmore.com/tracking-status-detail-en-268.html
	regexp.MustCompile(`(?i)^\s*LW\d{8,}|\s*LX\d{8,}|\s*1LS\d{12,}$`),
}

var DHLPatterns = []*regexp.Regexp{ // https://github.com/jkeen/tracking_number_data/blob/main/couriers/dhl.json
	regexp.MustCompile(`(?i)^\s*(?P<SerialNumber>([0-9]\s*){9})(?P<CheckDigit>([0-9]\s*))$`), // DHL Express
	regexp.MustCompile(`(?i)^\s*(?P<SerialNumber>([0-9]\s*){10})(?P<CheckDigit>[0-9]\s*)$`),  // DHL Express Air
}

var ShiptPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^\s*SHIPT\d{11}$`),
}

var AmazonPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^\s*TBA\d{10,12}$`),
}

func (t *Tracktry) recognizeCarrier() (carrier string, err error) {
	carriersPatterns := map[string][]*regexp.Regexp{
		"fedex":     FedexPatterns,
		"ups":       UPSPatterns,
		"usps":      USPSPatterns,
		"ontrac":    OnTracPatterns,
		"lasership": LaserShipPatterns,
		"dhl":       DHLPatterns,
		"shipt":     ShiptPatterns,
		"amazon":    AmazonPatterns,
	}

	for carrier, carrierPatterns := range carriersPatterns {
		for _, carrierPattern := range carrierPatterns {
			if carrierPattern.MatchString(t.code) {
				return carrier, err
			}
		}
	}

	return carrier, ErrCarrierUnrecognized
}

func (t *Tracktry) IsValid() (isValid bool) {
	_, err := t.recognizeCarrier()
	if err != nil {
		return false
	}

	return true

	//carriersPatterns := [][]*regexp.Regexp{
	//	FedexPatterns,
	//	UPSPatterns,
	//	USPSPatterns,
	//	OnTracPatterns,
	//	LaserShipPatterns,
	//	DHLPatterns,
	//	ShiptPatterns,
	//	AmazonPatterns,
	//}
	//
	//for _, carrierPatterns := range carriersPatterns {
	//	for _, carrierPattern := range carrierPatterns {
	//		if carrierPattern.MatchString(t.code) {
	//			return true
	//		}
	//	}
	//}
	//
	//return false
}
