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
	var data *statementData
	if data, err = createStatementData(invoice); err != nil {
		return "", err
	}
	return renderPlainText(data), nil
}

func enrichPerformance(aPerformance Performance) (enrichedPerformance, error) {
	result := enrichedPerformance{
		Performance: aPerformance,
		play:        playFor(aPerformance),
	}
	amount, err := amountFor(result)
	if err == nil {
		result.amount = amount
	}
	return result, err
}

func renderPlainText(data *statementData) string {
	result := fmt.Sprintf("Statement for %s\n", data.Customer)
	for _, perf := range data.Performances {
		result += fmt.Sprintf("  %s: %s (%d seats)\n", perf.play.Name, usd(perf.amount), perf.Audience)
	}
	result += fmt.Sprintf("Amount owed is %s\n", usd(data.TotalAmount))
	result += fmt.Sprintf("You earned %d credits\n", data.TotalCredits)
	return result
}

func htmlStatement(invoice Invoice) (result string, err error) {
	var data *statementData
	if data, err = createStatementData(invoice); err != nil {
		return "", err
	}
	return renderHTML(data), nil
}

func renderHTML(data *statementData) string {
	result := "<h1>Statement for ${data.customer}</h1>\n"
	result += "<tr><th>play</th><th>seats</th><th>cost</th></tr>"
	for _, perf := range data.Performances {
		result += fmt.Sprintf("<tr><td>%s</td><td>%d</td>", perf.play.Name, perf.Audience)
		result += fmt.Sprintf("<td>%s</td></tr>\n", usd(perf.amount))
	}
	result += "</table>\n"
	result += "<p>Amount owed is <em>${usd(data.totalAmount)}</em></p>\n"
	result += "<p>You earned <em>${data.totalVolumeCredits}</em> credits</p>\n"
	return result
}

func totalVolumeCredits(data *statementData) int {
	result := 0
	for _, perf := range data.Performances {
		result += volumeCreditsFor(perf)
	}
	return result
}

func volumeCreditsFor(perf enrichedPerformance) int {
	volumeCredits := int(math.Max(float64(perf.Audience-30), 0))
	if perf.play.Type == "comedy" {
		volumeCredits += int(math.Floor(float64(perf.Audience / 5)))
	}
	return volumeCredits
}

func totalAmount(data *statementData) int {
	result := 0
	for _, perf := range data.Performances {
		result += perf.amount
	}
	return result
}

func amountFor(perf enrichedPerformance) (int, error) {
	result := 0
	switch perf.play.Type {
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
