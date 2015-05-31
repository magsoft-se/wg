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

var ggsock *websocket.Conn = nil
var ggimg *image.RGBA = nil
var gmx, gmy int
var gkeys KeyMap
var gmlbtn, gmrbtn bool

var gloopfunc LoopFunc

func GetImage() *image.RGBA {
	return ggimg
}

func GetMX() int {
	return gmx
}

func GetMY() int {
	return gmy
}

func GetMLBtn() bool {
	return gmlbtn
}

func GetKeys() *KeyMap {
	return &gkeys
}

func GetKey(key int) bool {
	return gkeys[key]
}

func ClearImage(col color.Color) {
	draw.Draw(ggimg, ggimg.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
}

func receiver(ws *websocket.Conn) {
	var recv_data string
	var err error

	ggsock = ws

	for {
		if err = websocket.Message.Receive(ws, &recv_data); err != nil {
			fmt.Println("Can't receive" + err.Error())
			break
		}

		if strings.Index(recv_data, "mx") == 0 {
			fmt.Sscanf(recv_data, "mx%dmy%d", &gmx, &gmy)

		} else if strings.Index(recv_data, "mlbtn") == 0 {
			var btn int
			fmt.Sscanf(recv_data, "mlbtn%d", &btn)
			gmlbtn = (btn == 1)

		} else if strings.Index(recv_data, "kdn") == 0 {
			var key int
			fmt.Sscanf(recv_data, "kdn%d", &key)

			if key >= 0 && key < MAX_KEY {
				gkeys[key] = true
			}

		} else if strings.Index(recv_data, "kup") == 0 {
			var key int
			fmt.Sscanf(recv_data, "kup%d", &key)

			if key >= 0 && key < MAX_KEY {
				gkeys[key] = false
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
		if ggsock != nil {
			gloopfunc()

			buf.Reset()
			jpeg.Encode(buf, ggimg, nil)
			enc := base64.StdEncoding.EncodeToString(buf.Bytes())

			if err = websocket.Message.Send(ggsock, "data:image/jpeg;base64,"+enc); err != nil {
				fmt.Println("Can't send " + err.Error())
				break
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func Start(width, height, port int, loopfunc func()) {
	gloopfunc = loopfunc

	http.Handle("/ws", websocket.Handler(receiver))
	http.HandleFunc("/", serveStatic)

	ggimg = image.NewRGBA(image.Rect(0, 0, width, height))

	fmt.Println("Handle...")
	go sender()

	hosturl := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(hosturl, nil); err != nil {
		fmt.Println("Serving failed!!!")
	}
}
