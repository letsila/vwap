package coinbase

import (
	"context"
	"encoding/json"
	"log"

	"github.com/letsila/vwap/internal/websocket"
	ws "golang.org/x/net/websocket"
	"golang.org/x/xerrors"
)

type client struct {
	conn *ws.Conn
}

// NewClient returns a new websocket client.
func NewClient(url string) (websocket.Client, error) {
	conn, err := ws.Dial(url, "", "http://localhost/")
	if err != nil {
		return nil, err
	}

	log.Printf("websocket connected to: %s", url)

	return &client{
		conn: conn,
	}, nil
}

// wesocketResponse converts the coinbase response into a websocket response.
func wesocketResponse(res Response) websocket.Response {
	return websocket.Response{
		Type:      res.Type,
		Size:      res.Size,
		Price:     res.Price,
		ProductID: res.ProductID,
	}
}

// Subscribe subscribes to the websocket.
func (c *client) Subscribe(ctx context.Context, tradingPairs []string, receiver chan websocket.Response) error {
	if len(tradingPairs) == 0 {
		return xerrors.New("tradingPairs must be provided")
	}

	subscription := Request{
		Type:       RequestTypeSubscribe,
		ProductIDs: tradingPairs,
		Channels: []Channel{
			{Name: "matches"},
		},
	}

	payload, err := json.Marshal(subscription)
	if err != nil {
		return xerrors.Errorf("failed to marshal subscription: %w", err)
	}

	err = ws.Message.Send(c.conn, payload)
	if err != nil {
		return xerrors.Errorf("failed to send subscription: %w", err)
	}

	var subscriptionResponse Response

	err = ws.JSON.Receive(c.conn, &subscriptionResponse)
	if err != nil {
		return xerrors.Errorf("failed to receive subscription response: %w", err)
	}

	if subscriptionResponse.Type == "error" {
		return xerrors.Errorf("failed to subscribe: %s", subscriptionResponse.Message)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := c.conn.Close()
				if err != nil {
					log.Printf("failed closing ws connection: %s", err)
				}
			default:
				var message Response

				err := ws.JSON.Receive(c.conn, &message)
				if err != nil {
					log.Printf("failed receiving message: %s", err)

					break
				}

				receiver <- wesocketResponse(message)
			}
		}
	}()

	return nil
}
