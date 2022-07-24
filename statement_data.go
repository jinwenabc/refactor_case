package main

type enrichedPerformance struct {
	Performance
	play   Play
	amount int
}

type statementData struct {
	Customer     string
	Performances []enrichedPerformance
	TotalAmount  int
	TotalCredits int
}

func createStatementData(invoice Invoice) (*statementData, error) {
	data := &statementData{
		Customer: invoice.Customer,
	}
	for _, perf := range invoice.Performances {
		anEnrichPerf, enrichErr := enrichPerformance(perf)
		if enrichErr != nil {
			return nil, enrichErr
		}
		data.Performances = append(data.Performances, anEnrichPerf)
	}
	data.TotalAmount = totalAmount(data)
	data.TotalCredits = totalVolumeCredits(data)
	return data, nil
}
