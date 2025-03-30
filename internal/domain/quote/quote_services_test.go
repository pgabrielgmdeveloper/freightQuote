package quote

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockSimulatePort struct {
	mock.Mock
}

func (m *MockSimulatePort) Execute(req QuoteRequest) ([]Offer, error) {
	args := m.Called(req)
	return args.Get(0).([]Offer), args.Error(1)
}

type MockMetricsPort struct {
	mock.Mock
}

func (m *MockMetricsPort) Execute(lastQuotes int) ([]Metrics, error) {
	args := m.Called(lastQuotes)
	return args.Get(0).([]Metrics), args.Error(1)
}

func ValidRequest() QuoteRequest {
	validReq := QuoteRequest{
		Shipper: Shipper{
			RegisteredNumber: "25438296000158",
			Token:            "1d52a9b6b78cf07b08586152459a5c90",
			PlatformCode:     "5AKVkHqCn",
		},
		Recipient: Recipient{
			Type:             0,
			Country:          "BRA",
			Zipcode:          1311000,
			RegisteredNumber: "",
		},
		Dispatchers: []Dispatcher{
			{
				RegisteredNumber: "25438296000158",
				Zipcode:          49157021,
				Volumes: []Volume{
					{
						Category:      "7",
						Amount:        1,
						UnitaryWeight: 5,
						UnitaryPrice:  349,
						Height:        0.2,
						Width:         0.2,
						Length:        0.2,
					},
					{
						Category:      "7",
						Amount:        2,
						UnitaryWeight: 4,
						UnitaryPrice:  556,
						Height:        0.4,
						Width:         0.6,
						Length:        0.15,
					},
				},
			},
		},
	}
	return validReq
}

func InvalidRequest() QuoteRequest {
	invalidReq := QuoteRequest{
		Shipper: Shipper{
			RegisteredNumber: "123",
		},
	}
	return invalidReq
}

func TestSimulateQuote_Success(t *testing.T) {

	mockSimulate := new(MockSimulatePort)
	qs := NewQuoteService(mockSimulate, nil)
	validReq := ValidRequest()

	mockSimulate.On("Execute", validReq).Return([]Offer{
		{Carrier: "Correios", FinalPrice: 50.99, DeliveryTime: 1, Service: "SEDEX"},
	}, nil)

	offers, err := qs.SimulateQuote(validReq)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(offers))
	assert.Equal(t, "Correios", offers[0].Carrier)
	assert.Equal(t, "SEDEX", offers[0].Service)
	mockSimulate.AssertExpectations(t)
}

func TestSimulateQuote_ValidationError(t *testing.T) {
	qs := NewQuoteService(nil, nil) // Porta não será usada

	invalidReq := InvalidRequest()

	_, err := qs.SimulateQuote(invalidReq)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CNPJ do remetente inválido")
}

func TestGetQuoteMetrics_Success(t *testing.T) {
	mockMetrics := new(MockMetricsPort)
	qs := NewQuoteService(nil, mockMetrics)
	metricsCarrier := []CarrierMetrics{{Name: "Correios", AvgPrice: 50.99, MaxPrice: 50.99, MinPrice: 50.99, TotalPrice: 50.99 * 3, TotalOffer: 3}}
	expectedMetrics := []Metrics{
		{Carrier: metricsCarrier, GeneralMaxCarrierName: "Correios", GeneralMinCarrierName: "Correios", GeneralAvgPrice: 50.99, GeneralMaxPrice: 50.99, GeneralMinPrice: 50.99}}

	mockMetrics.On("Execute", 3).Return(expectedMetrics, nil)

	metrics, err := qs.GetQuoteMetrics(3)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(metrics))
	assert.Equal(t, "Correios", metrics[0].Carrier[0].Name)
	assert.Equal(t, 50.99, metrics[0].GeneralAvgPrice)
	mockMetrics.AssertExpectations(t)
}

func TestGetQuoteMetrics_Error(t *testing.T) {
	mockMetrics := new(MockMetricsPort)
	qs := NewQuoteService(nil, mockMetrics)

	mockMetrics.On("Execute", 5).Return([]Metrics{}, errors.New("falha no banco"))

	_, err := qs.GetQuoteMetrics(5)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "falha no banco")

}
