// dbus1.go

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/godbus/dbus"
	//"github.com/godbus/dbus/introspect"
)

type Callback func(s string, userdata interface{}) error

type MyDBus struct {
	Conn     *dbus.Conn
	F        Callback
	Name     string
	ch       chan os.Signal
	userdata interface{}
}

// Null callback function
func NullCall(s string, userdata interface{}) error { return nil }

// server mode remote callable with DBus,use callback function to deal string v
func (self *MyDBus) ClientMsg(v string) *dbus.Error {
	self.F(v, self.userdata)
	return nil
}

// server mode serve with block mode
func (self *MyDBus) Serv() error {
	defer self.Conn.Close()
	defer self.Conn.ReleaseName(self.Name)
	path1 := fmt.Sprintf("/%v", strings.Replace(self.Name, ".", "/", -1))
	err := self.Conn.Export(self, dbus.ObjectPath(path1), self.Name)
	if err != nil {
		return err
	}
	self.ch = make(chan os.Signal, 1)
	signal.Notify(self.ch, os.Interrupt)
	<-self.ch
	self.ch = nil
	return nil
}

// server mode shutdown
func (self *MyDBus) Shutdown() {
	if self.ch == nil {
		self.Conn.Close()
		self.Conn.ReleaseName(self.Name)
	} else {
		self.ch <- os.Interrupt
	}
}

// client mode send message
func (self *MyDBus) ClientSend(s string) error {
	bo1 := self.Conn.Object(self.Name, dbus.ObjectPath(fmt.Sprintf("/%v", strings.Replace(self.Name, ".", "/", -1))))
	c1 := bo1.Call(self.Name+".ClientMsg", 0, s)
	return c1.Err
}

// client mode close connection
func (self *MyDBus) ClientClose() {
	self.Conn.Close()
}

// try make new DBus server, or make client if server is exists
func TryNewMyDBus(name1 string, f Callback, userdata interface{}) (bus1 *MyDBus, unique bool) {
	conn, err := dbus.SessionBus()
	if err != nil {
		log.Println(err)
		return nil, false
	}
	rb1 := &MyDBus{conn, f, name1, nil, userdata}
	reply, err := conn.RequestName(name1,
		dbus.NameFlagDoNotQueue)
	if err != nil {
		log.Println(err)
		return rb1, false
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		log.Println("name already taken : ", name1)
		return rb1, false
	}
	return rb1, true
}

// client get channel to listen, 不要在一个进程中重复调用
func (self *MyDBus) GetChan() chan *dbus.Signal {
	ch := make(chan *dbus.Signal, 10)
	//path1 := fmt.Sprintf("/%v", strings.Replace(self.Name, ".", "/", -1))
	//sig1 := fmt.Sprintf("type='signal',interface='%v'", self.Name)
	o1 := self.Conn.BusObject().(*dbus.Object)
	o1.AddMatchSignal(self.Name, "SigMsg")
	self.Conn.Signal(ch)
	return ch
}

// server send signal
func (self *MyDBus) SendSignal(s string) {
	path1 := fmt.Sprintf("/%v", strings.Replace(self.Name, ".", "/", -1))
	self.Conn.Emit(dbus.ObjectPath(path1), self.Name+".SigMsg", s)
}

// set userdata
func (self *MyDBus) SetUserdata(data interface{}) {
	self.userdata = data
}
