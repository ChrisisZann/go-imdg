package comms

import (
	"errors"
	"strconv"
	"strings"
)

// <source|suid|dest|payload>
type Message struct {
	source      NodeAddr
	suid        int
	destination NodeAddr
	payload     *Payload
}

func NewMessage(source NodeAddr, suid int, destination NodeAddr, payload *Payload) *Message {
	return &Message{
		source:      source,
		suid:        suid,
		destination: destination,
		payload:     payload,
	}
}

func (m Message) ReadSender() NodeAddr {
	return m.source
}

func (m Message) ReadDest() string {
	return m.destination.String()
}

func (m Message) ReadSenderID() int {
	return m.suid
}

func (m Message) ReadPayloadData() string {
	return m.payload.ReadData()
}

func (m Message) GetPayloadType() PayloadType {
	return m.payload.ptype
}

func (m Message) ReadPayloadType() string {
	return m.payload.ptype.String()
}

func (m Message) String() string {
	return "source=" + m.source.String() + ", suid=" + strconv.Itoa(m.suid) + ", payload=" + m.payload.String()
}

func (m Message) Compile() (string, error) {
	pl, err := m.payload.Compile()
	if err != nil {
		return "", err
	}
	return m.source.String() + "|" + strconv.Itoa(m.suid) + "|" + m.destination.String() + "|" + pl, nil
}

func ParseMessage(input string) (*Message, error) {
	tok := strings.Split(input, "|")
	if len(tok) != 4 {
		return nil, errors.New("wrong Message structure: expected 4 parts, got " + strconv.Itoa(len(tok)) + " in input: " + input)
	}

	// DEBUGGING
	// for i, t := range tok {
	// 	fmt.Printf("tok %d=%s\n", i, t)
	// }

	sourceAddr, err := NewNodeAddr("tcp", tok[0])
	if err != nil {
		return nil, err
	}

	uid, err := strconv.Atoi(tok[1])
	if err != nil {
		return nil, err
	}

	destAddr, err := NewNodeAddr("tcp", tok[2])
	if err != nil {
		return nil, err
	}

	// TODO : PATCH payload first
	// payload, err := ParsePayload(tok[3])
	// if err != nil {
	// 	return nil, err
	// }

	return &Message{
		source:      sourceAddr,
		destination: destAddr,
		suid:        uid,
		payload:     ParsePayload(tok[3]),
		// payload:     payload,
	}, nil
}

// <source|suid|dest|payload>
func CompileHeader(source, suid, dest string) string {
	return source + "|" + suid + "|" + dest + "|"
}
