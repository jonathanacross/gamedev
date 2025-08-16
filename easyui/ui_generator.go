package main

import (
	"image/color"

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
	hover := b.generateButtonImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, label)
	disabled := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label) // Example: disabled state image

	return &Button{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idle, pressed, hover, disabled),
		onClick:              nil, // Will be set by SetClickHandler
	}
}

// generateButtonImage creates an Ebiten image for a button's specific state.
func (b *BareBonesUiGenerator) generateButtonImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw button background with rounded corners
	cornerRadius := float64(height) * 0.2 // 20% of height for radius
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.FillPreserve()
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(2)
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw button text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringAnchored(text, float64(width)/2, float64(height)/2, 0.5, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// generateDropdownImage creates an Ebiten image for a dropdown's specific state.
func (b *BareBonesUiGenerator) generateDropdownImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw the background of the dropdown button as a rectangle
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.FillPreserve()
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(2)
	dc.Stroke()

	// Draw the V arrow on the right side of the dropdown
	arrowHeight := float64(height) / 5
	arrowWidth := 2 * arrowHeight

	padding := float64(width) * 0.05
	arrowX := float64(width) - arrowWidth - padding
	arrowY := float64(height)/2 - arrowHeight/2

	dc.MoveTo(arrowX, arrowY)
	dc.LineTo(arrowX+arrowWidth/2, arrowY+arrowHeight)
	dc.LineTo(arrowX+arrowWidth, arrowY)
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Set text color and draw the text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringAnchored(text, float64(width)/2-arrowWidth/2, float64(height)/2, 0.5, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// generateMenuItemImage creates an Ebiten image for a menu item.
func (b *BareBonesUiGenerator) generateMenuItemImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw menu item background
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw menu item text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	textPadding := 10.0                                                   // Small padding from the left edge
	dc.DrawStringAnchored(text, textPadding, float64(height)/2, 0.0, 0.5) // Anchor left-center

	return ebiten.NewImageFromImage(dc.Image())
}

// New: generateMenuImage creates an Ebiten image for the menu background.
func (b *BareBonesUiGenerator) generateMenuImage(
	width, height int,
	bgColor color.RGBA,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw menu background (e.g., a simple rectangle)
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	return ebiten.NewImageFromImage(dc.Image())
}
