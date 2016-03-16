// Package socketiocli implements a client for SocketIO protocol
// as specified in https://github.com/LearnBoost/socket.io-spec
package socketiocli

import (
	"fmt"
	logging "github.com/jeppeter/go-logging"
	"strconv"
	"strings"
	"time"
)

type Socket struct {
	URL       string
	Session   *Session
	Transport transport
	probenum  int
	curnumber int
}

// Dial opens a new client connection to the socket.io server then connects
// to the given channel.
func DialAndConnect(url string, channel string, query string) (*Socket, error) {
	socket, err := Dial(url)
	if err != nil {
		logging.Errorf("can not connect (%s) error (%s)", url, err.Error())
		return nil, err
	}

	endpoint := NewEndpoint(channel, query)
	connectMsg := NewConnect(endpoint)
	socket.Send(connectMsg)

	return socket, nil
}

// Dial opens a new client connection to the socket.io server using one of
// the implemented and supported Transports.
func Dial(url string) (*Socket, error) {
	session, err := NewSession(url)
	if err != nil {
		return nil, err
	}

	transport, err := newTransport(session, url)
	if err != nil {
		return nil, err
	}

	sock := &Socket{url, session, transport, 2, 42}

	defer func() {
		if sock != nil {
			sock.Close()
		}
	}()

	s := fmt.Sprintf("%dprobe", sock.probenum)
	err = transport.Send(s)

	if err != nil {
		sock.Close()
		return nil, err
	}

	msg, err := sock.Transport.Receive()
	if err != nil {
		sock.Close()
	}

	e, _, err := sock.parseMessage(msg)
	if err != nil {
		return nil, err
	}
	if e != "" {
		err = fmt.Errorf("receive message %s error", msg)
		return nil, err
	}

	// Heartbeat goroutine
	go func() {
		//heartbeatMsg := NewHeartbeat()
		for {
			time.Sleep(session.HeartbeatTimeout - time.Second)
			s := fmt.Sprintf("%dprobe", sock.probenum)
			err := transport.Send(s)
			if err != nil {
				return
			}
			logging.Debugf("send %s ok", s)
		}
	}()

	nsock := sock
	/*make socket not closed*/
	sock = nil

	return nsock, nil
}

func trimString(s string, f, l byte) string {
	cs := s
	if cs[0] == f {
		cs = cs[1:]
	}
	slen := len(cs)
	if slen > 0 {
		if cs[(slen-1)] == l {
			cs = cs[:(slen - 1)]
		}
	}

	return cs
}

func (socket *Socket) parseMessage(msg string) (string, string, error) {
	snum := ""
	partmsg := ""
	state := ""
	buf := []byte(msg)
	for i, a := range buf {
		if a >= '0' && a <= '9' {
			snum = snum + string(a)
		} else {
			partmsg = msg[i:]
			break
		}
	}
	id := -1

	id, err := strconv.Atoi(snum)
	if err != nil {
		logging.Errorf("parse %s error (%s)", snum, err.Error())
		return "", "", err
	}
	if partmsg == "probe" {
		s := fmt.Sprintf("%d", id+socket.probenum)
		err := socket.Transport.Send(s)
		if err != nil {
			return "", "", err
		}
		return "", "", nil
	}
	if len(partmsg) <= 0 {
		return "", "", nil
	}

	/*get message for */
	partmsg = trimString(partmsg, byte('['), byte(']'))

	sessions := strings.SplitN(partmsg, ",", 2)
	if len(sessions) < 2 {
		return "", "", nil
	}

	state = sessions[0]
	partmsg = sessions[1]
	state = trimString(state, '"', '"')
	partmsg = trimString(partmsg, '"', '"')

	return state, partmsg, nil
}

// Receive receives the raw message from the underlying transport and
// converts it to the Message type.
func (socket *Socket) Receive() (string, string, error) {
	var state string
	var newmsg string
	state = ""
	newmsg = ""
	for {
		rawMsg, err := socket.Transport.Receive()
		if err != nil {
			return "", "", err
		}

		state, newmsg, err = socket.parseMessage(rawMsg)
		if err != nil {
			return "", "", err
		} else if state != "" {
			/*we do not handle internal ,so out will give*/
			break
		}
	}
	return state, newmsg, nil
}

func (sock *Socket) SendMessage(state string, msg string) error {
	s := fmt.Sprintf("%d[\"%s\",%s]", sock.curnumber, state, msg)
	sock.curnumber++
	logging.Infof("send %s", s)
	err := sock.Transport.Send(s)
	return err
}

// Send sends the given Message to the socket.io server using it's
// underlying transport.
func (socket *Socket) Send(msg *Message) error {
	return socket.Transport.Send(msg.String())
}

// Close underlying transport
func (socket *Socket) Close() error {
	return socket.Transport.Close()
}
