package main

import (
	"bufio"
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
			Foreground(lipgloss.Color("#00FF00")).
			Background(lipgloss.Color("#1a1a1a")).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00"))
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
}

type tickMsg time.Time
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
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m.handleInput()
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
		// Connection successful, start interactive mode
		go m.handleConnection(msg.conn)
		return m, tea.Quit

	case tickMsg:
		return m, nil
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

	// Bidirectional communication
	go func() {
		io.Copy(os.Stdout, conn)
	}()

	io.Copy(conn, os.Stdin)
}

func (m model) View() string {
	var s strings.Builder

	// Banner
	banner := `
    _   __     __  __     
   / | / /__  / /_/ /_  __
  /  |/ / _ \/ __/ / / / /
 / /|  /  __/ /_/ / /_/ / 
/_/ |_/\___/\__/_/\__, /  
                 /____/   
`
	s.WriteString(titleStyle.Render(banner))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render(fmt.Sprintf("Modern Netcat Alternative v%s | by %s", version, author)))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("✗ Error: %v", m.err)))
		s.WriteString("\n\n")
		s.WriteString(promptStyle.Render("Press 'q' to quit"))
		return s.String()
	}

	if m.loading {
		s.WriteString(fmt.Sprintf("%s Processing...\n", m.spinner.View()))
		return s.String()
	}

	switch m.mode {
	case modeMenu:
		s.WriteString(infoStyle.Render("╔═══════════════════════════════════════╗"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("║         SELECT OPERATION MODE         ║"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("╚═══════════════════════════════════════╝"))
		s.WriteString("\n\n")

		s.WriteString("  " + successStyle.Render("1.") + " Listen Mode (Server)    - Accept incoming connections\n")
		s.WriteString("  " + successStyle.Render("2.") + " Connect Mode (Client)   - Connect to remote host\n")
		s.WriteString("  " + errorStyle.Render("3.") + " Exit\n\n")

		s.WriteString(promptStyle.Render("→ ") + m.textInput.View())

	case modeServer:
		s.WriteString(infoStyle.Render("╔═══════════════════════════════════════╗"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("║          LISTEN MODE (SERVER)         ║"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("╚═══════════════════════════════════════╝"))
		s.WriteString("\n\n")
		s.WriteString(promptStyle.Render("→ ") + m.textInput.View())

	case modeClient:
		s.WriteString(infoStyle.Render("╔═══════════════════════════════════════╗"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("║         CONNECT MODE (CLIENT)         ║"))
		s.WriteString("\n")
		s.WriteString(infoStyle.Render("╚═══════════════════════════════════════╝"))
		s.WriteString("\n\n")
		if m.currentStep == 2 {
			s.WriteString(fmt.Sprintf("  Target: %s\n\n", infoStyle.Render(m.host)))
		}
		s.WriteString(promptStyle.Render("→ ") + m.textInput.View())
	}

	s.WriteString("\n\n")
	s.WriteString(subtitleStyle.Render("Press 'Ctrl+C' or 'q' to quit"))

	return s.String()
}

var rootCmd = &cobra.Command{
	Use:   "netly",
	Short: "Modern and fast netcat alternative",
	Long: `Netly - A modern, fast, and beautiful netcat alternative built with Go
	
Created by penguinshero for ethical hacking and network testing.`,
	Version: version,
}

var listenCmd = &cobra.Command{
	Use:   "listen [port]",
	Short: "Start listening mode (server)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		port := args[0]
		runDirectServer(port)
	},
}

var connectCmd = &cobra.Command{
	Use:   "connect [host] [port]",
	Short: "Connect to remote host (client)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		host := args[0]
		port := args[1]
		runDirectClient(host, port)
	},
}

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start interactive mode with GUI",
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

	fmt.Printf("%s Listening on port %s\n", successStyle.Render("✓"), port)
	fmt.Printf("%s Waiting for connection...\n\n", infoStyle.Render("→"))

	conn, err := listener.Accept()
	if err != nil {
		fmt.Printf("%s Error: %v\n", errorStyle.Render("✗"), err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("%s Connection from %s\n\n", successStyle.Render("✓"), conn.RemoteAddr())

	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
}

func runDirectClient(host, port string) {
	fmt.Printf("\n%s Connecting to %s:%s...\n", infoStyle.Render("→"), host, port)

	conn, err := net.DialTimeout("tcp", host+":"+port, 10*time.Second)
	if err != nil {
		fmt.Printf("%s Error: %v\n", errorStyle.Render("✗"), err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("%s Connected successfully!\n\n", successStyle.Render("✓"))

	go io.Copy(os.Stdout, conn)
	io.Copy(conn, os.Stdin)
}

func init() {
	rootCmd.AddCommand(listenCmd)
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(interactiveCmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if len(os.Args) == 1 {
		// No arguments, start interactive mode
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

