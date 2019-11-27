package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetFlyPaysB(t *testing.T) {
	var filter Filter
	c := make(chan *Response)
	var flyPay *FlyPayB
	go flyPay.GetFlyPays(filter, c)
	res := <-c
	assert.Equal(t, nil, res.Err, "error should be nil")
	assert.Greater(t, len(res.List.FlyPay), 0, "flypays list length should be greater than 0")
}

func TestIsAMatchForStatusCodeB(t *testing.T) {
	var filter Filter
	var flyPay FlyPayB
	var isMatch bool
	filter.BStatusCode = 1
	flyPay.StatusCode = 1
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, true, isMatch, "it should match")

	filter.BStatusCode = 2
	flyPay.StatusCode = 2
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, true, isMatch, "it should match")

	filter.BStatusCode = 3
	flyPay.StatusCode = 3
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, true, isMatch, "it should match")

	filter.BStatusCode = 2
	flyPay.StatusCode = 1
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, false, isMatch, "it should not match")
}

func TestIsAMatchForCurrencyB(t *testing.T) {
	var filter Filter
	var flyPay FlyPayB
	var isMatch bool
	filter.Currency = "EGP"
	flyPay.Currency = "EGP"
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, true, isMatch, "it should match")

	filter.Currency = "EGP"
	flyPay.Currency = "AUD"
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, false, isMatch, "it should not match")
}

func TestIsAMatchForAmountB(t *testing.T) {
	var filter Filter
	var flyPay FlyPayB
	var isMatch bool
	filter.AmountMin = 50
	flyPay.Amount = 1000
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, true, isMatch, "it should match")

	filter.AmountMin = 50
	flyPay.Amount = 1
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, false, isMatch, "it should not match")

	filter.AmountMin = 0
	filter.AmountMax = 50
	flyPay.Amount = 1
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, true, isMatch, "it should match")

	filter.AmountMax = 50
	flyPay.Amount = 1000
	isMatch = flyPay.IsAMatch(&filter)
	assert.Equal(t, false, isMatch, "it should not match")
}

/**************************BenchMark****************************************/

func BenchmarkGetFlyPaysB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var filter Filter
		c := make(chan *Response)
		var flyPay *FlyPayB
		flyPay.GetFlyPays(filter, c)
	}

}

func BenchmarkIsAMatchForStatusCodeB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var filter Filter
		var flyPay FlyPayB
		filter.BStatusCode = 1
		flyPay.StatusCode = 1
		flyPay.IsAMatch(&filter)
	}

}
