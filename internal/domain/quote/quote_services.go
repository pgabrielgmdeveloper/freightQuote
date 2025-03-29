package quote

type QuoteService struct {
	SmltPort    SimulateQuoteOutPutPort
	MetricsPort MetricsOutputPort
}

func NewQuoteService(portSmlt SimulateQuoteOutPutPort, portMetrics MetricsOutputPort) *QuoteService {
	return &QuoteService{
		SmltPort:    portSmlt,
		MetricsPort: portMetrics,
	}
}

func (qs *QuoteService) SimulateQuote(quote QuoteRequest) ([]Offer, error) {

	if err := quote.Validate(); err != nil {
		return nil, err
	}
	return qs.SmltPort.Execute(quote)

}

func (qs *QuoteService) GetQuoteMetrics(lastQuotes int) ([]Metrics, error) {
	return qs.MetricsPort.Execute(lastQuotes)
}
