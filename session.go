package socketiocli

import (
	"errors"
	"fmt"
	logging "github.com/jeppeter/go-logging"
	"io/ioutil"
	"net/http"
	"regexp"
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

func sessionCreate(resp *http.Response) (*Session, error) {
	var numexp *regexp.Regexp
	var matchstrs []string
	var sbody string
	var sessionVars []string
	var lefts string
	var slen int
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	sbody = string(body)

	numexp, err = regexp.Compile("^([0-9]+)\\{")
	if err != nil {
		return nil, err
	}

	/*now we should handle the body*/
	fmt.Printf("body [%s]\n", sbody)
	matchstrs = numexp.FindStringSubmatch(sbody)
	if len(matchstrs) > 1 {
		sessionVars = []string{}
		sessionVars = append(sessionVars, matchstrs[1])
		lefts = strings.Replace(sbody, sessionVars[0], "", 1)
		sessionVars = append(sessionVars, lefts)
	} else {
		sessionVars = strings.SplitN(sbody, ":", 2)
	}
	if len(sessionVars) != 2 {
		return nil, errors.New("Session variables is not valid")
	}

	slen, err = strconv.Atoi(sessionVars[0])
	if err != nil {
		return nil, err
	}

	if slen > len(sessionVars[1]) {
		err = fmt.Errorf("slen %d != len(%d)", slen, len(sessionVars[1]))
		logging.Errorf("%s", err.Error())
		return nil, err
	}

	indexjson := strings.Index(sessionVars[1], "{")
	if indexjson < 0 {
		err = fmt.Errorf(" %s not valid json format ", sessionVars[1])
		return nil, err
	}

	jsonstr := sessionVars[1][indexjson:]
	id := getJsonValueDefault(jsonstr, SOCKET_IO_SID, "")
	if id == "" {
		err = fmt.Errorf("%s not valid", SOCKET_IO_SID)
		return nil, err
	}
	heartbeattimeout := getJsonValueDefault(jsonstr, SOCKET_IO_HEART_TIMEOUT, "10000.0")
	connectiontimeout := getJsonValueDefault(jsonstr, SOCKET_IO_TIMEOUT, "30000.0")
	supportprotocol := getJsonValueDefault(jsonstr, SOCKET_IO_PROTOCOL, "")

	milsec, err := strconv.ParseFloat(heartbeattimeout, 64)
	if err != nil {
		logging.Errorf("%s set value value (%s) error(%s)", SOCKET_IO_HEART_TIMEOUT, heartbeattimeout, err.Error())
		return nil, err
	}
	intmilsec := int(milsec)
	heartbeatTimeout := time.Duration(intmilsec) * time.Millisecond

	milsec, err = strconv.ParseFloat(connectiontimeout, 64)
	if err != nil {
		logging.Errorf("%s set value value (%s) error(%s)", SOCKET_IO_TIMEOUT, connectiontimeout, err.Error())
		return nil, err
	}
	intmilsec = int(milsec)
	connectionTimeout := time.Duration(intmilsec) * time.Millisecond
	supportedProtocols := strings.Split(supportprotocol, ",")

	return &Session{id, heartbeatTimeout, connectionTimeout, supportedProtocols}, nil

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
	client := &http.Client{}
	req, err := http.NewRequest("GET", shakequery, nil)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Accept-Encoding", "identity")
	//req.Header.Del("Accept-Encoding")
	//req.Header.Set("User-Agent", "node-XMLHttpRequest")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "close")
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return sessionCreate(response)
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
