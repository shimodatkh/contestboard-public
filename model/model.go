package model

import (
	"time"
)

type Measurement struct {
	Measurementid    int64        `json:"measurementid" gorm:"primaryKey"`
	Score            float64      `json:"score"`
	MeasurementTime  *time.Time   `json:"measurementtime"`
	Param            string       `json:"param"`
	Memo             string       `json:"memo"`
	NginxAccUrl      string       `json:"nginxaccurl"`
	NginxErrUrl      string       `json:"nginxerrurl"`
	AppUrl           string       `json:"appurl"`
	MysqlUrl         string       `json:"mysqlurl"`
	AlpUrl           string       `json:"alpurl"`
	ProfUrl          string       `json:"profurl"`
	PtqUrl           string       `json:"ptqurl"`
	TopslowqStateUrl string       `json:"topslowqstateurl"`
	TopslowqPlanUrl  string       `json:"topslowqplanurl"`
	TimeResults      []TimeResult `json:"timeresults"  gorm:"foreignKey:UserRefer;references:Measurementid;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type DataType int

const (
	Alp DataType = iota
	Prof
	Ptq
)

func (t DataType) String() string {
	switch t {
	case Alp:
		return "Alp"
	case Prof:
		return "Prof"
	case Ptq:
		return "Ptq"
	default:
		return "Unknown"
	}
}

func (t DataType) InitialChar() string {
	switch t {
	case Alp:
		return "U" // URI
	case Prof:
		return "F" // Function
	case Ptq:
		return "Q" // Query
	default:
		return "Unknown"
	}
}

type TimeResult struct {
	UserRefer uint
	DataType  DataType

	// 共通
	Name  string  `json:"name"`
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Count int64   `json:"count"`

	// alp
	HTTPMethod string `json:"httpmethod"`

	// ptq
	Checksum    string  `json:"checksum"`
	Fingerprint string  `json:"fingerprint"`
	Rawsql      string  `json:"rawsql"`
	Rate        float64 `json:"rate"`
}

type TimeResultTable struct {
	Measurementids        []int64
	TimeResultTableRowMap map[string]*TimeResultTableRow
}

// html table表示用
type TimeResultTableRow struct {
	SortValue float64
	Sums      []string
	Avgs      []string
	Counts    []string

	// ptq
	QueryID     string
	Checksum    string
	Fingerprint string
	Rawsql      string
}

type URISort struct {
	URI       string
	SortValue float64
}

// サーバ側

// Chart.jsにJSON形式で渡すための構造
type PlotPoint struct {
	X string  `json:"x"`
	Y float64 `json:"y"`
}
type PlotProject struct {
	Projectname string      `json:"projectname"`
	PlotPoints  []PlotPoint `json:"plotpoints"`
	SortValue   float64
}
type PlotProjects struct {
	Projects []PlotProject `json:"projects"`
}

// クライアント側

type CSVReadIndex struct {
	CountIdx  int
	PrefixIdx int
	NameIdx   int
	SumIdx    int
	AvgIdx    int
	DataType  DataType
}
