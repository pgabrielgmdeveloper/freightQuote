package quote

import (
	"testing"
)

func TestShipper_Validate(t *testing.T) {
	tests := []struct {
		name          string
		shipper       Shipper
		expectedError bool
		errMsg        string
	}{
		{
			name: "Shipper válido",
			shipper: Shipper{
				RegisteredNumber: "12345678901234", // CNPJ válido
				Token:            "abcdefghijklmnopqrstuvwxyz123456",
				PlatformCode:     "PLAT123",
			},
			expectedError: false,
		},
		{
			name: "CNPJ inválido",
			shipper: Shipper{
				RegisteredNumber: "123", // CNPJ inválido
				Token:            "abcdefghijklmnopqrstuvwxyz123456",
				PlatformCode:     "PLAT123",
			},
			expectedError: true,
			errMsg:        "CNPJ do remetente inválido",
		},
		{
			name: "Token inválido",
			shipper: Shipper{
				RegisteredNumber: "12345678901234",
				Token:            "curto", // Token inválido
				PlatformCode:     "PLAT123",
			},
			expectedError: true,
			errMsg:        "token deve ter 32 caracteres",
		},
		{
			name: "PlatformCode vazio",
			shipper: Shipper{
				RegisteredNumber: "12345678901234",
				Token:            "abcdefghijklmnopqrstuvwxyz123456",
				PlatformCode:     "", // PlatformCode inválido
			},
			expectedError: true,
			errMsg:        "código da plataforma é obrigatório",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.shipper.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("Shipper.Validate() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if tt.expectedError && err.Error() != tt.errMsg {
				t.Errorf("Shipper.Validate() error = %v, expectedErrorMsg %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestRecipient_Validate(t *testing.T) {
	tests := []struct {
		name          string
		recipient     Recipient
		expectedError bool
		errMsg        string
	}{
		{
			name: "Recipient PJ válido",
			recipient: Recipient{
				Type:             1,
				Country:          "BRA",
				Zipcode:          12345678,
				RegisteredNumber: "12345678901234", // CNPJ válido
			},
			expectedError: false,
		},
		{
			name: "Recipient PF válido",
			recipient: Recipient{
				Type:             0,
				Country:          "BRA",
				Zipcode:          12345678,
				RegisteredNumber: "12345678901", // CPF válido
			},
			expectedError: false,
		},
		{
			name: "Tipo inválido",
			recipient: Recipient{
				Type:    5,
				Country: "BRA",
				Zipcode: 12345678,
			},
			expectedError: true,
			errMsg:        "tipo de destinatário deve ser 'PF(0)' ou 'PJ(1)'",
		},
		{
			name: "País inválido",
			recipient: Recipient{
				Type:    1,
				Country: "USA",
				Zipcode: 12345678,
			},
			expectedError: true,
			errMsg:        "país deve ser 'BRA'",
		},
		{
			name: "CEP inválido",
			recipient: Recipient{
				Type:    1,
				Country: "BRA",
				Zipcode: 123, // CEP inválido
			},
			expectedError: true,
			errMsg:        "CEP inválido",
		},
		{
			name: "CNPJ inválido para PJ",
			recipient: Recipient{
				Type:             1,
				Country:          "BRA",
				Zipcode:          12345678,
				RegisteredNumber: "123", // CNPJ inválido
			},
			expectedError: true,
			errMsg:        "CNPJ do destinatário inválido",
		},
		{
			name: "CPF inválido para PF",
			recipient: Recipient{
				Type:             0,
				Country:          "BRA",
				Zipcode:          12345678,
				RegisteredNumber: "123", // CPF inválido
			},
			expectedError: true,
			errMsg:        "CPF do destinatário inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.recipient.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("Recipient.Validate() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if tt.expectedError && err.Error() != tt.errMsg {
				t.Errorf("Recipient.Validate() error = %v, expectedErrorMsg %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestVolume_Validate(t *testing.T) {
	tests := []struct {
		name          string
		volume        Volume
		expectedError bool
		errMsg        string
	}{
		{
			name: "Volume válido",
			volume: Volume{
				Category:      "19",
				Amount:        2,
				UnitaryWeight: 1.5,
				UnitaryPrice:  10.99,
				Height:        10,
				Width:         20,
				Length:        30,
			},
			expectedError: false,
		},
		{
			name: "Categoria vazia",
			volume: Volume{
				Category:      "",
				Amount:        2,
				UnitaryWeight: 1.5,
				UnitaryPrice:  10.99,
				Height:        10,
				Width:         20,
				Length:        30,
			},
			expectedError: true,
			errMsg:        "categoria do volume é obrigatória",
		},
		{
			name: "Quantidade inválida",
			volume: Volume{
				Category:      "19",
				Amount:        0,
				UnitaryWeight: 1.5,
				UnitaryPrice:  10.99,
				Height:        10,
				Width:         20,
				Length:        30,
			},
			expectedError: true,
			errMsg:        "quantidade deve ser maior que zero",
		},
		{
			name: "Peso inválido",
			volume: Volume{
				Category:      "125",
				Amount:        2,
				UnitaryWeight: 0,
				UnitaryPrice:  10.99,
				Height:        10,
				Width:         20,
				Length:        30,
			},
			expectedError: true,
			errMsg:        "peso unitário deve ser maior que zero",
		},
		{
			name: "Preço negativo",
			volume: Volume{
				Category:      "125",
				Amount:        2,
				UnitaryWeight: 1.5,
				UnitaryPrice:  -1,
				Height:        10,
				Width:         20,
				Length:        30,
			},
			expectedError: true,
			errMsg:        "preço unitário não pode ser negativo",
		},
		{
			name: "Dimensões inválidas",
			volume: Volume{
				Category:      "125",
				Amount:        2,
				UnitaryWeight: 1.5,
				UnitaryPrice:  10.99,
				Height:        0,
				Width:         -1,
				Length:        0,
			},
			expectedError: true,
			errMsg:        "dimensões devem ser maiores que zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.volume.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("Volume.Validate() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if tt.expectedError && err.Error() != tt.errMsg {
				t.Errorf("Volume.Validate() error = %v, expectedErrorMsg %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDispatcher_Validate(t *testing.T) {
	validVolume := Volume{
		Category:      "19",
		Amount:        1,
		UnitaryWeight: 1,
		UnitaryPrice:  10,
		Height:        10,
		Width:         10,
		Length:        10,
	}

	invalidVolume := Volume{
		Category:      "125",
		Amount:        0,
		UnitaryWeight: 0,
		UnitaryPrice:  -1,
		Height:        0,
		Width:         0,
		Length:        0,
	}

	tests := []struct {
		name          string
		dispatcher    Dispatcher
		expectedError bool
		errMsg        string
	}{
		{
			name: "Dispatcher válido",
			dispatcher: Dispatcher{
				RegisteredNumber: "12345678901234",
				Zipcode:          12345678,
				Volumes:          []Volume{validVolume},
			},
			expectedError: false,
		},
		{
			name: "CNPJ inválido",
			dispatcher: Dispatcher{
				RegisteredNumber: "123",
				Zipcode:          12345678,
				Volumes:          []Volume{validVolume},
			},
			expectedError: true,
			errMsg:        "CNPJ do expedidor inválido",
		},
		{
			name: "CEP inválido",
			dispatcher: Dispatcher{
				RegisteredNumber: "12345678901234",
				Zipcode:          123,
				Volumes:          []Volume{validVolume},
			},
			expectedError: true,
			errMsg:        "CEP do expedidor inválido",
		},
		{
			name: "Sem volumes",
			dispatcher: Dispatcher{
				RegisteredNumber: "12345678901234",
				Zipcode:          12345678,
				Volumes:          []Volume{},
			},
			expectedError: true,
			errMsg:        "pelo menos um volume é obrigatório",
		},
		{
			name: "Volume inválido",
			dispatcher: Dispatcher{
				RegisteredNumber: "12345678901234",
				Zipcode:          12345678,
				Volumes:          []Volume{invalidVolume},
			},
			expectedError: true,
			// A mensagem de erro será a primeira validação que falhar no volume
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.dispatcher.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("Dispatcher.Validate() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if tt.expectedError && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Dispatcher.Validate() error = %v, expectedErrorMsg %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestQuoteRequest_Validate(t *testing.T) {
	validRequest := QuoteRequest{
		Shipper: Shipper{
			RegisteredNumber: "12345678901234",
			Token:            "abcdefghijklmnopqrstuvwxyz123456",
			PlatformCode:     "PLAT123",
		},
		Recipient: Recipient{
			Type:             1,
			Country:          "BRA",
			Zipcode:          49160000,
			RegisteredNumber: "12345678901234",
		},
		Dispatchers: []Dispatcher{
			{
				RegisteredNumber: "12345678901234",
				Zipcode:          12345678,
				Volumes: []Volume{
					{
						Category:      "125",
						Amount:        1,
						UnitaryWeight: 1,
						UnitaryPrice:  10,
						Height:        10,
						Width:         10,
						Length:        10,
					},
				},
			},
		},
	}

	tests := []struct {
		name          string
		request       QuoteRequest
		expectedError bool
		errMsg        string
	}{
		{
			name:          "Caso válido",
			request:       validRequest,
			expectedError: false,
		},
		{
			name: "Sem dispatchers",
			request: func() QuoteRequest {
				r := validRequest
				r.Dispatchers = []Dispatcher{}
				return r
			}(),
			expectedError: true,
			errMsg:        "pelo menos um expedidor é obrigatório",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.expectedError {
				t.Errorf("QuoteRequest.Validate() error = %v, expectedError %v", err, tt.expectedError)
			}
			if tt.expectedError && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Dispatcher.Validate() error = %v, expectedErrorMsg %v", err.Error(), tt.errMsg)
			}
		})
	}
}
