package wg

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"net/http"
	"strings"
	"time"
)

const MAX_KEY = 127

type LoopFunc func()
type KeyMap [MAX_KEY]bool

var wgsock *websocket.Conn = nil
var wgimg *image.RGBA = nil
var wgmx, wgmy int
var wgkeys KeyMap
var wgmlbtn, wgmrbtn bool

var wgLoopfunc LoopFunc

func GetImage() *image.RGBA {
	return wgimg
}

func GetMX() int {
	return wgmx
}

func GetMY() int {
	return wgmy
}

func GetMLBtn() bool {
	return wgmlbtn
}

func GetKeys() *KeyMap {
	return &wgkeys
}

func GetKey(key int) bool {
	return wgkeys[key]
}

func ClearImage(col color.Color) {
	draw.Draw(wgimg, wgimg.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
}

func receiver(ws *websocket.Conn) {
	var recv_data string
	var err error

	wgsock = ws

	for {
		if err = websocket.Message.Receive(ws, &recv_data); err != nil {
			fmt.Println("Can't receive" + err.Error())
			break
		}

		if strings.Index(recv_data, "mx") == 0 {
			fmt.Sscanf(recv_data, "mx%dmy%d", &wgmx, &wgmy)

		} else if strings.Index(recv_data, "mlbtn") == 0 {
			var btn int
			fmt.Sscanf(recv_data, "mlbtn%d", &btn)
			wgmlbtn = (btn == 1)

		} else if strings.Index(recv_data, "kdn") == 0 {
			var key int
			fmt.Sscanf(recv_data, "kdn%d", &key)

			if key >= 0 && key < MAX_KEY {
				wgkeys[key] = true
			}

		} else if strings.Index(recv_data, "kup") == 0 {
			var key int
			fmt.Sscanf(recv_data, "kup%d", &key)

			if key >= 0 && key < MAX_KEY {
				wgkeys[key] = false
			}
		}
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func sender() {
	var err error
	buf := new(bytes.Buffer)

	for {
		if wgsock != nil {
			wgLoopfunc()

			buf.Reset()
			jpeg.Encode(buf, wgimg, nil)
			enc := base64.StdEncoding.EncodeToString(buf.Bytes())

			if err = websocket.Message.Send(wgsock, "data:image/jpeg;base64,"+enc); err != nil {
				fmt.Println("Can't send " + err.Error())
				break
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func Start(width, height, port int, loopfunc func()) {
	wgLoopfunc = loopfunc

	http.Handle("/ws", websocket.Handler(receiver))
	http.HandleFunc("/", serveStatic)

	wgimg = image.NewRGBA(image.Rect(0, 0, width, height))

	fmt.Println("Handle...")
	go sender()

	hosturl := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(hosturl, nil); err != nil {
		fmt.Println("Serving failed!!!")
	}
}
