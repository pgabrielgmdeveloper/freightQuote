package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/database"
	http2 "github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/http"
	"io"
	"net/http"

	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
)

type FreteRapidoAdapter struct {
	client *http.Client
	repo   database.IQuoteRepository
}

func NewFreteRapidoAdapter(repo database.IQuoteRepository) *FreteRapidoAdapter {
	return &FreteRapidoAdapter{
		client: &http.Client{},
		repo:   repo,
	}
}

func (fra *FreteRapidoAdapter) Execute(quoteData quote.QuoteRequest) ([]quote.Offer, error) {
	freteApiRequest := http2.DomainToFreteRapidoContractRequest(quoteData)
	requestPayload, err := json.Marshal(freteApiRequest)
	if err != nil {
		return nil, err
	}

	response, err := fra.client.Post("https://sp.freterapido.com/api/v3/quote/simulate", "application/json", bytes.NewBuffer(requestPayload))
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("frete Rapido Contract returned %s", string(body))
	}
	defer response.Body.Close()

	var freteApiResponse http2.FreteRapidoApiResponse
	if err := json.NewDecoder(response.Body).Decode(&freteApiResponse); err != nil {
		return nil, err
	}
	offers := http2.FreteApiResponseToDomainOffer(freteApiResponse)
	if err = fra.repo.SaveAllOffers(offers); err != nil {
		return nil, err
	}

	return offers, nil
}
