package comms

import (
	"errors"
	"net"
	"strconv"
	"strings"
)

const LENGTH int = 4

type message struct {
	source   net.Addr
	respPort string

	sourceUUID int
	content    string
}

func (m message) Header() string {
	return m.source.String() + ";" + m.respPort + ";" + strconv.Itoa(m.sourceUUID) + ";"
}

func (m message) String() string {
	return "source=" + m.source.String() + ", respPort=" + m.respPort + ", sourceUUID=" + strconv.Itoa(m.sourceUUID) + ", content=" + m.content
}

func (m message) CreateMessageString() string {
	return m.Header() + m.content
}

func (m message) CompileMessage() []byte {
	return []byte(m.Header() + m.content)
}

func ParseMessage(input string) (*message, error) {
	tok := strings.Split(input, ";")
	if len(tok) != LENGTH {
		return nil, errors.New("wrong message structure")
	}
	uuid, err := strconv.Atoi(tok[2])
	if err != nil {
		return nil, err
	}

	temp := masterAddr{network: "tcp", addr: tok[0]}
	return &message{
		source:     temp,
		respPort:   tok[1],
		sourceUUID: uuid,
		content:    tok[3],
	}, nil
}
