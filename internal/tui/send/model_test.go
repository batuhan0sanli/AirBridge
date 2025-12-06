package send

import (
	"testing"
)

func TestInitialModel(t *testing.T) {
	// 1. Test with no arguments (Default)
	model1 := InitialModel("", "")
	if model1.rawPublicKey != "" {
		t.Error("Expected empty rawPublicKey")
	}
	if model1.selectedFile != "" {
		t.Error("Expected empty selectedFile")
	}

	// 2. Test with File and Public Key
	initialFile := "/path/to/file"
	initialKey := "some_public_key_string"

	model2 := InitialModel(initialFile, initialKey)
	if model2.selectedFile != initialFile {
		t.Errorf("Expected selectedFile %q, got %q", initialFile, model2.selectedFile)
	}
	if model2.rawPublicKey != initialKey {
		t.Errorf("Expected rawPublicKey %q, got %q", initialKey, model2.rawPublicKey)
	}
}
