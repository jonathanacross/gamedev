package main

import (
	"image/color"
	"log"
	"os" // Import the os package for os.ReadFile

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Demo struct {
	ui       *Ui
	button   *Button
	dropdown *DropDown
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

func NewDemo() *Demo {
	fontPath := "Go-Mono.ttf" // Make sure this file exists
	fontSize := 14.0          // Slightly smaller font for menu items

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

	theme := BareBonesTheme{
		BackgroundColor:    color.RGBA{100, 100, 100, 255},
		PrimaryColor:       color.RGBA{100, 150, 200, 255},
		OnPrimaryColor:     color.RGBA{255, 255, 255, 255},
		AccentColor:        color.RGBA{255, 255, 0, 255},
		MenuColor:          color.RGBA{80, 80, 80, 255},
		MenuItemHoverColor: color.RGBA{120, 120, 120, 255},
		Face:               face,
	}
	uiGenerator := &BareBonesUiGenerator{theme}

	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	button := uiGenerator.NewButton(100, 100, 200, 50, "Click me!")
	button.SetClickHandler(func() {
		log.Println("Regular button clicked!")
	})
	ui.AddChild(button)

	// --- Dropdown Menu Implementation ---
	menuWidth := 200
	animalMenu := NewMenu(350, 200, menuWidth, theme, uiGenerator, ui)

	animals := []string{"Lion", "Tiger", "Bear"}
	var dropdown *DropDown // Forward declaration for dropdown

	for _, animal := range animals {
		animalMenu.AddItem(animal, func() {
			// This placeholder handler will be overridden below.
		})
	}

	dropdownWidth := 200
	dropdownHeight := 50
	dropdownX := 350
	dropdownY := 100

	dropdown = NewDropDown(
		dropdownX,
		dropdownY,
		dropdownWidth,
		dropdownHeight,
		"Select an Animal", // Initial label for the dropdown button
		animalMenu,         // Pass the menu to the dropdown
		theme,
		uiGenerator,
	)

	// Update the menu items' handlers AFTER 'dropdown' is initialized
	for _, item := range animalMenu.items {
		itemLabel := item.Label // Capture loop variable for the closure
		item.SetClickHandler(func() {
			log.Printf("Selected animal: %s\n", itemLabel)
			dropdown.SelectedOption = itemLabel // Update the dropdown's displayed text
			animalMenu.Hide()                   // Hide the menu
		})
	}

	ui.AddChild(dropdown)

	return &Demo{
		ui:       ui,
		button:   button,
		dropdown: dropdown,
	}
}
