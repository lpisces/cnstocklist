package cnstocklist

import "testing"
import "sync"

func TestSSE(t *testing.T) {
	stockTypes := []string{"A", "B"}
	var wg sync.WaitGroup
	wg.Add(len(stockTypes))
	
	for _,st := range(stockTypes) {
		go func(tt string) {
			lst := SSE(tt)
			t.Logf("%d %s-type stocks in SSE.", len(lst), tt)
			t.Log("e.g.", lst[0])
			wg.Done()
		}(st)
	}

	wg.Wait()
}

func TestSZSE(t *testing.T) {
	stockTypes := []string{"A", "B"}
	var wg sync.WaitGroup
	wg.Add(len(stockTypes))
	
	for _,st := range(stockTypes) {
		go func(tt string) {
			lst := SZSE(tt)
			t.Logf("%d %s-type stocks in SZSE.", len(lst), tt)
			t.Log("e.g.", lst[0])
			wg.Done()
		}(st)
	}

	wg.Wait()
}

func TestAll(t *testing.T) {
	t.Logf("%d stocks in China Market", len(All()))
}

func TestA(t *testing.T) {
	t.Logf("%d stocks in China Market", len(A()))
}