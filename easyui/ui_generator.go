package main

import (
	"image"
	"image/color"
	"log" // Import for logging

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// BareBonesTheme defines the color and font theme for UI elements.
type BareBonesTheme struct {
	BackgroundColor    color.RGBA // e.g. dark gray
	PrimaryColor       color.RGBA // your color theme (for buttons, dropdowns)
	OnPrimaryColor     color.RGBA // probably white/near white. color of text/border on primary elements
	AccentColor        color.RGBA // for accents when pressed
	MenuColor          color.RGBA // Color for the menu background
	MenuItemHoverColor color.RGBA // Color for menu items on hover
	Face               font.Face  // The loaded font face
}

// BareBonesUiGenerator helps instantiate UI components with the defined theme.
type BareBonesUiGenerator struct {
	theme BareBonesTheme
}

// NewButton creates a new Button instance with the specified dimensions, label, and theme.
func (b *BareBonesUiGenerator) NewButton(x, y, width, height int, label string) *Button {
	idle := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label)
	pressed := b.generateButtonImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, label)
	hover := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label) // Could be a different hover color
	disabled := b.generateButtonImage(width, height, b.theme.BackgroundColor, b.theme.OnPrimaryColor, label)

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
		state:    ButtonIdle,
	}
}

// generateButtonImage draws a rounded rectangle button with text using the gg framework.
func (b *BareBonesUiGenerator) generateButtonImage(
	width, height int,
	buttonColor, textColor color.RGBA,
	buttonText string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	dc.SetRGBA255(int(buttonColor.R), int(buttonColor.G), int(buttonColor.B), int(buttonColor.A))
	cornerRadius := float64(height) / 4
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.FillPreserve() // Fill and preserve path for stroke

	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(2) // Border thickness
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringWrapped(buttonText, float64(width)/2, float64(height)/2, 0.5, 1.0, float64(width)-10, 1.5, gg.AlignCenter)

	img := ebiten.NewImageFromImage(dc.Image())
	log.Printf("Generated ButtonImage ('%s'): Bounds %v", buttonText, img.Bounds())
	return img
}

// generateMenuItemImage draws a single menu item with text and a background color that changes on hover/press.
func (b *BareBonesUiGenerator) generateMenuItemImage(
	width, height int,
	bgColor, textColor color.RGBA,
	itemText string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height)) // Menu items typically have sharp corners
	dc.Fill()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringWrapped(itemText, float64(width)/2, float64(height)/2, 0.5, 1.0, float64(width)-10, 1.5, gg.AlignCenter)

	img := ebiten.NewImageFromImage(dc.Image())
	log.Printf("Generated MenuItemImage ('%s'): Bounds %v", itemText, img.Bounds())
	return img
}

// generateMenuImage draws a solid background for the entire menu.
func (b *BareBonesUiGenerator) generateMenuImage(
	width, height int,
	bgColor color.RGBA,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	img := ebiten.NewImageFromImage(dc.Image())
	log.Printf("Generated MenuImage: Bounds %v (W:%d, H:%d)", img.Bounds(), width, height)
	return img
}

// generateDropdownImage draws the dropdown button with text and a small arrow.
func (b *BareBonesUiGenerator) generateDropdownImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw the background of the dropdown button as a rectangle
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height)) // Changed to DrawRectangle
	dc.FillPreserve()                                       // Fill and preserve path for stroke
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(2)
	dc.Stroke()

	// Draw the V arrow on the right side of the dropdown
	arrowHeight := float64(height) / 5 // A fifth of the button height for the arrow's "tallness"
	arrowWidth := 2 * arrowHeight      // Arrow is twice as wide as it is tall

	padding := float64(width) * 0.05 // Small padding from the right edge
	arrowX := float64(width) - arrowWidth - padding
	arrowY := float64(height)/2 - arrowHeight/2

	dc.MoveTo(arrowX, arrowY)                          // Top-left point of the arrow's base
	dc.LineTo(arrowX+arrowWidth/2, arrowY+arrowHeight) // Bottom tip of the arrow (middle of base)
	dc.LineTo(arrowX+arrowWidth, arrowY)               // Top-right point of the arrow's base
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Set text color and draw the text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	// Position text to the left of the arrow, with some padding
	textWidth, _ := dc.MeasureString(text)
	textX := padding                                          // Start text from left padding
	maxTextWidth := float64(width) - (arrowWidth + padding*2) // Max space for text before arrow

	// Ensure text fits and is centered within its available space
	// If the text is too wide, DrawStringWrapped will handle it, but for single line we adjust x.
	if textWidth > maxTextWidth {
		// If text is too wide, just left align with padding
		dc.DrawString(text, textX, float64(height)/2+dc.FontHeight()/3)
	} else {
		// Otherwise, center the text within the available space to the left of the arrow
		centeredTextX := (maxTextWidth-textWidth)/2 + padding
		dc.DrawString(text, centeredTextX, float64(height)/2+dc.FontHeight()/3)
	}

	img := ebiten.NewImageFromImage(dc.Image())
	log.Printf("Generated DropdownImage ('%s'): Bounds %v", text, img.Bounds())
	return img
}
