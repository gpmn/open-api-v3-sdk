package okex

/*
 OKEX websocket API agent
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import (
	"bytes"
	"compress/flate"
	"io/ioutil"

	"github.com/gorilla/websocket"

	"log"
	"sync"
	"time"
)

const (
	maxPongInterval = 35 * time.Second
)

type OKWSAgent struct {
	lastPongTm time.Time
	startHook  func() error
	baseUrl    string
	config     *Config
	conn       *websocket.Conn
	connLock   sync.Mutex

	wsEvtCh chan interface{}
	wsErrCh chan interface{}
	wsTbCh  chan interface{}

	subMap         map[string]ReceivedDataCallback
	activeChannels map[string]bool
	hotDepthsMap   map[string]*WSHotDepths

	processMut sync.Mutex
}

func (a *OKWSAgent) Start(config *Config, startHook func() error) error { // 没有restart、stop的必要，删除stop和finize
	a.baseUrl = config.WSEndpoint + "ws/v3?compress=true"
	log.Printf("Connecting to %s", a.baseUrl)
	c, _, err := websocket.DefaultDialer.Dial(a.baseUrl, nil)

	if err != nil {
		log.Fatalf("dial:%+v", err)
		return err
	}
	log.Printf("Connected to %s", a.baseUrl)
	a.lastPongTm = time.Now().Add(maxPongInterval)
	a.conn = c
	a.config = config
	a.startHook = startHook

	a.wsEvtCh = make(chan interface{})
	a.wsErrCh = make(chan interface{})
	a.wsTbCh = make(chan interface{})
	a.activeChannels = make(map[string]bool)
	a.subMap = make(map[string]ReceivedDataCallback)
	a.hotDepthsMap = make(map[string]*WSHotDepths)

	go a.work()
	go a.receive()
	if startHook != nil {
		return startHook()
	}
	return nil
}

func (a *OKWSAgent) Subscribe(channel, filter string, cb ReceivedDataCallback) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	st := SubscriptionTopic{channel, filter}
	bo, err := subscribeOp([]*SubscriptionTopic{&st})
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	log.Printf("Send Msg: %s", msg)
	a.connLock.Lock()
	err = a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	a.connLock.Unlock()
	if err != nil {
		return err
	}

	a.activeChannels[st.channel] = false
	a.subMap[st.channel] = cb

	return nil
}

func (a *OKWSAgent) SubscribeEx(channel string, filters []string, cb ReceivedDataCallback) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	var sts []*SubscriptionTopic
	for _, filter := range filters {
		sts = append(sts, &SubscriptionTopic{
			channel: channel,
			filter:  filter,
		})
	}

	bo, err := subscribeOp(sts)
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	log.Printf("Send Msg: %s", msg)
	a.connLock.Lock()
	err = a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	a.connLock.Unlock()
	if err != nil {
		return err
	}

	a.activeChannels[channel] = false
	a.subMap[channel] = cb
	return nil
}

func (a *OKWSAgent) UnSubscribe(channel, filter string) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	st := SubscriptionTopic{channel, filter}
	bo, err := unsubscribeOp([]*SubscriptionTopic{&st})
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	log.Printf("Send Msg: %s", msg)
	a.connLock.Lock()
	err = a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	a.connLock.Unlock()
	if err != nil {
		return err
	}

	a.subMap[channel] = nil
	a.activeChannels[channel] = false

	return nil
}

func (a *OKWSAgent) Login(apiKey, passphrase string) error {
	timestamp := EpochTime()
	preHash := PreHashString(timestamp, GET, "/users/self/verify", "")
	sign, err := HmacSha256Base64Signer(preHash, a.config.SecretKey)
	if err != nil {
		return err
	}
	op, err := loginOp(apiKey, passphrase, timestamp, sign)
	data, err := Struct2JsonString(op)
	a.connLock.Lock()
	err = a.conn.WriteMessage(websocket.TextMessage, []byte(data))
	a.connLock.Unlock()
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 100)

	return nil
}

func (a *OKWSAgent) keepalive() {
	a.ping()
}

func (a *OKWSAgent) ping() {
	msg := "ping"
	a.connLock.Lock()
	a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	a.connLock.Unlock()
}

func (a *OKWSAgent) GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func (a *OKWSAgent) handleErrResponse(r interface{}) error {
	log.Printf("handleErrResponse %+v \n", r)
	return nil
}

func (a *OKWSAgent) handleEventResponse(r interface{}) error {
	er := r.(*WSEventResponse)
	a.activeChannels[er.Channel] = (er.Event == CHNL_EVENT_SUBSCRIBE)
	return nil
}

func (a *OKWSAgent) handleTableResponse(r interface{}) error {
	tb := ""
	switch r.(type) {
	case *WSTableResponse:
		tb = r.(*WSTableResponse).Table
	case *WSDepthTableResponse:
		tb = r.(*WSDepthTableResponse).Table
	}

	cb := a.subMap[tb]
	if cb != nil {
		// for i := 0; i < len(cbs); i++ {
		// cb := cbs[i]
		if err := cb(r); err != nil {
			return err
		}
		// }
	}
	return nil
}

func (a *OKWSAgent) work() {
	ticker := time.NewTicker(9 * time.Second)

	a.keepalive()
	for {
		select {
		case <-ticker.C:
			if time.Now().Sub(a.lastPongTm) > maxPongInterval {
				log.Printf("lastPongTm %s timeout, reset connection", a.lastPongTm.Local().Format(time.RFC3339))
				a.connLock.Lock()
				a.conn.Close()
				a.connLock.Unlock()
				a.lastPongTm = time.Now().Add(maxPongInterval)
			}
			a.keepalive()
		case errR := <-a.wsErrCh:
			a.handleErrResponse(errR)
		case evtR := <-a.wsEvtCh:
			a.handleEventResponse(evtR)
		case tb := <-a.wsTbCh:
			a.handleTableResponse(tb)
		}
	}
}

func (a *OKWSAgent) receive() {
	for {
		messageType, message, err := a.conn.ReadMessage()
		if err != nil {
			log.Printf("a.conn.ReadMessage failed : %v", err)
			a.connLock.Lock()
			a.conn.Close()
			a.connLock.Unlock()
			conn, _, err := websocket.DefaultDialer.Dial(a.baseUrl, nil)
			if err != nil {
				log.Fatalf("a.receive : dial failed :%+v", err)
				time.Sleep(3 * time.Second)
				continue
			}
			a.connLock.Lock()
			log.Printf("a.receive - conn changed from %p -> %p", a.conn.UnderlyingConn(), conn.UnderlyingConn())
			a.conn = conn
			a.connLock.Unlock()
			if nil != a.startHook {
				if err = a.startHook(); nil != err {
					log.Printf("a.receive - a.startHook failed : %v, restart later", err)
					a.connLock.Lock()
					conn.Close()
					a.connLock.Unlock()
					time.Sleep(3 * time.Second)
				}
			}
			continue
		}

		txtMsg := message
		switch messageType {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			txtMsg, err = a.GzipDecode(message)
		}

		if string(txtMsg) == "pong" {
			a.lastPongTm = time.Now()
			continue
		}

		rsp, err := loadResponse(txtMsg)

		if err != nil {
			break
		}

		switch rsp.(type) {
		case *WSErrorResponse:
			a.wsErrCh <- rsp
		case *WSEventResponse:
			er := rsp.(*WSEventResponse)
			a.wsEvtCh <- er
		case *WSDepthTableResponse:

			dtr := rsp.(*WSDepthTableResponse)
			hotDepths := a.hotDepthsMap[dtr.Table]
			if hotDepths == nil {
				hotDepths = NewWSHotDepths(dtr.Table)
				hotDepths.loadWSDepthTableResponse(dtr)
				a.hotDepthsMap[dtr.Table] = hotDepths
			} else {
				hotDepths.loadWSDepthTableResponse(dtr)
			}
			a.wsTbCh <- dtr

		case *WSTableResponse:
			tb := rsp.(*WSTableResponse)
			a.wsTbCh <- tb
		default:
			log.Printf("LoadedRep: Warning - unknown response : %+v", string(txtMsg))
		}
	}
}
