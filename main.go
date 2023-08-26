package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/joho/godotenv"
	"github.com/shimodatkh/contestboard-public/env"
	"github.com/shimodatkh/contestboard-public/log"
	"github.com/shimodatkh/contestboard-public/model"
	"github.com/shimodatkh/contestboard-public/repository"
	"go.uber.org/zap/zapcore"
)

var mu sync.Mutex
var sr *repository.SqliteRepository

// user auth
var authUser = "user"
var authPw = "pass"

// 認証関数
func checkAuth(req *http.Request) bool {
	// req.BasicAuthメソッドで送信されたユーザ名・パスワードを取得
	user, pw, ok := req.BasicAuth()
	// 正しいユーザ名・パスワードと比較
	if !ok || user != authUser || pw != authPw {
		return false
	}
	return true
}

// Basic認証を行うハンドラ関数
func basicAuthHandler(w http.ResponseWriter, r *http.Request) {
	// 認証関数の実行
	if !checkAuth(r) {
		// 認証に失敗した場合のヘッダ情報付与
		w.Header().Add("WWW-Authenticate", `Basic realm="my private area"`)
		w.WriteHeader(http.StatusUnauthorized) // 401コード
		// 認証失敗時の出力内容
		w.Write([]byte("401 認証失敗\n"))
		return
	}
	// 認証成功時の出力内容
	mainHandler(w, r)
}

func main() {
	log := log.InitLogger(zapcore.DebugLevel)

	// env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// DB Setting
	useInmemoryDBStr := os.Getenv("USE_INMEMORY_DB")
	useInmemoryDB, _ := strconv.ParseBool(useInmemoryDBStr)
	if useInmemoryDB {
		log.Infow("use sqlite in in-memory mode")
	} else {
		log.Infow("use sqlite in file mode")
	}
	sr = repository.NewSqliteRepository(useInmemoryDB)

	// for only dev
	insertDebugDataStr := os.Getenv("INSERT_DEBUG_DATA")
	insertDebugData, _ := strconv.ParseBool(insertDebugDataStr)
	if insertDebugData {
		initialDataInsertForDev()
	}

	dir, _ := os.Getwd()
	log.Infow("Server start.", "static-files-dir", http.Dir(dir+"/static/"))
	// http.HandleFunc("/", basicAuthHandler)
	http.HandleFunc("/", mainHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(dir+"/static/"))))
	http.HandleFunc("/measurement", measurementHandler)
	http.HandleFunc("/entryMeasurement", entryMeasurementHandler)
	http.HandleFunc("/fixMeasurement", fixMeasurementHandler)
	http.HandleFunc("/deleteMeasurement", deleteMeasurementHandler)
	http.HandleFunc("/deleteAllMeasurements", deleteAllMeasurementsHandler)
	http.ListenAndServe(":8080", nil)
}

// メイン画面描画
func mainHandler(w http.ResponseWriter, r *http.Request) {
	log := log.InitLogger(zapcore.DebugLevel)
	log.Infow("mainHandler start")
	start := time.Now()

	username, password, ok := r.BasicAuth()
	if !ok {
		basicAuthHandler(w, r)
	}
	if !(username == authUser && password == authPw) {
		return
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		log.Errorw("template.ParseFiles index.html error:", "err", err.Error())
	}
	measurements := sr.GetAllMeasurements()

	// 表示用の加工
	uuidreg := "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}"
	for _, measurement := range measurements {
		// MeasurementTimeはJSTに変換
		jst := time.FixedZone("Asia/Tokyo", 9*60*60) // JSTのタイムゾーンを定義
		measurementTime := measurement.MeasurementTime.In(jst)
		measurement.MeasurementTime = &measurementTime

		// alpの見づらい正規表現を置換する
		newtr := make([]model.TimeResult, 0)
		for _, timeresult := range measurement.TimeResults {
			timeresult.Name = strings.Replace(timeresult.Name, uuidreg, "{UUID}", 1)
			newtr = append(newtr, timeresult)
		}
		measurement.TimeResults = newtr
	}

	if err := t.Execute(w, struct {
		Measurements          []*model.Measurement
		ScoreChartDataJSONStr string
		AlpChartDataJSONStr   string
		ProfChartDataJSONStr  string
		PtqChartDataJSONStr   string
		DividedAlpTables      []*model.TimeResultTable
		DividedProfTables     []*model.TimeResultTable
		DividedPtqTables      []*model.TimeResultTable
		AlpURINumPerTable     int
		AlpColnumPerTable     int
	}{
		Measurements:          measurements,
		ScoreChartDataJSONStr: makeScoreChartDataJSONStr(measurements),
		AlpChartDataJSONStr:   makeTimeResultChartDataJSONStr(measurements, model.Alp),
		ProfChartDataJSONStr:  makeTimeResultChartDataJSONStr(measurements, model.Prof),
		PtqChartDataJSONStr:   makeTimeResultChartDataJSONStr(measurements, model.Ptq),
		DividedAlpTables:      divideTimeResultTable(makeTimeResultTableForDisplay(measurements, model.Alp)),
		DividedProfTables:     divideTimeResultTable(makeTimeResultTableForDisplay(measurements, model.Prof)),
		DividedPtqTables:      divideTimeResultTable(makeTimeResultTableForDisplay(measurements, model.Ptq)),
		AlpURINumPerTable:     env.TimeResultNamePerTableNum,
		AlpColnumPerTable:     env.TimeResultNamePerTableNum * env.TimeResultParamNum,
	}); err != nil {
		panic(err.Error())
	}
	end := time.Now()
	log.Infow("mainHandler end", "totalExecuteTimeSeconds", (end.Sub(start)).Seconds())
}

