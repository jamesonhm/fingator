package edgar

import "github.com/jamesonhm/fingator/internal/sec/models"

func FilterDCF(cf *models.CompanyFactsResponse, d *DCFData) {

	var XBRLTags = map[string][]string{
		"CashFlow":         {"NetCashProvidedByUsedInOperatingActivities"},
		"CapEx":            {"PaymentsToAcquirePropertyPlantAndEquipment"},
		"Revenue":          {"Revenues"},
		"NetIncome":        {"NetIncomeLoss"},
		"OperatingExpense": {"OperatingExpenses"},
	}

	d.CashFlow = find(cf.Facts.Data, XBRLTags["CashFlow"])
}

func find(d map[string]models.FactData, tags []string) *models.UnitData {

}
