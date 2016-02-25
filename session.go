package socketiocli

import (
	"errors"
	logging "github.com/jeppeter/go-logging"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Session holds the configuration variables received from the socket.io
// server.
type Session struct {
	ID                 string
	HeartbeatTimeout   time.Duration
	ConnectionTimeout  time.Duration
	SupportedProtocols []string
}

// NewSession receives the configuraiton variables from the socket.io
// server.
func NewSession(urls string) (*Session, error) {
	urlParser, err := newURLParser(urls)
	if err != nil {
		logging.Errorf("error %s", err.Error())
		return nil, err
	}

	shakequery := urlParser.handshake()
	logging.Debugf("shakequery %s", shakequery)
	response, err := http.PostForm(shakequery,url.Values{''})
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	response.Body.Close()

	sessionVars := strings.Split(string(body), ":")
	if len(sessionVars) != 4 {
		return nil, errors.New("Session variables is not 4")
	}

	id := sessionVars[0]

	heartbeatTimeoutSec, _ := strconv.Atoi(sessionVars[1])
	connectionTimeoutSec, _ := strconv.Atoi(sessionVars[2])

	heartbeatTimeout := time.Duration(heartbeatTimeoutSec) * time.Second
	connectionTimeout := time.Duration(connectionTimeoutSec) * time.Second

	supportedProtocols := strings.Split(string(sessionVars[3]), ",")

	return &Session{id, heartbeatTimeout, connectionTimeout, supportedProtocols}, nil
}

// SupportProtocol checks if the given protocol is supported by the
// socket.io server.
func (session *Session) SupportProtocol(protocol string) bool {
	for _, supportedProtocol := range session.SupportedProtocols {
		if protocol == supportedProtocol {
			return true
		}
	}
	return false
}