func makeTimeResultTableForDisplay(measurements []*model.Measurement, dataType model.DataType) *model.TimeResultTable {
	log := log.InitLogger(zapcore.DebugLevel)
	var timeResultTable model.TimeResultTable = model.TimeResultTable{
		Measurementids:        []int64{},
		TimeResultTableRowMap: make(map[string]*model.TimeResultTableRow),
	}
	chartPlotProjectMap := extractDataToDisplayFromMeasurements(measurements, dataType)
	for k, v := range chartPlotProjectMap {
		timeResultTable.TimeResultTableRowMap[k] = &model.TimeResultTableRow{
			SortValue: v.SortValue,
		}
	}

	// 計測ごとに、
	for i := 0; i < len(measurements); i++ {
		timeResultTable.Measurementids = append(timeResultTable.Measurementids, measurements[i].Measurementid)
		othersSumTime := 0.0
		othersCounts := 0
		// 各URIの、（当該の計測で出てこないURIにも-を埋めるために掃過する）
		// 1.TimeResultsに入っているURIを埋める
		timeResults := selectTimeResultType(measurements[i], dataType)
		for _, timeResult := range timeResults {
			// 値があれば詰め替え
			v, ok := timeResultTable.TimeResultTableRowMap[timeResult.Name]
			if ok {
				v.Sums = append(v.Sums, fmt.Sprintf("%.1f", timeResult.Sum))
				v.Avgs = append(v.Avgs, fmt.Sprintf("%.3f", timeResult.Avg))
				v.Counts = append(v.Counts, fmt.Sprintf("%d", timeResult.Count))
				if i == 0 { // i=0が最新計測
					v.SortValue = timeResult.Sum

					// for ptq
					v.Checksum = timeResult.Checksum
					v.Fingerprint = timeResult.Fingerprint
					v.Rawsql = timeResult.Rawsql
				}
				// isNotExistInURIList = false
			} else { // mapに無い＝others判定（無駄に同じ計算を複数回している）
				othersSumTime += timeResult.Sum
				othersCounts += int(timeResult.Count)
			}
		}
		// 2.TimeResultsに入っていないURIを埋める
		for key, v := range timeResultTable.TimeResultTableRowMap {
			// othersではなく、かつまだsumsに追加されてないやつ
			if env.TimeResultOthersName != key && len(v.Sums) != i+1 {
				v, ok := timeResultTable.TimeResultTableRowMap[key]
				if ok {
					v.Sums = append(v.Sums, env.DisplayForInvalidValue)
					v.Avgs = append(v.Avgs, env.DisplayForInvalidValue)
					v.Counts = append(v.Counts, env.DisplayForInvalidValue)
				}
			}
		}
		// 3.othersを埋める
		v, ok := timeResultTable.TimeResultTableRowMap[env.TimeResultOthersName]
		if ok {
			v.Sums = append(v.Sums, fmt.Sprintf("%.1f", othersSumTime))
			v.Avgs = append(v.Avgs, env.DisplayForInvalidValue)
			v.Counts = append(v.Counts, fmt.Sprintf("%d", othersCounts))
			// othersは末尾にするためSortValueを更新しない
		} else { // mapに無い＝others判定
			log.Infow("others key is not in TableURIRowMap", "env.TimeResultOthersName", env.TimeResultOthersName)
		}
	}

	// URIを表示したい順に並べる
	urisorts := []model.URISort{}
	for k, v := range timeResultTable.TimeResultTableRowMap {
		urisorts = append(urisorts, model.URISort{
			URI:       k,
			SortValue: v.SortValue,
		})
	}
	sort.Slice(urisorts, func(i, j int) bool {
		if urisorts[i].SortValue != 0 || urisorts[j].SortValue != 0 {
			return urisorts[i].SortValue > urisorts[j].SortValue
		}
		return urisorts[i].URI < urisorts[j].URI
	})

	// log.Debugw("URI sort", "urisorts", urisorts)

	urinums := len(timeResultTable.TimeResultTableRowMap)

	// 複数テーブルに分割、入れ物用意
	renamedTimeResultTable := model.TimeResultTable{
		Measurementids:        timeResultTable.Measurementids,
		TimeResultTableRowMap: make(map[string]*model.TimeResultTableRow),
	}
	for i := 0; i < urinums; i++ {
		v := timeResultTable.TimeResultTableRowMap[urisorts[i].URI]
		initial := dataType.InitialChar()
		newKey := fmt.Sprintf("["+initial+"%02d] ", i+1) + urisorts[i].URI

		// for ptq
		v.QueryID = fmt.Sprintf("["+initial+"%02d] ", i+1)

		renamedTimeResultTable.TimeResultTableRowMap[newKey] = v
	}

	return &renamedTimeResultTable
}

