package comms

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

// <source|suid|dest|payload>
type message struct {
	source      NodeAddr
	suid        int
	destination NodeAddr
	payload     *Payload
}

func (m message) String() string {
	return "source=" + m.source.String() + ", suid=" + strconv.Itoa(m.suid) + ", payload=" + m.payload.String()
}

func (m message) Compile() string {

	pl, err := m.payload.Compile()
	if err != nil {
		log.Fatal("Failed to compile payload")
	}
	return m.source.String() + "|" + strconv.Itoa(m.suid) + "|" + m.destination.String() + "|" + pl
}

func ParseMessage(input string) (*message, error) {
	tok := strings.Split(input, "|")
	if len(tok) != 4 {
		return nil, errors.New("wrong message structure:" + strconv.Itoa(len(tok)))
	}

	// DEBUGGING
	// for i, t := range tok {
	// 	fmt.Printf("tok %d=%s\n", i, t)
	// }

	sourceAddr := NewNodeAddr("tcp", tok[0])

	uid, err := strconv.Atoi(tok[1])
	if err != nil {
		return nil, err
	}

	destAddr := NewNodeAddr("tcp", tok[2])

	return &message{
		source:      sourceAddr,
		destination: destAddr,
		suid:        uid,
		payload:     ParsePayload(tok[3]),
	}, nil
}

// <source|suid|dest|payload>
func CompileHeader(source, suid, dest string) string {
	return source + "|" + suid + "|" + dest + "|"
}
