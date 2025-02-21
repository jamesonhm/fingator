WITH findata AS (
	SELECT end_d,
			category,
			value
	FROM public.facts
	WHERE cik = 1868275
	AND category in ('Revenue', 'NetIncome', 'CashFlow', 'CapEx')
	AND form = '10-K'
),
crossfindata AS (
	SELECT end_d,
			max(value) filter (where category = 'Revenue') as revenue,
			max(value) filter (where category = 'NetIncome') as netincome,
			max(value) filter (where category = 'CashFlow') as cashflow,
			max(value) filter (where category = 'CapEx') as capex
	FROM findata
	group by end_d
	ORDER BY end_d
)
SELECT *,
		cashflow - capex as fcf,
		((cashflow - capex) - LAG(cashflow - capex) OVER (ORDER BY end_d)) / LAG(cashflow - capex) OVER (ORDER BY end_d) * 100 as pctgrowth 
FROM crossfindata
