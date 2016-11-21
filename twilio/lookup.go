package twilio

import (
	"net/url"
)

type LookupService struct {
	client *Client
}

type Carrier struct {
	MobileCountryCode string `json:"mobile_country_code"`
	MobileNetworkCode string `json:"mobile_network_code"`
	Name              string `json:"name"`
	Type              string `json:"type"`
}

type PhoneNumber struct {
	CountryCode    string  `json:"country_code"`
	PhoneNumber    string  `json:"phone_number"`
	NationalFormat string  `json:"national_format"`
	Carrier        Carrier `json:"carrier"`
}

func (ls *LookupService) LookupPhoneNumber(phoneNumber string) (*PhoneNumber, error) {
	data := url.Values{}
	data.Set("CountryCode", "US")
	data.Set("Type", "carrier")

	lookup := new(PhoneNumber)
	_, err := ls.client.LookupResource("PhoneNumbers", phoneNumber, data, lookup)
	if err != nil {
		return nil, err
	}
	return lookup, err
}
