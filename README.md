# go-skyremote

*IMPORTANT - This is not an official library by BSkyB* 

Go library for controlling Sky+HD and SkyQ boxes over IP

Largely adapted from https://github.com/dalhundal/sky-remote

## Usage

The bundled CLI is an example of how to use the library in your own applications. The CLI itself is very simple and let's you change the channel.

```
$ go run main.go --help
      --channel string   3 digit channel number (default "100")
      --ip string        IP of remote box
      --port int         Port on remote box (default 49160)
```
