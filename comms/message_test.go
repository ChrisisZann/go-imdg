package comms

// import (
// 	"testing"
// )

// func TestNewMessage(t *testing.T) {
// 	source, err := NewNodeAddr("tcp", "localhost:3333")
// 	if err != nil {
// 		t.Fatalf("Failed to create source NodeAddr: %v", err)
// 	}
// 	dest, err := NewNodeAddr("tcp", "localhost:3334")
// 	if err != nil {
// 		t.Fatalf("Failed to create destination NodeAddr: %v", err)
// 	}
// 	payload, err := NewPayload("cmd", "test")
// 	if err != nil {
// 		t.Fatalf("Failed to create Payload: %v", err)
// 	}

// 	msg := NewMessage(source, 1, dest, payload)

// 	if msg.source != source {
// 		t.Errorf("Expected source %v, got %v", source, msg.source)
// 	}
// 	if msg.suid != 1 {
// 		t.Errorf("Expected suid 1, got %d", msg.suid)
// 	}
// 	if msg.destination != dest {
// 		t.Errorf("Expected destination %v, got %v", dest, msg.destination)
// 	}
// 	if msg.payload != payload {
// 		t.Errorf("Expected payload %v, got %v", payload, msg.payload)
// 	}
// }

// func TestParseMessage_ValidInput(t *testing.T) {
// 	input := "localhost:3333|1|localhost:3334|cmd:test"
// 	msg, err := ParseMessage(input)

// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	if msg.source.String() != "localhost:3333" {
// 		t.Errorf("Expected source localhost:3333, got %s", msg.source.String())
// 	}
// 	if msg.suid != 1 {
// 		t.Errorf("Expected suid 1, got %d", msg.suid)
// 	}
// 	if msg.destination.String() != "localhost:3334" {
// 		t.Errorf("Expected destination localhost:3334, got %s", msg.destination.String())
// 	}
// 	if msg.payload.ReadData() != "test" {
// 		t.Errorf("Expected payload data 'test', got %s", msg.payload.ReadData())
// 	}
// }

// func TestParseMessage_InvalidInput(t *testing.T) {
// 	invalidInputs := []string{
// 		"localhost:3333|1|localhost:3334",            // Missing payload
// 		"localhost:3333|abc|localhost:3334|cmd:test", // Invalid suid
// 		"localhost:3333|1|localhost:3334|",           // Empty payload
// 	}

// 	for _, input := range invalidInputs {
// 		_, err := ParseMessage(input)
// 		if err == nil {
// 			t.Errorf("Expected error for input: %s", input)
// 		}
// 	}
// }

// func TestMessage_Compile(t *testing.T) {
// 	source := NodeAddr{Protocol: "tcp", Address: "localhost:3333"}
// 	dest := NodeAddr{Protocol: "tcp", Address: "localhost:3334"}
// 	payload := &Payload{ptype: PayloadType("cmd"), data: "test"}

// 	msg := NewMessage(source, 1, dest, payload)
// 	compiled, err := msg.Compile()

// 	if err != nil {
// 		t.Fatalf("Unexpected error: %v", err)
// 	}

// 	expected := "localhost:3333|1|localhost:3334|cmd:test"
// 	if compiled != expected {
// 		t.Errorf("Expected compiled message '%s', got '%s'", expected, compiled)
// 	}
// }

// func TestMessage_String(t *testing.T) {
// 	source := NodeAddr{Protocol: "tcp", Address: "localhost:3333"}
// 	dest := NodeAddr{Protocol: "tcp", Address: "localhost:3334"}
// 	payload := &Payload{ptype: PayloadType("cmd"), data: "test"}

// 	msg := NewMessage(source, 1, dest, payload)
// 	str := msg.String()

// 	expected := "source=localhost:3333, suid=1, payload=cmd:test"
// 	if str != expected {
// 		t.Errorf("Expected string representation '%s', got '%s'", expected, str)
// 	}
// }
