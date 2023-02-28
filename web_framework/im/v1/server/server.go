package server

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/sirupsen/logrus"
)

// ServerOptions ServerOptions
type ServerOptions struct {
	writewait time.Duration // 写超时时间
	readwait  time.Duration // 读超时时间
}

// Server is a websocket implement of the Server
type Server struct {
	once    sync.Once
	options ServerOptions
	id      string
	address string
	sync.Mutex
	// 会话列表，key=uid，value=connection
	users map[string]net.Conn
}

// NewServer NewServer
func NewServer(id, address string) *Server {
	return newServer(id, address)
}

func newServer(id, address string) *Server {
	return &Server{
		id:      id,
		address: address,
		users:   make(map[string]net.Conn, 100),
		options: ServerOptions{
			writewait: time.Second * 10,
			readwait:  time.Minute * 2,
		},
	}
}

// Start server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	log := logrus.WithFields(logrus.Fields{
		"module": "Server",
		"listen": s.address,
		"id":     s.id,
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			conn.Close()
			return
		}
		// 读取userId
		user := r.URL.Query().Get("user")
		if user == "" {
			conn.Write([]byte("no user param"))
			conn.Close()
			return
		}

		// 添加到会话管理中
		old, ok := s.addUser(user, conn)
		if ok {
			// 断开旧的连接
			old.Close()
			log.Infof("close old connection %v", old.RemoteAddr())
		}
		log.Infof("user %s in from %v", user, conn.RemoteAddr())

		go func(user string, conn net.Conn) {
			err := s.readloop(user, conn)
			if err != nil {
				log.Warn("readloop - ", err)
			}
			conn.Close()
			// 删除用户
			s.delUser(user)

			log.Infof("connection of %s closed", user)
		}(user, conn)
	})
	log.Infoln("started")
	return http.ListenAndServe(s.address, mux)
}

// addUser 用户上线
func (s *Server) addUser(user string, conn net.Conn) (net.Conn, bool) {
	s.Lock()
	defer s.Unlock()
	// 判断用户是否存在历史链接，存在则可以返回让后续 t 下线
	old, ok := s.users[user]
	s.users[user] = conn
	return old, ok
}

// delUser 剔除用户
func (s *Server) delUser(user string) {
	s.Lock()
	defer s.Unlock()

	// 直接从 map 移除对应的用户
	delete(s.users, user)
}

// Shutdown
func (s *Server) Shutdown() {
	s.once.Do(func() {
		s.Lock()
		defer s.Unlock()

		// 将所有客户端链接安全关闭
		for _, conn := range s.users {
			conn.Close()
		}
	})
}

func (s *Server) readloop(user string, conn net.Conn) error {
	for {
		// 要求：客户端必须在指定时间内发送一条消息过来，可以是ping，也可以是正常数据包
		_ = conn.SetReadDeadline(time.Now().Add(s.options.readwait))

		// 从 TCP 缓存中区读取一帧数据
		frame, err := ws.ReadFrame(conn)
		if err != nil {
			return err
		}

		// 判断帧头部是否为 ping 请求
		if frame.Header.OpCode == ws.OpPing {
			// 返回一个pong消息
			_ = wsutil.WriteServerMessage(conn, ws.OpPong, nil)
			logrus.Info("wirte a pong...")
			continue
		}

		// 判断是否为 close 请求的
		if frame.Header.OpCode == ws.OpClose {
			return errors.New("remote side close the conn")
		}
		logrus.Info(frame.Header)

		// ws协议规定客户端发送消息的时候必须使用随机的mask码对消息体做一次编码
		// 服务端此时得解码
		if frame.Header.Masked {
			// 解码数据
			ws.Cipher(frame.Payload, frame.Header.Mask, 0)
		}
		// 接收文本帧内容
		if frame.Header.OpCode == ws.OpText {
			go s.handleBoardCast(user, string(frame.Payload))
		} else if frame.Header.OpCode == ws.OpBinary {
			go s.handleBinary(user, frame.Payload)
		}
	}
}

// 广播消息
func (s *Server) handleBoardCast(user string, message string) {
	logrus.Infof("recv message %s from %s", message, user)
	s.Lock()
	defer s.Unlock()
	broadcast := fmt.Sprintf("%s -- FROM %s", message, user)
	for u, conn := range s.users {
		if u == user { // 跳过自己
			continue
		}
		logrus.Infof("send to %s : %s", u, broadcast)
		err := s.writeText(conn, broadcast)
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}

// 写入文本
func (s *Server) writeText(conn net.Conn, message string) error {
	// 创建文本帧数据
	f := ws.NewTextFrame([]byte(message))
	return ws.WriteFrame(conn, f)
}

// command of message
const (
	CommandPing = 100
	CommandPong = 101
)

func (s *Server) handleBinary(user string, message []byte) {
	logrus.Infof("recv message %v from %s", message, user)
	s.Lock()
	defer s.Unlock()
	// handle ping request
	i := 0
	command := binary.BigEndian.Uint16(message[i : i+2])
	i += 2
	payloadLen := binary.BigEndian.Uint32(message[i : i+4])
	logrus.Infof("command: %v payloadLen: %v", command, payloadLen)
	if command == CommandPing {
		u := s.users[user]
		// return pong
		err := wsutil.WriteServerBinary(u, []byte{0, CommandPong, 0, 0, 0, 0})
		if err != nil {
			logrus.Errorf("write to %s failed, error: %v", user, err)
		}
	}
}
