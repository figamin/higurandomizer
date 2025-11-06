package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sqweek/dialog"
)

type menu int

const (
	mainMenu menu = iota
	spriteMenu
	characterMenu
	meiVariantMenu
)

var spriteChoices = []string{
	"mion", "ooishi", "rena", "rika", "satoko",
	"takano", "chie", "tomitake", "kasai", "shion",
	"irie", "akane", "keiichi", "satoshi", "teppei",
	"rina", "akasaka", "hanyuu", "oko", "kameda",
	"mo", "mura", "tamura", "une",
}

func generateVariants(n int) []string {
	var result []string
	for i := 1; i <= n; i++ {
		result = append(result, fmt.Sprintf("v%03d", i))
	}
	return result
}

var meiOptions = append([]string{
	"Best Match",
	"Random Outfits",
	"Random Outfits & Expressions",
}, generateVariants(60)...)

const itemsPerPage = 5

type model struct {
	currentMenu       menu
	cursor            int
	page              int
	selectedCharacter string

	filePath   string
	spritePath string
	message    string
	quitting   bool
}

func initialModel() model {
	return model{
		currentMenu: mainMenu,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m *model) move(limit int, up bool) {
	if up && m.cursor > 0 {
		m.cursor--
	}
	if !up && m.cursor < limit-1 {
		m.cursor++
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		key := msg.String()

		if key == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		switch m.currentMenu {

		case mainMenu:
			switch key {
			case "q":
				m.quitting = true
				return m, tea.Quit

			case "up", "k":
				m.move(3, true)
			case "down", "j":
				m.move(3, false)

			case "enter", " ":
				switch m.cursor {

				case 0:
					path, err := dialog.File().
						Title("Select Higurashi Episode (Ep01–Ep03)").
						Filter("Higurashi Episodes", "exe").
						Load()
					if err != nil {
						m.message = "No file selected or error occurred."
						return m, nil
					}
					base := filepath.Base(path)
					allowed := map[string]bool{
						"HigurashiEp01.exe": true,
						"HigurashiEp02.exe": true,
						"HigurashiEp03.exe": true,
					}
					if !allowed[base] {
						m.message = fmt.Sprintf("Invalid file selected: %s", base)
						return m, nil
					}

					m.filePath = path
					dir := filepath.Dir(path)
					dataFolder := filepath.Join(dir, base[:len(base)-4]+"_Data")
					m.spritePath = filepath.Join(dataFolder, "StreamingAssets", "CGAlt", "sprite")
					m.message = "Game selected."

				case 1:
					m.currentMenu = spriteMenu
					m.cursor = 0
					m.page = 0

				case 2:
					m.quitting = true
					return m, tea.Quit
				}
			}

		case spriteMenu:
			switch key {
			case "q":
				m.currentMenu = mainMenu

			case "up", "k":
				m.move(itemsPerPage, true)
			case "down", "j":
				m.move(itemsPerPage, false)
			case "left", "h":
				if m.page > 0 {
					m.page--
					m.cursor = 0
				}
			case "right", "l":
				if (m.page+1)*itemsPerPage < len(spriteChoices) {
					m.page++
					m.cursor = 0
				}
			case "enter", " ":
				idx := m.page*itemsPerPage + m.cursor
				if idx < len(spriteChoices) {
					m.selectedCharacter = spriteChoices[idx]
					m.currentMenu = characterMenu
					m.cursor = 0
				}
			}

		case characterMenu:
			switch key {
			case "q":
				m.currentMenu = spriteMenu
			case "up", "k":
				m.move(2, true)
			case "down", "j":
				m.move(2, false)
			case "enter", " ":
				if m.cursor == 0 {
					m.currentMenu = meiVariantMenu
					m.cursor = 0
					m.page = 0
				} else {
					m.message = fmt.Sprintf("Selected Ace Attorney for %s", m.selectedCharacter)
					m.currentMenu = spriteMenu
				}
			}

		case meiVariantMenu:
			visible := meiOptions[m.page*itemsPerPage:]
			if len(visible) > itemsPerPage {
				visible = visible[:itemsPerPage]
			}

			switch key {
			case "q":
				m.currentMenu = characterMenu
				m.cursor = 0
				m.page = 0

			case "up", "k":
				m.move(len(visible), true)
			case "down", "j":
				m.move(len(visible), false)

			case "left", "h":
				if m.page > 0 {
					m.page--
					m.cursor = 0
				}
			case "right", "l":
				if (m.page+1)*itemsPerPage < len(meiOptions) {
					m.page++
					m.cursor = 0
				}

			case "enter", " ":
				chosen := meiOptions[m.page*itemsPerPage+m.cursor]
				m.message = fmt.Sprintf("Selected %s → Mei → %s", m.selectedCharacter, chosen)
				m.currentMenu = spriteMenu
			}
		}
	}

	return m, nil
}

func (m model) View() string {

	if m.quitting {
		return "Goodbye!\n"
	}

	switch m.currentMenu {

	case mainMenu:
		return fmt.Sprintf(
			"Main Menu\n\n%s Select Game\n%s Select Sprites\n%s Exit\n\n%s\n",
			cursor(m.cursor, 0), cursor(m.cursor, 1), cursor(m.cursor, 2), m.message,
		)

	case spriteMenu:
		start := m.page * itemsPerPage
		end := start + itemsPerPage
		if end > len(spriteChoices) {
			end = len(spriteChoices)
		}

		s := fmt.Sprintf("Select Character (Page %d)\n\n", m.page+1)
		for i, name := range spriteChoices[start:end] {
			s += fmt.Sprintf("%s %s\n", cursor(m.cursor, i), name)
		}
		return s + "\nUse ↑↓ ←→ Enter, q to return.\n"

	case characterMenu:
		return fmt.Sprintf(
			"Character: %s\n\n%s Mei\n%s Ace Attorney\n\nUse ↑↓ Enter, q to return.\n",
			m.selectedCharacter, cursor(m.cursor, 0), cursor(m.cursor, 1),
		)

	case meiVariantMenu:
		start := m.page * itemsPerPage
		end := start + itemsPerPage
		if end > len(meiOptions) {
			end = len(meiOptions)
		}
		s := fmt.Sprintf("Mei Variant (%s) Page %d\n\n", m.selectedCharacter, m.page+1)
		for i, name := range meiOptions[start:end] {
			s += fmt.Sprintf("%s %s\n", cursor(m.cursor, i), name)
		}
		return s + "\nUse ↑↓ ←→ Enter, q to return.\n"
	}

	return ""
}

func cursor(cur, i int) string {
	if cur == i {
		return ">"
	}
	return " "
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
