package tests

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary once for all tests
	if runtime.GOOS == "windows" {
		binaryPath = filepath.Join(os.TempDir(), "airbridge_test.exe")
	} else {
		binaryPath = filepath.Join(os.TempDir(), "airbridge_test")
	}

	// Build command
	// assuming tests/ is one level deep from root
	rootPath, _ := filepath.Abs("..")
	cmd := exec.Command("go", "build", "-o", binaryPath, rootPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to build binary: %v\nOutput: %s\n", err, output)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup binary
	_ = os.Remove(binaryPath)
	os.Exit(code)
}

func runCLI(dir string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestHeadlessHappyPath(t *testing.T) {
	// Setup Temp Dir
	tempDir, err := os.MkdirTemp("", "airbridge_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }() // CLEANUP

	// 1. Generate Keys
	_, err = runCLI(tempDir, "keygen", "-o", ".")
	if err != nil {
		t.Fatalf("Keygen failed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(tempDir, "private.pem")); os.IsNotExist(err) {
		t.Fatal("private.pem not created")
	}

	// 2. Create Dummy File
	testFileName := "testfile.dat"
	testContent := make([]byte, 1024*1024) // 1MB
	if _, err := rand.Read(testContent); err != nil {
		t.Fatalf("Failed to generate random content: %v", err)
	}
	err = os.WriteFile(filepath.Join(tempDir, testFileName), testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Calculate Hash of original
	hashOriginal := sha256.Sum256(testContent)

	// 3. Send (Headless)
	payloadName := "payload.abp"
	output, err := runCLI(tempDir, "send", testFileName, "-k", "public.pem", "-o", payloadName, "-H")
	if err != nil {
		t.Fatalf("Send failed: %v\nOutput: %s", err, output)
	}

	// 4. Receive (Headless)
	// Rename original file to avoid overwrite issues/confusion, or verify received file
	// The receive command uses the filename from metadata. Since we can't easily change the behavior
	// of where it writes (it writes to CWD with original name), let's create a sub-directory for receiver
	receiverDir := filepath.Join(tempDir, "receiver")
	if err := os.Mkdir(receiverDir, 0755); err != nil {
		t.Fatalf("Failed to create receiver dir: %v", err)
	}

	// Copy private key and payload to receiver dir
	copyFile(t, filepath.Join(tempDir, "private.pem"), filepath.Join(receiverDir, "private.pem"))
	copyFile(t, filepath.Join(tempDir, payloadName), filepath.Join(receiverDir, payloadName))

	output, err = runCLI(receiverDir, "receive", "-k", "private.pem", "-i", payloadName, "-H")
	if err != nil {
		t.Fatalf("Receive failed: %v\nOutput: %s", err, output)
	}

	// 5. Verify
	receivedContent, err := os.ReadFile(filepath.Join(receiverDir, testFileName))
	if err != nil {
		t.Fatalf("Failed to read received file: %v", err)
	}

	hashReceived := sha256.Sum256(receivedContent)
	if !bytes.Equal(hashOriginal[:], hashReceived[:]) {
		t.Fatal("Checksum mismatch! File corrupted during transfer.")
	}
}

func TestHeadlessDeletePayload(t *testing.T) {
	// Setup Temp Dir
	tempDir, err := os.MkdirTemp("", "airbridge_del_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Generate Keys & File (Simplified reuse)
	if _, err := runCLI(tempDir, "keygen", "-o", "."); err != nil {
		t.Fatalf("Keygen failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "small.txt"), []byte("restore"), 0644); err != nil {
		t.Fatalf("Failed to write small.txt: %v", err)
	}
	if _, err := runCLI(tempDir, "send", "small.txt", "-k", "public.pem", "-o", "payload.abp", "-H"); err != nil {
		t.Fatalf("Send failed: %v", err)
	}

	// Receive with -d
	output, err := runCLI(tempDir, "receive", "-k", "private.pem", "-i", "payload.abp", "-H", "-d")
	if err != nil {
		t.Fatalf("Receive -d failed: %v\nOutput: %s", err, output)
	}

	// Verify File Exists
	if _, err := os.Stat(filepath.Join(tempDir, "small.txt")); os.IsNotExist(err) {
		t.Fatal("Decrypted file not found")
	}

	// Verify Payload Deleted
	if _, err := os.Stat(filepath.Join(tempDir, "payload.abp")); !os.IsNotExist(err) {
		t.Fatal("Payload file was NOT deleted")
	}
}

func TestHeadlessErrorCases(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "airbridge_err_test_*")
	defer func() { _ = os.RemoveAll(tempDir) }()

	// Send without file
	_, err := runCLI(tempDir, "send", "-k", "missing.pem", "-H")
	if err == nil {
		t.Error("Send without file should fail")
	}

	// Receive without key
	_, err = runCLI(tempDir, "receive", "-i", "payload.abp", "-H")
	if err == nil {
		t.Error("Receive without key should fail")
	}
}

func copyFile(t *testing.T, src, dst string) {
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("Failed to read src %s: %v", src, err)
	}
	err = os.WriteFile(dst, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write dst %s: %v", dst, err)
	}
}
