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
	PrecipProb    string `json:"precipitation_probability"`
	PrecipAmmount string `json:"precipitation"`
	Humidity      string `json:"relative_humidity_2m"`
}
type WeatherState struct {
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
// These should be mapped together on creation of HourlyStruct for the sake of less formatting in the future. Once response is properly turned into structs, then do formatting there.

type HourlyStruct struct {
	Time          []string      `json:"time"`
	WeatherCore   []int         `json:"weather_code"`
	Temp          []json.Number `json:"temperature_2m"`
	FeelsTemp     []json.Number `json:"apparent_temperature"`
	PrecipProb    []json.Number `json:"precipitation_probability"`
	PrecipAmmount []json.Number `json:"precipitation"`
	Humidity      []json.Number `json:"relative_humidity_2m"`
}

// TODO: Same as above with daily data which is currently not fetched.

type Response struct {
	Latitude  json.Number  `json:"latitude"`
	Longitude json.Number  `json:"longitude"`
	Units     Units        `json:"current_units"`
	Current   WeatherState `json:"current"`
	Hourly    HourlyStruct `json:"hourly"`
}

func urlFormat() string {
	// TODO: Currently a plceholder, implement actual string builder logic with basic tui

	const timezone string = "America/New_York"
	const lat float32 = 52.52
	const long float32 = 13.41
	const currentSelect string = "weather_code,temperature_2m,apparent_temperature,precipitation,relative_humidity_2m"
	const hourlySelect string = "temperature_2m,precipitation_probability,precipitation,apparent_temperature,relative_humidity_2m,weather_code"
	const dailySelect string = ""

	output := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?timezone=%v&latitude=%v&longitude=%v&current=%v&hourly=%v", timezone, lat, long, currentSelect, hourlySelect)
	return output
}

func main() {
	var url = urlFormat()
	fmt.Printf("Url: %v\n", url)

	// - HTTP Request ---
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HTTP Response:\n%v\n", response.Body)

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject Response
	if err := json.Unmarshal(responseData, &responseObject); err != nil {
		log.Fatal(err)
	}
	// debug printout
	debug, _ := json.MarshalIndent(responseObject, "", "  ")
	fmt.Println(string(debug))

}
