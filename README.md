# scattare
Capture Twitch chats into logs


## Download and build

```
$ git clone https://github.com/So-chiru/scattare.git
$ cd scattare
$ go build
```


## Usage

```
Usage of ./scattare:
  -channel string
        channel to connect to
  -debug
        enable debug mode
  -interval int
        collect interval in miliseconds (default 3000)
  -output string
        output file (.csv, .json supported) (default "data.json")
 ```
 
 For example, if you want to capture a channel chats with streamer `woowakgood`, use:
 
 ```sh
 $ scattare -channel="woowakgood"
 ```
