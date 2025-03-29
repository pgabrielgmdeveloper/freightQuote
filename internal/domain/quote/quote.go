package quote

import (
	"errors"
	"fmt"
	"regexp"
)

type Shipper struct {
	RegisteredNumber string
	Token            string
	PlatformCode     string
}

func (s *Shipper) Validate() error {
	if !isValidCNPJ(s.RegisteredNumber) {
		return errors.New("CNPJ do remetente inválido")
	}
	if len(s.Token) != 32 {
		return errors.New("token deve ter 32 caracteres")
	}
	if s.PlatformCode == "" {
		return errors.New("código da plataforma é obrigatório")
	}
	return nil
}

type Recipient struct {
	Type             int
	Country          string
	Zipcode          int
	RegisteredNumber string
}

func (r *Recipient) Validate() error {
	if r.Type != 0 && r.Type != 1 {
		return errors.New("tipo de destinatário deve ser 'PF(0)' ou 'PJ(1)'")
	}
	if r.Country != "BRA" {
		return errors.New("país deve ser 'BRA'")
	}
	if !isValidCEP(r.Zipcode) {
		return errors.New("CEP inválido")
	}
	if (r.Type == 1 && r.RegisteredNumber != "") && !isValidCNPJ(r.RegisteredNumber) {
		return errors.New("CNPJ do destinatário inválido")
	}
	if (r.Type == 0 && r.RegisteredNumber != "") && !isValidCPF(r.RegisteredNumber) {
		return errors.New("CPF do destinatário inválido")
	}
	return nil
}

type Volume struct {
	Category      string
	Amount        int
	UnitaryWeight float64
	UnitaryPrice  float64
	Height        float64
	Width         float64
	Length        float64
}

func (v *Volume) Validate() error {
	if v.Category == "" {
		return errors.New("categoria do volume é obrigatória")
	}
	if v.Amount <= 0 {
		return errors.New("quantidade deve ser maior que zero")
	}
	if v.UnitaryWeight <= 0 {
		return errors.New("peso unitário deve ser maior que zero")
	}
	if v.UnitaryPrice < 0 {
		return errors.New("preço unitário não pode ser negativo")
	}
	if v.Height <= 0 || v.Width <= 0 || v.Length <= 0 {
		return errors.New("dimensões devem ser maiores que zero")
	}
	return nil
}

type Dispatcher struct {
	RegisteredNumber string
	Zipcode          int
	Volumes          []Volume
}

func (d *Dispatcher) Validate() error {
	if !isValidCNPJ(d.RegisteredNumber) {
		return errors.New("CNPJ do expedidor inválido")
	}
	if !isValidCEP(d.Zipcode) {
		return errors.New("CEP do expedidor inválido")
	}
	if len(d.Volumes) == 0 {
		return errors.New("pelo menos um volume é obrigatório")
	}
	for _, volume := range d.Volumes {
		if err := volume.Validate(); err != nil {
			return err
		}
	}
	return nil
}

type QuoteRequest struct {
	Shipper     Shipper
	Recipient   Recipient
	Dispatchers []Dispatcher
}

func (q *QuoteRequest) Validate() error {

	if err := q.Shipper.Validate(); err != nil {
		return err
	}
	if err := q.Recipient.Validate(); err != nil {
		return err
	}
	if len(q.Dispatchers) == 0 {
		return errors.New("pelo menos um expedidor é obrigatório")
	}
	for _, dispatcher := range q.Dispatchers {
		if err := dispatcher.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func isValidCNPJ(cnpj string) bool {
	matched, _ := regexp.MatchString(`^\d{14}$`, cnpj)
	return matched
}

func isValidCPF(cpf string) bool {
	matched, _ := regexp.MatchString(`^\d{11}$`, cpf)
	return matched
}

func isValidCEP(cep int) bool {
	cepStr := fmt.Sprintf("%d", cep)

	matched, _ := regexp.MatchString(`^\d{7,8}$`, cepStr)
	return matched
}

type Offer struct {
	FinalPrice   float64
	Carrier      string
	Service      string
	DeliveryTime int
}

type CarrierMetrics struct {
	Name       string
	AvgPrice   float64
	MaxPrice   float64
	MinPrice   float64
	TotalPrice float64
}

type Metrics struct {
	Carrier         CarrierMetrics
	GeneralAvgPrice float64
	GeneralMinPrice float64
	GeneralMaxPrice float64
}
