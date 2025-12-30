package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sqweek/dialog"
)
var SelectedVariants map[string]string

type Config struct {
	GamePath   string            `json:"game_path"`
	SpritePath string            `json:"sprite_path"`
	Selections map[string]string `json:"selections"`
}
type progressWriter struct {
	total     int64
	downloaded int64
	lastPrint int64
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n := len(b)
	p.downloaded += int64(n)

	// throttle redraws
	if p.downloaded-p.lastPrint > p.total/100 || p.downloaded == p.total {
		p.lastPrint = p.downloaded
		p.print()
	}
	return n, nil
}

func (p *progressWriter) print() {
	const barWidth = 40

	var percent float64
	if p.total > 0 {
		percent = float64(p.downloaded) / float64(p.total)
	}
	filled := int(percent * barWidth)

	fmt.Printf(
		"\r[%-*s] %3.0f%%",
		barWidth,
		strings.Repeat("=", filled),
		percent*100,
	)
}

func ensureSprites() {
	spritesDir := "sprites"
	meiDir := filepath.Join(spritesDir, "mei")

	if _, err := os.Stat(spritesDir); os.IsNotExist(err) {
		log.Println("[INFO] 'sprites' folder not found, downloading ZIPs...")

		urls := []string{
			"https://github.com/figamin/higurandomizer-assets/releases/download/mei1/mei_part_1.zip",
			"https://github.com/figamin/higurandomizer-assets/releases/download/mei1/mei_part_2.zip",
			"https://github.com/figamin/higurandomizer-assets/releases/download/mei1/mei_part_3.zip",
		}

		tmpDir := "tmp_sprites_download"
		os.MkdirAll(tmpDir, os.ModePerm)
		defer os.RemoveAll(tmpDir)

		os.MkdirAll(meiDir, os.ModePerm)

		for i, url := range urls {
			tmpFile := filepath.Join(tmpDir, "sprites_part_"+strconv.Itoa(i)+".zip")
			log.Printf("[INFO] Downloading %s\n", url)

			resp, err := http.Get(url)
			if err != nil {
				log.Fatalf("Failed to download %s: %v", url, err)
			}
			if resp.StatusCode != http.StatusOK {
				log.Fatalf("Failed to download %s: server returned %d", url, resp.StatusCode)
			}

			out, err := os.Create(tmpFile)
			if err != nil {
				log.Fatalf("Failed to create temp file: %v", err)
			}

			pw := &progressWriter{
				total: resp.ContentLength,
			}

			_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
			fmt.Print("\n") // finish progress line

			out.Close()
			resp.Body.Close()

			if err != nil {
				log.Fatalf("Failed to save %s: %v", url, err)
			}

			if err := unzip(tmpFile, meiDir); err != nil {
				log.Fatalf("Failed to extract %s: %v", tmpFile, err)
			}

			log.Printf("[INFO] Extracted %s\n", url)
		}

		log.Println("[INFO] All sprites downloaded and merged successfully.")
	} else {
		log.Println("[INFO] 'sprites' folder exists")
	}
}


// unzip extracts a ZIP file into a destination folder
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}


func extractVariant(selection string) string {
    if selection == "" || strings.ToLower(selection) == "best match" {
        return ""
    }
    return selection
}
func getVariantForKey(key string, selectedVariants map[string]string) string {
    folder := GetFolder(key)
    log.Printf("[DEBUG] Key: %s, Folder: %s", key, folder)

    sel, ok := selectedVariants[folder]
    if !ok {
        log.Printf("[DEBUG] No selection found for folder '%s', falling back to default variant", folder)
        return RawGameSprites[key][1]
    }

    log.Printf("[DEBUG] Selection for folder '%s': %s", folder, sel)

    if sel == "" || strings.ToLower(sel) == "best match" {
        log.Printf("[DEBUG] Selection is empty or Best Match, using default variant: %s", RawGameSprites[key][1])
        return RawGameSprites[key][1]
    }

    v := extractVariant(sel)
    if v == "" {
        log.Printf("[DEBUG] Could not extract variant from selection, using default variant: %s", RawGameSprites[key][1])
        return RawGameSprites[key][1]
    }

    log.Printf("[DEBUG] Using variant from selection: %s", v)
    return v
}


