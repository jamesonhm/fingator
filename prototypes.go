package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/jamesonhm/fingator/internal/database"
	"github.com/jamesonhm/fingator/internal/iter"
	"github.com/jamesonhm/fingator/internal/openfigi"
	fmodels "github.com/jamesonhm/fingator/internal/openfigi/models"
	"github.com/jamesonhm/fingator/internal/polygon"
	"github.com/jamesonhm/fingator/internal/polygon/models"
	edgar "github.com/jamesonhm/fingator/internal/sec"
	emodels "github.com/jamesonhm/fingator/internal/sec/models"
)

func runEdgarTickers(ctx context.Context, dbq *database.Queries, edgarClient *edgar.Client, logger *slog.Logger) {
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Tickers Started")
	companies, err := edgarClient.GetCompanyTickers(ctx)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Error getting OHLCStartEnd", slog.Any("Error", err))
	}
	logger.LogAttrs(ctx, slog.LevelDebug, "Edgar Tickers", slog.Int("no. of companies:", len(companies)))

	for _, comp := range companies {
		_, err := dbq.CreateCompany(ctx, database.CreateCompanyParams{
			Cik:      int32(comp.CIK),
			Name:     comp.Name,
			Ticker:   comp.Ticker,
			Exchange: comp.Exchange,
		})
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"DB error adding company",
				slog.Any("company", comp),
				slog.Any("Error", err),
			)
		}
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Tickers Complete", slog.Int("no. of companies:", len(companies)))
}

func runEdgar10k(
	ctx context.Context,
	//dbq *database.Queries,
	edgarClient *edgar.Client,
	logger *slog.Logger,
) {
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Company 10k Started")
	formType := "10-K"
	resCount := 100
	// TODO: Uncomment next 9 lines in prod
	// ciks, err := dbq.GetExchangeCiks(ctx)
	// if err != nil {
	// 	logger.LogAttrs(ctx, slog.LevelError, "Error getting company facts", slog.Any("Error", err))
	// 	return
	// }
	// if len(ciks) == 0 {
	// 	logger.LogAttrs(ctx, slog.LevelError, "No CIK's found")
	// 	return
	// }

	// TODO: Comment next line in prod
	ciks := []int32{1868275, 320193, 789019}
	for _, cik := range ciks {
		cik := emodels.NumericCIK(cik).Pad()
		params := &emodels.BrowseEdgarParams{
			Action: emodels.GetCompany,
			Type:   &formType,
			Count:  &resCount,
			CIK:    &cik,
			Output: emodels.Atom,
		}
		res, err := edgarClient.FetchFilings(ctx, params)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "error getting company filings", slog.Any("Error", err))
		}

		for _, e := range res.Entries {
			//fmt.Printf("Fetch 10-K entry: %+v\n", e)
			//fmt.Println()
			fmt.Println("html page:", e.Link.Href.String())
			fmt.Println()
			link, err := edgarClient.File10kURLFromHTML(ctx, e)
			if err != nil {
				logger.LogAttrs(ctx, slog.LevelError, "error getting 10k URL", slog.Any("Error", err))
			}
			fmt.Println(link)
			r, err := edgarClient.Fetch10k(ctx, link)
			if err != nil {
				logger.LogAttrs(ctx, slog.LevelError, "error getting 10k", slog.Any("Error", err))
			}
			err = process10k(r)
			if err != nil {
				logger.LogAttrs(ctx, slog.LevelError, "error processing 10k", slog.Any("Error", err))
			}
			break
		}
		break
	}
}

