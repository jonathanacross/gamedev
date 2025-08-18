package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Demo struct {
	ui *Ui
}

func (g *Demo) Update() error {
	g.ui.Update()
	return nil
}

func (g *Demo) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *Demo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func loadFontFace() font.Face {
	fontPath := "Go-Mono.ttf"
	fontSize := 14.0

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("Error loading font %s: %v", fontPath, err)
	}

	tt, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("Error parsing font: %v", err)
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Error creating font face: %v", err)
	}

	return face
}

func NewDemo() *Demo {
	face := loadFontFace()

	theme := ShapeTheme{
		BackgroundColor:    color.RGBA{20, 20, 20, 255},
		PrimaryColor:       color.RGBA{120, 120, 120, 255},
		OnPrimaryColor:     color.RGBA{255, 255, 255, 255},
		AccentColor:        color.RGBA{220, 120, 120, 255},
		MenuColor:          color.RGBA{100, 100, 100, 255},
		MenuItemHoverColor: color.RGBA{120, 120, 120, 255},
		Face:               face,
	}
	uiGenerator := &ShapeRenderer{theme}

	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	// --- A button ---
	button := NewButton(100, 100, 200, 50, "Click me!", uiGenerator)
	button.SetClickHandler(func() {
		log.Println("Regular button clicked!")
		button.SetText("Clicked!")
	})
	ui.AddChild(button)

	// --- Dropdown Menu ---
	menuWidth := 200
	// The menu's initial position will be set absolutely by the dropdown.
	animalMenu := NewMenu(0, 0, menuWidth, theme, uiGenerator, ui)

	dropdown := NewDropDown(350, 150, 200, 40, "Select an Animal", animalMenu, uiGenerator)
	ui.AddChild(dropdown)

	animals := []string{"Lion", "Tiger", "Bear", "Elephant"}

	for _, animal := range animals {
		currentAnimal := animal // Capture loop variable
		animalMenu.AddItem(currentAnimal, func() {
			log.Printf("Dropdown: %s selected!", currentAnimal)
			dropdown.SetSelectedOption(currentAnimal)
		})
	}

	// --- Checkbox ---
	checkbox := NewCheckbox(100, 200, 150, 30, "Enable Feature", false, uiGenerator)
	checkbox.OnCheckChanged = func(checked bool) {
		log.Printf("Checkbox 'Enable Feature' state changed to: %t", checked)
		if checked {
			button.SetText("Feature Enabled!")
		} else {
			button.SetText("Feature Disabled.")
		}
	}
	ui.AddChild(checkbox)

	// --- TextField ---
	nameField := NewTextField(100, 250, 300, 30, "Enter your name", uiGenerator)
	ui.AddChild(nameField)

	// --- Container ---
	container := NewContainer(50, 350, 700, 100, uiGenerator)
	ui.AddChild(container)

	// Create a label and add it to the container
	// Label's x,y are relative to the container's top-left corner
	infoLabel := NewLabel(10, 10, 380, 20, "This label is inside a container!", uiGenerator)
	container.AddChild(infoLabel)

	// Add another label to the main UI to show it's separate
	globalLabel := NewLabel(50, 470, 400, 20, "This label is directly on the UI.", uiGenerator)
	ui.AddChild(globalLabel)

	return &Demo{ui: ui}
}

func main() {
	demo := NewDemo()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("EasyUi Demo")
	if err := ebiten.RunGame(demo); err != nil {
		log.Fatal(err)
	}
}
