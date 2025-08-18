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
	button := NewButton(50, 50, 150, 40, "Click me!", uiGenerator)
	button.SetClickHandler(func() {
		log.Println("Button clicked!")
		button.SetText("Clicked!")
	})
	container.AddChild(button)

	button2 := NewButton(250, 50, 150, 40, "Disabled", uiGenerator)
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

	// --- Radio buttons ---
	radioLabels := []string{"Peanuts", "Crackers", "Cookies"}
	radioButtons := make([]*RadioButton, len(radioLabels))
	initialY := 150

	for i, label := range radioLabels {
		rb := NewRadioButton(400, initialY+(i*30), 150, 30, label, i == 0, uiGenerator) // i==0 sets the first button as checked
		radioButtons[i] = rb
		container.AddChild(rb)

		// Create a closure to capture the correct button for the handler
		currentButton := rb
		currentButton.OnCheckChanged = func(checked bool) {
			if checked {
				log.Printf("Radio button '%s' checked.", currentButton.Label)
				// Uncheck all other radio buttons
				for _, otherRb := range radioButtons {
					if otherRb != currentButton {
						otherRb.SetChecked(false)
					}
				}
			}
		}
	}

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
