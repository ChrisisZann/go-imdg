package comms

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

type payloadType int

const (
	cmd payloadType = iota
	dat
	bad
	def
)

/* ********************************************************************************
 *
 * any change in types need to update functions: ParsePayloadType, String
 * ********************************************************************************/
func (pt payloadType) String() string {
	if pt == cmd {
		return "0"
	} else if pt == dat {
		return "1"
	} else if pt == def {
		return "3"
	}
	return ""
}

func ParsePayloadType(s string) payloadType {
	si, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}

	if si == 0 {
		return cmd
	} else if si == 1 {
		return dat
	} else if si == 3 {
		return def
	}
	return bad
}

type Payload struct {
	ptype payloadType
	data  string
	// delim string
}

func NewPayload(s string, pt payloadType) *Payload {

	// DEBUGGING
	// temp := &Payload{
	// 	ptype: pt,
	// 	data:  s,
	// }
	// fmt.Println("New Payload created:", temp.String())

	return &Payload{
		ptype: pt,
		data:  s,
	}
}

func (p Payload) String() string {
	switch p.ptype {
	case cmd:
		return "cmd:" + p.data
	case dat:
		return "data:" + p.data
	case def:
		return "def:" + p.data
	default:
		return "bad payload"
	}
}

func (p Payload) Compile() (string, error) {
	if !p.validate() {
		return "", errors.New("failed to compile payload validations")
	}
	return p.ptype.String() + ":" + p.data, nil
	// return []byte(p.ptype.String() + p.delim + p.data), nil

}

func (p Payload) validate() bool {
	// TODO
	return true
}

func ParsePayload(s string) *Payload {
	tok := strings.Split(s, ":")
	if len(tok) != 2 {
		log.Fatal("wrong payload structure")
		return nil
	}
	return &Payload{
		ptype: ParsePayloadType(tok[0]),
		data:  tok[1],
	}
}
