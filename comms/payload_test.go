package comms

import "testing"

func TestEmptyPayloadData(t *testing.T) {
	_, err := NewPayload("", "cmd")
	if err == nil {
		t.Error(`NewPayload("", "cmd") : created payload with empty data`)
	}
}

func TestEmptyPayloadType(t *testing.T) {
	_, err := NewPayload("my_data", "")
	if err == nil {
		t.Error(`NewPayload("my_data", "") : created payload with empty type`)
	}
}

func TestInvalidPayloadType(t *testing.T) {
	_, err := NewPayload("my_data", "invalid")
	if err == nil {
		t.Error(`NewPayload("my_data", "invalid") : created payload with invalid type`)
	}
}

func TestPayloadString(t *testing.T) {
	payload, err := NewPayload("my_data", "cmd")
	if err != nil {
		t.Error(`NewPayload("my_data", "cmd") failed to create payload`)
	}
	expected := "cmd:my_data"
	if payload.String() != expected {
		t.Errorf("Expected payload string to be %s, but got %s", expected, payload.String())
	}
}

func TestPayloadReadType(t *testing.T) {
	payload, err := NewPayload("my_data", "cmd")
	if err != nil {
		t.Error(`NewPayload("my_data", "cmd") failed to create payload`)
	}
	if payload == nil {
		t.Error("ParsePayload(\"cmd:my_data\") returned nil payload")
	} else if payload.ReadType() != cmd {
		t.Errorf("Expected payload type to be %v, but got %v", cmd, payload.ReadType())
	}
}

func TestPayloadReadData(t *testing.T) {
	payload, err := NewPayload("my_data", "cmd")
	if err != nil {
		t.Error(`NewPayload("my_data", "cmd") failed to create payload`)
	}
	if payload.ReadData() != "my_data" {
		t.Errorf("Expected payload data to be %s, but got %s", "my_data", payload.ReadData())
	}
}

func TestParsePayload(t *testing.T) {
	payloadStr := "cmd:my_data"
	payload := ParsePayload(payloadStr)

	if payload == nil {
		t.Error("ParsePayload(\"cmd:my_data\") returned nil payload")
	} else if payload.ReadType() != cmd {
		t.Errorf("Expected payload type to be %v, but got %v", cmd, payload.ReadType())
	}
	if payload.ReadData() != "my_data" {
		t.Errorf("Expected payload data to be %s, but got %s", "my_data", payload.ReadData())
	}
}
