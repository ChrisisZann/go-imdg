package data

import (
	"errors"
	"fmt"
	"strconv"
)

//Dummy store for testing

const PAGE_SIZE = 10
const LINE_SIZE = 50

type memLine struct {
	Line [LINE_SIZE]byte
}

type MemPage struct {
	Page [PAGE_SIZE]memLine
	pos  int
	full bool
}

func (m *MemPage) Init() {

	fmt.Println("memPage len=", len(m.Page))
	fmt.Println("memPage len=", len(m.Page[0].Line))
	m.full = false
	m.pos = 0
}

func (m *MemPage) Save(dat []byte) error {
	if len(dat) >= LINE_SIZE {
		return errors.New("data >" + strconv.Itoa(LINE_SIZE))
	}

	if !m.full {
		copy(m.Page[m.pos].Line[:], dat)
	} else {
		fmt.Println("WARNING - Page is full")
		return errors.New("Page is full")
	}

	m.pos++
	if m.pos >= PAGE_SIZE {
		m.full = true
		fmt.Println("Page is now full")
	} else {
		fmt.Println("Remaining lines:", PAGE_SIZE-m.pos)
	}

	return nil
}

func (m MemPage) Read(pos int) []byte {
	return m.Page[pos].Line[:]
}
