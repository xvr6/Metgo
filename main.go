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
	Time          string      `json:"time"`
	WeatherCode   int         `json:"weather_code"`
	Temp          json.Number `json:"temperature_2m"`
	FeelsTemp     json.Number `json:"apparent_temperature"`
	PrecipAmmount json.Number `json:"precipitation"`
	Humidity      json.Number `json:"relative_humidity_2m"`
	// Precip Prob must be inherited from time 0 in hourly data
}

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
	Moonset       []string      `json:"moonset"`
}

type Response struct {
	Latitude  json.Number `json:"latitude"`
	Longitude json.Number `json:"longitude"`
	Units     Units       `json:"current_units"`
	Current   WeatherJson `json:"current"`
	Hourly    HourlyJson  `json:"hourly"`
	Daily     DailyJson   `json:"daily"`
}

// - Reformat of data for passing to frontend ---
// Just the above data formatted into a more logical structure for the frontend to parse through.

// An instance of weather data. Used for current weather and hourly weather breakdowns. Helps keep data predictible if the same structure is used EVERYWHERE.
type WeatherInstance struct {
	Time          string
	WeatherCode   int
	Temp          int
	FeelsTemp     int
	PrecipAmmount int
	Humidity      int
}

type ParsedHour struct {
}

// Other information that only makes sense to view on a 'day scale', as well as the hourly reports and averages.
type ParsedDay struct {
	Hourly   [24]ParsedHour
	TempMax  string
	TempMin  string
	Sunrise  string
	Sunset   string
	Moonrise string
	Moonset  string
}

// instead of storing timestamps as a 'string', could i do it with some sort of existing timestamp type or otherwise make my own?

// The actual datastructure to be sent to the frontend, includes other infrmation about session and specific location
type ParsedData struct {
	Lat            float32
	Long           float32
	Timezone       string //maybe move to units; effectively acts like one
	CurrentWeather WeatherInstance
	Units          Units
	Day            [7]ParsedDay
}

func reformat(res Response) ParsedData {
	var formatted ParsedData

	return formatted
}

func urlFormat() string {
	// TODO: Currently a placeholder, implement actual string builder logic with basic tui.
	// Regardless, likely have these be the default that are always enabled, and any settings would be atop of these.

	const timezone string = "America/New_York"
	const lat float32 = 42.3584
	const long float32 = -71.0598
	const tempUnit string = "celsius" // celsius or fahrenheit
	const precipUnit string = "mm"    // mm or inch
	const speedUnit string = "kmh"    // kmh mph ms(meters/sec) kn (knots)
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

	response, err := io.ReadAll(responseRaw.Body)
	if err != nil {
		log.Fatal(err)
	}

	var res Response
	if err := json.Unmarshal(response, &res); err != nil {
		log.Fatal(err)
	}
	// debug printout
	debug, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(debug))

	/* TODO: mass reformatting of data to fit something more useable for an actual app.
	   type MetgoData struct{
		 	current weather struct (pull out precip percentage from current hour)
		  daily weathers
		  units (likely not needed but is good for sanity)
	} */
	data := reformat(res)
	fmt.Printf("Struct:\n%#v", data)

}
