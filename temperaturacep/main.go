package main

import (
	"encoding/json"
	"io"
	http "net/http"
	"net/url"
	"regexp"
	"unicode"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type ViaCep struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

type WeatherApi struct {
	Current struct {
		TempC float64 `json:"temp_c"`
		TempF float64 `json:"temp_f"`
		TempK float64 `json:"temp_k"`
	} `json:"current"`
}

type Response struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/", ProcuraCepHandler)
	})

	err := http.ListenAndServe(":8082", r)
	if err != nil {
		return
	}
}

func ProcuraCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")

	validate := regexp.MustCompile(`^[0-9]{8}$`)
	if !validate.MatchString(cepParam) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("CEP inválido"))
		return
	}

	cep, err := EncontraCep(cepParam)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorStr := err.Error()
		w.Write([]byte("erro ao pesquisar CEP: " + errorStr))
		return
	}

	if cep.Erro {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("CEP não encontrado"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	weather, err := EncontraTemperatura(cep.Localidade)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		errorStr := err.Error()
		w.Write([]byte("erro ao procurar temperatura: " + errorStr))
		return
	}
	_ = json.NewEncoder(w).Encode(weather)
}

func EncontraCep(cep string) (*ViaCep, error) {
	req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(req.Body)

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var data ViaCep
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func EncontraTemperatura(city string) (*Response, error) {
	urlWeatherApi := "http://api.weatherapi.com/v1/current.json?key=12969ce544064451ab2103040240905&aqi=no&q=" + transCidade(city)
	req, err := http.Get(urlWeatherApi)

	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(req.Body)

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	var data WeatherApi
	err = json.Unmarshal(res, &data)
	if err != nil {
		return nil, err
	}

	data.Current.TempF = data.Current.TempC*1.8 + 32
	data.Current.TempK = data.Current.TempC + 273

	return &Response{TempC: data.Current.TempC, TempF: data.Current.TempF, TempK: data.Current.TempK}, nil
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func transCidade(s string) string {
	t := transform.Chain(
		norm.NFD,
		transform.RemoveFunc(isMn),
		norm.NFC,
	)
	result, _, _ := transform.String(t, s)
	result = url.QueryEscape(result)
	return result
}
