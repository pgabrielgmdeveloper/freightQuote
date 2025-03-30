package infra

import (
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/database"
)

type MetricsAdapter struct {
	repo *database.QuoteRepository
}

func NewMetricsAdapter(repo *database.QuoteRepository) *MetricsAdapter {
	return &MetricsAdapter{
		repo: repo,
	}
}

func (m MetricsAdapter) Execute(lastQuotes int) (*quote.Metrics, error) {
	return m.repo.GetMetricsQuotes(lastQuotes)
}
