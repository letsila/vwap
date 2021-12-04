package internal

import (
	"context"
	"fmt"

	"github.com/letsila/vwap/internal/vwap"
	"github.com/letsila/vwap/internal/websocket"
	"github.com/shopspring/decimal"
	"golang.org/x/xerrors"
)

// Service is our main service.
type Service struct {
	wsClient     websocket.Client
	tradingPairs []string
	list         *vwap.List
}

// NewService returns a new service.
func NewService(ws websocket.Client, tradingPairs []string, list *vwap.List) *Service {
	return &Service{
		wsClient:     ws,
		tradingPairs: tradingPairs,
		list:         list,
	}
}

func (s *Service) Run(ctx context.Context) error {
	receiver := make(chan websocket.Response)

	err := s.wsClient.Subscribe(ctx, s.tradingPairs, receiver)
	if err != nil {
		return xerrors.Errorf("service subscription err: %w", err)
	}

	for data := range receiver {
		if data.Price == "" {
			continue
		}

		decimalPrice, err := decimal.NewFromString(data.Price)
		if err != nil {
			return xerrors.Errorf("decimalPrice %s: %w", data.Price, err)
		}

		decimalSize, err := decimal.NewFromString(data.Size)
		if err != nil {
			return xerrors.Errorf("decimalSize %s: %w", data.Size, err)
		}

		s.list.Push(vwap.DataPoint{
			Price:     decimalPrice,
			Volume:    decimalSize,
			ProductID: data.ProductID,
		})

		// Print to sdout.
		fmt.Println(s.list.VWAP)
	}

	return nil
}
