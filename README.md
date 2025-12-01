# AirBridge

![AirBridgeMain.jpg](assets/AirBridgeMain.jpg)

![GitHub Release](https://img.shields.io/github/v/release/batuhan0sanli/AirBridge?color=blue)
![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/batuhan0sanli/AirBridge/release.yml)

**Secure, serverless file transfer right from your terminal.**

AirBridge is a CLI tool that allows you to securely transfer small files between computers using a text-based interface.
It converts your files into encrypted text payloads, enabling you to send sensitive data over any channel that supports
text (chat apps, email, pastebins, etc.) without trusting the intermediary.

## 🚀 Why AirBridge?

- **🔒 End-to-End Encryption:** Your files are encrypted locally before they ever leave your machine. Only the intended
  recipient can decrypt them.
- **🚫 Serverless:** No third-party servers, no clouds, no tracking. Just you and the recipient.
- **💻 Cross-Platform:** Works on macOS, Linux, and Windows.
- **✨ Beautiful TUI:** Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for a delightful interactive
  experience.

![AirBridgeAnimation.gif](assets/AirBridgeAnimation.gif)

## 🛠 Installation

###  macOS

**Using Homebrew (Recommended):**

```bash
brew tap batuhan0sanli/tap
brew install airbridge
```

**Manual Install:**

1. Download the `airbridge_Darwin_x86_64.tar.gz` (Intel) or `airbridge_Darwin_arm64.tar.gz` (Apple Silicon) from
   [Releases](https://github.com/batuhan0sanli/AirBridge/releases/).
2. Extract and move to your path:

```bash
tar -xvf airbridge_Darwin_arm64.tar.gz
sudo mv airbridge /usr/local/bin/
```

### ⊞ Windows

**Using Scoop (Recommended):**

```powershell
scoop bucket add airbridge https://github.com/batuhan0sanli/scoop-bucket
scoop install airbridge
```

**Manual Install:**

1. Download the `airbridge_Windows_x86_64.zip` from [Releases](https://github.com/batuhan0sanli/AirBridge/releases/).
2. Extract the zip file.
3. Open PowerShell or Command Prompt in that folder.
4. Run `./airbridge.exe`. (Optional: Add the folder to your System PATH environment variable to run it from anywhere)

### 🐧 Linux

**Using Homebrew (Linuxbrew):**

```bash
brew tap batuhan0sanli/tap
brew install airbridge
```

**Using DEB/RPM Packages:** Download the appropriate file from Releases and run:

- **Debian/Ubuntu:** `sudo dpkg -i airbridge_linux_amd64.deb`
- **Fedora/RHEL:** `sudo rpm -i airbridge_linux_amd64.rpm`

**Manual Install:**

1. Download the `airbridge_Linux_x86_64.tar.gz` from [Releases](https://github.com/batuhan0sanli/AirBridge/releases/).
2. Extract and move to your path:

```bash
tar -xvf airbridge_Linux_x86_64.tar.gz
sudo mv airbridge /usr/local/bin/
```

### 🏗 From Source (All Platforms)

If you prefer to build from source:

**Prerequisites:** [Go](https://go.dev/dl/) 1.21 or higher.

```bash
# Clone the repository
git clone https://github.com/batuhan0sanli/AirBridge.git
cd AirBridge

# Build the binary
make build

# (Optional) Move to your PATH
sudo mv airbridge /usr/local/bin/
```

## 📖 Usage

AirBridge has two main modes: **Send** and **Receive**.

### 📥 Receiving a File

1. Run the receive command:
   ```bash
   airbridge receive
   ```
2. AirBridge will generate a **Public Key**. Copy this key and send it to the sender.
3. Wait for the sender to give you the **Encrypted Payload**.
4. Paste the payload into the terminal.
5. The file will be decrypted and saved to your current directory.

### 📤 Sending a File

1. Run the send command:
   ```bash
   airbridge send [optional-file-path]
   ```
2. Paste the **Public Key** provided by the receiver.
3. Select the file you want to send (if you didn't provide a path).
4. AirBridge will generate an **Encrypted Payload**.
5. Copy this payload and send it to the receiver.

## 🔐 Technical Details

AirBridge uses a robust hybrid encryption scheme to ensure security:

1. **Key Exchange:** The receiver generates an ephemeral **RSA-2048** key pair.
2. **Symmetric Encryption:** The sender generates a random **AES-256** key and a random Initialization Vector (IV).
3. **Data Encryption:** The file is encrypted using **AES-256-CTR**.
4. **Key Encapsulation:** The AES key is encrypted with the receiver's RSA Public Key using **RSA-OAEP** (with SHA-256).
5. **Payload:** The encrypted AES key, IV, and encrypted file data are bundled into a JSON object and Base64 encoded for
   easy transport.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

Copyright © 2025 Batuhan Sanli.
