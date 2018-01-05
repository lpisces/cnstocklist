package cnstocklist

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/tealeg/xlsx"
)

const (
	// StockListDownloadLink A+B股列表下载链接
	SZSEStockListDownloadLink = "http://www.szse.cn/szseWeb/ShowReport.szse?SHOWTYPE=xlsx&CATALOGID=1110&tab1PAGENO=1&ENCODE=1&TABKEY=%s"
)

// StockTypeMap AB股对应类型参数
var SZSEStockTypeMap = map[string]string{"A": "tab1", "B": "tab3"}

// GetStockList 获取股票列表
func SZSE(t string) (stocks []Stock) {
	// 下载文件
	url := fmt.Sprintf(SZSEStockListDownloadLink, SZSEStockTypeMap[t])
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "http://www.szse.cn/main/marketdata/jypz/colist/")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("Status Code is %d", resp.StatusCode))
	}

	tmpfile, err := ioutil.TempFile("", "szse")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	io.Copy(tmpfile, reader)

	//解析文件
	xlFile, err := xlsx.OpenFile(tmpfile.Name())
	if err != nil {
		panic(err)
	}

	cnt := 0
	for _, sheet := range xlFile.Sheets {
		for _, row := range sheet.Rows {
			if cnt == 0 {
				cnt++
				continue
			}
			var stock Stock
			if t == "A" {
				if row.Cells[8].String() == "0" {
					continue
				}
				stock = Stock{Symbol: row.Cells[0].String(), ShortName: row.Cells[1].String(), Market: "sz", TimeToMarket: row.Cells[7].String(), T: t}
			} else {
				stock = Stock{Symbol: row.Cells[10].String(), ShortName: row.Cells[11].String(), Market: "sz", TimeToMarket: row.Cells[12].String(), T: t}
			}
			stocks = append(stocks, stock)

		}
	}
	return
}
