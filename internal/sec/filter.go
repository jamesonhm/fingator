package edgar

import (
	"fmt"

	"github.com/jamesonhm/fingator/internal/sec/models"
)

func FilterDCF(cf *models.CompanyFactsResponse, d *models.DCFData) []*models.FilteredFact {

	var XBRLTags = map[string][]string{
		"CashFlow":         {"NetCashProvidedByUsedInOperatingActivities"},
		"CapEx":            {"PaymentsToAcquirePropertyPlantAndEquipment"},
		"Revenue":          {"Revenues"},
		"NetIncome":        {"NetIncomeLoss"},
		"OperatingExpense": {"OperatingExpenses"},
		"MadeUpCategory":   {"MadeUpTag"},
	}

	var filteredFacts []*models.FilteredFact
	for key, tags := range XBRLTags {
		factData, err := findFact(cf.Facts.Data, key, tags)
		if err != nil {
			fmt.Printf("%w\n", err)
			continue
		}
		filteredFacts = append(filteredFacts, factData)
	}
	return filteredFacts
}

func findFact(d map[string]models.FactData, key string, tags []string) (*models.FilteredFact, error) {
	for i := 0; i < len(tags); i++ {
		fact, ok := d[tags[i]]
		if !ok {
			continue
		}

		return &models.FilteredFact{
			Category: key,
			Tag:      tags[i],
			FactData: fact,
		}, nil
	}
	return nil, fmt.Errorf("No fact found for %s", key)
}
