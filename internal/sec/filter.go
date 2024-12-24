package edgar

import (
	"fmt"

	"github.com/jamesonhm/fingator/internal/sec/models"
)

func FilterDCF(cf *models.CompanyFactsResponse, d *models.DCFData) []*models.FilteredFact {

	var XBRLTags = map[string][]string{
		"CashFlow": {"NetCashProvidedByUsedInOperatingActivities"},
		//"NetIncome":          {"NetIncomeLoss", "NetIncomeLossAvailableToCommonStockholdersBasic", "ProfitLoss"},
		//"NonCashExpense":     {"DepreciationAndAmortization", "DepreciationDepletionAndAmortization"},
		//"AccountsReceivable": {"IncreaseDecreaseInAccountsReceivable", "AccountsReceivableNetCurrent"},
		//"Inventory":          {"IncreaseDecreaseInInventories", "InventoryNet"},
		//"AccountsPayable":    {"IncreaseDecreaseInAccountsPayable", "AccountsPayableCurrent"},
		//"OperatingExpense":   {"OperatingExpenses"},
		"CapEx":           {"PaymentsToAcquirePropertyPlantAndEquipment", "CapitalExpenditures"},
		"InterestPaid":    {"InterestPaid", "InterestPaidNet", "InterestPaidCapitalized"},
		"DebtRepayment":   {"RepaymentsOfDebtAndCapitalLeaseObligations", "RepaymentsOfDebt", "RepaymentsOfLongTermDebt"},
		"DebtIssuance":    {"ProceedsFromIssuanceOfDebt", "ProceedsFromIssuanceOfLongTermDebt"},
		"EquityValue":     {"StockholdersEquity", "MarketCapitalization"},
		"DebtValue":       {"LongTermDebt", "DebtCurrent"},
		"InterestExpense": {"InterestExpense", "InterestExpenseDebt"},
		"TaxRate":         {"EffectiveIncomeTaxRateContinuingOperations", "IncomeTaxExpenseBenefit"},
		"MadeUpCategory":  {"MadeUpTag"},
	}

	var filteredFacts []*models.FilteredFact
	for key, tags := range XBRLTags {
		factData, err := findFact(cf.Facts.USGAAP, key, tags)
		if err != nil {
			fmt.Printf("%v\n", err)
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
	return nil, fmt.Errorf("\tNo fact found for %s", key)
}
