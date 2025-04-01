package http

import "github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"

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
