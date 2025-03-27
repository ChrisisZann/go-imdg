package data

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// Dummy store for testing

const PAGE_SIZE = 10
const LINE_SIZE = 50

type memLine struct {
	Line [LINE_SIZE]byte
}

type MemPage struct {
	Page       [PAGE_SIZE]memLine
	DataExists [PAGE_SIZE]bool // Parallel table to track if data is present
	pos        int
	full       bool
	pageLock   sync.RWMutex // Mutex for thread-safe access
}

func (m *MemPage) Init() {
	m.pageLock.Lock()
	defer m.pageLock.Unlock()

	fmt.Println("memPage len=", len(m.Page))
	fmt.Println("memPage len=", len(m.Page[0].Line))
	m.full = false
	m.pos = 0
}

func (m *MemPage) Save(dat []byte) error {
	m.pageLock.Lock()
	defer m.pageLock.Unlock()

	if len(dat) >= LINE_SIZE {
		return errors.New("data >" + strconv.Itoa(LINE_SIZE))
	}

	if !m.full {
		copy(m.Page[m.pos].Line[:], dat)
		m.DataExists[m.pos] = true // Mark this position as containing data
	} else {
		fmt.Println("WARNING - Page is full")
		return errors.New("page is full")
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

func (m *MemPage) Read(pos int) ([]byte, error) {

	if pos < 0 || pos >= PAGE_SIZE {
		panic("position out of bounds")
	}
	m.pageLock.RLock()
	defer m.pageLock.RUnlock()

	if pos < 0 || pos >= PAGE_SIZE {
		return nil, errors.New("invalid position")
	}

	if !m.DataExists[pos] {
		return nil, errors.New("no data exists at this position")
	}

	return m.Page[pos].Line[:], nil
}
