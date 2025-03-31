package drivers

import (
	"github.com/gin-gonic/gin"
	"github.com/pgabrielgmdeveloper/freightQuote/internal/domain/quote"
	"net/http"
	"strconv"
)

type QuoteAdapterHandler struct {
	inputSimulate quote.SimulateInputPort
	inputMetrics  quote.MetricsInputPort
}

func NewQuoteAdapterHandler(inputSimulate quote.SimulateInputPort, inputMetrics quote.MetricsInputPort) *QuoteAdapterHandler {
	return &QuoteAdapterHandler{
		inputSimulate: inputSimulate,
		inputMetrics:  inputMetrics,
	}
}

func (q *QuoteAdapterHandler) SimulateQuote(c *gin.Context) {
	var simulateRequest SimulateQuoteRequest

	if err := c.ShouldBindJSON(&simulateRequest); err != nil {
		JSONErrorResponse(http.StatusBadRequest, "Error ao converter json em struct", err, c)
		return
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

	offersResponse := DomainToSimulateQuoteResponse(offers)

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
	c.JSON(http.StatusOK, metrics)
}

func JSONErrorResponse(statusCode int, ErrorMessage string, error error, c *gin.Context) {
	c.JSON(statusCode, struct {
		ErrorMessage string `json:"errorMessage"`
		ErrorDetails string `json:"errorDetails"`
	}{ErrorMessage: ErrorMessage, ErrorDetails: error.Error()})
}
