package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ViaCepApiResponse struct {
	Cep        string `json:"cep"`
	Estado     string `json:"estado"`
	Localidade string `json:"localidade"`
}

type BrasilCepApiResponse struct {
	Cep   string `json:"cep"`
	State string `json:"state"`
	City  string `json:"city"`
}

type CepApiData struct {
	Cep         string
	State       string
	City        string
	ServiceName string
}

func main() {
	cep := os.Args[1]
	cepChan := make(chan CepApiData)

	go viaCepWorker(&cep, cepChan)
	go brasilCepWorker(&cep, cepChan)

	var fastestWorker CepApiData
	select {
	case fastestWorker = <-cepChan:
	case <-time.After(1 * time.Second):
		fmt.Printf("Request timeout\n")
		return
	}

	if fastestWorker != (CepApiData{}) {
		fmt.Printf("Fastest API response is: %s, city: %s, state: %s, CEP: %s\n", fastestWorker.ServiceName, fastestWorker.City, fastestWorker.State, fastestWorker.Cep)
	}
}

func viaCepWorker(cep *string, viaCepChan chan CepApiData) {
	const cepDataUrl = "https://viacep.com.br/ws/%s/json/"
	url := fmt.Sprintf(cepDataUrl, *cep)

	var viaCepResponse ViaCepApiResponse
	err := fetchAndDecode(url, &viaCepResponse)

	if err != nil {
		fmt.Printf("error fetching data from ViaCep: %s\n", err)
		return
	}

	cepApiData := CepApiData{
		Cep:         viaCepResponse.Cep,
		State:       viaCepResponse.Estado,
		City:        viaCepResponse.Localidade,
		ServiceName: "ViaCep",
	}

	viaCepChan <- cepApiData
}

func brasilCepWorker(cep *string, brasilCepChan chan CepApiData) {
	const cepDataUrl = "https://brasilapi.com.br/api/cep/v1/"
	url := cepDataUrl + *cep

	var brasilCepApiResponse BrasilCepApiResponse
	err := fetchAndDecode(url, &brasilCepApiResponse)

	if err != nil {
		fmt.Printf("error fetching data from BrasilCep: %s\n", err)
		return
	}

	cepApiData := CepApiData{
		Cep:         brasilCepApiResponse.Cep,
		State:       brasilCepApiResponse.State,
		City:        brasilCepApiResponse.City,
		ServiceName: "BrasilCep",
	}

	brasilCepChan <- cepApiData
}

func fetchAndDecode(url string, target interface{}) error {
	resp, err := fetchCepData(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch target.(type) {
	case *ViaCepApiResponse, *BrasilCepApiResponse:
		return json.NewDecoder(resp.Body).Decode(target)
	default:
		return fmt.Errorf("unsupported target type")
	}
}

func fetchCepData(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
