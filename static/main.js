// Goから渡されたscoreplotdataをChart.jsが読めるchartjsdatasetに詰め直す
var scoreChartData = JSON.parse(document.querySelector('#ScoreChartDataJSONStr').value)
var alpChartData = JSON.parse(document.querySelector('#AlpChartDataJSONStr').value)
var profChartData = JSON.parse(document.querySelector('#ProfChartDataJSONStr').value)
var ptqChartData = JSON.parse(document.querySelector('#PtqChartDataJSONStr').value)

// https://developer.mozilla.org/ja/docs/Web/CSS/CSS_Colors/Color_picker_tool
var scoreColor = 'rgba(255, 255, 255, 0)'
var alpColor = 'rgba(244, 4, 8, 0.01)'
var profColor = 'rgba(4, 236, 244, 0.01)'
var ptqColor = 'rgba(127, 191, 63, 0.01)'

makeNewChart(scoreChart, scoreChartData, 'スコア', false, '', scoreColor)
makeNewChart(scoreChart2, scoreChartData, 'スコア', false, '', scoreColor)
makeNewChart(alpChart, alpChartData, 'URIごと累計処理時間[s]', false, 'right', alpColor)
makeNewChart(alpChart2, alpChartData, 'URIごと累計処理時間[s]', true, 'right', alpColor)
makeNewChart(profChart, profChartData, '関数プロファイル累計処理時間[s]', false, 'false', profColor)
makeNewChart(profChart2, profChartData, '関数プロファイル累計処理時間[s]', true, 'right', profColor)
makeNewChart(ptqChart, ptqChartData, 'SQLごと累計処理時間[s]', false, 'false', ptqColor)
makeNewChart(ptqChart2, ptqChartData, 'SQLごと累計処理時間[s]', true, 'right', ptqColor)

function golangJSONToChartjsData(chartData) {
    var chartjsdataset = []
    if (!(chartData.projects === null)) {
        for (let i = 0; i < chartData.projects.length; i++) {
            var project = chartData.projects[i]
            chartjsdataset.push({
                label: project.projectname,
                data: project.plotpoints,
                "fill": false,
                "lineTension": 0
            })
        }
    } else {
        console.log('chartData.projectsはnullです');
    }
    return chartjsdataset
}

function makeNewChart(chart, chartData, title, legendDisplay, legendPos, bgcolor) {
    Chart.pluginService.register({
        beforeDraw: function (c) {
            if (c.config.options.chartArea && c.config.options.chartArea.backgroundColor) {
                var ctx = c.chart.ctx;
                var chartArea = c.chartArea;
                ctx.save();
                ctx.fillStyle = c.config.options.chartArea.backgroundColor;
                ctx.fillRect(chartArea.left, chartArea.top, chartArea.right - chartArea.left, chartArea.bottom - chartArea.top);
                ctx.restore();
            }
        }
    });
    var chartjsdataset = golangJSONToChartjsData(chartData)
    new Chart(chart, {
        "type": "line",
        "data": {
            datasets: chartjsdataset
        },
        "options": {
            title: {
                display: true,
                text: title,
                fontSize: 15
            },
            scales: {
                xAxes: [{
                    type: 'time',
                }],
                yAxes: [{
                    ticks: {
                        suggestedMin: 0,
                    }
                }]
            },
            chartArea: {
                backgroundColor: bgcolor
            },
            plugins: {
                colorschemes: {
                    scheme: 'brewer.Paired12'
                }
            },
            legend: {
                display: legendDisplay,
                position: legendPos,
            }
        }
    })
}




function check() {
    if (window.confirm('計測データを削除しますか？')) {
        return true;
    }
    else {
        window.alert('削除を中止しました');
        return false;
    }
}
