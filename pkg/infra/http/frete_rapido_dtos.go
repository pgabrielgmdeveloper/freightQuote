package http

import "github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"

type FreteRapidoApiRequest struct {
	Shipper        Shipper      `json:"shipper"`
	Recipient      Recipient    `json:"recipient"`
	Dispatchers    []Dispatcher `json:"dispatchers"`
	SimulationType []int        `json:"simulation_type"`
}

type FreteRapidoApiResponse struct {
	Dispatchers []DispatcherResponse `json:"dispatchers"`
}

type Shipper struct {
	RegisteredNumber string `json:"registered_number"`
	Token            string `json:"token"`
	PlatformCode     string `json:"platform_code"`
}

type Recipient struct {
	Type    int    `json:"type"`
	Zipcode int    `json:"zipcode"`
	Country string `json:"country"`
}

type Dispatcher struct {
	RegisteredNumber string             `json:"registered_number"`
	Zipcode          int                `json:"zipcode"`
	Volumes          []VolumeApiRequest `json:"volumes"`
}
type DispatcherResponse struct {
	Offer []OfferResponse `json:"offers"`
}

type OfferResponse struct {
	FinalPrice   float64              `json:"final_price"`
	Carrier      CarrierResponse      `json:"carrier"`
	Service      string               `json:"service"`
	DeliveryTime DeliveryTimeResponse `json:"delivery_time"`
}

type DeliveryTimeResponse struct {
	Days    int `json:"days"`
	Hours   int `json:"hours"`
	Minutes int `json:"minutes"`
}

type CarrierResponse struct {
	Name string `json:"name"`
}

type VolumeApiRequest struct {
	Amount        int     `json:"amount"`
	Category      string  `json:"category"`
	Height        float64 `json:"height"`
	Width         float64 `json:"width"`
	Length        float64 `json:"length"`
	UnitaryPrice  float64 `json:"unitary_price"`
	UnitaryWeight float64 `json:"unitary_weight"`
}

func DomainToVolumeContract(volumeDomain quote.Volume) VolumeApiRequest {
	return VolumeApiRequest{
		Category:      volumeDomain.Category,
		Amount:        volumeDomain.Amount,
		Height:        volumeDomain.Height,
		Length:        volumeDomain.Length,
		UnitaryPrice:  volumeDomain.UnitaryPrice,
		UnitaryWeight: volumeDomain.UnitaryWeight,
		Width:         volumeDomain.Width,
	}
}

func DomainToDispacherContract(dispacherDomain quote.Dispatcher) Dispatcher {
	return Dispatcher{
		RegisteredNumber: dispacherDomain.RegisteredNumber,
		Zipcode:          dispacherDomain.Zipcode,
		Volumes: func() []VolumeApiRequest {
			var volumes []VolumeApiRequest
			for _, volume := range dispacherDomain.Volumes {
				volumes = append(volumes, DomainToVolumeContract(volume))
			}
			return volumes
		}(),
	}
}

func DomainToFreteRapidoContractRequest(request quote.QuoteRequest) FreteRapidoApiRequest {
	return FreteRapidoApiRequest{
		Shipper: Shipper{
			RegisteredNumber: request.Shipper.RegisteredNumber,
			Token:            request.Shipper.Token,
			PlatformCode:     request.Shipper.PlatformCode,
		},
		Recipient: Recipient{
			Type:    request.Recipient.Type,
			Country: request.Recipient.Country,
			Zipcode: request.Recipient.Zipcode,
		},
		Dispatchers: func() []Dispatcher {
			var dispatchers []Dispatcher
			for _, d := range request.Dispatchers {
				dispatcher := DomainToDispacherContract(d)
				dispatchers = append(dispatchers, dispatcher)
			}
			return dispatchers
		}(),
		SimulationType: []int{0},
	}
}

func FreteApiResponseToDomainOffer(response FreteRapidoApiResponse) []quote.Offer {
	var offers []quote.Offer
	for _, d := range response.Dispatchers {
		for _, offer := range d.Offer {
			offerDomain := quote.Offer{
				Carrier:    offer.Carrier.Name,
				Service:    offer.Service,
				FinalPrice: offer.FinalPrice,
				DeliveryTime: func() int {
					if offer.DeliveryTime.Days == 0 {
						return 1
					}
					return offer.DeliveryTime.Days
				}(),
			}
			offers = append(offers, offerDomain)
		}
	}
	return offers
}
