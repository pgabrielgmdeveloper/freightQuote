package infra

import (
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/database"
)

type MetricsAdapter struct {
	repo database.IQuoteRepository
}

func NewMetricsAdapter(repo database.IQuoteRepository) *MetricsAdapter {
	return &MetricsAdapter{
		repo: repo,
	}
}

func (m MetricsAdapter) Execute(lastQuotes int) (*quote.Metrics, error) {
	return m.repo.GetMetricsQuotes(lastQuotes)
}
