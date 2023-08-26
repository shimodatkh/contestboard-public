package model

type PTQJSON struct {
	Classes []PTQJSONClass `json:"classes"`
	Global  PTQJSONGlobal  `json:"global"`
}

type PTQJSONClass struct {
	Attribute   string                 `json:"attribute"`
	Checksum    string                 `json:"checksum"`
	Distillate  string                 `json:"distillate"`
	Example     PTQJSONClassExample    `json:"example"`
	Fingerprint string                 `json:"fingerprint"`
	Histograms  PTQJSONClassHistograms `json:"histograms"`
	Metrics     PTQJSONClassMetrics    `json:"metrics"`
	QueryCount  int                    `json:"query_count"`
	Tables      []PTQJSONClassTable    `json:"tables"`
	TsMax       string                 `json:"ts_max"`
	TsMin       string                 `json:"ts_min"`
}

type PTQJSONClassExample struct {
	QueryTime string `json:"Query_time"`
	Query     string `json:"query"`
	Ts        string `json:"ts"`
}

type PTQJSONClassHistograms struct {
	QueryTime []int `json:"Query_time"`
}

type PTQJSONClassMetrics struct {
	LockTime     PTQJSONClassMetricsLockTime     `json:"Lock_time"`
	QueryLength  PTQJSONClassMetricsQueryLength  `json:"Query_length"`
	QueryTime    PTQJSONClassMetricsQueryTime    `json:"Query_time"`
	RowsExamined PTQJSONClassMetricsRowsExamined `json:"Rows_examined"`
	RowsSent     PTQJSONClassMetricsRowsSent     `json:"Rows_sent"`
	Host         PTQJSONClassMetricsHost         `json:"host"`
	User         PTQJSONClassMetricsUser         `json:"user"`
}

type PTQJSONClassMetricsLockTime struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct    string `json:"pct"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONClassMetricsQueryLength struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct    string `json:"pct"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONClassMetricsQueryTime struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct    string `json:"pct"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONClassMetricsRowsExamined struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct    string `json:"pct"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONClassMetricsRowsSent struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct    string `json:"pct"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONClassMetricsHost struct {
	Value string `json:"value"`
}

type PTQJSONClassMetricsUser struct {
	Value string `json:"value"`
}

type PTQJSONClassTable struct {
	Create string `json:"create"`
	Status string `json:"status"`
}

type PTQJSONGlobal struct {
	Files            []PTQJSONGlobalFile  `json:"files"`
	Metrics          PTQJSONGlobalMetrics `json:"metrics"`
	QueryCount       int                  `json:"query_count"`
	UniqueQueryCount int                  `json:"unique_query_count"`
}

type PTQJSONGlobalFile struct {
	Name string  `json:"name"`
	Size float64 `json:"size"`
}

type PTQJSONGlobalMetrics struct {
	LockTime     PTQJSONGlobalMetricsLockTime     `json:"Lock_time"`
	QueryLength  PTQJSONGlobalMetricsQueryLength  `json:"Query_length"`
	QueryTime    PTQJSONGlobalMetricsQueryTime    `json:"Query_time"`
	RowsExamined PTQJSONGlobalMetricsRowsExamined `json:"Rows_examined"`
	RowsSent     PTQJSONGlobalMetricsRowsSent     `json:"Rows_sent"`
}

type PTQJSONGlobalMetricsLockTime struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONGlobalMetricsQueryLength struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONGlobalMetricsQueryTime struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONGlobalMetricsRowsExamined struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}

type PTQJSONGlobalMetricsRowsSent struct {
	Avg    string `json:"avg"`
	Max    string `json:"max"`
	Median string `json:"median"`
	Min    string `json:"min"`
	Pct95  string `json:"pct_95"`
	Stddev string `json:"stddev"`
	Sum    string `json:"sum"`
}
