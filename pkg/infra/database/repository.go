package database

import (
	"database/sql"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
)

type QuoteRepository struct {
	db *sql.DB
}

func NewQuoteRepository(db *sql.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (q *QuoteRepository) SaveAllOffers(offers []quote.Offer) error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO offers(final_price, carrier, service, delivery_time) VALUES ($1, $2, $3, $4)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()
	for _, offer := range offers {
		_, err := stmt.Exec(offer.FinalPrice, offer.Carrier, offer.Service, offer.DeliveryTime)
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	return tx.Commit()

}
