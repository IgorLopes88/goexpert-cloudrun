package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type WeatherApi struct {
	Location struct {
		Name           string  `json:"name"`
		Region         string  `json:"region"`
		Country        string  `json:"country"`
		Lat            float64 `json:"lat"`
		Lon            float64 `json:"lon"`
		TzID           string  `json:"tz_id"`
		LocaltimeEpoch int     `json:"localtime_epoch"`
		Localtime      string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdatedEpoch int     `json:"last_updated_epoch"`
		LastUpdated      string  `json:"last_updated"`
		TempC            float64 `json:"temp_c"`
		TempF            float64 `json:"temp_f"`
		IsDay            int     `json:"is_day"`
		Condition        struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph    float64 `json:"wind_mph"`
		WindKph    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMb float64 `json:"pressure_mb"`
		PressureIn float64 `json:"pressure_in"`
		PrecipMm   float64 `json:"precip_mm"`
		PrecipIn   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelslikeC float64 `json:"feelslike_c"`
		FeelslikeF float64 `json:"feelslike_f"`
		VisKm      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		Uv         float64 `json:"uv"`
		GustMph    float64 `json:"gust_mph"`
		GustKph    float64 `json:"gust_kph"`
	} `json:"current"`
}

type ViaCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"uf"`
	City         string `json:"localidade"`
	Neighborhood string `json:"bairro"`
	Street       string `json:"logradouro"`
	Complemento  string `json:"complemento"`
	Ibge         string `json:"ibge"`
	Gia          string `json:"gia"`
	DDD          string `json:"ddd"`
	Siafi        string `json:"siafi"`
}

var ApiKey = "8a1e75434bfd4056852172426241307"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /temperature/{zipcode}", handlerTemperature)
	http.ListenAndServe(":8080", mux)
}

func handlerTemperature(w http.ResponseWriter, r *http.Request) {
	request := r.PathValue("zipcode")
	zipcode, err := convertZipcode(request)
	if err != nil || zipcode == "" {
		w.WriteHeader(422)
		w.Write([]byte("invalid zipcode"))
		log.Printf("Invalid Zipcode Request: %s", request)
		return
	}

	city, err := SearchLocation(zipcode)
	if err != nil || city == "" {
		w.WriteHeader(404)
		w.Write([]byte("can not find zipcode"))
		log.Printf("City Not Located: %s", zipcode)
		return
	}

	response, err := GetTemperature(city)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("can not find zipcode"))
		log.Print(err)
		return
	}

	celsius := strconv.FormatFloat(response, 'f', 0, 64)
	fahrenheit := strconv.FormatFloat(response*1.8+32, 'f', 0, 64)
	kelvin := strconv.FormatFloat(response+273.15, 'f', 0, 64)

	responseJson := []byte(`{ "temp_C": ` + celsius + `, "temp_F": ` + fahrenheit + `, "temp_K": ` + kelvin + ` }`)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJson)
	log.Printf("Zipcode Request: %s (%s): %vC / %vF / %vK ", zipcode, city, celsius, fahrenheit, kelvin)
}

func SearchLocation(cep string) (string, error) {
	result, err := requestGetUrl("http://viacep.com.br/ws/" + cep + "/json")
	if err != nil {
		return "", err
	}

	var data ViaCEP
	err = json.Unmarshal(result, &data)
	if err != nil {
		return "", err
	}

	// CASO NÃO ENCONTRE O CEP
	if data.Cep != "" {
		return data.City, nil
	} else {
		return "", nil
	}
}

func GetTemperature(city string) (float64, error) {
	city = convertName(city)

	result, err := requestGetUrl("http://api.weatherapi.com/v1/current.json?key=" + ApiKey + "&q=" + city + "&aqi=no")
	if err != nil {
		return 0, err
	}

	var u WeatherApi
	err = json.Unmarshal(result, &u)
	if err != nil {
		return 0, err
	}
	return u.Current.TempC, err
}

func requestGetUrl(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*3))
	defer cancel()

	req = req.WithContext(ctx)
	c := &http.Client{}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func convertZipcode(zipcode string) (string, error) {
	zipcode = strings.Replace(zipcode, "-", "", 1)
	match, err := regexp.MatchString("[0-9]", zipcode)
	if err != nil {
		return "", err
	}
	if !match {
		return "", err
	}
	if zipcode == "" || len(zipcode) != 8 {
		return "", err
	}
	return zipcode, nil
}

func convertName(name string) string {
	format := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(format, name)
	if err != nil {
		return ""
	}
	result = strings.Replace(result, " ", "%20", -1)
	return result
}