func selectTimeResultType(measurement *model.Measurement, dataType model.DataType) []model.TimeResult {
	var tr []model.TimeResult
	for _, v := range measurement.TimeResults {
		if v.DataType == dataType {
			tr = append(tr, v)
		}
	}
	return tr
}

// 画面の横幅制約により、テーブルを分割するのに合わせてデータも切っておく
func divideTimeResultTable(timeResultTable *model.TimeResultTable) []*model.TimeResultTable {
	urinums := len(timeResultTable.TimeResultTableRowMap)

	var dividedTimeResultTables = []*model.TimeResultTable{}
	tempTimeResultTable := model.TimeResultTable{
		Measurementids:        timeResultTable.Measurementids,
		TimeResultTableRowMap: make(map[string]*model.TimeResultTableRow),
	}

	urisorts := []model.URISort{}
	for k, v := range timeResultTable.TimeResultTableRowMap {
		urisorts = append(urisorts, model.URISort{
			URI:       k,
			SortValue: v.SortValue,
		})
	}

	sort.Slice(urisorts, func(i, j int) bool {
		if urisorts[i].SortValue != 0 || urisorts[j].SortValue != 0 {
			return urisorts[i].SortValue > urisorts[j].SortValue
		}
		return urisorts[i].URI > urisorts[j].URI
	})
	for i := 0; i < urinums; i++ {
		tempTimeResultTable.TimeResultTableRowMap[urisorts[i].URI] = timeResultTable.TimeResultTableRowMap[urisorts[i].URI]
		if i%env.TimeResultNamePerTableNum == env.TimeResultNamePerTableNum-1 || i == urinums-1 {
			// 分割分詰め込み
			copyTimeResultTable := &model.TimeResultTable{
				Measurementids:        tempTimeResultTable.Measurementids,
				TimeResultTableRowMap: tempTimeResultTable.TimeResultTableRowMap,
			}
			dividedTimeResultTables = append(dividedTimeResultTables, copyTimeResultTable)
			// 初期化
			tempTimeResultTable = model.TimeResultTable{
				Measurementids:        timeResultTable.Measurementids,
				TimeResultTableRowMap: make(map[string]*model.TimeResultTableRow),
			}
		}
	}
	return dividedTimeResultTables
}

func makeScoreChartDataJSONStr(measurements []*model.Measurement) string {
	p := model.PlotProjects{}
	var mpn model.PlotProject
	jst := time.FixedZone("Asia/Tokyo", 9*60*60) // JSTのタイムゾーンを定義
	for _, measurement := range measurements {
		jstTime := measurement.MeasurementTime.In(jst) // JSTに変換
		mpn.PlotPoints = append(mpn.PlotPoints, model.PlotPoint{
			X: jstTime.Format(env.Timeformat),
			Y: measurement.Score,
		})
	}
	p.Projects = append(p.Projects, mpn)
	return jsonMarshalToString(p)
}

func jsonMarshalToString(p model.PlotProjects) string {
	log := log.InitLogger(zapcore.DebugLevel)
	bytes, err := json.Marshal(p)
	if err != nil {
		log.Errorw("JSON marshal error.", "err", err.Error())
		return ""
	}
	return string(bytes)
}

