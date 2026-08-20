package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Units struct {
	Time          string `json:"time"`
	Temp          string `json:"temperature_2m"`
	FeelsTemp     string `json:"apparent_temperature"`
	PrecipAmmount string `json:"precipitation"`
	Humidity      string `json:"relative_humidity_2m"`
}
type WeatherJson struct {
	Time        string      `json:"time"`
	WeatherCode int         `json:"weather_code"`
	Temp        json.Number `json:"temperature_2m"`
	FeelsTemp   json.Number `json:"apparent_temperature"`
	PrecipProb  json.Number ``
	// Precip Prob must be inherited from time 0 in hourly data
	PrecipAmmount json.Number `json:"precipitation"`
	Humidity      json.Number `json:"relative_humidity_2m"`
}

// TODO: The hourly data is listed as "time": [xxxx-xx-xxTxx:xx], "temperature_2m": [<list of corelated temps>], etc
// These should be mapped together on creation of HourlyJson for the sake of less formatting in the future. Once response is properly turned into structs, then do formatting there.

type HourlyJson struct {
	Time          []string      `json:"time"`
	WeatherCode   []int         `json:"weather_code"`
	Temp          []json.Number `json:"temperature_2m"`
	FeelsTemp     []json.Number `json:"apparent_temperature"`
	PrecipProb    []json.Number `json:"precipitation_probability"`
	PrecipAmmount []json.Number `json:"precipitation"`
	Humidity      []json.Number `json:"relative_humidity_2m"`
}

type DailyJson struct {
	Time          []string      `json:"time"`
	WeatherCode   []int         `json:"weather_code"`
	TempMax       []json.Number `json:"temperature_2m_max"`
	TempMin       []json.Number `json:"temperature_2m_min"`
	PrecipProbMax []int         `json:"precipitation_probability_max"`
	Sunrise       []string      `json:"sunrise"`
	Sunset        []string      `json:"sunset"`
	Moonrise      []string      `json:"moonrise"`
	Moonset       []string      `json"moonset"`
}

type Response struct {
	Latitude  json.Number `json:"latitude"`
	Longitude json.Number `json:"longitude"`
	Units     Units       `json:"current_units"`
	Current   WeatherJson `json:"current"`
	Hourly    HourlyJson  `json:"hourly"`
	Daily     DailyJson   `json:"daily"`
}

func urlFormat() string {
	// TODO: Currently a plceholder, implement actual string builder logic with basic tui

	const timezone string = "America/New_York"
	const lat float32 = 52.52
	const long float32 = 13.41
	const tempUnit = "celsius" // celsius or fahrenheit&
	const precipUnit = "mm"    // mm or inch
	const speedUnit = "kmh"    // kmh mph ms(meters/sec) kn (knots)
	const currentSelect string = "weather_code,temperature_2m,apparent_temperature,precipitation,relative_humidity_2m"
	const hourlySelect string = "temperature_2m,precipitation_probability,precipitation,apparent_temperature,relative_humidity_2m,weather_code"
	const dailySelect string = "weather_code,temperature_2m_max,temperature_2m_min,apparent_temperature_max,apparent_temperature_min,precipitation_sum,precipitation_hours,precipitation_probability_max,sunrise,sunset,moonset,moonrise"

	output := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?timezone=%v&latitude=%v&longitude=%v&temperature_unit=%v&precipitation_unit=%v&wind_speed_unit=%v&current=%v&hourly=%v&daily=%v", timezone, lat, long, tempUnit, precipUnit, speedUnit, currentSelect, hourlySelect, dailySelect)
	return output
}

func main() {
	var url = urlFormat()
	fmt.Printf("Url: %v\n", url)

	// - HTTP Request ---
	responseRaw, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	responseParsed, err := io.ReadAll(responseRaw.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	if err := json.Unmarshal(responseParsed, &responseObject); err != nil {
		log.Fatal(err)
	}
	// debug printout
	debug, _ := json.MarshalIndent(responseObject, "", "  ")
	fmt.Println(string(debug))

	/* TODO: mass reformatting of data to fit something more useable for an actual app.
	   type MetgoData struct{
		 	current weather struct (pull out precip percentage from current hour)
		  daily weathers
		  units (likely not needed but is good for sanity)
	} */
}
