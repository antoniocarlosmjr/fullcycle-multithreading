package models

type BrasilAPIResponse struct {
	CEP     string `json:"cep"`
	State   string `json:"state"`
	City    string `json:"city"`
	Bairro  string `json:"neighborhood"`
	Address string `json:"street"`
}

type ViaCepAPIResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	UF          string `json:"uf"`
	IBGE        string `json:"ibge"`
	Gia         string `json:"gia"`
	DDD         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}