func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		return Config{Selections: make(map[string]string)}
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return Config{Selections: make(map[string]string)}
	}
	SelectedVariants = make(map[string]string)
    for character, selection := range cfg.Selections {
        SelectedVariants[character] = extractVariant(selection)
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
	enc.SetIndent("", "  ")
	enc.Encode(cfg)
}

type menu int

const (
	mainMenu menu = iota
	episodeMenu
	manualEpisodeMenu
	spriteMenu
	characterMenu
	meiVariantMenu
	checkSelectionsMenu
	helpMenu
)

var spriteChoices = []string{
	"rena", "mion", "shion", "satoko", 
	"rika", "keiichi", "hanyuu", "satoshi", "chie",
	"ooishi", "tomitake", "takano", "irie", "kasai", 
	"akasaka", "teppei", "akane", "rina", "oko",
	"kameda", "mo", "mura", "tamura", "une",
}

const itemsPerPage = 5

type model struct {
	currentMenu       menu
	cursor            int
	page              int
	selectedCharacter string
	selectedEpisode	  HigurashiEpisode

	filePath   string
	spritePath string
	message    string
	quitting   bool
	meiOptions []string

	selections map[string]string 
}
type HigurashiEpisode int

const (
	Ep01 HigurashiEpisode = iota + 1
	Ep02
	Ep03
	Ep04
	Ep05
	Ep06
	Ep07
	Ep08
	Ep09
	Ep10
)

func (e HigurashiEpisode) ExeName() string {
	return fmt.Sprintf("HigurashiEp%02d", e)
}
func (e HigurashiEpisode) Label() string {
	switch e {
	case Ep01:
		return "Episode 1 — Onikakushi"
	case Ep02:
		return "Episode 2 — Watanagashi"
	case Ep03:
		return "Episode 3 — Tatarigoroshi"
	case Ep04:
		return "Episode 4 — Himatsubushi"
	case Ep05:
		return "Episode 5 — Meakashi"
	case Ep06:
		return "Episode 6 — Tsumihoroboshi"
	case Ep07:
		return "Episode 7 — Minagoroshi"
	case Ep08:
		return "Episode 8 — Matsuribayashi"
	case Ep09:
		return "Episode 9 — Higurashi Rei"
	case Ep10:
		return "Episode 10 — Higurashi Hou+"
	default:
		return string(e)
	}
}

 
func cursor(cur, i int) string {
	if cur == i {
		return ">"
	}
	return " "
}

func loadMeiOptions(charKey string) []string {
	data, ok := Characters[charKey]
	if !ok {
		// fallback
		return []string{"Best Match", "Random Outfits", "Random Outfits & Expressions", "Fully Random",}
	}

	opts := []string{
		"Best Match",
		"Random Outfits",
		"Random Outfits & Expressions",
		"Fully Random",
	}

	for _, o := range data.OutfitsMei {
		opts = append(opts, o.Name)
	}
	return opts
}

