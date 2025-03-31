package database

import "github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"

type IQuoteRepository interface {
	SaveAllOffers(offers []quote.Offer) error
	GetMetricsQuotes(lastQuotes int) (*quote.Metrics, error)
}
