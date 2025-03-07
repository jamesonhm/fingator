package edgar

import (
	"context"
	//"encoding/json"
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
			"SalesRevenueNet",
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
		factData, err := findFact(ctx, cf.Facts.USGAAP, key, tags, logger)
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"Unable to find Fact",
				slog.Int("CIK", int(cf.CIK)),
				slog.String("Key", key),
			)
			continue
			//break
		}
		filteredFacts = append(filteredFacts, factData)
	}
	return filteredFacts
}

func findFact(
	ctx context.Context,
	d map[string]models.FactData,
	key string,
	tags []string,
	logger *slog.Logger,
) (*models.FilteredFact, error) {
	for i := 0; i < len(tags); i++ {
		fact, ok := d[tags[i]]
		if !ok {
			continue
		}
		err := fact.Filter()
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"Fact skipped, Filter error",
				slog.Any("Error", err),
				slog.String("Key", key),
				slog.String("Tag", tags[i]),
			)
			continue
		}
		if fact.Age() > 1 {
			logger.LogAttrs(
				ctx,
				slog.LevelInfo,
				"Fact skipped, age > 1",
				slog.String("Key", key),
				slog.String("Tag", tags[i]),
				slog.Int("FY", fact.LastFY()),
			)
			continue
			//break
		}

		return &models.FilteredFact{
			Category: key,
			Tag:      tags[i],
			FactData: fact,
		}, nil
	}
	return nil, fmt.Errorf("\tNo fact found for %s", key)
}
