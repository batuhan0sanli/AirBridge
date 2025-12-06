package pkg

import (
	"encoding/json"
	"testing"
)

func TestSmallFilePayload_JSON(t *testing.T) {
	payload := SmallFilePayload{
		Key:   "testkey",
		Data:  "testdata",
		Nonce: "testnonce",
		Metadata: FileMetadata{
			Name: "test.txt",
			Size: 123,
			Hash: "testhash",
		},
	}

	// Marshal
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal
	var decoded SmallFilePayload
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Key != payload.Key {
		t.Errorf("Expected Key %s, got %s", payload.Key, decoded.Key)
	}
	if decoded.Data != payload.Data {
		t.Errorf("Expected Data %s, got %s", payload.Data, decoded.Data)
	}
	if decoded.Nonce != payload.Nonce {
		t.Errorf("Expected Nonce %s, got %s", payload.Nonce, decoded.Nonce)
	}
	if decoded.Metadata.Name != payload.Metadata.Name {
		t.Errorf("Expected Metadata Name %s, got %s", payload.Metadata.Name, decoded.Metadata.Name)
	}
}

func TestLargeFilePayload_JSON(t *testing.T) {
	payload := LargeFilePayload{
		Key:      "largekey",
		Name:     "large.iso",
		Hash:     "largehash",
		ChunkNum: 5,
	}

	// Marshal
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal
	var decoded LargeFilePayload
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.Key != payload.Key {
		t.Errorf("Expected Key %s, got %s", payload.Key, decoded.Key)
	}
	if decoded.Name != payload.Name {
		t.Errorf("Expected Name %s, got %s", payload.Name, decoded.Name)
	}
	if decoded.Hash != payload.Hash {
		t.Errorf("Expected Hash %s, got %s", payload.Hash, decoded.Hash)
	}
	if decoded.ChunkNum != payload.ChunkNum {
		t.Errorf("Expected ChunkNum %d, got %d", payload.ChunkNum, decoded.ChunkNum)
	}
}
