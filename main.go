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

var (
	botToken       = os.Getenv("botToken")
	channelID      = os.Getenv("channelID")
	telegramBotAPI = "https://api.telegram.org/bot" + botToken + "/"
)

const (
	submitHour                = 10
	retryInterval             = 5 * time.Minute
	notificationCheckInterval = 15 * time.Minute
)

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

func toFloat(s string) float64 {
	res, _ := strconv.ParseFloat(s, 32)
	return res
}

func (c Course) String() string {
	isNegativeGrow := c.PercentChange24h[0] == '-'
	var courseGrowEmoji string
	if isNegativeGrow {
		courseGrowEmoji = "↘️️"
	} else {
		courseGrowEmoji = "↗️"
	}
	utcSeconds, _ := strconv.ParseInt(c.LastUpdated, 10, 64)
	c.LastUpdated = time.Unix(utcSeconds, 0).String()
	c.PriceRUB = fmt.Sprintf("%.2f", toFloat(c.PriceRUB))
	c.PriceUSD = fmt.Sprintf("%.2f", toFloat(c.PriceUSD))
	return fmt.Sprintf("%v\nPrice: %v$\n            %v₽\nChange 24h: %v%%%v\nUpdated: %v",
		c.Name, c.PriceUSD, c.PriceRUB, c.PercentChange24h, courseGrowEmoji, c.LastUpdated)
}

func getCourse(currency string) Course {
	resp, _ := http.Get("https://api.coinmarketcap.com/v1/ticker/" + currency + "/?convert=RUB")
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var courses []Course
	err = json.Unmarshal(bodyBytes, &courses)
	if err != nil {
		log.Fatal(err)
	}
	return courses[0]
}

func sendTelegramMsg(channel string, msg string) error {
	channel = netUrl.QueryEscape(channel)
	msg = netUrl.QueryEscape(msg)
	url := fmt.Sprintf("%vsendMessage?chat_id=%v&text=%v", telegramBotAPI, channel, msg)
	_, err := http.Get(url)
	return err
}

func main() {
	var lastSubmitDay int
	for {
		currentTime := time.Now()
		if lastSubmitDay != currentTime.Day() && currentTime.Hour() >= submitHour {
			log.Println("Sending course...")

			covestingCourse := getCourse("covesting")
			err := sendTelegramMsg(channelID, covestingCourse.String())
			if err != nil {
				log.Println("Error:", err)
				time.Sleep(retryInterval)
				continue
			}
			log.Println("Message has been sent.")

			lastSubmitDay = currentTime.Day()
		}
		time.Sleep(notificationCheckInterval)
	}
}
