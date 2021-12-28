package main

import (
	"crypto/rand"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lemon-mint/authproxy/packet"
	"github.com/lemon-mint/vbox"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/curve25519"
)

// powered by Curve25519 ( https://cr.yp.to/ecdh.html )

var keyFileBin []byte

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var authNonceBufferPool = sync.Pool{
	New: func() interface{} {
		buffer := make([]byte, 128)
		return &buffer
	},
}

var x25519KeyBufferPool = sync.Pool{
	New: func() interface{} {
		buffer := make([]byte, 32)
		return &buffer
	},
}

var proxy4kBufPool = sync.Pool{
	New: func() interface{} {
		buffer := make([]byte, 4096)
		return &buffer
	},
}

func WsConnection(ws *websocket.Conn) {
	defer ws.Close()
	nonce := authNonceBufferPool.Get().(*[]byte)
	defer authNonceBufferPool.Put(nonce)
	_, err := io.ReadFull(rand.Reader, *nonce)
	if err != nil {
		return
	}
	ws.WriteMessage(websocket.BinaryMessage, *nonce)
	h, err := blake2b.New(32, nil)
	if err != nil {
		return
	}
	h.Write(*nonce)
	h.Write(keyFileBin)
	h.Write(*nonce)
	var EncryptionKey = h.Sum(nil)
	blackbox := vbox.NewBlackBox(EncryptionKey)

	privkey := x25519KeyBufferPool.Get().(*[]byte)
	defer func() {
		for i := range *privkey {
			(*privkey)[i] = 0
		}
		x25519KeyBufferPool.Put(privkey)
	}()

	pubkey := x25519KeyBufferPool.Get().(*[]byte)
	defer func() {
		for i := range *pubkey {
			(*pubkey)[i] = 0
		}
		x25519KeyBufferPool.Put(pubkey)
	}()

	_, err = io.ReadFull(rand.Reader, *privkey)
	if err != nil {
		return
	}

	curve25519.ScalarBaseMult((*[32]byte)(*pubkey), (*[32]byte)(*privkey))

	err = ws.SetWriteDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return
	}
	err = ws.WriteMessage(websocket.BinaryMessage, blackbox.Seal(*pubkey))
	if err != nil {
		return
	}
	err = ws.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil {
		return
	}
	_, message, err := ws.ReadMessage()
	if err != nil {
		return
	}
	clientPubKey, ok := blackbox.OpenOverWrite(message)
	if !ok {
		return
	}
	skey, err := curve25519.X25519(*privkey, clientPubKey)
	if err != nil {
		return
	}
	blackbox = vbox.NewBlackBox(skey)

	var upstream net.Conn
	var wsMutex sync.Mutex

	closeAll := func() {
		if upstream != nil {
			upstream.Close()
		}
		ws.Close()
	}

	go func() {
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				return
			}
			rawmsg, ok := blackbox.OpenOverWrite(message)
			if !ok {
				return
			}
			pkt := packet.Packet(rawmsg)
			if !pkt.Vstruct_Validate() {
				closeAll()
				return
			}
			switch pkt.Type() {
			case packet.PacketType_Connect:
				var connectionInfo packet.Connect = pkt.Payload()
				if !connectionInfo.Vstruct_Validate() {
					closeAll()
					return
				}
				if upstream != nil {
					closeAll()
					return
				}

				upstream, err = net.Dial("tcp", connectionInfo.Host())
				if err != nil {
					wsMutex.Lock()
					err = ws.WriteMessage(websocket.BinaryMessage, blackbox.Seal(packet.New_Packet(packet.PacketType_ConnectResponse, packet.New_ConnectResponse(false))))
					if err != nil {
						wsMutex.Unlock()
						closeAll()
						return
					}
					wsMutex.Unlock()
					continue
				}
				wsMutex.Lock()
				err = ws.WriteMessage(websocket.BinaryMessage, blackbox.Seal(packet.New_Packet(packet.PacketType_ConnectResponse, packet.New_ConnectResponse(true))))
				if err != nil {
					wsMutex.Unlock()
					closeAll()
					return
				}
				wsMutex.Unlock()
				go func() {
					defer closeAll()
					buffer := proxy4kBufPool.Get().(*[]byte)
					defer proxy4kBufPool.Put(buffer)
					for {
						n, err := upstream.Read(*buffer)
						if err != nil {
							closeAll()
							return
						}
						wsMutex.Lock()
						err = ws.WriteMessage(websocket.BinaryMessage, blackbox.Seal(packet.New_Packet(packet.PacketType_Downstream, packet.New_Stream(false, (*buffer)[:n]))))
						if err != nil {
							wsMutex.Unlock()
							closeAll()
							return
						}
						wsMutex.Unlock()
					}
				}()
			}
		}
	}()
}

func WsHandle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	go WsConnection(ws)
}

func main() {
	const keyFile = "key.bin"

	f, err := os.Open(keyFile)
	if err != nil {
		panic(err)
	}
	keyFileBin, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":8080", http.HandlerFunc(WsHandle))
}
