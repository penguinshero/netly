package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	version = "1.0.0"
	author  = "penguinshero"

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF00"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Bold(true)

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")).
			Bold(true)

	highlightStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF00FF")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))
)

type mode int

const (
	modeMenu mode = iota
	modeServer
	modeClient
)

type model struct {
	mode        mode
	spinner     spinner.Model
	textInput   textinput.Model
	loading     bool
	message     string
	err         error
	host        string
	port        string
	currentStep int
	showHelp    bool
	width       int
	height      int
}

type connectionMsg struct {
	conn net.Conn
	err  error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	ti := textinput.New()
	ti.Placeholder = "Enter choice..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return model{
		mode:      modeMenu,
		spinner:   s,
		textInput: ti,
		loading:   false,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if !m.showHelp {
				return m, tea.Quit
			}
		case "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
			if m.mode != modeMenu && !m.loading {
				m.mode = modeMenu
				m.currentStep = 0
				m.textInput.SetValue("")
				m.textInput.Placeholder = "Enter choice..."
				m.err = nil
				return m, nil
			}
		case "h", "?":
			if m.mode == modeMenu && !m.loading {
				m.showHelp = !m.showHelp
				return m, nil
			}
		case "enter":
			if !m.showHelp {
				return m.handleInput()
			}
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case connectionMsg:
		if msg.err != nil {
			m.err = msg.err
			m.loading = false
			return m, nil
		}
		go m.handleConnection(msg.conn)
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) handleInput() (tea.Model, tea.Cmd) {
	value := strings.TrimSpace(m.textInput.Value())

	if m.mode == modeMenu {
		switch value {
		case "1":
			m.mode = modeServer
			m.currentStep = 1
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Enter port (e.g., 4444)..."
			return m, nil
		case "2":
			m.mode = modeClient
			m.currentStep = 1
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Enter host (e.g., 192.168.1.100)..."
			return m, nil
		case "3":
			return m, tea.Quit
		}
	}

	if m.mode == modeServer {
		if m.currentStep == 1 {
			m.port = value
			m.loading = true
			return m, m.startServer()
		}
	}

	if m.mode == modeClient {
		if m.currentStep == 1 {
			m.host = value
			m.currentStep = 2
			m.textInput.SetValue("")
			m.textInput.Placeholder = "Enter port (e.g., 4444)..."
			return m, nil
		} else if m.currentStep == 2 {
			m.port = value
			m.loading = true
			return m, m.startClient()
		}
	}

	return m, nil
}

func (m model) startServer() tea.Cmd {
	return func() tea.Msg {
		listener, err := net.Listen("tcp", ":"+m.port)
		if err != nil {
			return connectionMsg{nil, err}
		}

		fmt.Printf("\n%s Listening on port %s...\n", successStyle.Render("✓"), m.port)
		fmt.Printf("%s Waiting for connection...\n\n", infoStyle.Render("→"))

		conn, err := listener.Accept()
		if err != nil {
			listener.Close()
			return connectionMsg{nil, err}
		}

		fmt.Printf("%s Connection established from %s\n\n", successStyle.Render("✓"), conn.RemoteAddr())
		return connectionMsg{conn, nil}
	}
}

func (m model) startClient() tea.Cmd {
	return func() tea.Msg {
		fmt.Printf("\n%s Connecting to %s:%s...\n", infoStyle.Render("→"), m.host, m.port)

		conn, err := net.DialTimeout("tcp", m.host+":"+m.port, 10*time.Second)
		if err != nil {
			return connectionMsg{nil, err}
		}

		fmt.Printf("%s Connected successfully!\n\n", successStyle.Render("✓"))
		return connectionMsg{conn, nil}
	}
}

func (m model) handleConnection(conn net.Conn) {
	defer conn.Close()

	go func() {
		io.Copy(os.Stdout, conn)
	}()

	io.Copy(conn, os.Stdin)
}

