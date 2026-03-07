package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type FunctionCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
}

type Part struct {
	FunctionCall FunctionCall `json:"functionCall"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
}

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

var cities = []string{
	"New York", "London", "Paris", "Tokyo", "Berlin",
	"San Francisco", "Singapore", "Sydney", "Toronto", "Dubai",
}

var stocks = []string{
	"AAPL", "GOOGL", "AMZN", "TSLA", "MSFT", "META", "NVDA",
}

var currencies = []string{
	"USD", "EUR", "GBP", "JPY", "INR", "CAD",
}

var musicGenres = []string{
	"jazz", "lofi", "classical", "ambient", "rock",
}

var rooms = []string{
	"living room", "bedroom", "kitchen", "office",
}

var languages = []string{
	"French", "Spanish", "German", "Japanese", "Hindi",
}

var restaurants = []string{
	"Italian", "Japanese", "Mexican", "Indian", "Thai",
}

func buildResponse(function string, args map[string]interface{}) string {

	resp := Response{
		Candidates: []Candidate{
			{
				Content: Content{
					Role: "model",
					Parts: []Part{
						{
							FunctionCall: FunctionCall{
								Name: function,
								Args: args,
							},
						},
					},
				},
				FinishReason: "STOP",
			},
		},
	}

	b, _ := json.MarshalIndent(resp, "", "  ")
	return string(b)
}

func randomCity() string {
	return cities[rand.Intn(len(cities))]
}

func randomStock() string {
	return stocks[rand.Intn(len(stocks))]
}

func randomCurrency() string {
	return currencies[rand.Intn(len(currencies))]
}

func randomGenre() string {
	return musicGenres[rand.Intn(len(musicGenres))]
}

func randomRoom() string {
	return rooms[rand.Intn(len(rooms))]
}

func randomLanguage() string {
	return languages[rand.Intn(len(languages))]
}

func randomRestaurant() string {
	return restaurants[rand.Intn(len(restaurants))]
}

func GenerateToolCallDataset(count int) map[string]string {

	_ = rand.New(rand.NewSource(time.Now().UnixNano()))
	data := make(map[string]string)

	for len(data) < count {

		switch rand.Intn(12) {

		case 0:

			city := randomCity()
			prompt := fmt.Sprintf("Get the current weather in %s.", city)

			data[prompt] = buildResponse(
				"get_weather",
				map[string]interface{}{
					"city": city,
				},
			)

		case 1:

			a := rand.Intn(100)
			b := rand.Intn(100)

			prompt := fmt.Sprintf("Calculate %d * %d.", a, b)

			data[prompt] = buildResponse(
				"calculator_multiply",
				map[string]interface{}{
					"a": a,
					"b": b,
				},
			)

		case 2:

			currencyA := randomCurrency()
			currencyB := randomCurrency()
			amount := rand.Intn(1000)

			prompt := fmt.Sprintf("Convert %d %s to %s.", amount, currencyA, currencyB)

			data[prompt] = buildResponse(
				"currency_convert",
				map[string]interface{}{
					"amount": amount,
					"from":   currencyA,
					"to":     currencyB,
				},
			)

		case 3:

			city := randomCity()
			prompt := fmt.Sprintf("Find flights to %s next week.", city)

			data[prompt] = buildResponse(
				"search_flights",
				map[string]interface{}{
					"destination": city,
					"date_range":  "next_week",
				},
			)

		case 4:

			ticker := randomStock()

			prompt := fmt.Sprintf("What is the stock price of %s?", ticker)

			data[prompt] = buildResponse(
				"get_stock_price",
				map[string]interface{}{
					"ticker": ticker,
				},
			)

		case 5:

			lang := randomLanguage()

			prompt := fmt.Sprintf("Translate this webpage to %s.", lang)

			data[prompt] = buildResponse(
				"translate_webpage",
				map[string]interface{}{
					"target_language": lang,
				},
			)

		case 6:

			cuisine := randomRestaurant()

			prompt := fmt.Sprintf("Book a table at a %s restaurant for two people.", cuisine)

			data[prompt] = buildResponse(
				"book_restaurant",
				map[string]interface{}{
					"cuisine": cuisine,
					"people":  2,
				},
			)

		case 7:

			genre := randomGenre()

			prompt := fmt.Sprintf("Play some %s music.", genre)

			data[prompt] = buildResponse(
				"play_music",
				map[string]interface{}{
					"genre": genre,
				},
			)

		case 8:

			room := randomRoom()

			prompt := fmt.Sprintf("Turn off the lights in the %s.", room)

			data[prompt] = buildResponse(
				"set_light_state",
				map[string]interface{}{
					"room":  room,
					"state": "off",
				},
			)

		case 9:

			room := randomRoom()

			prompt := fmt.Sprintf("Turn on the lights in the %s.", room)

			data[prompt] = buildResponse(
				"set_light_state",
				map[string]interface{}{
					"room":  room,
					"state": "on",
				},
			)

		case 10:

			time := rand.Intn(12) + 1

			prompt := fmt.Sprintf("Set an alarm for %d AM tomorrow.", time)

			data[prompt] = buildResponse(
				"set_alarm",
				map[string]interface{}{
					"time": fmt.Sprintf("%02d:00", time),
					"day":  "tomorrow",
				},
			)

		case 11:

			cityA := randomCity()
			cityB := randomCity()

			prompt := fmt.Sprintf("Get directions from %s to %s.", cityA, cityB)

			data[prompt] = buildResponse(
				"maps_directions",
				map[string]interface{}{
					"from": cityA,
					"to":   cityB,
				},
			)
		}
	}

	return data
}
