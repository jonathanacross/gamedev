package main

import (
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
		BackgroundColor:    color.RGBA{20, 20, 20, 255},
		PrimaryColor:       color.RGBA{120, 120, 120, 255},
		OnPrimaryColor:     color.RGBA{255, 255, 255, 255},
		AccentColor:        color.RGBA{220, 120, 120, 255},
		MenuColor:          color.RGBA{100, 100, 100, 255},
		MenuItemHoverColor: color.RGBA{120, 120, 120, 255},
		Face:               face,
		// MenuItemHeight:     30, // Removed for this specific step to resolve a potential conflict
	}
	// Initialize BareBonesUiGenerator, which now implements UiRenderer
	uiGenerator := &BareBonesUiGenerator{theme}

	// The root UI has absolute coordinates (0,0)
	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	// Example: A regular button
	// Pass the uiGenerator (as UiRenderer) to the component constructor
	button := NewButton(100, 100, 200, 50, "Click me!", uiGenerator)
	button.SetClickHandler(func() {
		log.Println("Regular button clicked!")
		// Example of changing button text dynamically:
		button.SetText("Clicked!")
	})
	ui.AddChild(button)

	// --- Dropdown Menu Implementation ---
	menuWidth := 200
	// Pass the uiGenerator (as UiRenderer) to NewMenu
	// The menu's initial position will be set absolutely by the dropdown, and then it knows its parent (the UI)
	animalMenu := NewMenu(0, 0, menuWidth, theme, uiGenerator, ui) // Initial pos (0,0) will be overridden by dropdown

	// Pass the uiGenerator (as UiRenderer) to NewDropDown
	dropdown := NewDropDown(350, 150, 200, 40, "Select an Animal", animalMenu, uiGenerator)
	ui.AddChild(dropdown)

	animals := []string{"Lion", "Tiger", "Bear", "Elephant"}

	for _, animal := range animals {
		currentAnimal := animal // Capture loop variable
		animalMenu.AddItem(currentAnimal, func() {
			log.Printf("Dropdown: %s selected!", currentAnimal)
			// Update the dropdown's displayed text
			dropdown.SetSelectedOption(currentAnimal)
		})
	}

	// --- Checkbox Implementation ---
	// Pass the uiGenerator (as UiRenderer) to NewCheckbox
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

	// --- TextField Implementation ---
	// Pass the uiGenerator (as UiRenderer) to NewTextField
	nameField := NewTextField(100, 250, 300, 30, "Enter your name", uiGenerator)
	ui.AddChild(nameField)

	// --- Container Implementation ---
	// Create a container (position is relative to 'ui' which is 0,0)
	container := NewContainer(50, 350, 700, 100, uiGenerator) // x, y, width, height, renderer
	ui.AddChild(container)                                    // Add the container to the root UI

	// Create a label and add it to the container
	// Label's x,y (10,10) are now RELATIVE to the container's top-left corner
	infoLabel := NewLabel(10, 10, 380, 20, "This label is inside a container!", uiGenerator)
	container.AddChild(infoLabel) // Add the label as a child of the container

	// Add another label to the main UI to show it's separate
	globalLabel := NewLabel(50, 470, 400, 20, "This label is directly on the UI.", uiGenerator)
	ui.AddChild(globalLabel)

	return &Demo{ui: ui}
}
