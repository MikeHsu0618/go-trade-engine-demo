package util

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go_trade_engine_demo/internal/trade/constants"
	"math/rand"
	"time"
)

func FormatDecimal2String(d decimal.Decimal, digit int) string {
	f, _ := d.Float64()
	format := "%." + fmt.Sprintf("%d", digit) + "f"
	return fmt.Sprintf(format, f)
}

func quickSort(nums []string, asc_desc string) []string {
	if len(nums) <= 1 {
		return nums
	}

	spilt := nums[0]
	left := []string{}
	right := []string{}
	mid := []string{}

	for _, v := range nums {
		vv, _ := decimal.NewFromString(v)
		sp, _ := decimal.NewFromString(spilt)
		if vv.Cmp(sp) == -1 {
			left = append(left, v)
		} else if vv.Cmp(sp) == 1 {
			right = append(right, v)
		} else {
			mid = append(mid, v)
		}
	}

	left = quickSort(left, asc_desc)
	right = quickSort(right, asc_desc)

	if asc_desc == "asc" {
		return append(append(left, mid...), right...)
	}
	return append(append(right, mid...), left...)
}

func SortMap2Slice(m map[string]string, askBid constants.OrderSide) [][2]string {
	var res [][2]string
	var keys []string
	for k, _ := range m {
		keys = append(keys, k)
	}

	if askBid == constants.OrderSideSell {
		keys = quickSort(keys, "asc")
	} else {
		keys = quickSort(keys, "desc")
	}

	for _, k := range keys {
		res = append(res, [2]string{k, m[k]})
	}
	return res
}

func String2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}

func RandDecimal(min, max int64) decimal.Decimal {
	rand.Seed(time.Now().UnixNano())

	d := decimal.New(rand.Int63n(max-min)+min, 0)
	return d
}
