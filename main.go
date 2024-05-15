package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fullcycle-multithreading/models"
)

const (
	baseURLBrasilAPI = "https://brasilapi.com.br/api/cep/v1/"
	baseURLViaCEP    = "https://viacep.com.br/ws/"
	timeoutSeconds   = 5
)

func main() {
	flag.Parse()
	cep := flag.Arg(0)
	if cep == "" {
		fmt.Println("Por favor, informe um CEP.")
		os.Exit(1)
	}

	cep = cleanCEP(cep)

	if err := validateCEP(cep); err != nil {
		fmt.Println("CEP inválido:", err)
		os.Exit(1)
	}

	channelBrasilAPI := make(chan models.BrasilAPIResponse)
	channelViaCEPAPI := make(chan models.ViaCepAPIResponse)

	go getBrasilAPI(cep, channelBrasilAPI)
	go getViaCEPAPI(cep, channelViaCEPAPI)

	select {
	case brasilAPI := <-channelBrasilAPI:
		log.Printf("BrasilAPI Response: %#v", brasilAPI)
	case viaCEP := <-channelViaCEPAPI:
		log.Printf("ViaCEP Response: %#v", viaCEP)
	case <-time.After(timeoutSeconds * time.Second):
		fmt.Println("Limite excedido!")
	}
}

func cleanCEP(cep string) string {
	return strings.ReplaceAll(cep, "-", "")
}

func validateCEP(cep string) error {
	if len(cep) != 8 {
		return errors.New("o CEP deve conter 8 dígitos")
	}
	return nil
}

func getBrasilAPI(cep string, ch chan<- models.BrasilAPIResponse) {
	req, err := http.Get(baseURLBrasilAPI + cep)
	if err != nil {
		fmt.Println("Erro ao buscar dados da BrasilAPI:", err)
		ch <- models.BrasilAPIResponse{}
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Erro ao fechar o corpo da resposta da BrasilAPI:", err)
			ch <- models.BrasilAPIResponse{}
		}
	}(req.Body)

	res, _ := io.ReadAll(req.Body)

	var data models.BrasilAPIResponse
	if err := json.Unmarshal(res, &data); err != nil {
		fmt.Println("Erro ao decodificar resposta da BrasilAPI:", err)
		ch <- models.BrasilAPIResponse{}
		return
	}

	ch <- data
}

func getViaCEPAPI(cep string, ch chan<- models.ViaCepAPIResponse) {
	req, err := http.Get(baseURLViaCEP + cep + "/json")
	if err != nil {
		fmt.Println("Erro ao buscar dados da ViaCEP:", err)
		ch <- models.ViaCepAPIResponse{}
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Erro ao fechar o corpo da resposta da ViaCEP:", err)
			ch <- models.ViaCepAPIResponse{}
		}
	}(req.Body)

	var viaCEPData models.ViaCepAPIResponse
	if err := json.NewDecoder(req.Body).Decode(&viaCEPData); err != nil {
		fmt.Println("Erro ao decodificar resposta da ViaCEP:", err)
		ch <- models.ViaCepAPIResponse{}
		return
	}

	ch <- viaCEPData
}
