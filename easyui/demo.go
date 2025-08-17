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
	uiGenerator := &BareBonesUiGenerator{theme}

	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	// Example: A regular button
	button := NewButton(100, 100, 200, 50, "Click me!", uiGenerator) // Updated call
	button.SetClickHandler(func() {
		log.Println("Regular button clicked!")
		// Example of changing button text dynamically:
		button.SetText("Clicked!")
	})
	ui.AddChild(button)

	// --- Dropdown Menu Implementation ---
	menuWidth := 200
	animalMenu := NewMenu(350, 200, menuWidth, theme, uiGenerator, ui) // Still uses theme and uiGenerator directly for now

	dropdown := NewDropDown(350, 150, 200, 40, "Select an Animal", animalMenu, uiGenerator) // Updated call
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
	checkbox := NewCheckbox(100, 200, 150, 30, "Enable Feature", false, uiGenerator) // Updated call
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
	nameField := NewTextField(100, 250, 300, 30, "Enter your name", uiGenerator) // Updated call
	ui.AddChild(nameField)

	// --- Label Implementation ---
	infoLabel := NewLabel(100, 300, 400, 20, "Welcome to the UI Demo!", uiGenerator) // Updated call
	ui.AddChild(infoLabel)

	// Example of changing label text dynamically (e.g., after 5 seconds)
	// For a real application, you'd trigger this based on game state or user action.
	// This simple example runs once at startup.
	go func() {
		// time.Sleep(5 * time.Second) // In a real Ebiten game, don't use time.Sleep in main goroutine.
		// Instead, use an internal timer in the Demo's Update loop.
		// For this simple demo, we'll just log it.
		// log.Println("Changing label text...")
		// infoLabel.SetText("Demo Ready! Interact above.")
	}()

	return &Demo{ui: ui}
}
