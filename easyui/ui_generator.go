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
	dc.FillPreserve()

	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(2)
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringWrapped(buttonText, float64(width)/2, float64(height)/2, 0.5, 1.0, float64(width)-10, 1.5, gg.AlignCenter)

	img := ebiten.NewImageFromImage(dc.Image())
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
	textWidth, _ := dc.MeasureString(text)
	textX := padding
	maxTextWidth := float64(width) - (arrowWidth + padding*2)

	if textWidth > maxTextWidth {
		dc.DrawString(text, textX, float64(height)/2+dc.FontHeight()/3)
	} else {
		centeredTextX := (maxTextWidth-textWidth)/2 + padding
		dc.DrawString(text, centeredTextX, float64(height)/2+dc.FontHeight()/3)
	}

	img := ebiten.NewImageFromImage(dc.Image())
	return img
}