func runEdgarFacts(ctx context.Context, dbq *database.Queries, edgarClient *edgar.Client, logger *slog.Logger) {
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Facts Started")
	// TODO: Uncomment next 9 lines in prod
	//ciks, err := dbq.GetExchangeCiks(ctx)
	//if err != nil {
	//	logger.LogAttrs(ctx, slog.LevelError, "Error getting company facts", slog.Any("Error", err))
	//	return
	//}
	//if len(ciks) == 0 {
	//	logger.LogAttrs(ctx, slog.LevelError, "No CIK's found")
	//	return
	//}

	// TODO: Comment next line in prod
	ciks := []int32{320193, 789019, 1868275}
	for i, cik := range ciks {
		if i >= 5 {
			break
		}
		params := &emodels.CompanyFactsParams{
			CIK: emodels.NumericCIK(cik),
		}
		res, err := edgarClient.GetCompanyFacts(ctx, params)
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"Error Getting Company Facts",
				slog.Int("cik", int(cik)),
				slog.Any("Error", err),
			)
		}

		facts := edgar.FilterDCF(ctx, res, logger)
		for _, fact := range facts {
			units, err := fact.LabelUnits()
			if err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"Unknow Units Label",
					slog.Any("Error", err),
				)
				continue
			}
			entries := fact.UnitEntries()
			for _, entry := range entries {
				err = dbq.CreateFact(ctx, database.CreateFactParams{
					Cik:          cik,
					Category:     fact.Category,
					Tag:          fact.Tag,
					Label:        fact.Label,
					Description:  fact.Description,
					Units:        units,
					EndD:         time.Time(entry.End),
					Value:        entry.Value.String(),
					FiscalYear:   int32(entry.FiscalYear),
					FiscalPeriod: entry.FiscalPeriod,
					Form:         entry.Form,
				})
				if err != nil {
					logger.LogAttrs(
						ctx,
						slog.LevelError,
						"Error adding fact data to DB",
						slog.Int("cik", int(cik)),
						slog.Any("Error", err),
					)
				}
			}
		}
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Facts Complete")
}

func runEdgarFilings(
	ctx context.Context,
	dbq *database.Queries,
	edgarClient *edgar.Client,
	logger *slog.Logger,
) {
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Filings Started")

	filers, err := dbq.GetFilers(ctx)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Unable to get 13F filers", slog.Any("Error", err))
		return
	}

	formType := "13F-HR"
	resCount := 100
	params := &emodels.BrowseEdgarParams{
		Action: emodels.GetCurrent,
		Type:   &formType,
		Count:  &resCount,
		Output: emodels.Atom,
	}
	res, err := edgarClient.FetchFilings(ctx, params)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"Fetch Filings Err",
			slog.Any("Error", err),
		)
		return
	}

	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Filings Found",
		slog.Int("Count", len(res.Entries)),
	)
	for _, entry := range res.Entries {
		cik, err := res.CIK()
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"Error saving holdings",
				slog.Any("Error", err),
			)
		}

		for _, filer := range filers {
			if filer.Cik == cik {
				logger.LogAttrs(
					ctx,
					slog.LevelInfo,
					"Filing for Tracked Filer",
					slog.String("Company", entry.Title),
					slog.Int("cik", int(cik)),
				)
				err := saveHoldings(ctx, edgarClient, logger, dbq, entry)
				if err != nil {
					logger.LogAttrs(
						ctx,
						slog.LevelError,
						"Error saving holdings",
						slog.Any("Error", err),
						slog.String("Company", entry.Title),
						slog.Int("cik", int(cik)),
					)
				}
			}
		}
	}
}

func runEdgarCompanyFilings(
	ctx context.Context,
	dbq *database.Queries,
	edgarClient *edgar.Client,
	logger *slog.Logger,
) {
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Company Filings Started")
	formType := "13F-HR"
	resCount := 100
	//cik := "0001471384"
	filers, err := dbq.GetFilers(ctx)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Unable to get 13F filers", slog.Any("Error", err))
		return
	}
	for _, filer := range filers {
		cik := emodels.NumericCIK(filer.Cik).Pad()
		params := &emodels.BrowseEdgarParams{
			Action: emodels.GetCompany,
			Type:   &formType,
			Count:  &resCount,
			CIK:    &cik,
			Output: emodels.Atom,
		}
		compRes, err := edgarClient.FetchFilings(ctx, params)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "error getting company filings", slog.Any("Error", err))
		}

		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"Filings Found",
			slog.Int("Count", len(compRes.Entries)),
			slog.String("Name", compRes.CompanyInfo.ConformedName),
		)
		for j, e := range compRes.Entries {
			cik, err := compRes.CIK()
			if err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"unable to convert Filer CIK to int",
					slog.String("Str CIK", compRes.CompanyInfo.CIK),
					slog.String("Name", compRes.CompanyInfo.ConformedName),
				)
				continue
			}

			cutoff := time.Now().AddDate(-5, 0, 0)
			if e.FilingDate().Before(cutoff) {
				logger.LogAttrs(
					ctx,
					slog.LevelInfo,
					"reached cutoff date",
					slog.String("FilingDate", e.FilingDate().String()),
					slog.String("Count", fmt.Sprintf("%d / %d", j+1, len(compRes.Entries))),
					slog.String("Name", compRes.CompanyInfo.ConformedName),
				)
				break
			}
			err = dbq.CreateFiling(ctx, database.CreateFilingParams{
				AccessionNo: e.Content.AccessionNumber,
				FilmNo:      e.FilmNo(),
				Cik:         cik,
				FilingDate:  e.FilingDate(),
			})
			if err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"Unable to create filing entry",
					slog.Any("Error", err),
					slog.String("Accession", e.AccessionNo()),
					slog.String("Link", e.Link.Href.String()),
					slog.String("Name", compRes.CompanyInfo.ConformedName),
				)
				continue
			}

			err = saveHoldings(ctx, edgarClient, logger, dbq, e)
			if err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"Error saving holdings",
					slog.Any("Error", err),
					slog.String("Company", e.Title),
					slog.Int("cik", int(cik)),
				)
			}

		}
	}
	logger.LogAttrs(ctx, slog.LevelInfo, "Edgar Company Filings Complete")
}

