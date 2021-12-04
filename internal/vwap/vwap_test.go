package vwap_test

import (
	"sync"
	"testing"

	"github.com/letsila/vwap/internal/vwap"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Parallel()

	list, err := vwap.NewList([]vwap.DataPoint{}, 1)
	require.NoError(t, err)

	first := vwap.DataPoint{Price: decimal.NewFromInt(1), Volume: decimal.NewFromInt(1)}

	second := vwap.DataPoint{Price: decimal.NewFromInt(2), Volume: decimal.NewFromInt(2)}

	third := vwap.DataPoint{Price: decimal.NewFromInt(3), Volume: decimal.NewFromInt(3)}

	list.Push(first)
	require.Equal(t, 1, list.Len())
	require.Equal(t, first, list.DataPoints[0])

	list.Push(second)
	require.Equal(t, 1, list.Len())
	require.Equal(t, second, list.DataPoints[0])

	list.Push(third)
	require.Equal(t, 1, list.Len())
	require.Equal(t, third, list.DataPoints[0])
}

func TestListConcurrentPush(t *testing.T) {
	t.Parallel()

	list, err := vwap.NewList([]vwap.DataPoint{}, 2)
	require.NoError(t, err)

	first := vwap.DataPoint{Price: decimal.NewFromInt(1), Volume: decimal.NewFromInt(1)}

	second := vwap.DataPoint{Price: decimal.NewFromInt(2), Volume: decimal.NewFromInt(2)}

	third := vwap.DataPoint{Price: decimal.NewFromInt(3), Volume: decimal.NewFromInt(3)}

	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		list.Push(first)
		wg.Done()
	}()

	go func() {
		list.Push(second)
		wg.Done()
	}()

	go func() {
		list.Push(third)
		wg.Done()
	}()

	wg.Wait()

	require.Len(t, list.DataPoints, 2)
}

func TestVWAP(t *testing.T) {
	t.Parallel()

	tests := []struct {
		Name       string
		DataPoints []vwap.DataPoint
		WantVwap   map[string]decimal.Decimal
		MaxSize    uint
	}{
		{
			Name:       "EmptyDataPoints",
			DataPoints: []vwap.DataPoint{},
			WantVwap: map[string]decimal.Decimal{
				"BTC-USD": decimal.Zero,
				"ETH-USD": decimal.Zero,
			},
		},
		{
			Name: "FullDataPoints1",
			DataPoints: []vwap.DataPoint{
				{Price: decimal.NewFromInt(10), Volume: decimal.NewFromInt(10), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(10), Volume: decimal.NewFromInt(10), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(31), Volume: decimal.NewFromInt(30), ProductID: "ETH-USD"},
				{Price: decimal.NewFromInt(21), Volume: decimal.NewFromInt(20), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(41), Volume: decimal.NewFromInt(33), ProductID: "ETH-USD"},
			},
			MaxSize: 4,
			WantVwap: map[string]decimal.Decimal{
				"BTC-USD": decimal.RequireFromString("17.3333333333333333"),
				"ETH-USD": decimal.RequireFromString("36.2380952380952381"),
			},
		},
		{
			Name: "FullDataPoints2",
			DataPoints: []vwap.DataPoint{
				{Price: decimal.NewFromInt(10), Volume: decimal.RequireFromString("10.1"), ProductID: "BTC-USD"},
				{Price: decimal.NewFromInt(10), Volume: decimal.RequireFromString("10.1"), ProductID: "BTC-USD"},
			},
			WantVwap: map[string]decimal.Decimal{
				"BTC-USD": decimal.RequireFromString("10"),
			},
			MaxSize: 4,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			list, err := vwap.NewList([]vwap.DataPoint{}, tt.MaxSize)
			require.NoError(t, err)

			for _, d := range tt.DataPoints {
				list.Push(d)
			}

			for k := range tt.WantVwap {
				require.Equal(t, tt.WantVwap[k].String(), list.VWAP[k].String())
			}
		})
	}
}