func makeTimeResultChartDataJSONStr(measurements []*model.Measurement, dataType model.DataType) string {
	p := model.PlotProjects{}
	TimeResultTable := makeTimeResultTableForDisplay(measurements, dataType)
	jst := time.FixedZone("Asia/Tokyo", 9*60*60) // JSTのタイムゾーンを定義
	for k, v := range TimeResultTable.TimeResultTableRowMap {
		plotpoints := []model.PlotPoint{}
		for i := 0; i < len(TimeResultTable.Measurementids); i++ {
			if v.Sums[i] != env.DisplayForInvalidValue {
				f, _ := strconv.ParseFloat(v.Sums[i], 64)
				newPlotPoint := model.PlotPoint{X: measurements[i].MeasurementTime.In(jst).Format(env.Timeformat), Y: f}
				plotpoints = append(plotpoints, newPlotPoint)
			}
		}
		p.Projects = append(p.Projects, model.PlotProject{
			Projectname: k,
			PlotPoints:  plotpoints,
			SortValue:   v.SortValue,
		})
	}

	sort.Slice(p.Projects, func(i, j int) bool {
		if p.Projects[i].SortValue != 0 || p.Projects[j].SortValue != 0 {
			return p.Projects[i].SortValue > p.Projects[j].SortValue
		}
		return p.Projects[i].Projectname < p.Projects[j].Projectname
	})

	return jsonMarshalToString(p)
}

// 画面に表示するデータの選別(others判定を含む)
func extractDataToDisplayFromMeasurements(measurements []*model.Measurement, dataType model.DataType) map[string]*model.PlotProject {
	plotProjectMap := make(map[string]*model.PlotProject)
	for _, measurement := range measurements {
		// timeResultデータの蓄積、結果の行ごとに処理
		timeResults := selectTimeResultType(measurement, dataType)
		totalTime := calcTotalTime(timeResults)
		thresholdTime := totalTime * env.DisplayThresholdRate
		for _, timeResult := range timeResults {
			_, ok := plotProjectMap[timeResult.Name]
			if !ok {
				if timeResult.Sum > thresholdTime {
					plotProjectMap[timeResult.Name] = &model.PlotProject{
						Projectname: timeResult.Name,
						SortValue:   0,
						PlotPoints:  []model.PlotPoint{}}
				}
			}
		}
	}
	// othersは必ず入れる
	plotProjectMap[env.TimeResultOthersName] = &model.PlotProject{
		Projectname: env.TimeResultOthersName,
		SortValue:   env.OthersSortValue,
		PlotPoints:  []model.PlotPoint{}}
	return plotProjectMap
}

func calcTotalTime(timeResult []model.TimeResult) float64 {
	totalTime := 0.0
	for _, v := range timeResult {
		totalTime += v.Sum
	}
	return totalTime
}

func measurementHandler(w http.ResponseWriter, r *http.Request) {
	log := log.InitLogger(zapcore.DebugLevel)
	log.Infow("measurementHandler start")
	start := time.Now()

	// バリデーション
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// データパース
	buffer, err := ioutil.ReadAll(r.Body)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusConflict)
		return
	}
	defer r.Body.Close()

	var postedMeasurement model.Measurement
	err = json.Unmarshal(buffer, &postedMeasurement)
	if err != nil {
		log.Errorw("json unmarshal error.", "err", err.Error())
		return
	}

	mu.Lock()
	addMeasurement(postedMeasurement)
	mu.Unlock()
	end := time.Now()
	log.Infow("measurementHandler end", "totalExecuteTimeSeconds", (end.Sub(start)).Seconds())
}

func addMeasurement(postedMeasurement model.Measurement) {
	sr.InsertMeasurement(&postedMeasurement)
}

func fixMeasurementHandler(w http.ResponseWriter, r *http.Request) {
	log := log.InitLogger(zapcore.DebugLevel)
	log.Infow("fixMeasurementHandler start")

	// バリデーション
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// log.Infow("data", "score", r.FormValue("score"), "mid", r.FormValue("mid"), "memo", r.FormValue("memo"))

	measurementid, err := strconv.ParseInt(r.FormValue("mid"), 10, 64)
	if err != nil {
		log.Errorw("mid int parse error", "err", err.Error(), "r.FormValue(mid)", r.FormValue("mid"))
	}

	formScore := strings.Trim(r.FormValue("score"), " ")
	if formScore != "" {
		score, err := strconv.ParseFloat(formScore, 64)
		if err != nil {
			log.Errorw("ParseFloat failed.", "err", err.Error(), "scoreValue", formScore)
		}
		sr.UpdateMeasurementByScore(measurementid, score)
	} else if r.FormValue("memo") != "" {
		sr.UpdateMeasurementByMemo(measurementid, r.FormValue("memo"))
	} else if r.FormValue("param") != "" {
		sr.UpdateMeasurementByParam(measurementid, r.FormValue("param"))
	}

}