func initialModel() model {
	cfg := loadConfig()

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


func (m model) Init() tea.Cmd { return tea.ClearScreen }

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
				m.move(8, true)
			case "down", "j":
				m.move(8, false)
			case "enter", " ":
				switch m.cursor {
				case 0:
					m.currentMenu = episodeMenu
    				m.cursor = int(m.selectedEpisode) - 1
    				return m, nil
				case 1:
					m.currentMenu = manualEpisodeMenu
    				m.cursor = int(m.selectedEpisode) - 1
    				return m, nil
				case 2:
					m.currentMenu = spriteMenu
					m.cursor = 0
					m.page = 0
				case 3:
					m.currentMenu = checkSelectionsMenu
					m.cursor = 0
				case 4:
					return m.randomizeSprites()
				case 5:
    				return m.restoreOriginalSprites()
				case 6:
					m.currentMenu = helpMenu
					m.cursor = 0
				case 7:
					m.quitting = true
					return m, tea.Quit
				}
			}
		case episodeMenu:
	switch key {
	case "q":
		m.currentMenu = mainMenu

	case "up", "k":
		m.move(10, true)

	case "down", "j":
		m.move(10, false)

	case "enter", " ":
		m.selectedEpisode = HigurashiEpisode(m.cursor + 1)

		// Try autodetect
		path, err := autodetectHigurashi(m.selectedEpisode)
		if err != nil {
			path, err = dialog.File().
				Title("Select Higurashi executable").
				Filter("Higurashi0X.exe").
				Load()
			if err != nil {
				m.message = "No file selected."
				m.currentMenu = mainMenu
				return m, nil
			}
		}

		m.filePath = path

		base := filepath.Base(path)
		dir := filepath.Dir(path)
		dataFolder := filepath.Join(
			dir,
			strings.TrimSuffix(base, filepath.Ext(base)),
		)

		m.spritePath = filepath.Join(
			dataFolder,
			m.selectedEpisode.ExeName() + "_Data",
			"StreamingAssets",
			"CGAlt",
			"sprite",
		)

		saveConfig(Config{
			GamePath:   m.filePath,
			SpritePath: m.spritePath,
			Selections: m.selections,
		})

		m.message = fmt.Sprintf("Selected %s", m.selectedEpisode.Label())
		m.currentMenu = mainMenu
	}
case manualEpisodeMenu:
	switch key {
	case "q":
		m.currentMenu = mainMenu

	case "up", "k":
		m.move(10, true)

	case "down", "j":
		m.move(10, false)

	case "enter", " ":
		m.selectedEpisode = HigurashiEpisode(m.cursor + 1)

			path, err := dialog.File().
				Title("Select Higurashi executable").
				Filter("Higurashi Episodes").
				Load()
			if err != nil {
				m.message = "No file selected."
				m.currentMenu = mainMenu
				return m, nil
			}

		m.filePath = path
		base := filepath.Base(path)
		dir := filepath.Dir(path)
		dataFolder := filepath.Join(dir, base[:len(base)-4]+"_Data")
		m.spritePath = filepath.Join(dataFolder, "StreamingAssets", "CGAlt", "sprite")

		saveConfig(Config{
			GamePath:   m.filePath,
			SpritePath: m.spritePath,
			Selections: m.selections,
		})

		m.message = fmt.Sprintf("Selected %s", m.selectedEpisode.Label())
		m.currentMenu = mainMenu
	}
