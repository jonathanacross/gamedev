package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
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

func NewDemo() *Demo {
	face := basicfont.Face7x13

	theme := ShapeTheme{
		PrimaryAccentColor: color.RGBA{R: 0x40, G: 0x70, B: 0xB0, A: 0xFF},
		BackgroundColor:    color.RGBA{R: 0x1A, G: 0x20, B: 0x2C, A: 0xFF},
		SurfaceColor:       color.RGBA{R: 0x2D, G: 0x37, B: 0x48, A: 0xFF},
		TextColor:          color.RGBA{R: 0xF8, G: 0xF8, B: 0xF8, A: 0xFF},
		BorderColor:        color.RGBA{R: 0x4A, G: 0x55, B: 0x68, A: 0xFF},
		Face:               face,
	}

	uiGenerator := &ShapeRenderer{theme}

	ui := NewUi(0, 0, ScreenWidth, ScreenHeight)

	// --- Buttons ---
	button := NewButton(100, 100, 150, 40, "Click me!", uiGenerator)
	button.SetClickHandler(func() {
		log.Println("Button clicked!")
		button.SetText("Clicked!")
	})
	ui.AddChild(button)

	button2 := NewButton(300, 100, 150, 40, "Disabled", uiGenerator)
	button2.state = ButtonDisabled
	button.SetClickHandler(func() {})
	ui.AddChild(button2)

	// --- Dropdown Menu ---
	menuWidth := 200
	// The menu's initial position will be set absolutely by the dropdown.
	animalMenu := NewMenu(0, 0, menuWidth, uiGenerator, ui)

	dropdown := NewDropDown(350, 150, 200, 40, "Select an Animal", animalMenu, uiGenerator)
	ui.AddChild(dropdown)

	animals := []string{"Lion", "Tiger", "Bear", "Elephant"}

	for _, animal := range animals {
		currentAnimal := animal
		animalMenu.AddItem(currentAnimal, func() {
			log.Printf("Dropdown: %s selected!", currentAnimal)
			dropdown.SetSelectedOption(currentAnimal)
			animalMenu.Hide()
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
