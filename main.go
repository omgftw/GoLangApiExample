package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"text/template"
	"time"
)

type StockApiResponse struct {
	MetaData   MetaData              `json:"Meta Data"`
	TimeSeries map[string]TimeSeries `json:"Time Series (Daily)"`
}

type ApiResponse struct {
	Symbol string `json:"symbol"`
	Data []float64 `json:"data"`
	Average float64 `json:"average"`
}

type MetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type TimeSeries struct {
	Open             float64 `json:"1. open,string"`
	High             float64 `json:"2. high,string"`
	Low              float64 `json:"3. low,string"`
	Close            float64 `json:"4. close,string"`
	AdjustedClose    float64 `json:"5. adjusted close,string"`
	Volume           int     `json:"6. volume,string"`
	DividendAmount   float64 `json:"7. dividend amount,string"`
	SplitCoefficient float64 `json:"8. split coefficient,string"`
	Date time.Time
}

func getStocks(c *gin.Context) {
	var days []TimeSeries

	// Create slice from map of TimeSeries to allow deterministic sorting
	for _, timeSeries := range apiData.TimeSeries {
		days = append(days, timeSeries)
	}

	// Sort by date
	sort.Slice(days, func(a, b int) bool {
		return days[a].Date.After(days[b].Date)
	})

	// Create response
	days = days[:envVars.NDays]
	resp := ApiResponse{
		Symbol: envVars.Symbol,
	}

	var average float64
	for _, timeSeries := range days {
		resp.Data = append(resp.Data, timeSeries.Close)
		average += timeSeries.Close
	}
	average /= float64(len(days))
	// TODO Round to specific number of decimal places?
	resp.Average = average

	c.IndentedJSON(http.StatusOK, resp)
}

// Generic error handler to keep things DRY
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func GetData() StockApiResponse {
	data := StockApiResponse{}
	var err error

	// If no APIKEY is passed in, load test data from a local file
	if envVars.ApiKey == "" {
		jsonData, err := ioutil.ReadFile("data.json")
		handleError(err)
		err = json.Unmarshal(jsonData, &data)
		handleError(err)
	} else {
		// if APIKEY is passed, fetch the data from the API
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(GetBaseUrl())
		handleError(err)
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(&data)
		handleError(err)
	}

	// Parse dates for easy filtering
	layout := "2006-01-02"
	for key, timeSeries := range data.TimeSeries {
		timeSeries.Date, err = time.Parse(layout, key)
		data.TimeSeries[key] = timeSeries
		handleError(err)
	}

	return data
}

func GetEnvVars() {
	envVars.ApiKey = os.Getenv("APIKEY")

	envVars.Symbol = os.Getenv("SYMBOL")
	if envVars.Symbol == "" {
		envVars.Symbol = "MSFT"
	}

	envVars.BaseUrl = os.Getenv("BASEURL")
	if envVars.BaseUrl == "" {
		envVars.BaseUrl = "https://www.alphavantage.co/query?apikey={{ .ApiKey }}&function=TIME_SERIES_DAILY_ADJUSTED&symbol={{ .Symbol }}"
	}

	ndaysString := os.Getenv("NDAYS")
	if ndaysString == "" {
		envVars.NDays = 3
	} else {
		var err error
		envVars.NDays, err = strconv.Atoi(ndaysString)
		handleError(err)
	}
}

// GetBaseUrl Templates the base URL for the API
func GetBaseUrl() string {
	t, err :=template.New("baseUrl").Parse(envVars.BaseUrl)
	handleError(err)
	var output bytes.Buffer
	err = t.Execute(&output, envVars)
	handleError(err)
	return output.String()
}

// Store data here for caching purposes
var apiData StockApiResponse

type EnvVars struct {
	ApiKey string
	Symbol string
	NDays int
	BaseUrl string
}

var envVars = EnvVars{}

func main() {
	GetEnvVars()
	apiData = GetData()

	router := gin.Default()
	router.GET("/", getStocks)

	err := router.Run("0.0.0.0:8080")
	handleError(err)
}
