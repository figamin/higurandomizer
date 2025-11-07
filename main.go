package main

import (
	"fmt"
	"math/rand"
	"log"
	"os"
	"path/filepath"
	"encoding/json"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sqweek/dialog"
)

type Config struct {
	GamePath   string            `json:"game_path"`
	SpritePath string            `json:"sprite_path"`
	Selections map[string]string `json:"selections"`
}
func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		// no config → return defaults
		return Config{Selections: make(map[string]string)}
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{Selections: make(map[string]string)}
	}
	return cfg
}

func saveConfig(cfg Config) {
	f, err := os.Create("config.json")
	if err != nil {
		log.Printf("Failed to write config: %v\n", err)
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ") // pretty-print
	enc.Encode(cfg)
}

type menu int

const (
	mainMenu menu = iota
	spriteMenu
	characterMenu
	meiVariantMenu
	checkSelectionsMenu
)

var spriteChoices = []string{
	"mion", "ooishi", "rena", "rika", "satoko",
	"takano", "chie", "tomitake", "kasai", "shion",
	"irie", "akane", "keiichi", "satoshi", "teppei",
	"rina", "akasaka", "hanyuu", "oko", "kameda",
	"mo", "mura", "tamura", "une",
}

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
	meiOptions []string

	selections map[string]string // new: character → selected option
}

// cursor helper
func cursor(cur, i int) string {
	if cur == i {
		return ">"
	}
	return " "
}

// loadMeiOptions builds the options for a character
func loadMeiOptions(charKey string) []string {
	data, ok := Characters[charKey]
	if !ok {
		// fallback
		return []string{"Best Match", "Random Outfits", "Random Outfits & Expressions"}
	}

	opts := []string{
		"Best Match",
		"Random Outfits",
		"Random Outfits & Expressions",
	}

	for _, o := range data.OutfitsMei {
		opts = append(opts, o.Name)
	}
	return opts
}