case helpMenu:
	switch key {
	case "q", "esc", "enter":
		m.currentMenu = mainMenu
		m.cursor = 0
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
				if len(m.meiOptions) == 0 {
					break
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

				
var variant string
switch chosen {
case "Best Match":
    data := Characters[m.selectedCharacter]
    if len(data.OutfitsMei) > 0 {
        variant = data.OutfitsMei[0].SpriteSet
        chosen = data.OutfitsMei[0].Name
    } else {
        variant = spriteSets[0]
    }
case "Random Outfits", "Random Outfits & Expressions", "Fully Random":
    variant = "" 
default:
    data := Characters[m.selectedCharacter]
    for _, o := range data.OutfitsMei {
        if o.Name == chosen {
            variant = o.SpriteSet
            break
        }
    }
}

if variant != "" {
    m.selections[m.selectedCharacter] = fmt.Sprintf("%s (variant: %s)", chosen, variant)
} else {
    m.selections[m.selectedCharacter] = chosen
}

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

func (m model) randomizeSprites() (tea.Model, tea.Cmd) {
    if m.spritePath == "" {
        m.message = "Select a game first."
        return m, nil
    }
	//log.Println("SPRITE PATH")
	//log.Println(m.spritePath)
	//log.Println("SPRITE PATH")
    spriteDir := m.spritePath
    backupDir := filepath.Join(filepath.Dir(spriteDir), "sprite_backup")
	//log.Println("BACKUP PATH")
	//log.Println(backupDir)
	//log.Println("BACKUP PATH")

    if _, err := os.Stat(backupDir); os.IsNotExist(err) {
        log.Println("Creating backup at:", backupDir)
        filepath.Walk(spriteDir, func(path string, info os.FileInfo, err error) error {
            if err != nil {
                return nil
            }
            if !info.IsDir() && filepath.Ext(path) == ".png" {
                rel, _ := filepath.Rel(spriteDir, path)
                dst := filepath.Join(backupDir, rel)
                os.MkdirAll(filepath.Dir(dst), 0755)
                data, _ := os.ReadFile(path)
                os.WriteFile(dst, data, 0644)
            }
            return nil
        })
    }

for key := range RawGameSprites {
    dst := filepath.Join(spriteDir, key+".png")
    if _, err := os.Stat(dst); err != nil {
        continue 
    }

    folder := GetFolder(key)
    selection := m.selections[folder]

    var chosenVariant string
    var chosenExpression string

    switch selection {
    case "Random Outfits":
        data := Characters[folder]
        if len(data.OutfitsMei) > 0 {
            o := data.OutfitsMei[rand.Intn(len(data.OutfitsMei))]
            chosenVariant = o.SpriteSet
            chosenExpression = RawGameSprites[key][0] // preserve the original expression
        } else {
            chosenVariant = spriteSets[0]
            chosenExpression = RawGameSprites[key][0]
        }
    case "Random Outfits & Expressions":
        data := Characters[folder]
        if len(data.OutfitsMei) > 0 {
            o := data.OutfitsMei[rand.Intn(len(data.OutfitsMei))]
            chosenVariant = o.SpriteSet

            variantFolder := filepath.Join("sprites", "mei", folder, chosenVariant)
            files, err := os.ReadDir(variantFolder)
            if err != nil || len(files) == 0 {
                chosenExpression = RawGameSprites[key][0] // fallback
            } else {
                pngFiles := []string{}
                for _, f := range files {
                    if !f.IsDir() && filepath.Ext(f.Name()) == ".png" {
                        pngFiles = append(pngFiles, strings.TrimSuffix(f.Name(), ".png"))
                    }
                }
                if len(pngFiles) > 0 {
                    chosenExpression = pngFiles[rand.Intn(len(pngFiles))]
                } else {
                    chosenExpression = RawGameSprites[key][0]
                }
            }
        } else {
            chosenVariant = spriteSets[0]
            chosenExpression = RawGameSprites[key][0]
        }
	case "Fully Random":
    const maxAttempts = 50
    attempt := 0

    var folders []string
    for k := range Characters {
        folders = append(folders, k)
    }

    if len(folders) == 0 {
        chosenVariant = spriteSets[0]
        chosenExpression = RawGameSprites[key][0]
        break
    }

    for attempt < maxAttempts {
        attempt++

        randomFolder := folders[rand.Intn(len(folders))]
        data := Characters[randomFolder]

        if len(data.OutfitsMei) == 0 {
            continue
        }

        o := data.OutfitsMei[rand.Intn(len(data.OutfitsMei))]
        variant := o.SpriteSet

        variantDir := filepath.Join("sprites", "mei", randomFolder, variant)
        entries, err := os.ReadDir(variantDir)
        if err != nil {
            continue
        }

        var candidates []string
        for _, e := range entries {
            if !e.IsDir() && strings.HasSuffix(e.Name(), ".png") {
                candidates = append(candidates, strings.TrimSuffix(e.Name(), ".png"))
            }
        }

        if len(candidates) == 0 {
            continue
        }

        expr := candidates[rand.Intn(len(candidates))]
        src := filepath.Join("sprites", "mei", randomFolder, variant, expr+".png")

        if fileExists(src) {
            folder = randomFolder
            chosenVariant = variant
            chosenExpression = expr
            break
        }
    }

    if chosenVariant == "" || chosenExpression == "" {
        chosenVariant = spriteSets[0]
        chosenExpression = RawGameSprites[key][0]
    }


    default:
        if start := strings.LastIndex(selection, "(variant: "); start != -1 {
            end := strings.Index(selection[start:], ")")
            if end != -1 {
                chosenVariant = selection[start+10 : start+end]
            }
        }
        if chosenVariant == "" {
            chosenVariant = spriteSets[0]
        }
        chosenExpression = RawGameSprites[key][0]
    }

    src := filepath.Join("sprites", "mei", folder, chosenVariant, chosenExpression+".png")
    data, err := os.ReadFile(src)
    if err != nil {
        log.Printf("Could not read Mei sprite: %s", src)
        continue
    }

    err = os.WriteFile(dst, data, 0644)
    if err != nil {
        log.Printf("Could not write sprite: %s", dst)
        continue
    }

    //log.Printf("Replaced: %s → %s (variant: %s, expression: %s)", key, dst, chosenVariant, chosenExpression)
}


    m.message = fmt.Sprintf(
    	"\nSprites randomized successfully. [%s]",
    	time.Now().Format("15:04:05"),
	)
    return m, nil
}


func (m model) restoreOriginalSprites() (tea.Model, tea.Cmd) {
    if m.spritePath == "" {
        m.message = "Select a game first."
        return m, nil
    }

    spriteDir := m.spritePath
    backupDir := filepath.Join(filepath.Dir(spriteDir), "sprite_backup")

    if _, err := os.Stat(backupDir); os.IsNotExist(err) {
        m.message = "No backup found. You must randomize once before restoring."
        return m, nil
    }

    filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
        if err != nil { return nil }
        if !info.IsDir() && filepath.Ext(path) == ".png" {
            rel, _ := filepath.Rel(backupDir, path)
            dst := filepath.Join(spriteDir, rel)
            os.MkdirAll(filepath.Dir(dst), 0755)
            data, _ := os.ReadFile(path)
            os.WriteFile(dst, data, 0644)
        }
        return nil
    })

    m.message = "Original sprites restored successfully."
    return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	switch m.currentMenu {
	case mainMenu:
		return fmt.Sprintf(
			"Higurandomizer 1.0\nBy figamin\n\n%s Auto Select Game\n%s Manually Select Game\n%s Select Sprites\n%s Check Selections\n%s Randomize\n%s Restore Original Sprites\n%s Help\n%s Exit\n\n%s\n",
			cursor(m.cursor, 0), cursor(m.cursor, 1), cursor(m.cursor, 2), cursor(m.cursor, 3), cursor(m.cursor, 4), cursor(m.cursor, 5), cursor(m.cursor, 6), cursor(m.cursor, 7),
			m.message,
		)



	case spriteMenu:
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
			//"Character: %s\n\n%s Mei\n%s Ace Attorney\n\nUse ↑↓ Enter, q to return.\n",
			"Character: %s\n\n%s Mei\n\nUse ↑↓ Enter, q to return.\n",
			m.selectedCharacter, cursor(m.cursor, 0),// cursor(m.cursor, 1),
		)
	case meiVariantMenu:
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
	case episodeMenu:
		s := "Select Episode\n\n"
		for i := 1; i <= 10; i++ {
			ep := HigurashiEpisode(i)
			s += fmt.Sprintf(
				"%s %s\n",
				cursor(m.cursor, i-1),
				ep.Label(),
			)
		}
		return s + "\nUse ↑↓ Enter, q to cancel.\n"
	case manualEpisodeMenu:
		s := "Select Episode\n\n"
		for i := 1; i <= 10; i++ {
			ep := HigurashiEpisode(i)
			s += fmt.Sprintf(
				"%s %s\n",
				cursor(m.cursor, i-1),
				ep.Label(),
			)
		}
		return s + "\nUse ↑↓ Enter, q to cancel.\n"
case helpMenu:
	return `
This program swaps the Mangagamer Higurashi sprites with those of alternative sets, with the option to randomize the sprites. Currently only the Mei sprites are included.

For now you will need 07th-Mod to use this. Later on I will add support for the vanilla game (it will just take a while because the vanilla sprite names are different)

The program saves your previous settings to config.json.
It should be able to autodetect Steam and GOG installs on Windows and Linux. If it does not work for whatever reason use the "Manually Select Game" option.

In terms of the sprite options:

- Best Match uses the closest sprites to the Ryukishi originals. Note that there were never Mei blinking sprites so the expressions may be more basic than the original. In the future I will try and make custom Mei sprites for certain expressions like that.

- Random Outfits uses random sprites from across the Mei outfit variations for the given character, but with the same expressions as the Best Match

- Random Outfits + Expresisons uses random sprites from across the Mei outfit variations for the given character.

- Fully Random uses entirely random sprites from Mei. It will make zero sense.

- The rest of the options change all the character sprites to a specific Mei variant. Note that there is currently no way to assign different variants to different character outfit sets (such as school vs casual vs gym), so the characters will always be wearing the same outfit if one of these options is used.
`

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
func fileExists(p string) bool {
	i, err := os.Stat(p)
	return err == nil && !i.IsDir()
}

func dirExists(p string) bool {
	i, err := os.Stat(p)
	return err == nil && i.IsDir()
}

func autodetectHigurashi(ep HigurashiEpisode) (string, error) {
	exeBase := ep.ExeName()
	folderName := ep.InstallFolderName()

	if folderName == "" {
		return "", fmt.Errorf("unknown episode %q", ep)
	}

	for _, root := range possibleGameRoots() {
		gameDir := filepath.Join(root, folderName)

		info, err := os.Stat(gameDir)
		if err != nil || !info.IsDir() {
			continue
		}
		if runtime.GOOS == "windows" {
			exe := filepath.Join(gameDir, exeBase+".exe")
			if fileExists(exe) {
				return gameDir, nil
			}
		}

		if runtime.GOOS == "darwin" {
			app := filepath.Join(gameDir, exeBase+".app")
			if dirExists(app) {
				return app, nil
			}
		}

		if runtime.GOOS == "linux" {
			if fileExists(filepath.Join(gameDir, exeBase+".x86")) ||
				fileExists(filepath.Join(gameDir, exeBase+".x86_64")) ||
				fileExists(filepath.Join(gameDir, exeBase+".exe")) {
					return gameDir, nil
			}
		}
	}

	return "", fmt.Errorf(
		"could not autodetect %s (%s)",
		exeBase,
		folderName,
	)
}


func possibleGameRoots() []string {
	var roots []string

	switch runtime.GOOS {
	case "windows":
		roots = append(roots,
			os.Getenv("PROGRAMFILES(X86)")+"\\Steam\\steamapps\\common",
			os.Getenv("PROGRAMFILES")+"\\Steam\\steamapps\\common",
			"C:\\GOG Games",
			"C:\\games\\Mangagamer",
		)
	case "darwin":
		roots = append(roots,
			filepath.Join(os.Getenv("HOME"), "Library/Application Support/Steam/steamapps/common"),
			"/Applications",
			filepath.Join(os.Getenv("HOME"), "GOG Games"),
		)
	case "linux":
		roots = append(roots,
			filepath.Join(os.Getenv("HOME"), ".steam/steam/steamapps/common"),
			filepath.Join(os.Getenv("HOME"), ".steam/steambeta/steamapps/common"),
			filepath.Join(os.Getenv("HOME"), ".var/app/com.valvesoftware.Steam/data/Steam/steamapps/common"),
			filepath.Join(os.Getenv("HOME"), "GOG Games"),
		)
	}

	seen := map[string]bool{}
	var out []string
	for _, r := range roots {
		r = filepath.Clean(r)
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}
func (e HigurashiEpisode) InstallFolderName() string {
	switch e {
	case Ep01:
		return "Higurashi When They Cry"
	case Ep02:
		return "Higurashi 02 - Watanagashi"
	case Ep03:
		return "Higurashi 03 - Tatarigoroshi"
	case Ep04:
		return "Higurashi 04 - Himatsubushi"
	case Ep05:
		return "Higurashi When They Cry Hou - Ch. 5 Meakashi"
	case Ep06:
		return "Higurashi When They Cry Hou - Ch.6 Tsumihoroboshi"
	case Ep07:
		return "Higurashi When They Cry Hou - Ch.7 Minagoroshi"
	case Ep08:
		return "Higurashi When They Cry Hou - Ch.8 Matsuribayashi"
	case Ep09:
		return "Higurashi When They Cry Hou - Rei"
	case Ep10:
		return "Higurashi When They Cry Hou+"
	default:
		return ""
	}
}

func main() {
	ensureSprites()
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
