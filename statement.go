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

var (
	Plays  map[string]Play
	format accounting.Accounting
)

func init() {
	Plays = map[string]Play{
		"hamlet":  {Name: "Hamlet", Type: "tragedy"},
		"as-like": {Name: "As You Like It", Type: "comedy"},
		"othello": {Name: "Othello", Type: "tragedy"},
	}
	format = accounting.Accounting{Symbol: "$", Precision: 2}
}

func statement(invoice Invoice) (result string, err error) {
	result = fmt.Sprintf("Statement for %s\n", invoice.Customer)
	for _, perf := range invoice.Performances {
		thisAmount, err2 := amountFor(perf)
		if err2 != nil {
			return "", err2
		}
		result += fmt.Sprintf("  %s: %s (%d seats)\n", playFor(perf).Name, usd(thisAmount), perf.Audience)
	}

	amount, err3 := totalAmount(invoice)
	if err3 != nil {
		return "", err3
	}

	result += fmt.Sprintf("Amount owed is %s\n", usd(amount))
	result += fmt.Sprintf("You earned %d credits\n", totalVolumeCredits(invoice))
	return
}

func totalVolumeCredits(invoice Invoice) int {
	result := 0
	for _, perf := range invoice.Performances {
		result += volumeCreditsFor(perf)
	}
	return result
}

func volumeCreditsFor(perf Performance) int {
	volumeCredits := int(math.Max(float64(perf.Audience-30), 0))
	if playFor(perf).Type == "comedy" {
		volumeCredits += int(math.Floor(float64(perf.Audience / 5)))
	}
	return volumeCredits
}

func totalAmount(invoice Invoice) (int, error) {
	result := 0
	for _, perf := range invoice.Performances {
		amount, err := amountFor(perf)
		if err != nil {
			return 0, err
		}
		result += amount
	}
	return result, nil
}

func amountFor(perf Performance) (int, error) {
	result := 0
	switch playFor(perf).Type {
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

func playFor(aPerformance Performance) Play {
	return Plays[aPerformance.PlayID]
}

func usd(amount int) string {
	return format.FormatMoney(amount / 100)
}
