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
		PerfCalculator: NewPerfCalculator(aPerformance, playFor(aPerformance)),
	}
	if result.PerfCalculator == nil {
		return result, errors.New("construct perf calculator fail")
	}
	result.amount = result.PerfCalculator.getAmount()
	return result, nil
}

func renderPlainText(data *statementData) string {
	result := fmt.Sprintf("Statement for %s\n", data.Customer)
	for _, perf := range data.Performances {
		result += fmt.Sprintf("  %s: %s (%d seats)\n", perf.getPlay().Name, usd(perf.amount), perf.getPerf().Audience)
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
		result += fmt.Sprintf("<tr><td>%s</td><td>%d</td>", perf.getPlay().Name, perf.getPerf().Audience)
		result += fmt.Sprintf("<td>%s</td></tr>\n", usd(perf.amount))
	}
	result += "</table>\n"
	result += "<p>Amount owed is <em>${usd(data.totalAmount)}</em></p>\n"
	result += "<p>You earned <em>${data.totalVolumeCredits}</em> credits</p>\n"
	return result
}

type PerfCalculator interface {
	getPlay() Play
	getPerf() Performance
	getAmount() int
	getCredits() int
}

type PerformanceBasicInfo struct {
	perf Performance
	play Play
}

func (p PerformanceBasicInfo) getPlay() Play {
	return p.play
}

func (p PerformanceBasicInfo) getPerf() Performance {
	return p.perf
}

func NewPerfCalculator(aPerf Performance, aPlay Play) PerfCalculator {
	switch aPlay.Type {
	case "comedy":
		return newComedyCalculator(aPerf, aPlay)
	case "tragedy":
		return newTragedyCalculator(aPerf, aPlay)
	default:
		return nil
	}
}

type ComedyCalculator struct {
	PerformanceBasicInfo
}

func newComedyCalculator(aPerf Performance, aPlay Play) *ComedyCalculator {
	return &ComedyCalculator{
		PerformanceBasicInfo{
			perf: aPerf,
			play: aPlay,
		},
	}
}

func (c *ComedyCalculator) getAmount() int {
	result := 30000
	if c.perf.Audience > 20 {
		result += 10000 + 500*(c.perf.Audience-20)
	}
	result += 300 * c.perf.Audience
	return result
}

func (c *ComedyCalculator) getCredits() int {
	return int(math.Max(float64(c.perf.Audience-30), 0)) + int(math.Floor(float64(c.perf.Audience/5)))
}

type TragedyCalculator struct {
	PerformanceBasicInfo
}

func newTragedyCalculator(aPerf Performance, aPlay Play) *TragedyCalculator {
	return &TragedyCalculator{
		PerformanceBasicInfo{
			perf: aPerf,
			play: aPlay,
		},
	}
}

func (t *TragedyCalculator) getAmount() int {
	result := 40000
	if t.perf.Audience > 30 {
		result += 1000 * (t.perf.Audience - 30)
	}
	return result
}

func (t *TragedyCalculator) getCredits() int {
	return int(math.Max(float64(t.perf.Audience-30), 0))
}

func totalVolumeCredits(data *statementData) int {
	result := 0
	for _, perf := range data.Performances {
		result += perf.getCredits()
	}
	return result
}

func totalAmount(data *statementData) int {
	result := 0
	for _, perf := range data.Performances {
		result += perf.amount
	}
	return result
}

func playFor(aPerformance Performance) Play {
	return Plays[aPerformance.PlayID]
}

func usd(amount int) string {
	return format.FormatMoney(amount / 100)
}
