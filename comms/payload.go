package comms

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type PayloadType int

const (
	cmd PayloadType = iota
	dat
	network
	def
	bad
)

/* ********************************************************************************
 *
 * any change in types need to update functions: ParsePayloadType, String
 * ********************************************************************************/
func (pt PayloadType) String() string {
	if pt == cmd {
		return "cmd"
	} else if pt == dat {
		return "dat"
	} else if pt == def {
		return "def"
	} else if pt == network {
		return "network"
	}
	return ""
}

func StringToPayloadType(s string) PayloadType {

	if strings.Compare(s, "cmd") == 0 {
		return cmd
	} else if strings.Compare(s, "dat") == 0 {
		return dat
	} else if strings.Compare(s, "def") == 0 {
		return def
	} else if strings.Compare(s, "network") == 0 {
		return network
	}
	return bad
}

func ParsePayloadType(s string) PayloadType {
	si, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("Yep this is the error")
		log.Fatal(err)
	}

	if si == 0 {
		return cmd
	} else if si == 1 {
		return dat
	} else if si == 2 {
		return network
	} else if si == 3 {
		return def
	}
	return bad
}

type Payload struct {
	ptype PayloadType
	data  string
	// delim string
}

func NewPayload(d, pt string) (*Payload, error) {

	if !validatePayloadData(d) {
		return nil, errors.New("bad payload data")
	}

	pType := StringToPayloadType(pt)
	if pType == bad {
		return nil, errors.New("bad payload type")
	}

	return &Payload{
		ptype: pType,
		data:  d,
	}, nil
}

func (p Payload) ReadType() PayloadType {
	return p.ptype
}

func (p Payload) ReadData() string {

	// Emit zeros in buffer of Message
	return strings.Trim(p.data, "\x00")
}

func (p Payload) String() string {

	// Emit zeros in buffer the Message
	return p.ReadType().String() + ":" + strings.Trim(p.data, "\x00")
}

func (p Payload) Compile() (string, error) {
	return p.ptype.String() + ":" + p.data, nil
	// return []byte(p.ptype.String() + p.delim + p.data), nil
}

func validatePayloadData(s string) bool {

	if strings.Contains(s, " ") {
		return false
	} else if strings.Contains(s, ":") {
		return false
	}
	return true
}

func ParsePayload(s string) *Payload {
	tok := strings.Split(s, ":")
	if len(tok) != 2 {
		log.Fatal("wrong payload structure")
		return nil
	}
	return &Payload{
		ptype: StringToPayloadType(tok[0]),
		data:  tok[1],
	}
}
