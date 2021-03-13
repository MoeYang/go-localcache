package common

import (
	"strconv"
	"testing"
)

func TestGetShardIndex(t *testing.T) {
	arr := make([]int, 4)
	for i := 0; i < 10000; i++ {
		idx := GetShardIndex(strconv.Itoa(i), 4)
		arr[idx]++
	}
	for _, cnt := range arr {
		if cnt > 2600 {
			t.Error("Err_TestGetShardIndex")
		}
	}
}
