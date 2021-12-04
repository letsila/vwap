package websocket

import "context"

// Response is the response received after a request submission.
type Response struct {
	Type      string `json:"type"`
	Size      string `json:"size"`
	Price     string `json:"price"`
	ProductID string `json:"product_id"`
}

type Client interface {
	Subscribe(ctx context.Context, tradingPairs []string, receiver chan Response) error
}
