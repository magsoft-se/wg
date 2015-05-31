# wg

### About
A simple "Web Graphics" interface to use for golang programs.
The index.html is hosted with the built in web server, from there a
bi-directional websocket is started. From go backend to browser is each frame passed
as base64 encoded jpeg and immediately inline'd in the DOM (drawn). From browser to go
are mouse and key events passed for the backend code to react on.

Use as a base for a go game or something else with real time graphics.

### Install
```
go get github.com/magsoft-se/wg
```
For an example, `go run main.go` and point a browser to localhost.
You should se a blue square controllable with the arrow keys, and be able to draw a dotted line with the mouse. (you may have to focus click the graphics area)

### Basic Usage
```
wg.Start(WIDTH, HEIGHT, PORT, GameLoop)
```
See more in example/main.go

### Disclaimer
Please consider that this is one of my first golang applications, so there are probably better ways of doing things.
Also this is pretty cpu intensive and of course far from optimal if compared to more direct/native graphics.
In any case I think this will provide an easy to use graphics interface that I myself missed when I started with go.
