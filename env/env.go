package env

// server
// 画面に表示する閾値（各計測値のランキングでothersにまとめるかどうか）
var DisplayThresholdRate float64 = 0.001
var TimeResultOthersName string = "その他"
var TimeResultNamePerTableNum int = 5
var TimeResultParamNum int = 3 // sum,avg,count
var Timeformat string = "2006-01-02T15:04:05.000000000+09:00"
var OthersSortValue float64 = -1000
var PtqNum string = "20"
var PtqSQLChecksumShortenLength int = 4
var PtqSQLStatementShortenLength int = 10
var PtqSQLStatementDisplayLength int = 1000
var DisplayForInvalidValue string = "-"

// client
var ServerHost string = "localhost"
var DefaultUrl string = "-"
