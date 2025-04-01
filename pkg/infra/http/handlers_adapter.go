package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"github.com/pgabrielgmdeveloper/freightQuote/pkg/infra/cache"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type QuoteAdapterHandler struct {
	inputSimulate quote.SimulateInputPort
	inputMetrics  quote.MetricsInputPort
	redisCache    cache.IRedisCache
}

func NewQuoteAdapterHandler(inputSimulate quote.SimulateInputPort, inputMetrics quote.MetricsInputPort, redis cache.IRedisCache) *QuoteAdapterHandler {
	return &QuoteAdapterHandler{
		inputSimulate: inputSimulate,
		inputMetrics:  inputMetrics,
		redisCache:    redis,
	}
}

func (q *QuoteAdapterHandler) SimulateQuote(c *gin.Context) {
	var simulateRequest SimulateQuoteRequest
	var offersResponse SimulateQuoteResponse
	if err := c.ShouldBindJSON(&simulateRequest); err != nil {
		JSONErrorResponse(http.StatusBadRequest, "Error ao converter json em struct", err, c)
		return
	}

	zipcode, err := ConverterStrinToInZipcode(simulateRequest.Recipient.Address.Zipcode)
	if err != nil {
		JSONErrorResponse(http.StatusBadRequest, "Error ao converter json em struct", err, c)
		return
	}
	var skuAmounts []string
	for _, v := range simulateRequest.Volumes {
		skuAmounts = append(skuAmounts, fmt.Sprintf("%s-%d", v.Sku, v.Amount))
	}
	sort.Strings(skuAmounts)
	cachedKey := fmt.Sprintf("%d-%s", zipcode, strings.Join(skuAmounts, "-"))
	ctx := context.Background()
	resultCached, err := q.redisCache.Get(ctx, cachedKey)
	if err == redis.Nil {
		log.Println("Cache não encontrado para a key:", cachedKey)
	}

	if resultCached != "" {
		err = json.Unmarshal([]byte(resultCached), &offersResponse)
		if err != nil {
			log.Println("Não foi possivel converter o cache me json. Error: ", err.Error())
		} else {
			log.Println("Resultado retornado em cache")
			c.JSON(http.StatusOK, offersResponse)
			return
		}
	}

	quoteRequest, err := RequestToDomainQuote(simulateRequest)
	if err != nil {
		JSONErrorResponse(http.StatusBadRequest, "Error processar dados", err, c)
		return
	}

	offers, err := q.inputSimulate.Simulate(*quoteRequest)
	if err != nil {
		JSONErrorResponse(http.StatusInternalServerError, "Error ao Simular cotações", err, c)
		return
	}

	offersResponse = DomainToSimulateQuoteResponse(offers)
	if err = q.redisCache.Set(ctx, cachedKey, offersResponse, time.Minute*30); err != nil {
		log.Println("Não foi possivel salvar retorno em cache err: ", err.Error())
	}

	c.JSON(http.StatusOK, offersResponse)
	return
}

func (q *QuoteAdapterHandler) GetMetrics(c *gin.Context) {
	lastQuotes, _ := strconv.Atoi(c.Query("last_quotes"))

	metrics, err := q.inputMetrics.GetMetrics(lastQuotes)
	if err != nil {
		JSONErrorResponse(http.StatusInternalServerError, "Error ao gerar metricas", err, c)
		return
	}
	c.JSON(http.StatusOK, DomainMetricsToRequest(*metrics))
}

func JSONErrorResponse(statusCode int, ErrorMessage string, error error, c *gin.Context) {
	c.JSON(statusCode, struct {
		ErrorMessage string `json:"errorMessage"`
		ErrorDetails string `json:"errorDetails"`
	}{ErrorMessage: ErrorMessage, ErrorDetails: error.Error()})
}
