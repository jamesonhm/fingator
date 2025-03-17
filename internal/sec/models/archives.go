package models

import ()

//IndexPath = "https://www.sec.gov/Archives/edgar/full-index/2020/QTR1/{index_name}"

type BrowseArchiveParams struct {
	Year      int     `path:"year"`
	Qtr       *string `path:"qtr"`
	IndexName *string `path:"index_name"`
}