func saveHoldings(
	ctx context.Context,
	edgarClient *edgar.Client,
	logger *slog.Logger,
	dbq *database.Queries,
	entry emodels.FilingEntry,
) error {
	path, err := edgarClient.InfotableURLFromHTML(ctx, entry)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"no holding url found",
			slog.Any("Accession", entry.AccessionNo()),
			slog.String("Date", entry.FilingDate().String()),
			slog.String("Title", entry.Title),
		)
		return fmt.Errorf("No Holding URL found, Accession %s, Date %s", entry.AccessionNo(), entry.FilingDate().String())
	}

	holdings, err := edgarClient.FetchHoldings(ctx, path)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"Error getting holdings",
			slog.Any("Error", err),
		)
		return fmt.Errorf("Error getting holdings, %v", err)
	}

	logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		"Holdings Found",
		slog.Int("Count", len(holdings.InfoTable)),
		slog.String("Date", entry.FilingDate().String()),
		slog.String("URL", path),
	)
	for k, h := range holdings.InfoTable {
		err = dbq.CreateHolding(ctx, database.CreateHoldingParams{
			AccessionNo:  entry.Content.AccessionNumber,
			NameOfIssuer: h.NameOfIssuer,
			ClassTitle:   h.TitleOfClass,
			Cusip:        h.CUSIP,
			Value:        int64(h.Value),
			Shares:       int32(h.SharesOrPrnAmt.Amount),
			PutCall: sql.NullString{
				String: h.PutCall,
				Valid:  h.PutCall != "",
			},
		})
		if err != nil {
			logger.LogAttrs(
				ctx,
				slog.LevelError,
				"Unable to create holding entry",
				slog.Any("Error", err),
				slog.String("Count", fmt.Sprintf("%d / %d", k, len(holdings.InfoTable))),
				slog.Any("holding", h),
			)
			continue
		}
	}
	return nil
}

