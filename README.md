# scattare

Capture Twitch chat logs into JSON, CSV with webhook supports.

## Build

```
$ git clone https://github.com/So-chiru/scattare.git
$ cd scattare
$ go build
```

## Usage

```
Usage of ./scattare:
  -c string
        alias of channel, twitch channel id that wishes to connect with
  -channel string
        twitch channel id that wishes to connect with
  -d    alias of debug, enable debug mode
  -debug
        enable debug mode
  -e string
        alias of endpoint, transport endpoint (accept https-http, leave empty to disable)
  -endpoint string
        transport endpoint (accept https-http, leave empty to disable)
  -h string
        alias of headers, transport headers (serialized JSON)
  -headers string
        transport headers (serialized JSON)
  -i int
        alias of interval, collect interval in miliseconds (default 3000)
  -interval int
        collect interval in miliseconds (default 3000)
  -o string
        alias of output, output file (.csv, .json) (default "data.json")
  -output string
        output file (.csv, .json) (default "data.json")
```

An example to capture chat logs of a channel `woowakgood`:

```sh
$ scattare -c woowakgood
```

An example to capture chat logs of a channel `woowakgood` and transport over http with headers:

```sh
$ scattare -c woowakgood -e "http://localhost:3000/webhooks" -h "{\"Authorization\": \"Bearer token\"}"
```
