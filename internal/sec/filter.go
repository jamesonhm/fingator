package edgar

import (
	"context"
	//"encoding/json"
	"fmt"
	"log/slog"

	"github.com/jamesonhm/fingator/internal/sec/models"
)

func FilterBasicFinancials(
	ctx context.Context,
	cf *models.CompanyFactsResponse,
	logger *slog.Logger,
) []*models.FilteredFact {
	var cashFlowTags = map[string][]string{
		"NetIncome": {"NetIncomeLoss", "ProfitLoss", "NetIncomeLossAvailableToCommonStockholdersBasic"},
		"DandA": {
			"DepreciationAmortizationAndOther",
			"DepreciationDepletionAndAmortization",
			"DepreciationAmortizationAndAccretionNet",
		},
		"Depreciation": {"Depreciation"},
		"Amortization": {
			"Amortization",
			"AmortizationOfIntangibleAssets",
		},
		"NetCashOps": {
			"NetCashProvidedByUsedInOperatingActivities",
			"NetCashProvidedByUsedInOperatingActivitiesContinuingOperations",
		},
		"CapEx": {
			"PaymentsToAcquirePropertyPlantAndEquipment",
			"CapitalExpenditures",
			"PaymentsToAcquireProductiveAssets",
		},
		"DebtIssuance": {
			"ProceedsFromIssuanceOfLongTermDebt",
			"ProceedsFromDebtMaturingInMoreThanThreeMonths",
			"ProceedsFromIssuanceOfDebt",
			"ProceedsFromConvertibleDebt",
		},
		"DebtRepayment": {
			"RepaymentsOfLongTermDebt",
			"RepaymentsOfDebtMaturingInMoreThanThreeMonths",
			"RepaymentsOfDebt",
			"RepaymentsOfDebtAndCapitalLeaseObligations",
			"RepaymentsOfConvertibleDebt",
		},
	}
	var incomeTags = map[string][]string{
		"Revenue": {
			"Revenues",
			"RevenueFromContractWithCustomerExcludingAssessedTax",
		},
		"EBIT": {
			"OperatingIncomeLoss",
		},
		"TaxExpense": {
			"IncomeTaxExpenseBenefit",
		},
		"PreTaxIncome": {
			"IncomeLossFromContinuingOperationsBeforeIncomeTaxesMinorityInterestAndIncomeLossFromEquityMethodInvestments",
			"IncomeLossFromContinuingOperationsBeforeIncomeTaxesExtraordinaryItemsNoncontrollingInterest",
		},
		"EPS":    {"EarningsPerShareBasic"},
		"Shares": {"WeightedAverageNumberOfSharesOutstandingBasic"},
	}
	var balanceTags = map[string][]string{
		"TotalCurrentAssets":      {"AssetsCurrent"},
		"TotalCurrentLiabilities": {"LiabilitiesCurrent"},
		"ShareholderEquity":       {"StockholdersEquity"},
	}
	var aggBalanceTags = map[string][]string{
		"NonOpAssets": {
			"CashAndCashEquivalentsAtCarryingValue",
			"CashCashEquivalentsAndShortTermInvestments",
			"MarketableSecuritiesCurrent",
			"RestrictedCashAndCashEquivalentsAtCarryingValue",
			"ShortTermInvestments",
			"DerivativeAssetsCurrent",
		},
		"OpAssets": {
			"AccountsReceivableNetCurrent",
			"InventoryNet",
			"InventoryRawMaterialsAndSuppliesNetOfReserves",
			"PrepaidExpenseAndOtherAssetsCurrent",
			"NontradeReceivablesCurrent",
			"OtherReceivablesNetCurrent",
			"OtherAssetsCurrent",
			"EnergyRelatedInventoryNaturalGasInStorage",
			"InventoryRawMaterialsAndSuppliesNetOfReserves",
			"RenewableEnergyCreditsCurrent",
		},
		"OpLiabilities": {
			"AccountsPayableCurrent",
			"AccountsPayableAndAccruedLiabilitiesCurrent",
			"AccountsPayableTradeCurrent",
			"OtherLiabilitiesCurrent",
			"AccruedLiabilitiesCurrent",
			"EmployeeRelatedLiabilitiesCurrent",
			"AccruedIncomeTaxesCurrent",
			"ContractWithCustomerLiabilityCurrent",
			"OtherLiabilitiesCurrent",
			"EnergyMarketingAccountsPayable",
			"OperatingLeaseLiabilityCurrent",
		},
		"NonOpLiabilities": {
			"CommercialPaper",
			"LongTermDebtCurrent",
			"ShortTermBorrowings",
			"LongTermDebtAndCapitalLeaseObligationsCurrent",
			"DerivativeLiabilitiesCurrent",
		},
	}
	var aggCashFlowTags = map[string][]string{
		"CFChangeNWC": {
			"CashAndCashEquivalentsAtCarryingValue",
			"CashCashEquivalentsAndShortTermInvestments",
			"MarketableSecuritiesCurrent",
			"RestrictedCashAndCashEquivalentsAtCarryingValue",
			"ShortTermInvestments",
			"DerivativeAssetsCurrent",
		},
	}

	var filteredFacts []*models.FilteredFact
	for key, tags := range cashFlowTags {
		factData, err := findFact(ctx, cf.Facts.USGAAP, "CashFlow", key, tags, logger)
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
	for key, tags := range incomeTags {
		factData, err := findFact(ctx, cf.Facts.USGAAP, "Income", key, tags, logger)
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
	for key, tags := range balanceTags {
		factData, err := findFact(ctx, cf.Facts.USGAAP, "Balance", key, tags, logger)
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
	for key, tags := range aggBalanceTags {
		for _, foundFact := range findAllFact(ctx, cf.Facts.USGAAP, "Balance", key, tags, logger) {
			filteredFacts = append(filteredFacts, foundFact)
		}
	}
	return filteredFacts
}

func findFact(
	ctx context.Context,
	d map[string]models.FactData,
	stmnt string,
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
			//logger.LogAttrs(
			//	ctx,
			//	slog.LevelInfo,
			//	"Fact skipped, age > 1",
			//	slog.String("Key", key),
			//	slog.String("Tag", tags[i]),
			//	slog.Int("FY", fact.LastFY()),
			//)
			continue
		}

		return &models.FilteredFact{
			Statement: stmnt,
			Category:  key,
			Tag:       tags[i],
			FactData:  fact,
		}, nil
	}
	return nil, fmt.Errorf("\tNo fact found for %s", key)
}

func findAllFact(
	ctx context.Context,
	d map[string]models.FactData,
	stmnt string,
	key string,
	tags []string,
	logger *slog.Logger,
) []*models.FilteredFact {
	var foundFacts []*models.FilteredFact
	for i := 0; i < len(tags); i++ {
		fact, ok := d[tags[i]]
		if !ok {
			continue
		}
		err := fact.Filter()
		if err != nil {
			continue
		}
		if fact.Age() > 1 {
			continue
		}

		foundFacts = append(foundFacts, &models.FilteredFact{
			Statement: stmnt,
			Category:  key,
			Tag:       tags[i],
			FactData:  fact,
		})
	}
	return foundFacts
}
