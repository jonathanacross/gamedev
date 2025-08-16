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
	ui       Ui
	button   *Button
	dropdown *DropDown // New field for the dropdown
}

func (g *Demo) Update() error {
	// Update all UI elements. The Ui struct now handles modal components.
	g.ui.Update()
	return nil
}

func (g *Demo) Draw(screen *ebiten.Image) {
	// Draw all UI elements. The Ui struct now handles modal components.
	g.ui.Draw(screen)
}

func (g *Demo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewDemo() *Demo {
	// Define font path and size
	fontPath := "Go-Mono.ttf" // Make sure this file exists
	fontSize := 14.0          // Slightly smaller font for menu items

	// Load the font using os.ReadFile
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

	// Define your theme with the loaded font face
	theme := BareBonesTheme{
		BackgroundColor:    color.RGBA{100, 100, 100, 255}, // Dark gray background
		PrimaryColor:       color.RGBA{100, 150, 200, 255}, // Blue-ish primary color for buttons/dropdown
		OnPrimaryColor:     color.RGBA{255, 255, 255, 255}, // White text/border on primary
		AccentColor:        color.RGBA{255, 255, 0, 255},   // Yellow accent for pressed state
		MenuColor:          color.RGBA{80, 80, 80, 255},    // Darker gray for menu background
		MenuItemHoverColor: color.RGBA{120, 120, 120, 255}, // Slightly lighter gray for menu item hover
		Face:               face,                           // Assign the loaded font face
	}
	uiGenerator := &BareBonesUiGenerator{theme}

	// Create the root UI container
	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	// Create a regular button and add it to the UI
	button := uiGenerator.NewButton(100, 100, 200, 50, "Click me!")
	button.SetClickHandler(func() {
		log.Println("Regular button clicked!")
	})
	ui.AddChild(button)

	// --- Dropdown Menu Implementation ---
	// 1. Create the Menu component
	// The Menu's position will be set by the DropDown when it opens.
	// Its width will match the dropdown. Initial Y is just a placeholder.
	menuWidth := 200
	animalMenu := NewMenu(350, 200, menuWidth, theme, uiGenerator, ui)

	// 2. Add animal items to the menu
	animals := []string{"Lion", "Tiger", "Bear", "Wolf", "Deer", "Fox", "Eagle", "Shark"}
	// Forward declaration for dropdown to be used in menu item handlers
	var dropdown *DropDown

	for _, animal := range animals {
		currentAnimal := animal // Capture loop variable
		animalMenu.AddItem(currentAnimal, func() {
			log.Printf("Selected animal: %s\n", currentAnimal)
			// When an item is clicked, it will hide the menu automatically (handled in Menu.AddItem handler)
			// The dropdown's text needs to update to the selected animal.
			if dropdown != nil { // Ensure dropdown is not nil before using
				dropdown.SelectedOption = currentAnimal // Update the dropdown's displayed text
			}
		})
	}

	// 3. Create the DropDown component
	dropdownWidth := 200
	dropdownHeight := 50
	dropdownX := 350
	dropdownY := 100

	dropdown = NewDropDown( // Assign to the declared variable
		dropdownX,
		dropdownY,
		dropdownWidth,
		dropdownHeight,
		"Select an Animal", // Initial label for the dropdown button
		animalMenu,         // Pass the menu to the dropdown
		theme,
		uiGenerator,
	)

	// 4. Add the dropdown to the root UI
	ui.AddChild(dropdown)

	return &Demo{
		ui:       *ui,
		button:   button,
		dropdown: dropdown,
	}
}
