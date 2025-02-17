package edgar

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jamesonhm/fingator/internal/sec/models"
)

func FilterDCF(ctx context.Context, cf *models.CompanyFactsResponse, logger *slog.Logger) []*models.FilteredFact {

	var XBRLTags = map[string][]string{
		"CashFlow": {
			"NetCashProvidedByUsedInOperatingActivitiesContinuingOperations",
			"NetCashProvidedByUsedInOperatingActivities",
		},
		"CapEx": {
			"PaymentsToAcquirePropertyPlantAndEquipment",
			"CapitalExpenditures",
			"PaymentsToAcquireProductiveAssets",
		},
		"NetIncome": {"NetIncomeLoss", "NetIncomeLossAvailableToCommonStockholdersBasic", "ProfitLoss"},
		"Revenue": {
			"Revenues",
			"RevenueFromContractWithCustomerExcludingAssessedTax",
			"SalesRevenueGoodsNet",
			"SalesRevenueServicesNet",
			"SalesRevenueEnergyServices",
			"OperatingLeasesIncomeStatementLeaseRevenue",
			"SalesTypeLeaseRevenue"},
		"DebtRepayment": {
			"RepaymentsOfDebtAndCapitalLeaseObligations",
			"RepaymentsOfDebt",
			"RepaymentsOfConvertibleDebt",
			"RepaymentsOfLongTermDebt",
		},
		"DebtIssuance": {
			"ProceedsFromIssuanceOfDebt",
			"ProceedsFromIssuanceOfLongTermDebt",
			"ProceedsFromConvertibleDebt",
		},
		"EquityValue": {"StockholdersEquity", "MarketCapitalization"},
		"DebtValue": {
			"Liabilities",
			"LiabilitiesFairValueDisclosure",
			"LiabilitiesCurrent",
			"DebtCurrent",
			"LongTermDebt",
		},
		"InterestExpense": {"InterestExpense", "InterestExpenseDebt"},
		"TaxRate":         {"EffectiveIncomeTaxRateContinuingOperations", "IncomeTaxExpenseBenefit"},
		"Shares":          {"WeightedAverageNumberOfDilutedSharesOutstanding"},
		//"NonCashExpense":     {"DepreciationAndAmortization", "DepreciationDepletionAndAmortization"},
		//"AccountsReceivable": {"IncreaseDecreaseInAccountsReceivable", "AccountsReceivableNetCurrent"},
		//"Inventory":          {"IncreaseDecreaseInInventories", "InventoryNet"},
		//"AccountsPayable":    {"IncreaseDecreaseInAccountsPayable", "AccountsPayableCurrent"},
		//"OperatingExpense":   {"OperatingExpenses"},
		//"InterestPaid": {"InterestPaid", "InterestPaidNet", "InterestPaidCapitalized"},
		//"MadeUpCategory":  {"MadeUpTag"},
	}

	var filteredFacts []*models.FilteredFact
	for key, tags := range XBRLTags {
		factData, err := findFact(cf.Facts.USGAAP, key, tags)
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"Unable to find Fact",
				slog.Int("CIK", int(cf.CIK)),
				slog.String("Key", key),
			)
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
