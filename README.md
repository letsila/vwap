# VWAP calculator
A realtime VWAP calculator of crypto currencies. It uses coinbase as its default provider for real time data over websocket.

## Design
The service in `./internal/service.go` is composed of two main components:
* A websocket client that pulls data off an exchange.
  * The default choice is coinbase.
  * Any exchange can be used as long as it implements the client interface defined in the websocket package.
* A list of data points defined in the VWAP package.
  * The VWAP calculation is performed each time a data point is pushed to the list and saved in a hash map for each trading pairs.
  * We don't loop over the datapoints so the VWAP calculation is done in constant time, O(1).

## Configuration
The following flags are available while running the project through CLI using the binary.
* `trading-pairs`: a comma separated strings of crypto currencies pairs, default is set to `BTC-USD,ETH-USD,ETH-BTC`
* `ws-url`: the URL of the websocket server to use, default is coinbase websocket URL.
* `window-size`: the sliding window used for the VWAP calculation, default is set to **200**.

## Decimal
For precision sake we used https://github.com/shopspring/decimal for all calculation.

## Run it
First, make sure that you have go version 0.17 installed on your machine. Then ...
```
make run
```
or 
```
make build
``` 
then ...
```
./vwap -ws-url "<coinbase_valid_ws_url>" -trading-pairs "<trading_pairs>" -window-size <window_size>
```

## Tests
* Runs all the tests.
```
make test
``` 
* Runs the unit tests.
```
make test-unit
```
* Runs the integration test.
```
make test-intergration
``` 