package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font/basicfont"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

var (
	pencil = createIcon([16]int{
		0x00f0, 0x0088, 0x0108, 0x0190, 0x0270, 0x0220, 0x0420, 0x0440,
		0x0840, 0x0880, 0x1080, 0x1100, 0x1e00, 0x1c00, 0x1800, 0x1000,
	})

	brush = createIcon([16]int{
		0x01c0, 0x0140, 0x01c0, 0x01c0, 0x01c0, 0x01c0, 0x0ff8, 0x0808,
		0x0ff8, 0x0808, 0x0808, 0x0808, 0x0808, 0x0aa8, 0x1558, 0x3ff0,
	})

	bucket = createIcon([16]int{
		0x0700, 0x0880, 0x0980, 0x0ac0, 0x0cb0, 0x089c, 0x108e, 0x2147,
		0x4087, 0x800f, 0x8017, 0x4027, 0x2047, 0x1087, 0x0906, 0x0604,
	})

	spraycan = createIcon([16]int{
		0x4000, 0x1000, 0x4540, 0x10e0, 0x4110, 0x03f8, 0x0208, 0x0278,
		0x0248, 0x0278, 0x0248, 0x0278, 0x0278, 0x0208, 0x0208, 0x03f8,
	})
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
	button := NewButton(50, 50, 150, 40, "Click me!", nil, false, uiGenerator)
	button.SetClickHandler(func() {
		log.Println("Button clicked!")
		button.SetText("Clicked!")
	})
	container.AddChild(button)

	button2 := NewButton(250, 50, 150, 40, "Disabled", nil, false, uiGenerator)
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
	b1 := NewButton(0, 0, 30, 30, "", pencil, true, uiGenerator)
	b2 := NewButton(0, 0, 30, 30, "", brush, true, uiGenerator)
	b3 := NewButton(0, 0, 30, 30, "", bucket, true, uiGenerator)
	b4 := NewButton(0, 0, 30, 30, "", spraycan, true, uiGenerator)
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

// Creates a 16x16 binary icon from data.
func createIcon(data [16]int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))

	for y, val := range data {
		for x := range 16 {
			isBitSet := (val & (1 << (15 - x))) != 0
			rgbaIndex := (y*16 + x) * 4

			if isBitSet {
				// Set the pixel to opaque white
				img.Pix[rgbaIndex+0] = 0xff
				img.Pix[rgbaIndex+1] = 0xff
				img.Pix[rgbaIndex+2] = 0xff
				img.Pix[rgbaIndex+3] = 0xff
			} else {
				// Set the pixel to transparent black
				img.Pix[rgbaIndex+0] = 0x00
				img.Pix[rgbaIndex+1] = 0x00
				img.Pix[rgbaIndex+2] = 0x00
				img.Pix[rgbaIndex+3] = 0x00
			}
		}
	}

	return img
}

func main() {
	demo := NewDemo()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("EasyUi Demo")
	if err := ebiten.RunGame(demo); err != nil {
		log.Fatal(err)
	}
}
