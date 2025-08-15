package main

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// BareBonesTheme defines the color and font theme for UI elements.
type BareBonesTheme struct {
	BackgroundColor color.RGBA // e.g. dark gray
	PrimaryColor    color.RGBA // your color theme
	OnPrimaryColor  color.RGBA // probably white/near white. color of text/border
	AccentColor     color.RGBA // for accents when pressed
	Face            font.Face  // The loaded font face
}

// BareBonesUiGenerator helps instantiate UI components with the defined theme.
type BareBonesUiGenerator struct {
	theme BareBonesTheme
}

// NewButton creates a new Button instance with the specified dimensions, label, and theme.
func (b *BareBonesUiGenerator) NewButton(x, y, width, height int, label string) *Button {
	// Generate images for different button states
	idle := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label)
	// For simplicity, using the same image for pressed and hover states for now.
	// You can create distinct images for each state later if desired.
	pressed := b.generateButtonImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, label) // Use accent color when pressed
	hover := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label)
	disabled := b.generateButtonImage(width, height, b.theme.BackgroundColor, b.theme.OnPrimaryColor, label) // Greyed out for disabled

	return &Button{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y + height},
			},
		},
		idle:     idle,
		pressed:  pressed,
		hover:    hover,
		disabled: disabled,
		state:    ButtonIdle, // Initial state
	}
}

// generateButtonImage draws a rounded rectangle button with text using the gg framework.
func (b *BareBonesUiGenerator) generateButtonImage(
	width, height int,
	buttonColor, textColor color.RGBA,
	buttonText string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Set background color for the button
	dc.SetRGBA255(int(buttonColor.R), int(buttonColor.G), int(buttonColor.B), int(buttonColor.A))
	// Draw a rounded rectangle for the button background
	cornerRadius := float64(height) / 4 // Adjust as needed for desired roundness
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.Fill()

	// Draw a frame (border) around the button
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(2) // Border thickness
	dc.Stroke()

	// Use the font face provided in the theme
	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Set text color
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	// Draw text centered on the button
	dc.DrawStringWrapped(buttonText, float64(width)/2, float64(height)/2, 0.5, 1.0, float64(width)-10, 1.5, gg.AlignCenter)

	// Convert gg.Context image to ebiten.Image
	return ebiten.NewImageFromImage(dc.Image())
}
