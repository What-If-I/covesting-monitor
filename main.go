package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	netUrl "net/url"
	"os"
	"strconv"
	"time"
)

const retryInterval = 5 * time.Minute

type Course struct {
	ID               string `json:"id"`
	Name             string `json:"names"`
	Symbol           string `json:"symbol"`
	Rank             string `json:"rank"`
	PriceUSD         string `json:"price_usd"`
	PriceRUB         string `json:"price_rub"`
	PriceBTC         string `json:"price_btc"`
	PercentChange1h  string `json:"percent_change_1h"`
	PercentChange24h string `json:"percent_change_24h"`
	PercentChange7d  string `json:"percent_change_7d"`
	LastUpdated      string `json:"last_updated"`
}

type TelegramClient struct {
	token string
}

func (t TelegramClient) sendMsg(channel string, msg string) error {
	channel = netUrl.QueryEscape(channel)
	msg = netUrl.QueryEscape(msg)
	url := fmt.Sprintf("%vsendMessage?chat_id=%v&text=%v", t.token, channel, msg)
	_, err := http.Get(url)
	return err
}

func (c Course) String() string {
	isNegativeGrow := c.PercentChange24h[0] == '-'
	var courseGrowEmoji string
	if isNegativeGrow {
		courseGrowEmoji = "😭️️️"
	} else {
		courseGrowEmoji = "🤩️"
	}
	utcSeconds, _ := strconv.ParseInt(c.LastUpdated, 10, 64)
	c.LastUpdated = time.Unix(utcSeconds, 0).String()
	c.PriceRUB = fmt.Sprintf("%.2f", toFloat(c.PriceRUB))
	c.PriceUSD = fmt.Sprintf("%.2f", toFloat(c.PriceUSD))
	return fmt.Sprintf("%v\nPrice: %v$\n            %v₽\nChange 24h: %v%% %v\nChange 7d:   %v%%\nUpdated: %v",
		c.Name, c.PriceUSD, c.PriceRUB, c.PercentChange24h, courseGrowEmoji, c.PercentChange7d, c.LastUpdated)
}

func getCourse(currency string) (Course, error) {
	resp, _ := http.Get("https://api.coinmarketcap.com/v1/ticker/" + currency + "/?convert=RUB")
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Course{}, err
	}
	var courses []Course
	err = json.Unmarshal(bodyBytes, &courses)
	if err != nil {
		return Course{}, err
	}

	return courses[0], nil
}

func findSecondsUntil(future time.Time) time.Duration {
	return time.Duration(future.Sub(time.Now()).Seconds()) * time.Second
}

func main() {
	submitHour := convertToInt(os.Getenv("submitTime"))
	channelID := os.Getenv("channelID")
	currencyName := os.Getenv("currency")
	botToken := os.Getenv("botToken")

	telegramClient := TelegramClient{token: botToken}

	for {
		log.Printf("Getting %s course...", currencyName)
		course, err := getCourse(currencyName)
		if err != nil {
			log.Printf("Failed to fetch course: %v\n", err)
			log.Printf("Will try again after: %d minutes\n", retryInterval/time.Minute)
			time.Sleep(retryInterval)
			continue
		}

		log.Println("Course is:\n", course)
		err = telegramClient.sendMsg(channelID, course.String())
		if err != nil {
			log.Printf("Failed to send message course: %v\n", err)
			log.Printf("Will try again after: %d minutes\n", retryInterval/time.Minute)
			time.Sleep(retryInterval)
			continue
		}

		log.Println("Message has been sent.")

		now := time.Now()
		nextTick := time.Date(
			now.Year(), now.Month(), now.Day()+1, submitHour,
			0, 0, 0, now.Location())
		secondsTillNextTick := findSecondsUntil(nextTick)

		log.Printf("Sleeping for %v.", secondsTillNextTick)
		time.Sleep(secondsTillNextTick)
	}
}

func convertToInt(string string) int {
	hour, err := strconv.ParseInt(string, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	return int(hour)
}

func toFloat(s string) float64 {
	res, _ := strconv.ParseFloat(s, 32)
	return res
}
