package cnstocklist

import "sync"

const userAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.91 Safari/537.36"

// Stock 类型
type Stock struct {
	Symbol       string //证券代码
	ShortName    string //公司简称
	Market       string //市场
	TimeToMarket string //上市时间
	T            string // A B 股
}

func _merge(types []string, exchanges []func(string) []Stock)(stocks []Stock) {
	var wg sync.WaitGroup

	//types := []string{"A", "B"}
	//methods := []func(string) []Stock{SSE, SZSE}
	methods := exchanges

	buffSize := len(types) * len(methods)
	c := make(chan []Stock, buffSize)

	wg.Add(buffSize)

	for _, tt := range types {
		for _, f := range methods {
			go func(g func(t string) []Stock, tt string, c chan []Stock) {
				defer wg.Done()
				s := g(tt)
				c <- s
			}(f, tt, c)
		}
	}

	wg.Wait()
	close(c)

	for s := range c {
		stocks = append(stocks, s...)
	}
	return
}

// All 获取全部股票列表
func All() (stocks []Stock) {
	types := []string{"A", "B"}
	exchanges := []func(string) []Stock{SSE, SZSE}
	return _merge(types, exchanges)
}

// A 获取全部A股列表
func A()(stocks []Stock) {
	types := []string{"A"}
	exchanges := []func(string) []Stock{SSE, SZSE}
	return _merge(types, exchanges)
}

// B 获取全部B股列表
func B()(stocks []Stock) {
	types := []string{"B"}
	exchanges := []func(string) []Stock{SSE, SZSE}
	return _merge(types, exchanges)
}