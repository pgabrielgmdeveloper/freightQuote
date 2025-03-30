package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/database"
	"net/http"

	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
)

type FreteRapidoAdapter struct {
	client *http.Client
	repo   *database.QuoteRepository
}

func NewFreteRapidoAdapter(repo *database.QuoteRepository) *FreteRapidoAdapter {
	return &FreteRapidoAdapter{
		client: &http.Client{},
		repo:   repo,
	}
}

func (fra *FreteRapidoAdapter) Execute(quoteData quote.QuoteRequest) ([]quote.Offer, error) {
	freteApiRequest := DomainToFreteRapidoContractRequest(quoteData)
	requestPayload, err := json.Marshal(freteApiRequest)
	if err != nil {
		return nil, err
	}
	response, err := fra.client.Post("https://sp.freterapido.com/api/v3/quote/simulate", "application/json", bytes.NewBuffer(requestPayload))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("frete Rapido Contract returned %s", response.Status)
	}
	defer response.Body.Close()

	var freteApiResponse FreteRapidoResponse
	if err := json.NewDecoder(response.Body).Decode(&freteApiResponse); err != nil {
		return nil, err
	}
	offers := FreteApiResponseToDomainOffer(freteApiResponse)
	if err = fra.repo.SaveAllOffers(offers); err != nil {
		return nil, err
	}

	return offers, nil
}
