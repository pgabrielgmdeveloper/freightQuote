package database

import (
	"database/sql"
	"fmt"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"strings"
)

type QuoteRepository struct {
	db *sql.DB
}

func NewQuoteRepository(db *sql.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

func (q *QuoteRepository) GetMetricsQuotes(lastQuotes int) (*quote.Metrics, error) {
	var queryBuilder strings.Builder
	var from string
	if lastQuotes > 0 {
		queryBuilder.WriteString(`
		with limited_offers as (
			select * from offers order by created_at desc limit $1
		)
		`)
		from = "limited_offers"
	} else {
		from = "offers"
	}

	queryBuilder.WriteString(fmt.Sprintf(`
		select
			carrier,
			count(*) as total_offer,
			round(sum(final_price),2) as total_price,
			round(avg(final_price),2) as avg_price,
			round(min(final_price),2) as min_price,
			round(max(final_price),2) as max_price,
			(select min(final_price) from %s) as min_general_price,
			(select max(final_price) from %s) as max_general_price,
			(select avg(final_price) from %s) as avg_general_price,
    		(select carrier from %s where final_price = (select min(final_price) from %s) limit 1) AS carrier_min_general_price,
    		(select carrier from %s where final_price = (select max(final_price) from %s) limit 1) AS carrier_max_general_price
		from %s group by carrier`, from, from, from, from, from, from, from, from))

	query := queryBuilder.String()
	stmt, err := q.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var rows *sql.Rows
	if lastQuotes > 0 {
		rows, err = stmt.Query(lastQuotes)
		if err != nil {
			return nil, err
		}
	} else {
		rows, err = stmt.Query()
		if err != nil {
			return nil, err
		}
	}
	var metrics quote.Metrics

	for rows.Next() {
		var carrierMetrics quote.CarrierMetrics
		err = rows.Scan(&carrierMetrics.Name,
			&carrierMetrics.TotalOffer,
			&carrierMetrics.TotalPrice,
			&carrierMetrics.AvgPrice,
			&carrierMetrics.MinPrice,
			&carrierMetrics.MaxPrice,
			&metrics.GeneralMinPrice,
			&metrics.GeneralMaxPrice,
			&metrics.GeneralAvgPrice,
			&metrics.GeneralMinCarrierName,
			&metrics.GeneralMaxCarrierName,
		)
		if err != nil {
			return nil, err
		}
		metrics.Carrier = append(metrics.Carrier, carrierMetrics)
	}

	return &metrics, nil
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
