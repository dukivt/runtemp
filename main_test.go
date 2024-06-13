package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRecebeCepETemperaturaCep(t *testing.T) {
	cep := "24230126"

	urlRecebeCep := "http://localhost:8080"
	payload := fmt.Sprintf(`{"cep": "%s"}`, cep)
	respRecebeCep, err := http.Post(urlRecebeCep, "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatalf("Erro ao fazer requisição para recebecep: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(respRecebeCep.Body)

	if respRecebeCep.StatusCode != http.StatusOK {
		t.Fatalf("recebecep retornou status %d, esperava %d", respRecebeCep.StatusCode, http.StatusOK)
	}

	bodyRecebeCep, err := ioutil.ReadAll(respRecebeCep.Body)
	if err != nil {
		t.Fatalf("Erro ao ler corpo da resposta de recebecep: %v", err)
	}

	var tempResponse struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		TempK float64 `json:"temp_k"`
	}
	err = json.Unmarshal(bodyRecebeCep, &tempResponse)
	if err != nil {
		t.Fatalf("Erro ao fazer parse da resposta JSON de recebecep: %v", err)
	}

	urlTemperaturaCep := fmt.Sprintf("http://localhost:8082/?cep=%s", cep)
	respTemperaturaCep, err := http.Get(urlTemperaturaCep)
	if err != nil {
		t.Fatalf("Erro ao fazer requisição para temperaturacep: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(respTemperaturaCep.Body)

	if respTemperaturaCep.StatusCode != http.StatusOK {
		t.Fatalf("temperaturacep retornou status %d, esperava %d", respTemperaturaCep.StatusCode, http.StatusOK)
	}

	bodyTemperaturaCep, err := ioutil.ReadAll(respTemperaturaCep.Body)
	if err != nil {
		t.Fatalf("Erro ao ler corpo da resposta de temperaturacep: %v", err)
	}

	var temperatura struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		TempK float64 `json:"temp_k"`
	}
	err = json.Unmarshal(bodyTemperaturaCep, &temperatura)
	if err != nil {
		t.Fatalf("Erro ao fazer parse da resposta JSON de temperaturacep: %v", err)
	}

	if temperatura.TempC != tempResponse.TempC || temperatura.TempF != tempResponse.TempF || temperatura.TempK != tempResponse.TempK {
		t.Errorf("Valores de temperatura inconsistentes. Recebido de recebecep: %+v, recebido de temperaturacep: %+v", tempResponse, temperatura)
	}
}
