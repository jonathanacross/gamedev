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
	ui     Ui
	button *Button
}

func (g *Demo) Update() error {
	// Update all UI elements
	g.ui.Update()
	return nil
}

func (g *Demo) Draw(screen *ebiten.Image) {
	// Draw all UI elements
	g.ui.Draw(screen)
}

func (g *Demo) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func NewDemo() *Demo {
	// Define font path and size
	fontPath := "Go-Mono.ttf" // Make sure this file exists
	fontSize := 12.0

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
		BackgroundColor: color.RGBA{100, 100, 100, 255}, // Dark gray background
		PrimaryColor:    color.RGBA{100, 150, 200, 255}, // Blue-ish primary color for buttons
		OnPrimaryColor:  color.RGBA{255, 255, 255, 255}, // White text/border on primary
		AccentColor:     color.RGBA{255, 255, 0, 255},   // Yellow accent for pressed state
		Face:            face,                           // Assign the loaded font face
	}
	uiGenerator := &BareBonesUiGenerator{theme}

	// Create the root UI container
	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	// Create a button and add it to the UI
	// Position: (100, 100), Size: (200, 50), Label: "Click me!"
	button := uiGenerator.NewButton(100, 100, 200, 50, "Click me!")
	ui.AddChild(button)

	return &Demo{
		ui:     *ui,
		button: button,
	}
}
