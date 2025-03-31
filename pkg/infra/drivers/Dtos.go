package drivers

import (
	"fmt"
	"github.com/pgabrielgmdeveloper/freightQuote/configs"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"strconv"
	"strings"
)

type Address struct {
	Zipcode string `json:"zipcode" binding:"required"`
}

type Recipient struct {
	Address Address `json:"address"`
}

type Volume struct {
	Category      int     `json:"category" binding:"required"`
	Amount        int     `json:"amount" binding:"required, gt=0"`
	UnitaryWeight float64 `json:"unitary_weight" binding:"required, gt=0"`
	Price         float64 `json:"price" binding:"required, gt=0"`
	Sku           string  `json:"sku" binding:"required"`
	Height        float64 `json:"height" binding:"required, gt=0"`
	Width         float64 `json:"width" binding:"required, gt=0"`
	Length        float64 `json:"length" binding:"required, gt=0"`
}

type SimulateQuoteRequest struct {
	Recipient Recipient `json:"recipient" binding:"required"`
	Volumes   []Volume  `json:"volumes" binding:"required"`
}

type Carrier struct {
	Name     string  `json:"name"`
	Service  string  `json:"service"`
	Deadline int     `json:"deadline"`
	Price    float64 `json:"price"`
}

type SimulateQuoteResponse struct {
	Carrier []Carrier `json:"carrier"`
}

type CarrierMetricsResponse struct {
	Name       string
	AvgPrice   float64
	MaxPrice   float64
	MinPrice   float64
	TotalPrice float64
	TotalOffer int
}

type MetricsResponse struct {
	Carrier               []CarrierMetricsResponse
	GeneralAvgPrice       float64
	GeneralMinPrice       float64
	GeneralMaxPrice       float64
	GeneralMinCarrierName string
	GeneralMaxCarrierName string
}

func DomainMetricsToRequest(metrics quote.Metrics) MetricsResponse {
	return MetricsResponse{
		Carrier: func() []CarrierMetricsResponse {
			var carrierMetris []CarrierMetricsResponse
			for _, m := range metrics.Carrier {
				metricDomain := CarrierMetricsResponse{
					Name:       m.Name,
					AvgPrice:   m.AvgPrice,
					MaxPrice:   m.MaxPrice,
					MinPrice:   m.MinPrice,
					TotalPrice: m.TotalPrice,
					TotalOffer: m.TotalOffer,
				}
				carrierMetris = append(carrierMetris, metricDomain)

			}
			return carrierMetris
		}(),
		GeneralMinCarrierName: metrics.GeneralMinCarrierName,
		GeneralMaxCarrierName: metrics.GeneralMaxCarrierName,
		GeneralMinPrice:       metrics.GeneralMinPrice,
		GeneralMaxPrice:       metrics.GeneralMaxPrice,
		GeneralAvgPrice:       metrics.GeneralAvgPrice,
	}
}

func ConverterStrinToInZipcode(zipcode string) (int, error) {
	zipcodeResponse, err := strconv.Atoi(strings.TrimLeft(zipcode, "0"))
	if err != nil {
		return 0, fmt.Errorf("zipcode deve ser um valor apenas numerico e sem pontos mas foi enviado %s", zipcode)
	}
	return zipcodeResponse, nil
}

func RequestToDomainQuote(request SimulateQuoteRequest) (*quote.QuoteRequest, error) {
	cfg, err := configs.LoadConfig()
	if err != nil {
		return nil, err
	}

	zipcode, err := ConverterStrinToInZipcode(request.Recipient.Address.Zipcode)
	if err != nil {
		return nil, err
	}
	return &quote.QuoteRequest{
		Shipper: quote.Shipper{
			RegisteredNumber: cfg.RegisteredNumber,
			Token:            cfg.TokenAPI,
			PlatformCode:     cfg.PlatformCode,
		},
		Recipient: quote.Recipient{
			Type:    0,
			Country: "BRA",
			Zipcode: zipcode,
		},
		Dispatchers: []quote.Dispatcher{{
			RegisteredNumber: cfg.RegisteredNumber,
			Zipcode:          1311000,
			Volumes: func() []quote.Volume {
				var volumes []quote.Volume
				for _, v := range request.Volumes {
					volume := quote.Volume{
						Category:      strconv.Itoa(v.Category),
						Amount:        v.Amount,
						UnitaryWeight: v.UnitaryWeight,
						Width:         v.Width,
						Height:        v.Height,
						Length:        v.Length,
						UnitaryPrice:  v.Price,
					}
					volumes = append(volumes, volume)
				}
				return volumes
			}(),
		}},
	}, nil
}

func DomainToSimulateQuoteResponse(offer []quote.Offer) SimulateQuoteResponse {
	return SimulateQuoteResponse{
		Carrier: func() []Carrier {
			var carriers []Carrier
			for _, o := range offer {
				carrier := Carrier{
					Name:     o.Carrier,
					Service:  o.Service,
					Deadline: o.DeliveryTime,
					Price:    o.FinalPrice,
				}
				carriers = append(carriers, carrier)
			}
			return carriers
		}(),
	}
}
