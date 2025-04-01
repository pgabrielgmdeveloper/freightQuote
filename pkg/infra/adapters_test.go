package infra

import (
	"encoding/json"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	http2 "github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"strings"
	"testing"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) GetMetricsQuotes(lastQuotes int) (*quote.Metrics, error) {
	args := m.Called(lastQuotes)
	return args.Get(0).(*quote.Metrics), args.Error(1)
}

func (m *MockRepo) SaveAllOffers(offers []quote.Offer) error {
	args := m.Called(offers)
	return args.Error(0)
}

func ResponseMockFreteRapidoApi(request *http.Request) (*http.Response, error) {
	var payload http2.FreteRapidoApiRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		return httpmock.NewStringResponse(400, "invalid request"), nil
	}
	if payload.Shipper.RegisteredNumber != "25438296000158" {
		responseBody, _ := json.Marshal(struct {
			Details []string `json:"details"`
			Error   string   `json:"error"`
		}{
			Details: []string{},
			Error:   "Shipper not found",
		})
		return httpmock.NewBytesResponse(400, responseBody), nil
	}

	simpleResponse := http2.FreteRapidoApiResponse{
		Dispatchers: []http2.DispatcherResponse{
			{
				Offer: []http2.OfferResponse{
					{
						FinalPrice: 30,
						Carrier: http2.CarrierResponse{
							Name: "CORREIO - SEDEX",
						},
						Service: "SEDEX",
						DeliveryTime: http2.DeliveryTimeResponse{
							Days:    1,
							Hours:   0,
							Minutes: 0,
						},
					},
				},
			},
		},
	}
	bytes, _ := json.Marshal(simpleResponse)
	return httpmock.NewBytesResponse(200, bytes), nil
}

func ValidRequest() quote.QuoteRequest {
	validReq := quote.QuoteRequest{
		Shipper: quote.Shipper{
			RegisteredNumber: "25438296000158",
			Token:            "1d52a9b6b78cf07b08586152459a5c90",
			PlatformCode:     "5AKVkHqCn",
		},
		Recipient: quote.Recipient{
			Type:             0,
			Country:          "BRA",
			Zipcode:          1311000,
			RegisteredNumber: "",
		},
		Dispatchers: []quote.Dispatcher{
			{
				RegisteredNumber: "25438296000158",
				Zipcode:          49157021,
				Volumes: []quote.Volume{
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

func InvalidRequest() quote.QuoteRequest {
	invalidReq := quote.QuoteRequest{
		Shipper: quote.Shipper{
			RegisteredNumber: "123",
		},
	}
	return invalidReq
}

func TestFreteRapidoAdaterSimulateSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST",
		"https://sp.freterapido.com/api/v3/quote/simulate",
		ResponseMockFreteRapidoApi,
	)
	request := ValidRequest()
	mockRepo := new(MockRepo)
	mockRepo.On("SaveAllOffers", mock.Anything).Return(nil)

	adapter := NewFreteRapidoAdapter(mockRepo)

	offers, err := adapter.Execute(request)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(offers))
}

func TestFreteRapidoAdaterSimulateFailureResponseApi(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST",
		"https://sp.freterapido.com/api/v3/quote/simulate",
		ResponseMockFreteRapidoApi,
	)
	request := InvalidRequest()
	mockRepo := new(MockRepo)
	adapter := NewFreteRapidoAdapter(mockRepo)
	_, err := adapter.Execute(request)
	assert.NotNil(t, err)
	assert.True(t, true, strings.Contains(err.Error(), "frete Rapido Contract returned"))

}

func TestFreteRapidoAdaterSimulateFailureSaveDb(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST",
		"https://sp.freterapido.com/api/v3/quote/simulate",
		ResponseMockFreteRapidoApi,
	)
	request := ValidRequest()
	mockRepo := new(MockRepo)
	mockRepo.On("SaveAllOffers", mock.Anything).Return(fmt.Errorf("Error ao salvar no banco"))

	adapter := NewFreteRapidoAdapter(mockRepo)

	_, err := adapter.Execute(request)
	assert.NotNil(t, err)
	assert.Equal(t, "Error ao salvar no banco", err.Error())
}

func TestGetMetricsQuotes(t *testing.T) {
	mockRepo := new(MockRepo)
	mockRepo.On("GetMetricsQuotes", mock.Anything).Return(&quote.Metrics{}, nil)

	adapter := NewMetricsAdapter(mockRepo)
	_, err := adapter.Execute(0)
	assert.Nil(t, err)
}
