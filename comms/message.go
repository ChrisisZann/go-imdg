package comms

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// <source|suid|dest|payload>
/* ********************************************************************************
 *
 *
 * ********************************************************************************/
type message struct {
	source      nodeAddr
	suid        int
	destination nodeAddr
	respPort    string
	payload     *Payload
}

/* ********************************************************************************
 *
 *
 * ********************************************************************************/
// func NewMessage() {
// 	return &message{}
// }

/* ********************************************************************************
 *
 *
 * ********************************************************************************/
func (m message) String() string {
	return "source=" + m.source.String() + ", respPort=" + m.respPort + ", suid=" + strconv.Itoa(m.suid) + ", payload=" + m.payload.String()
}

/* ********************************************************************************
 *
 *
 * ********************************************************************************/
func (m message) Compile() string {

	pl, err := m.payload.Compile()
	if err != nil {
		log.Fatal("Failed to compile payload")
	}

	fmt.Println("\nm.source.String()=", m.source.String())
	fmt.Println("\nstrconv.Itoa(m.suid)=", strconv.Itoa(m.suid))
	fmt.Println("\nm.destination.String()", m.destination.String())
	fmt.Println("\npl=", pl)

	tempMessage := m.source.String() + "|" + strconv.Itoa(m.suid) + "|" + m.destination.String() + "|" + pl
	fmt.Println("\ntempMessage=", tempMessage)
	return tempMessage
}

/* ********************************************************************************
 *
 *
 * ********************************************************************************/
func ParseMessage(input string) (*message, error) {

	tok := strings.Split(input, "|")
	if len(tok) != 4 {
		return nil, errors.New("wrong message structure:" + strconv.Itoa(len(tok)))
	}

	sourceAddr := NewNodeAddr("tcp", tok[0])

	uid, err := strconv.Atoi(tok[1])
	if err != nil {
		return nil, err
	}

	destAddr := NewNodeAddr("tcp", tok[2])

	return &message{
		source:      sourceAddr,
		destination: destAddr,
		respPort:    "",
		suid:        uid,
		payload:     ParsePayload(tok[3]),
	}, nil
}
