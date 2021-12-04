package coinbase_test

import (
	"context"
	"testing"

	"github.com/letsila/vwap/internal/websocket"
	"github.com/letsila/vwap/internal/websocket/exchanges/coinbase"
	"github.com/stretchr/testify/require"
)

func TestWebNew(t *testing.T) {
	t.Parallel()

	_, err := coinbase.NewClient(coinbase.DefaultURL)
	require.NoError(t, err)
}

func TestWebsocketSubscribe(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		Name         string
		TradingPairs []string
		WantErr      bool
	}{
		{
			Name:         "ValidTradingPairs",
			TradingPairs: []string{"BTC-USD"},
			WantErr:      false,
		},
		{
			Name:         "InvalidPairs",
			TradingPairs: []string{"xxx-USD"},
			WantErr:      true,
		},
	}

	receiver := make(chan websocket.Response)

	for _, tt := range tests {
		tt := tt

		t.Run(tt.Name, func(t *testing.T) {
			t.Parallel()

			ws, err := coinbase.NewClient(coinbase.DefaultURL)
			require.NoError(t, err)

			err = ws.Subscribe(ctx, tt.TradingPairs, receiver)
			if tt.WantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var counter int

				// Check the first two messages.
				for m := range receiver {
					if counter >= 2 {
						break
					}

					if m.Type == "last_match" {
						require.Equal(t, "BTC-USD", m.ProductID)
					}

					counter++
				}
			}
		})
	}
}
