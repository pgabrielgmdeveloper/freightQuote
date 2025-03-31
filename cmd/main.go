package main

import (
	"github.com/gin-gonic/gin"
	"github.com/pgabrielgmdeveloper/freightQuote/configs"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/cache"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/database"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/drivers"
)

func main() {

	cfg, err := configs.LoadConfig()

	db, err := database.NewDbInstance(
		cfg.DBDriver,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		"disable",
		cfg.DBPassword,
		cfg.DBUser)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	redis := cache.NewRedisInstance(cfg.RedisHost, cfg.RedisPort)
	redisCache := cache.NewRedisCache(redis)
	repo := database.NewQuoteRepository(db)
	adapterMetrics := infra.NewMetricsAdapter(repo)
	adapterSimulateQuote := infra.NewFreteRapidoAdapter(repo)

	quoteService := quote.NewQuoteService(adapterSimulateQuote, adapterMetrics)

	handlerQuoteServices := drivers.NewQuoteAdapterHandler(quoteService, quoteService, redisCache)

	r := gin.Default()
	r.POST("/simulate", handlerQuoteServices.SimulateQuote)
	r.GET("/metrics", handlerQuoteServices.GetMetrics)
	r.Run(":8000")

}
