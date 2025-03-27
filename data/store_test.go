package data

import (
	"strings"
	"testing"
)

func TestMemPage_Init(t *testing.T) {
	var memPage MemPage
	memPage.Init()

	if memPage.pos != 0 {
		t.Errorf("Expected pos to be 0, got %d", memPage.pos)
	}
	if memPage.full {
		t.Error("Expected full to be false, got true")
	}
}

func TestMemPage_Save(t *testing.T) {
	var memPage MemPage
	memPage.Init()

	// Test saving valid data
	err := memPage.Save([]byte("test_data"))
	if err != nil {
		t.Errorf("Failed to save valid data: %v", err)
	}

	// Test saving data exceeding LINE_SIZE
	largeData := strings.Repeat("a", LINE_SIZE+1)
	err = memPage.Save([]byte(largeData))
	if err == nil {
		t.Error("Expected error when saving data exceeding LINE_SIZE, but got none")
	}

	// Fill the page to test saving when full
	for i := 1; i < PAGE_SIZE; i++ {
		err = memPage.Save([]byte("fill"))
		if err != nil {
			t.Errorf("Failed to save data at position %d: %v", i, err)
		}
	}

	// Test saving when page is full
	err = memPage.Save([]byte("overflow"))
	if err == nil {
		t.Error("Expected error when saving data to a full page, but got none")
	}
}

func TestMemPage_Read(t *testing.T) {
	var memPage MemPage
	memPage.Init()

	// Save some data
	testData := []string{"data1", "data2", "data3"}
	for _, d := range testData {
		err := memPage.Save([]byte(d))
		if err != nil {
			t.Errorf("Failed to save data '%s': %v", d, err)
		}
	}

	// Read and verify the data
	for i, d := range testData {
		if !memPage.DataExists[i] {
			t.Errorf("Expected data to exist at position %d, but it does not", i)
			continue
		}
		data, err := memPage.Read(i)
		if err != nil {
			t.Errorf("Failed to read data at position %d: %v", i, err)
			continue
		}
		readData := strings.TrimRight(string(data), "\x00")
		if readData != d {
			t.Errorf("Expected '%s', got '%s'", d, readData)
		}
	}

	// Test reading from an invalid position
	invalidPos := PAGE_SIZE + 1
	memPage.Read(invalidPos)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when reading from invalid position %d, but got none", invalidPos)
		}
	}()
	_, _ = memPage.Read(invalidPos)
}
