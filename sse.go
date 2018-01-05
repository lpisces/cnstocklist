package cnstocklist

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"bytes"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

const (
	// SseDownloadLink 上交所股票列表下载链接
	StockListDownloadLink = "http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=%s"
)

// StockTypeMap AB股对应类型参数
var StockTypeMap = map[string]string{"A": "1", "B": "2"}

// SSE 获取上交所股票列表
func SSE(t string) (stocks []Stock) {

	// 下载xls文件
	url := fmt.Sprintf(StockListDownloadLink, StockTypeMap[t])
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "http://www.sse.com.cn/assortment/stock/list/share/")
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

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	s, err := GbkToUtf8(bodyBytes)

	r := csv.NewReader(strings.NewReader(string(s)))
	ss, _ := r.ReadAll()
	cnt := 0
	for _, i := range ss {
		for _, j := range i {
			if cnt == 0 {
				cnt++
				continue
			}
			row := string(j)
			line := strings.Split(row, "\t")

			var stock Stock
			if t == "A"{
				stock = Stock{Symbol:string(line[0]),ShortName: strings.Trim(string(line[1]), " "), Market: "sh", TimeToMarket: string(line[4]), T: t}
			} else {
				stock = Stock{Symbol:strings.Trim(string(line[2]), " "),ShortName: strings.Trim(string(line[3]), " "), Market: "sh", TimeToMarket: string(line[4]), T: t}
			}
			
			//t.Log(stock)
			stocks = append(stocks, stock)
		}
	}
	return
}

// GbkToUtf8 字符串转换
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// Utf8ToGbk 字符串转换
func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}