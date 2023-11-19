package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

type Questions struct {
	question string
	answer   string
}

func NewQuestion(question string) Questions {
	return Questions{question: question}
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

type model struct {
	width       int
	height      int
	index       int
	styles      *Styles
	answerField textinput.Model
	questions   []Questions
}

func New(questions []Questions) *model {
	styles := DefaultStyles()
	answerField := textinput.New()
	answerField.Placeholder = "Your answer here"
	answerField.Focus()
	return &model{
		questions:   questions,
		answerField: answerField,
		styles:      styles,
	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	current := &m.questions[m.index]

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			current.answer = m.answerField.Value()
			m.answerField.SetValue("")
			m.Next()
			return m, nil
		}
	}
	m.answerField, cmd = m.answerField.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	if m.width == 0 {
		return "loading..."
	}
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			m.questions[m.index].question,
			m.styles.InputField.Render(m.answerField.View()),
		),
	)
}

func (m *model) Next() {
	if m.index < len(m.questions)-1 {
		m.index++
	} else {
		m.index = 0
	}
}

func main() {
	questions := []Questions{
		NewQuestion("Eet jy kaas?"),
		NewQuestion("Eet jy grond?"),
		NewQuestion("Eet jy gras?"),
	}
	m := New(questions)

	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// fmt.Println("Znox V2 - Enjoy!")

	// if len(os.Args) == 1 {
	// 	showHelp()
	// 	os.Exit(0)
	// }

	// enc := flag.NewFlagSet("enc", flag.ExitOnError)
	// enci := enc.String("i", "", "Provide an input file to encrypt.")
	// enco := enc.String("o", "", "Provide an output filename.")

	// dec := flag.NewFlagSet("dec", flag.ExitOnError)
	// deci := dec.String("i", "", "Provide an input file to decrypt.")
	// deco := dec.String("o", "", "Provide an output filename.")

	// pw := flag.NewFlagSet("pw", flag.ExitOnError)
	// pwsize := pw.Int("s", 15, "Generate password of given length.")

	// switch os.Args[1] {
	// case "enc":
	// 	if err := enc.Parse(os.Args[2:]); err != nil {
	// 		log.Println("Error when parsing arguments to enc")
	// 		panic(err)
	// 	}
	// 	if *enci == "" {
	// 		fmt.Println("Provide an input file to encrypt.")
	// 		os.Exit(1)
	// 	}
	// 	if *enco != "" {
	// 		znox.Encryption(*enci, *enco)
	// 	} else {
	// 		znox.Encryption(*enci, *enci+".enc")
	// 	}

	// case "dec":
	// 	if err := dec.Parse(os.Args[2:]); err != nil {
	// 		log.Println("Error when parsing arguments to dec")
	// 		panic(err)
	// 	}
	// 	if *deci == "" {
	// 		fmt.Println("Provide an input file to decrypt.")
	// 		os.Exit(1)
	// 	}
	// 	if *deco != "" {
	// 		znox.Decryption(*deci, *deco)
	// 	} else {
	// 		dd := *deci
	// 		o := "decrypted-" + *deci
	// 		if dd[len(dd)-4:] == ".enc" {
	// 			o = "decrypted-" + dd[:len(dd)-4]
	// 		}
	// 		znox.Decryption(*deci, o)
	// 	}

	// case "pw":
	// 	if err := pw.Parse(os.Args[2:]); err != nil {
	// 		log.Println("Error when parsing arguments to pw")
	// 		panic(err)
	// 	}
	// 	fmt.Println("Password :", znox.GeneratePassword(*pwsize))

	// default:
	// 	showHelp()
	// }

	// log.Println("Done!")
}

func showHelp() {
	fmt.Println("Example commands:")
	fmt.Println("Encrypt a file : crypto_demo enc -i plaintext.txt -o ciphertext.enc")
	fmt.Println("Decrypt a file : crypto_demo dec -i ciphertext.enc -o decrypted-plaintext.txt")
	fmt.Println("Generate a password : crypto_demo pw -s 15")
}