func initialModel() model {
	cfg := loadConfig()

	// Fill selections with defaults where missing
	if cfg.Selections == nil {
		cfg.Selections = make(map[string]string)
	}
	for _, c := range spriteChoices {
		if _, ok := cfg.Selections[c]; !ok {
			cfg.Selections[c] = "Best Match"
		}
	}

	return model{
		currentMenu: mainMenu,
		filePath:    cfg.GamePath,
		spritePath:  cfg.SpritePath,
		selections:  cfg.Selections,
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
				m.move(4, true)
			case "down", "j":
				m.move(4, false)
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
					saveConfig(Config{
						GamePath:   m.filePath,
						SpritePath: m.spritePath,
						Selections: m.selections,
					})

					m.message = "Game selected."
				case 1:
					m.currentMenu = spriteMenu
					m.cursor = 0
					m.page = 0
				case 2:
					m.currentMenu = checkSelectionsMenu
					m.cursor = 0
				case 3:
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
					m.meiOptions = loadMeiOptions(m.selectedCharacter)
					m.currentMenu = meiVariantMenu
					m.cursor = 0
					m.page = 0
				} else {
					m.message = fmt.Sprintf("Selected Ace Attorney for %s", m.selectedCharacter)
					m.currentMenu = spriteMenu
				}
			}

		case meiVariantMenu:
			//log.Printf("Mei options length 1 = %d\n", len(m.meiOptions))
			if len(m.meiOptions) == 0 {
				break // prevent panic if somehow empty
			}
			maxPage := (len(m.meiOptions) - 1) / itemsPerPage
			if m.page > maxPage {
				m.page = maxPage
			}

			// Correct cursor if out of bounds
			start := m.page * itemsPerPage
			end := start + itemsPerPage
			if end > len(m.meiOptions) {
				end = len(m.meiOptions)
			}
			visible := m.meiOptions[start:end]
			if m.cursor >= len(visible) {
				m.cursor = len(visible) - 1
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
				if (m.page+1)*itemsPerPage < len(m.meiOptions) {
					m.page++
					m.cursor = 0
				}
			case "enter", " ":
				//log.Printf("Mei options length 2 = %d\n", len(m.meiOptions))
				if len(m.meiOptions) == 0 {
					break // safety check
				}
				if m.page*itemsPerPage >= len(m.meiOptions) {
					m.page = len(m.meiOptions) / itemsPerPage
					if m.page*itemsPerPage >= len(m.meiOptions) {
						m.page = 0
					}
				}

				idx := m.page*itemsPerPage + m.cursor
				if idx >= len(m.meiOptions) {
					idx = len(m.meiOptions) - 1
				}
				chosen := m.meiOptions[idx]

				// Determine variant to use
				var variant string
				switch chosen {
				case "Best Match":
					variant = spriteSets[0]
				case "Random Outfits":
					data := Characters[m.selectedCharacter]
					if len(data.OutfitsMei) > 0 {
						o := data.OutfitsMei[rand.Intn(len(data.OutfitsMei))]
						chosen = o.Name       // update display name
						variant = o.SpriteSet // store sprite set
					} else {
						variant = spriteSets[0]
					}
				case "Random Outfits & Expressions":
					data := Characters[m.selectedCharacter]
					if len(data.OutfitsMei) > 0 {
						o := data.OutfitsMei[rand.Intn(len(data.OutfitsMei))]
						chosen = o.Name
						variant = o.SpriteSet
					} else {
						variant = spriteSets[0]
					}
				default:
					data := Characters[m.selectedCharacter]
					for _, o := range data.OutfitsMei {
						if o.Name == chosen {
							variant = o.SpriteSet
							break
						}
					}
				}

				// Save the display name + variant in selections
				m.selections[m.selectedCharacter] = fmt.Sprintf("%s (variant: %s)", chosen, variant)
				saveConfig(Config{
					GamePath:   m.filePath,
					SpritePath: m.spritePath,
					Selections: m.selections,
				})
				m.message = fmt.Sprintf("Selected %s → Mei → %s", m.selectedCharacter, m.selections[m.selectedCharacter])
				m.currentMenu = spriteMenu
			}
		case checkSelectionsMenu:
    total := len(spriteChoices)
    maxPage := (total - 1) / itemsPerPage

    switch key {
		case "q":
			m.currentMenu = mainMenu
			m.cursor = 0
			m.page = 0

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
			if m.page < maxPage {
				m.page++
				m.cursor = 0
			}
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
			"Main Menu\n\n%s Select Game\n%s Select Sprites\n%s Check Selections\n%s Exit\n\n%s\n",
			cursor(m.cursor, 0), cursor(m.cursor, 1), cursor(m.cursor, 2), cursor(m.cursor, 3), m.message,
		)

	case spriteMenu:
		// Ensure page is within valid range for spriteChoices
		if len(spriteChoices) == 0 {
			return "No characters available.\n"
		}
		maxPage := (len(spriteChoices)-1) / itemsPerPage
		if m.page < 0 {
			m.page = 0
		}
		if m.page > maxPage {
			m.page = maxPage
		}

		start := m.page * itemsPerPage
		end := start + itemsPerPage
		if end > len(spriteChoices) {
			end = len(spriteChoices)
		}

		// Correct cursor so it never points outside the visible slice
		visibleCount := end - start
		if m.cursor < 0 {
			m.cursor = 0
		}
		if m.cursor >= visibleCount {
			m.cursor = visibleCount - 1
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
		//log.Printf("Mei options length 3 = %d\n", len(m.meiOptions))
		if len(m.meiOptions) == 0 {
			return "No options available.\n"
		}

		start := m.page * itemsPerPage
		if start >= len(m.meiOptions) {
			start = (len(m.meiOptions) / itemsPerPage) * itemsPerPage
			if start >= len(m.meiOptions) {
				start = 0
			}
			m.page = start / itemsPerPage
		}

		end := start + itemsPerPage
		if end > len(m.meiOptions) {
			end = len(m.meiOptions)
		}

		s := fmt.Sprintf("Mei Variant (%s) Page %d\n\n", m.selectedCharacter, m.page+1)
		for i, name := range m.meiOptions[start:end] {
			s += fmt.Sprintf("%s %s\n", cursor(m.cursor, i), name)
		}
		return s + "\nUse ↑↓ ←→ Enter, q to return.\n"

	case checkSelectionsMenu:
		start := m.page * itemsPerPage
		end := start + itemsPerPage
		if end > len(spriteChoices) {
			end = len(spriteChoices)
		}

		s := fmt.Sprintf("Current Selections (Page %d)\n\n", m.page+1)
		for i, c := range spriteChoices[start:end] {
			selection := m.selections[c]
			s += fmt.Sprintf("%s %s → %s\n", cursor(m.cursor, i), c, selection)
		}
		return s + "\nUse ↑↓ ←→ q to return.\n"

	}

	return ""
}

func main() {
	// Create or append to debug log file
/*
logFile, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
if err != nil {
    panic(err)
}
log.SetOutput(logFile)
log.SetFlags(log.LstdFlags | log.Lshortfile) // Include date/time + file/line number
*/

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