func runPolyGrouped(
	ctx context.Context,
	dbq *database.Queries,
	polyClient polygon.Client,
	days int,
	logger *slog.Logger,
) {
	logger.LogAttrs(ctx, slog.LevelInfo, "Polygon Grouped Daily Bars Started", slog.Int("History Days", days))

	startEnd, err := dbq.OHLCStartEnd(ctx)
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Error getting OHLCStartEnd", slog.Any("Error", err))
		return
	}
	logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		"runPolyGrouped",
		slog.Any("Start:", startEnd.Min),
		slog.Any("End:", startEnd.Max),
	)
	var minDate, maxDate *time.Time
	if start, ok := startEnd.Min.(time.Time); !ok {
		minDate = nil
	} else {
		minDate = &start
	}
	if end, ok := startEnd.Max.(time.Time); !ok {
		maxDate = nil
	} else {
		maxDate = &end
	}
	di := NewDateIter(days, minDate, maxDate, time.Now())
	for di.Next() {
		logger.LogAttrs(
			ctx,
			slog.LevelDebug,
			"GroupedDailyBars next",
			slog.Time("date", di.Date),
			slog.Int("range", di.Range()),
		)
		params := &models.GroupedDailyParams{
			Date: models.Date(di.Date),
		}
		res, err := polyClient.GroupedDailyBars(ctx, params)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "GroupedDailyBars method call", slog.Any("Error", err))
			break
		}
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"GroupedDailyBars Result",
			slog.Any("date", di.Date),
			slog.Int("result count", res.ResultCount),
			slog.String("status:", res.Status),
		)
		for i, tickerDay := range res.Results {
			// TODO: Remove in Production
			if i >= 20 {
				break
			}
			logger.LogAttrs(
				ctx,
				slog.LevelDebug,
				"GroupedDailyBars Results",
				slog.Any("TickerDay", tickerDay),
			)
			_, err := dbq.CreateTickerTimestamp(ctx, database.CreateTickerTimestampParams{
				Ticker: tickerDay.Ticker,
				Ts:     time.Time(tickerDay.Timestamp),
				Open:   strconv.FormatFloat(tickerDay.Open, 'f', 2, 64),
				High:   strconv.FormatFloat(tickerDay.High, 'f', 2, 64),
				Low:    strconv.FormatFloat(tickerDay.Low, 'f', 2, 64),
				Close:  strconv.FormatFloat(tickerDay.Close, 'f', 2, 64),
				Volume: strconv.FormatFloat(tickerDay.Volume, 'f', 2, 64),
			})
			if err != nil {
				logger.LogAttrs(ctx, slog.LevelError, "GroupedDailyBars Error adding ticker/timestamp to db", slog.Any("Error", err))
			}
		}
	}
}

func runOpenFigiCusips(
	ctx context.Context,
	dbq *database.Queries,
	figiClient openfigi.Client,
	logger *slog.Logger,
) {
	logger.LogAttrs(ctx, slog.LevelInfo, "OpenFigi CUSIP Mapping Started")
	// CUSIPS: Abbvie, Alphabet Class C, Amazon
	//cusips := []string{"00287Y109", "02079K107", "023135106"}
	cusips, err := dbq.GetUnmatchedCusips(ctx)
	if err != nil {
		logger.LogAttrs(
			ctx,
			slog.LevelError,
			"Unable to get unmatched cusips",
			slog.Any("Error", err),
		)
	}

	slIter := iter.NewSliceIter(len(cusips), figiClient.Batchsize)
	for slIter.Next() {
		logger.LogAttrs(
			ctx,
			slog.LevelInfo,
			"Unmatched CUSIPS",
			slog.Int("Total", len(cusips)),
			slog.Int("Batch Start", slIter.Start),
			slog.Int("Batch End", slIter.End),
		)
		params := []fmodels.MappingRequest{}
		for i := slIter.Start; i < slIter.End; i += 1 {
			mr := fmodels.MappingRequest{
				IDType:   fmodels.TypeCUSIP,
				IDValue:  cusips[i],
				ExchCode: fmodels.ExchUS,
			}
			//logger.LogAttrs(ctx, slog.LevelInfo, "OpenFigi Param Batch", slog.Any("MapReq", mr))
			params = append(params, mr)
		}

		res, err := figiClient.Mapping(ctx, params)
		if err != nil {
			logger.LogAttrs(ctx, slog.LevelError, "Error with mapping request", slog.Any("Error", err))
		}

		for j, obj := range *res {
			//for _, d := range obj.Data {
			//for j := 0; j < len(*res); j++ {
			fmt.Printf("%+v\n", cusips[slIter.Start+j])
			fmt.Printf("Obj Data: %+v\n", obj.Data)
			if len(obj.Data) == 0 {
				fmt.Printf("Skipping length 0\n")
				continue
			}
			err := dbq.AddCusip(ctx, database.AddCusipParams{
				Cusip:        cusips[slIter.Start+j],
				SecurityName: obj.Data[0].Name,
				Ticker:       obj.Data[0].Ticker,
				ExchangeCode: obj.Data[0].ExchangeCode,
				SecurityType: sql.NullString{
					String: obj.Data[0].SecurityType,
					Valid:  obj.Data[0].SecurityType != "",
				},
			})
			if err != nil {
				logger.LogAttrs(
					ctx,
					slog.LevelError,
					"Unable to add Cusip",
					slog.Any("Error", err),
					slog.String("CUSIP", cusips[slIter.Start+j]),
					slog.Any("Data", obj.Data[0]),
				)
			}
		}
		//}
	}
}
