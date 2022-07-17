package main

import (
	"errors"
	"fmt"
	"github.com/leekchan/accounting"
	"math"
)

type Performance struct {
	PlayID   string `json:"playID"`
	Audience int    `json:"audience"`
}

type Invoice struct {
	Customer     string        `json:"customer"`
	Performances []Performance `json:"performances"`
}

type Play struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

var Plays = map[string]Play{
	"hamlet":  {Name: "Hamlet", Type: "tragedy"},
	"as-like": {Name: "As You Like It", Type: "comedy"},
	"othello": {Name: "Othello", Type: "tragedy"},
}

func statement(invoice Invoice, plays map[string]Play) (result string, err error) {
	totalAmount := 0
	volumeCredits := 0
	result = fmt.Sprintf("Statement for %s\n", invoice.Customer)
	format := accounting.Accounting{Symbol: "$", Precision: 2}
	for _, perf := range invoice.Performances {
		play := plays[perf.PlayID]
		if play.Name == "" {
			continue
		}
		thisAmount, err2 := amountFor(perf, play)
		if err2 != nil {
			return "", err2
		}
		volumeCredits += int(math.Max(float64(perf.Audience-30), 0))
		if play.Type == "comedy" {
			volumeCredits += int(math.Floor(float64(perf.Audience / 5)))
		}
		result += fmt.Sprintf("  %s: %s (%d seats)\n", play.Name, format.FormatMoney(thisAmount/100), perf.Audience)
		totalAmount += thisAmount
	}
	result += fmt.Sprintf("Amount owed is %s\n", format.FormatMoney(totalAmount/100))
	result += fmt.Sprintf("You earned %d credits\n", volumeCredits)
	return
}

func amountFor(perf Performance, play Play) (int, error) {
	result := 0
	switch play.Type {
	case "tragedy":
		result = 40000
		if perf.Audience > 30 {
			result += 1000 * (perf.Audience - 30)
		}
	case "comedy":
		result = 30000
		if perf.Audience > 20 {
			result += 10000 + 500*(perf.Audience-20)
		}
		result += 300 * perf.Audience
	default:
		return 0, errors.New("invalid play type")
	}
	return result, nil
}