func (m model) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	var s strings.Builder

	banner := `
  ╔╗╔╔═╗╔╦╗╦  ╦ ╦
  ║║║║╣  ║ ║  ╚╦╝
  ╝╚╝╚═╝ ╩ ╩═╝ ╩ `
	
	s.WriteString(titleStyle.Render(banner))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render(fmt.Sprintf("   Modern Netcat Alternative v%s | by %s", version, author)))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("✗ Error: %v", m.err)))
		s.WriteString("\n\n")
		s.WriteString(warningStyle.Render("Press 'ESC' to return to menu or 'q' to quit"))
		return s.String()
	}

	if m.loading {
		s.WriteString(fmt.Sprintf("%s Processing...\n", m.spinner.View()))
		return s.String()
	}

	switch m.mode {
	case modeMenu:
		s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("          SELECT OPERATION MODE"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		s.WriteString("\n\n")

		s.WriteString("  " + successStyle.Render("1.") + " " + highlightStyle.Render("Listen Mode") + dimStyle.Render(" (Server)") + "   - Accept incoming connections\n")
		s.WriteString("  " + successStyle.Render("2.") + " " + highlightStyle.Render("Connect Mode") + dimStyle.Render(" (Client)") + " - Connect to remote host\n")
		s.WriteString("  " + errorStyle.Render("3.") + " " + errorStyle.Render("Exit") + "\n\n")

		s.WriteString(promptStyle.Render("→ ") + m.textInput.View())
		s.WriteString("\n\n")
		s.WriteString(dimStyle.Render("Press 'h' or '?' for help | 'Ctrl+C' or 'q' to quit"))

	case modeServer:
		s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		s.WriteString("\n")
		s.WriteString(successStyle.Render("          LISTEN MODE") + dimStyle.Render(" (SERVER)"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		s.WriteString("\n\n")
		s.WriteString(dimStyle.Render("  Will listen for incoming TCP connections\n"))
		s.WriteString(dimStyle.Render("  Example: 4444, 8080, 9999\n\n"))
		s.WriteString(promptStyle.Render("→ ") + m.textInput.View())
		s.WriteString("\n\n")
		s.WriteString(dimStyle.Render("Press 'ESC' for menu | 'Ctrl+C' to quit"))

	case modeClient:
		s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		s.WriteString("\n")
		s.WriteString(successStyle.Render("          CONNECT MODE") + dimStyle.Render(" (CLIENT)"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"))
		s.WriteString("\n\n")
		if m.currentStep == 1 {
			s.WriteString(dimStyle.Render("  Enter target host IP or domain\n"))
			s.WriteString(dimStyle.Render("  Example: 192.168.1.100, example.com\n\n"))
		} else if m.currentStep == 2 {
			s.WriteString(fmt.Sprintf("  Target: %s\n", highlightStyle.Render(m.host)))
			s.WriteString(dimStyle.Render("  Enter target port\n"))
			s.WriteString(dimStyle.Render("  Example: 4444, 8080, 9999\n\n"))
		}
		s.WriteString(promptStyle.Render("→ ") + m.textInput.View())
		s.WriteString("\n\n")
		s.WriteString(dimStyle.Render("Press 'ESC' for menu | 'Ctrl+C' to quit"))
	}

	return s.String()
}

func (m model) renderHelp() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("\n  NETLY - HELP & USAGE GUIDE\n"))
	s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n"))

	s.WriteString(highlightStyle.Render("📖 INTERACTIVE MODE COMMANDS:\n\n"))
	
	s.WriteString(successStyle.Render("  1-3      ") + "Select menu option\n")
	s.WriteString(successStyle.Render("  h / ?    ") + "Show this help screen\n")
	s.WriteString(successStyle.Render("  ESC      ") + "Return to main menu\n")
	s.WriteString(successStyle.Render("  Ctrl+C/q ") + "Quit application\n\n")

	s.WriteString(highlightStyle.Render("💻 CLI MODE COMMANDS:\n\n"))
	
	s.WriteString(successStyle.Render("  netly listen [port]\n"))
	s.WriteString(dimStyle.Render("    Start server on specified port\n"))
	s.WriteString(dimStyle.Render("    Example: netly listen 4444\n\n"))
	
	s.WriteString(successStyle.Render("  netly connect [host] [port]\n"))
	s.WriteString(dimStyle.Render("    Connect to remote host\n"))
	s.WriteString(dimStyle.Render("    Example: netly connect 192.168.1.100 4444\n\n"))

	s.WriteString(successStyle.Render("  netly interactive\n"))
	s.WriteString(dimStyle.Render("    Start GUI mode (default)\n\n"))

	s.WriteString(highlightStyle.Render("🔧 COMMON USE CASES:\n\n"))
	
	s.WriteString(warningStyle.Render("  File Transfer:\n"))
	s.WriteString(dimStyle.Render("    Receiver: netly listen 4444 > file.txt\n"))
	s.WriteString(dimStyle.Render("    Sender:   cat file.txt | netly connect IP 4444\n\n"))
	
	s.WriteString(warningStyle.Render("  Chat Server:\n"))
	s.WriteString(dimStyle.Render("    Server:   netly listen 8080\n"))
	s.WriteString(dimStyle.Render("    Client:   netly connect server-ip 8080\n\n"))

	s.WriteString(warningStyle.Render("  Port Testing:\n"))
	s.WriteString(dimStyle.Render("    Test:     echo \"test\" | netly connect target-ip 80\n\n"))

	s.WriteString(infoStyle.Render("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"))
	s.WriteString(dimStyle.Render("\nPress 'ESC' to return to menu\n"))

	return s.String()
}

var rootCmd = &cobra.Command{
	Use:   "netly",
	Short: "Modern and fast netcat alternative",
	Long: `Netly - A modern, fast, and beautiful netcat alternative built with Go

Created by penguinshero for ethical hacking and network testing.
Netly provides an intuitive interface for TCP connections with both
interactive GUI and direct command-line modes.

Features:
  • Fast and efficient TCP connections
  • Beautiful terminal UI powered by Bubble Tea
  • Server (Listen) and Client (Connect) modes
  • Interactive and direct command modes
  • Cross-platform support`,
	Version: version,
}

var listenCmd = &cobra.Command{
	Use:   "listen [port]",
	Short: "Start listening mode (server)",
	Long: `Start a TCP server that listens on the specified port.
Once a connection is established, you can send and receive data.

Example:
  netly listen 4444`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]
		runDirectServer(port)
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect [host] [port]",
	Short: "Connect to remote host (client)",
	Long: `Connect to a remote TCP server using the specified host and port.
Once connected, you can send and receive data.

Example:
  netly connect 192.168.1.100 4444
  netly connect example.com 8080`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port := args[1]
		runDirectClient(host, port)
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive mode with GUI",
	Long: `Launch Netly in interactive mode with a beautiful terminal GUI.
This mode provides a menu-driven interface to easily switch between
server and client modes.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runDirectServer(port string) {
	fmt.Printf("\n%s Starting server on port %s...\n", infoStyle.Render("→"), port)

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("%s Error: %v\n", errorStyle.Render("✗"), err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("%s Listening on 0.0.0.0:%s\n", successStyle.Render("✓"), port)
	fmt.Printf("%s Waiting for connection...\n", infoStyle.Render("⏳"))
	fmt.Printf("%s Press Ctrl+C to stop\n\n", dimStyle.Render("ℹ"))

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("%s Error: %v\n", errorStyle.Render("✗"), err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("%s Connection established!\n", successStyle.Render("✓"))
	fmt.Printf("%s Remote: %s\n", infoStyle.Render("→"), conn.RemoteAddr())
	fmt.Printf("%s Local:  %s\n\n", infoStyle.Render("→"), conn.LocalAddr())
	fmt.Printf("%s Session started. Type to send data...\n\n", highlightStyle.Render("━━━"))

	done := make(chan bool)

	// Remote to local
	go func() {
		io.Copy(os.Stdout, conn)
		done <- true
	}()

	// Local to remote
	go func() {
		io.Copy(conn, os.Stdin)
		done <- true
	}()

	<-done
	fmt.Printf("\n%s Connection closed.\n", warningStyle.Render("✗"))
}

func runDirectClient(host, port string) {
	fmt.Printf("\n%s Connecting to %s:%s...\n", infoStyle.Render("→"), host, port)

	conn, err := net.DialTimeout("tcp", host+":"+port, 10*time.Second)
	if err != nil {
		fmt.Printf("%s Connection failed: %v\n", errorStyle.Render("✗"), err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("%s Connected successfully!\n", successStyle.Render("✓"))
	fmt.Printf("%s Remote: %s\n", infoStyle.Render("→"), conn.RemoteAddr())
	fmt.Printf("%s Local:  %s\n\n", infoStyle.Render("→"), conn.LocalAddr())
	fmt.Printf("%s Session started. Type to send data...\n\n", highlightStyle.Render("━━━"))

	done := make(chan bool)

	// Remote to local
	go func() {
		io.Copy(os.Stdout, conn)
		done <- true
	}()

	// Local to remote
	go func() {
		io.Copy(conn, os.Stdin)
		done <- true
	}()

	<-done
	fmt.Printf("\n%s Connection closed.\n", warningStyle.Render("✗"))
}

func init() {
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if len(os.Args) == 1 {
		p := tea.NewProgram(initialModel())
		if _, err := p.Run(); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
