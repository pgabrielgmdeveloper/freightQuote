package quote

type SimulateQuoteOutPutPort interface {
	Execute(quoteData QuoteRequest) ([]Offer, error)
}

type SimulateInputPort interface {
	Simulate(request QuoteRequest) ([]Offer, error)
}

type MetricsOutputPort interface {
	Execute(lastQuotes int) ([]Metrics, error)
}

type MetricsInputPort interface {
	GetMetrics(lastQuotes int) ([]Metrics, error)
}
