package microSocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"errors"
	"net/http"
	"log"
	"time"
	"microSocket/util"
)

var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// websocket事件接口
type WebsktEventer interface {
	OnHandel(fd string) bool
	OnClose(fd string)
	OnMessage(fd string, msg string) bool
}

// websocket结构体
type Webskt struct {
	EventPool     *RouterMap
	SessionMaster *SessionM
	WebsktEvent   WebsktEventer
}

// websocket连接结构体
type WebsktConn struct {
	// 存放websocket连接
	wsConn *websocket.Conn
	// 用于存放数据
	inChan chan []byte
	// 用于读取数据
	outChan chan []byte
	closeChan chan byte
	mutex sync.Mutex
	// chan是否被关闭
	isClosed bool
	fd string
	Webskt	*Webskt
}

// 创建websocket
func NewWebskt(WebsktEvent WebsktEventer) *Webskt {
	return &Webskt{
		SessionMaster: NewSessonM(),
		WebsktEvent:   WebsktEvent,
		EventPool:     NewRouterMap(),
	}
}

// websocket监听地址
func (webskt *Webskt) Listening(address string) bool {
	go http.HandleFunc("/", webskt.connHandle)
	// 监听地址
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Println("websocket监听地址失败(websocket Failed to monitor address)")
		return false
	}
	return true
}

// websocket连接处理
func (webskt *Webskt) connHandle(w http.ResponseWriter, r *http.Request) {
	var (
		wsConn *websocket.Conn
		err    error
		data   []byte
		conn   *WebsktConn
	)
	// 完成http应答，在httpheader中放下如下参数
	if wsConn, err = upgrader.Upgrade(w, r, nil); err != nil {
		return // 获取连接失败直接返回
	}
	if conn, err = webskt.initWebsktConn(wsConn); err != nil {
		goto ERR
	}

	// 启动一个协程发送心跳包
	go func() {
		var (
			err error
		)
		for {
			log.Printf(conn.fd)
			// 每隔一秒发送一次心跳
			if err = conn.WriteMessage([]byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}

	}()
	// 调用握手事件
	webskt.WebsktEvent.OnHandel(conn.fd)
	for {
		if data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		// 数据转换为map(data conversion map)
		requestData := string(data)
		webskt.WebsktEvent.OnMessage(conn.fd, requestData)
		//webskt.EventPool.Hook(requestData["module"], requestData["action"], requestData)
		//if err = conn.WriteMessage(data); err != nil {
		//	goto ERR
		//}
		conn.WriteMessage(data)
	}

ERR:
	// 关闭当前连接
	conn.Close()
}

// 读取Api
func (conn *WebsktConn) ReadMessage() (data []byte, err error) {
	//select是Go中的一个控制结构，类似于用于通信的switch语句。每个case必须是一个通信操作，要么是发送要么是接收。 select随机执行一个可运行的case。如果没有case可运行，它将阻塞，直到有case可运行。一个默认的子句应该总是可运行的。
	select {
	case data = <- conn.inChan:
	case <- conn.closeChan:
		log.Println("关闭Chan失败(WebsktConn Chan failed)")
		return nil, errors.New("关闭Chan失败(WebsktConn Chan failed)")
	}
	return data, nil
}

// 发送Api
func (conn *WebsktConn) WriteMessage(data []byte) (err error) {
	select {
	case conn.outChan <- data:
	case <- conn.closeChan:
		log.Println("连接失败(WebsktConn failed)")
		return errors.New("WebsktConn is closed")
	}
	return nil
}

// 关闭连接的Api
func (conn *WebsktConn) Close()  {
	// 线程安全的Close，可以并发多次调用也叫做可重入的Close
	conn.wsConn.Close()
	conn.mutex.Lock()
	if !conn.isClosed {
		// 关闭chan,但是chan只能关闭一次
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()
}

// 初始化长连接
func (webskt *Webskt) initWebsktConn(wsConn *websocket.Conn) (*WebsktConn, error)  {
	conn := &WebsktConn{
			wsConn: wsConn,
			inChan: make(chan []byte, 1000),
			outChan: make(chan []byte, 1000),
			closeChan: make(chan byte, 1),
			fd: util.Transmitter(),
	}
	// 启动读协程
	go conn.readLoop()
	// 启动写协程
	go conn.writeLoop()
	return conn, nil
}

// 读协程
func (conn *WebsktConn) readLoop() {
	var (
		data []byte
		err error
	)
	for {
		if _, data, err = conn.wsConn.ReadMessage(); err != nil {
			goto ERR
		}
		// 容易阻塞到这里，等待inChan有空闲的位置
		select {
		case conn.inChan <- data:
		case <- conn.closeChan: // closeChan关闭的时候执行
			goto ERR
		}
	}

ERR:
	conn.Close()
}

// 写协程
func (conn *WebsktConn) writeLoop() {
	var (
		data []byte
		err error
	)
	for {
		select {
		case data = <- conn.outChan:
		case <- conn.closeChan:
			goto ERR
		}
		data = <- conn.outChan
		if err = conn.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	conn.Close()
}

