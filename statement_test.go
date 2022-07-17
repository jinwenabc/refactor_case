package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func initData() (map[string]Play, Invoice) {
	playsStr := []byte(`
		{
		  "hamlet": { "name": "Hamlet", "type": "tragedy" },
		  "as-like": { "name": "As You Like It", "type": "comedy" },
		  "othello": { "name": "Othello", "type": "tragedy" }
		}`)
	plays := make(map[string]Play)
	_ = json.Unmarshal(playsStr, &plays)
	invoiceStr := []byte(`
		{
		"customer": "BigCo",
		"performances": [
		  {
			"playID": "hamlet",
			"audience": 55
		  },
		  {
			"playID": "as-like",
			"audience": 35
		  },
		  {
			"playID": "othello",
			"audience": 40
		  }
		]
	  }`)
	var invoice Invoice
	_ = json.Unmarshal(invoiceStr, &invoice)
	return plays, invoice
}

func Test_statement(t *testing.T) {
	type args struct {
		invoice Invoice
		plays   map[string]Play
	}
	plays, invoice := initData()
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{
			"normal case",
			args{
				invoice: invoice,
				plays:   plays,
			},
			"Statement for BigCo\n  Hamlet: $650.00 (55 seats)\n  As You Like It: $580.00 (35 seats)\n  Othello: $500.00 (40 seats)\nAmount owed is $1,730.00\nYou earned 47 credits\n",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := statement(tt.args.invoice, tt.args.plays)
			if (err != nil) != tt.wantErr {
				t.Errorf("statement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("statement() gotResult = %v, want %v", gotResult, tt.wantResult)
				compareStrs(gotResult, tt.wantResult)
			}
		})
	}
}

func compareStrs(str1, str2 string) {
	len1, len2 := len(str1), len(str2)
	if len1 != len2 {
		fmt.Printf("len1:%d, len2:%d\n", len1, len2)
		return
	}
	for i := 0; i < len1 && i < len2; i++ {
		if str1[i] != str2[i] {
			fmt.Printf("str1[%d]=%v, str2[%d]=%v\n", i, rune(str1[i]), i, rune(str2[i]))
		}
	}
}
