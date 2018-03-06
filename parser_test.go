package main

import (
	"encoding/json"
	"testing"
)

/* RESPONSE JSON EXAMPLE
[
	{
		"id": "covesting",
		"name": "Covesting",
		"symbol": "COV",
		"rank": "317",
		"price_usd": "1.39593",
		"price_rub": "70.560382266",
		"price_btc": "0.00012872",
		"24h_volume_usd": "173512.0",
		"market_cap_usd": "24428775.0",
		"available_supply": "17500000.0",
		"total_supply": "20000000.0",
		"max_supply": null,
		"percent_change_1h": "0.33",
		"percent_change_24h": "2.28",
		"percent_change_7d": "-1.05",
		"last_updated": "1520002462"
	}
]
*/
var testBody = []byte(`[
	{
		"id": "covesting",
		"name": "Covesting",
		"symbol": "COV",
		"rank": "317",
		"price_usd": "1.40058",
		"price_rub": "70.560382266", 
		"price_btc": "0.00012932",
		"24h_volume_usd": "173253.0",
		"market_cap_usd": "24510150.0",
		"available_supply": "17500000.0",
		"total_supply": "20000000.0",
		"max_supply": null,
		"percent_change_1h": "1.05",
		"percent_change_24h": "2.62",
		"percent_change_7d": "-0.56",
		"last_updated": "1520002161"
	}
	]`)

func TestParser(t *testing.T) {
	bodyBytes := testBody
	var courses []Course
	err := json.Unmarshal(bodyBytes, &courses)
	if err != nil {
		t.Error("Failed to parse response:", err)
	}
}
