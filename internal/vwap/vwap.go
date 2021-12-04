package vwap

import (
	"sync"

	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

const defaultMaxSize = 200

// DataPoint represents a single data point from coinbase.
type DataPoint struct {
	Price     decimal.Decimal
	Volume    decimal.Decimal
	ProductID string
}

// List represents a queue of DataPoints.
type List struct {
	mu                sync.Mutex
	DataPoints        []DataPoint
	SumVolumeWeighted map[string]decimal.Decimal
	SumVolume         map[string]decimal.Decimal
	VWAP              map[string]decimal.Decimal

	MaxSize uint
}

// NewList creates a new queue.
func NewList(dataPoint []DataPoint, maxSize uint) (List, error) {
	if maxSize == 0 {
		maxSize = defaultMaxSize
	}

	if len(dataPoint) > int(maxSize) {
		return List{}, xerrors.New("initial datapoints exceeds maxSize")
	}

	return List{
		DataPoints:        dataPoint,
		MaxSize:           maxSize,
		SumVolumeWeighted: make(map[string]decimal.Decimal),
		SumVolume:         make(map[string]decimal.Decimal),
		VWAP:              make(map[string]decimal.Decimal),
	}, nil
}

// Len returns the length of the queue.
func (q *List) Len() int {
	return len(q.DataPoints)
}

// Push pushes an element onto the queue, drops the first one when MaxSize is reached.
func (q *List) Push(d DataPoint) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.DataPoints) == int(q.MaxSize) {
		d := q.DataPoints[0]
		q.DataPoints = q.DataPoints[1:]

		// Substract the datapoint values from the VWAP calculation.
		q.SumVolumeWeighted[d.ProductID] = q.SumVolumeWeighted[d.ProductID].Sub(d.Price.Mul(d.Volume))
		q.SumVolume[d.ProductID] = q.SumVolume[d.ProductID].Sub(d.Volume)
		if !q.SumVolume[d.ProductID].IsZero() {
			q.VWAP[d.ProductID] = q.SumVolumeWeighted[d.ProductID].Div(q.SumVolume[d.ProductID])
		}
	}

	if _, ok := q.VWAP[d.ProductID]; ok {
		q.SumVolumeWeighted[d.ProductID] = q.SumVolumeWeighted[d.ProductID].Add(d.Price.Mul(d.Volume))
		q.SumVolume[d.ProductID] = q.SumVolume[d.ProductID].Add(d.Volume)
		q.VWAP[d.ProductID] = q.SumVolumeWeighted[d.ProductID].Div(q.SumVolume[d.ProductID])
	} else {
		initialVW := d.Price.Mul(d.Volume)

		q.SumVolumeWeighted[d.ProductID] = initialVW
		q.SumVolume[d.ProductID] = d.Volume
		q.VWAP[d.ProductID] = initialVW.Div(d.Volume)
	}

	q.DataPoints = append(q.DataPoints, d)
}
