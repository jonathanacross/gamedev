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

	// --- Container ---
	container := NewContainer(50, 50, 700, 400, uiGenerator)
	ui.AddChild(container)

	// --- Buttons ---
	button := NewButton(50, 50, 150, 40, "Click me!", false, uiGenerator)
	button.SetClickHandler(func() {
		log.Println("Button clicked!")
		button.SetText("Clicked!")
	})
	container.AddChild(button)

	button2 := NewButton(250, 50, 150, 40, "Disabled", false, uiGenerator)
	button2.state = ButtonDisabled
	button2.SetClickHandler(func() { button2.SetText("ack! clicked!") })
	container.AddChild(button2)

	// --- Dropdown Menu ---
	menuWidth := 200
	// The menu's initial position will be set absolutely by the dropdown.
	animalMenu := NewMenu(0, 0, menuWidth, uiGenerator, ui)

	dropdown := NewDropDown(300, 100, 200, 40, "Select an Animal", animalMenu, uiGenerator)
	container.AddChild(dropdown)

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
	checkbox := NewCheckbox(50, 150, 150, 30, "Enable Feature", false, uiGenerator)
	checkbox.OnCheckChanged = func(checked bool) {
		log.Printf("Checkbox 'Enable Feature' state changed to: %t", checked)
		if checked {
			button.SetText("Feature Enabled!")
		} else {
			button.SetText("Feature Disabled.")
		}
	}
	container.AddChild(checkbox)

	// --- Radio Buttons managed by a ButtonGroup ---
	radioGroup := NewButtonGroup(400, 150, 200, 120, LayoutVertical, SingleSelection, 5)
	container.AddChild(radioGroup)

	// Now, add the radio buttons directly to the group. The group handles their positioning and exclusivity.
	rb1 := NewRadioButton(0, 0, 150, 20, "Peanuts", true, uiGenerator)
	rb2 := NewRadioButton(0, 0, 150, 20, "Cracker", false, uiGenerator)
	rb3 := NewRadioButton(0, 0, 150, 20, "Cookies", false, uiGenerator)
	radioGroup.AddChild(rb1)
	radioGroup.AddChild(rb2)
	radioGroup.AddChild(rb3)

	// -- Toggle button bar --
	toggleButtonGroup := NewButtonGroup(450, 50, 1, 1, LayoutHorizontal, SingleSelection, 2)
	b1 := NewButton(0, 0, 30, 30, "A", true, uiGenerator)
	b2 := NewButton(0, 0, 30, 30, "B", true, uiGenerator)
	b3 := NewButton(0, 0, 30, 30, "C", true, uiGenerator)
	b4 := NewButton(0, 0, 30, 30, "D", true, uiGenerator)
	toggleButtonGroup.AddChild(b1)
	toggleButtonGroup.AddChild(b2)
	toggleButtonGroup.AddChild(b3)
	toggleButtonGroup.AddChild(b4)
	container.AddChild(toggleButtonGroup)

	// --- TextField ---
	nameField := NewTextField(50, 200, 300, 30, "Enter your name", uiGenerator)
	container.AddChild(nameField)

	textField2 := NewTextField(50, 240, 300, 30, "another field", uiGenerator)
	container.AddChild(textField2)

	// Create a label and add it to the container
	// Label's x,y are relative to the container's top-left corner
	infoLabel := NewLabel(20, 20, 380, 20, "This label is inside a container!", uiGenerator)
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
