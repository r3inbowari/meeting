package web

import (
	"github.com/sirupsen/logrus"
	"meeting/utils"
	"net"
	"strings"
	"sync"
	"time"
)

func GetDTUSessionsKey() []string {
	var ret []string
	SessionsMap.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(string))
		return true
	})
	return ret
}

type DTUSession struct {
	readChan  chan []byte
	writeChan chan []byte
	stopChan  chan bool
	conn      net.Conn

	//tw *TimeWheel
}

func (ds *DTUSession) readConn() {
	for {
		data := make([]byte, 64)
		n, err := ds.conn.Read(data)
		if err != nil {
			break
		}
		ds.readChan <- data[:n]
	}
	ds.stopChan <- true
}

/**
 * 写
 */
func (ds *DTUSession) Write(b []byte) error {
	if _, err := ds.conn.Write(b); err != nil {
		return err
	}
	return nil
}

/**
 * 释放一个session
 */
func (ds *DTUSession) Release() {
	SessionsMap.Delete(GetIP(ds.conn))
	utils.Info("session release", logrus.Fields{"addr": GetIP(ds.conn)})
}

func DtuHandle(conn net.Conn) {
	defer func() { _ = conn.Close() }()
	session := RegDTUSession(conn)
	go session.readConn()

	time.Sleep(time.Second)
	session.Write([]byte{0xff, 0x05, 0xa1, 0xf1, 0x3a, 0x38, 0x66})

	for {
		select {
		case read := <-session.readChan:
			println(string(read))
			// session.Write([]byte("123FF4"))
		case stop := <-session.stopChan:
			if stop {
				utils.Info("disconnected", logrus.Fields{"addr": GetIP(conn)})
				session.Release()
				return
			}
		}
	}
}

var SessionsMap sync.Map

func RegDTUSession(conn net.Conn) DTUSession {
	var ds DTUSession
	ds.readChan = make(chan []byte)  // 读
	ds.writeChan = make(chan []byte) // 写
	ds.stopChan = make(chan bool)    // 停
	ds.conn = conn                   // 连接
	addr := GetIP(conn)
	SessionsMap.Store(addr, ds)
	utils.Info("iot connected", logrus.Fields{"addr": addr})
	return ds
}

func GetIP(conn net.Conn) string {
	return strings.Split(conn.RemoteAddr().String(), ":")[0]
}