func entryMeasurementHandler(w http.ResponseWriter, r *http.Request) {
	log := log.InitLogger(zapcore.DebugLevel)
	log.Infow("entryMeasurementHandler start")

	newIDint := sr.GetNewMeasurementID() + 1
	newID := fmt.Sprintf("%d", newIDint)
	log.Infow("next ID is", "newID", newID)
	w.Write([]byte(newID))
}

func deleteMeasurementHandler(w http.ResponseWriter, r *http.Request) {
	log := log.InitLogger(zapcore.DebugLevel)
	log.Infow("deleteMeasurementHandler start")

	// バリデーション
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// log.Infow("data", "score", r.FormValue("score"), "mid", r.FormValue("mid"), "memo", r.FormValue("memo"))

	measurementid, err := strconv.ParseInt(r.FormValue("mid"), 10, 64)
	if err != nil {
		log.Errorw("mid int parse error", "err", err.Error(), "r.FormValue(mid)", r.FormValue("mid"))
	}

	sr.DeleteMeasurement(measurementid)
}

func deleteAllMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	log := log.InitLogger(zapcore.DebugLevel)
	log.Infow("deleteAllMeasurementsHandler start")

	// バリデーション
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO:

}

func initialDataInsertForDev() {
	uri1 := "GET /api/get1/FirstBottleNeck"
	uri2 := "GET /api/get2/SecondBottleNeck"
	uri3 := "POST /api/post1/London"
	uri4 := "POST /api/post2/Birmingham-Glasgow-Liverpool"
	uri5 := "POST /api/post3/Leeds"
	func1 := "main.function1Paris"
	func2 := "main.function1Paris: part1"
	func3 := "main.function1Paris: part2"
	func4 := "main.function2Bordeaux "
	func5 := "main.function3Marseille"
	mtime1 := time.Now().Add(-25 * time.Minute)
	mtime2 := mtime1.Add(5 * time.Minute)
	mtime3 := mtime1.Add(10 * time.Minute)
	mtime4 := mtime1.Add(12 * time.Minute)
	mtime5 := mtime1.Add(13 * time.Minute)
	mea1 := model.Measurement{
		Measurementid: 1,
		Score:         100,
		Memo:          "テスト1",
		TimeResults: []model.TimeResult{
			{DataType: model.Alp, Name: uri1, Sum: 800, Avg: 800 / 100, Count: 100},
			{DataType: model.Alp, Name: uri2, Sum: 500, Avg: 500 / 200, Count: 200},
			{DataType: model.Alp, Name: uri3, Sum: 100, Avg: 100 / 400, Count: 400},
			{DataType: model.Alp, Name: uri4, Sum: 80, Avg: 80 / 4200, Count: 4200},
			{DataType: model.Alp, Name: uri5, Sum: 10, Avg: 10 / 100, Count: 100},
			{DataType: model.Prof, Name: func1, Sum: 80, Avg: 800 / 10, Count: 100},
			{DataType: model.Prof, Name: func2, Sum: 500, Avg: 500 / 200, Count: 200},
			{DataType: model.Prof, Name: func3, Sum: 100, Avg: 100 / 400, Count: 400},
			{DataType: model.Prof, Name: func4, Sum: 380, Avg: 380 / 4200, Count: 4200},
			{DataType: model.Prof, Name: func5, Sum: 10, Avg: 10 / 100, Count: 100},
		},
		MeasurementTime: &mtime1}
	mea2 := model.Measurement{
		Measurementid: 2,
		Score:         150,
		Memo:          "テスト2",
		TimeResults: []model.TimeResult{
			{DataType: model.Alp, Name: uri1, Sum: 200, Avg: 200 / 100, Count: 100},
			{DataType: model.Alp, Name: uri2, Sum: 700, Avg: 700 / 200, Count: 280},
			{DataType: model.Alp, Name: uri3, Sum: 150, Avg: 150 / 400, Count: 600},
			{DataType: model.Alp, Name: uri4, Sum: 10, Avg: 10 / 4200, Count: 4200},
			{DataType: model.Alp, Name: uri5, Sum: 120, Avg: 120 / 2230, Count: 2230},
			{DataType: model.Prof, Name: func1, Sum: 200, Avg: 200 / 100, Count: 100},
			{DataType: model.Prof, Name: func2, Sum: 170, Avg: 170 / 200, Count: 280},
			{DataType: model.Prof, Name: func3, Sum: 150, Avg: 150 / 400, Count: 600},
			{DataType: model.Prof, Name: func4, Sum: 10, Avg: 10 / 4200, Count: 4200},
			{DataType: model.Prof, Name: func5, Sum: 120, Avg: 120 / 2230, Count: 2230},
		},
		MeasurementTime: &mtime2}
	mea3 := model.Measurement{
		Measurementid:   3,
		Score:           250,
		MeasurementTime: &mtime3,
		Memo:            "テスト3",
		TimeResults: []model.TimeResult{
			{DataType: model.Alp, Name: uri1, Sum: 305, Avg: 305 / 100, Count: 150},
			{DataType: model.Alp, Name: uri2, Sum: 10, Avg: 10 / 200, Count: 200},
			{DataType: model.Alp, Name: uri3, Sum: 105, Avg: 105 / 400, Count: 400},
			{DataType: model.Alp, Name: uri5, Sum: 120, Avg: 120 / 2230, Count: 2230},
			{DataType: model.Prof, Name: func1, Sum: 35, Avg: 35 / 100, Count: 150},
			{DataType: model.Prof, Name: func2, Sum: 10, Avg: 10 / 200, Count: 200},
			{DataType: model.Prof, Name: func3, Sum: 15, Avg: 15 / 400, Count: 400},
			{DataType: model.Prof, Name: func5, Sum: 220, Avg: 220 / 2230, Count: 2230},
		},
		NginxAccUrl: "https://storage.cloud.google.com/contestboard/alp.log"}
	mea4 := model.Measurement{
		Measurementid:   4,
		Score:           252,
		MeasurementTime: &mtime4,
		Memo:            "テスト4",
		TimeResults: []model.TimeResult{
			{DataType: model.Alp, Name: uri1, Sum: 305, Avg: 305 / 100, Count: 150},
			{DataType: model.Alp, Name: uri2, Sum: 10, Avg: 10 / 200, Count: 200},
			{DataType: model.Alp, Name: uri3, Sum: 105, Avg: 105 / 400, Count: 400},
			{DataType: model.Alp, Name: uri5, Sum: 120, Avg: 120 / 2230, Count: 2230},
			{DataType: model.Prof, Name: func1, Sum: 305, Avg: 305 / 100, Count: 150},
			{DataType: model.Prof, Name: func2, Sum: 10, Avg: 10 / 200, Count: 200},
			{DataType: model.Prof, Name: func3, Sum: 5, Avg: 5 / 400, Count: 400},
			{DataType: model.Prof, Name: func5, Sum: 120, Avg: 120 / 2230, Count: 2230},
		},
		NginxAccUrl: "https://storage.cloud.google.com/contestboard/alp.log"}
	mea5 := model.Measurement{
		Measurementid:   5,
		Score:           240,
		MeasurementTime: &mtime5,
		Memo:            "テスト5",
		TimeResults: []model.TimeResult{
			{DataType: model.Alp, Name: uri1, Sum: 304, Avg: 304 / 100, Count: 150},
			{DataType: model.Alp, Name: uri2, Sum: 10, Avg: 10 / 200, Count: 200},
			{DataType: model.Alp, Name: uri3, Sum: 105, Avg: 105 / 400, Count: 400},
			{DataType: model.Alp, Name: uri5, Sum: 110, Avg: 110 / 2230, Count: 2230},
			{DataType: model.Prof, Name: func1, Sum: 304, Avg: 304 / 100, Count: 150},
			{DataType: model.Prof, Name: func2, Sum: 110, Avg: 110 / 200, Count: 200},
			{DataType: model.Prof, Name: func3, Sum: 105, Avg: 105 / 400, Count: 400},
			{DataType: model.Prof, Name: func5, Sum: 10, Avg: 10 / 2230, Count: 2230},
		},
		NginxAccUrl: "https://storage.cloud.google.com/contestboard/alp.log"}
	addMeasurement(mea1)
	addMeasurement(mea2)
	addMeasurement(mea3)
	addMeasurement(mea4)
	addMeasurement(mea5)
}
