package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// ShapeTheme defines the color and font theme for UI elements.
type ShapeTheme struct {
	BackgroundColor    color.RGBA // e.g. dark gray
	PrimaryColor       color.RGBA // your color theme (for buttons, dropdowns)
	OnPrimaryColor     color.RGBA // probably white/near white. color of text/border on primary elements
	AccentColor        color.RGBA // for accents when pressed
	MenuColor          color.RGBA // Color for the menu background
	MenuItemHoverColor color.RGBA // Color for menu items on hover
	Face               font.Face  // The loaded font face
}

// ShapeRenderer renders UI elements using graphic primitives without
// relying on any external image data.
type ShapeRenderer struct {
	theme ShapeTheme
}

// Ensure ShapeRnderer implements the UiRenderer interface
var _ UiRenderer = (*ShapeRenderer)(nil)

// GenerateButtonImage draws a button
func (b *ShapeRenderer) GenerateButtonImage(width, height int, text string, state ButtonState) *ebiten.Image {
	var bgColor, textColor color.RGBA
	switch state {
	case ButtonIdle:
		bgColor = b.theme.PrimaryColor
		textColor = b.theme.OnPrimaryColor
	case ButtonPressed, ButtonHover:
		bgColor = b.theme.AccentColor
		textColor = b.theme.OnPrimaryColor
	case ButtonDisabled:
		bgColor = b.theme.PrimaryColor
		textColor = b.theme.OnPrimaryColor
	}

	dc := gg.NewContext(width, height)

	// Draw button background with rounded corners
	cornerRadius := float64(height) * 0.2
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.FillPreserve()
	// Apply a stroke/border around the rounded rectangle
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(1)
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw button text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringAnchored(text, float64(width)/2, float64(height)/2, 0.5, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateDropdownImage implements UiRenderer.GenerateDropdownImage
func (b *ShapeRenderer) GenerateDropdownImage(width, height int, text string, state ButtonState) *ebiten.Image {
	var bgColor, textColor color.RGBA
	switch state {
	case ButtonIdle:
		bgColor = b.theme.PrimaryColor
		textColor = b.theme.OnPrimaryColor
	case ButtonPressed, ButtonHover:
		bgColor = b.theme.AccentColor
		textColor = b.theme.OnPrimaryColor
	case ButtonDisabled:
		bgColor = b.theme.PrimaryColor
		textColor = b.theme.OnPrimaryColor
	}

	dc := gg.NewContext(width, height)

	// Draw the background of the dropdown button as a simple rectangle
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height)) // No rounded corners here
	dc.FillPreserve()
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(1)
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

// GenerateMenuItemImage implements UiRenderer.GenerateMenuItemImage
func (b *ShapeRenderer) GenerateMenuItemImage(width, height int, text string, state ButtonState) *ebiten.Image {
	var bgColor, textColor color.RGBA
	switch state {
	case ButtonIdle:
		bgColor = b.theme.MenuColor
		textColor = b.theme.OnPrimaryColor
	case ButtonHover:
		bgColor = b.theme.MenuItemHoverColor
		textColor = b.theme.OnPrimaryColor
	case ButtonPressed:
		bgColor = b.theme.AccentColor
		textColor = b.theme.OnPrimaryColor
	case ButtonDisabled:
		bgColor = b.theme.MenuColor // Disabled menu item same as idle
		textColor = b.theme.OnPrimaryColor
	}

	dc := gg.NewContext(width, height)

	cornerRadius := float64(height) * 0.1
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.Fill()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw menu item text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	textPadding := 10.0
	dc.DrawStringAnchored(text, textPadding, float64(height)/2, 0.0, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateMenuImage implements UiRenderer.GenerateMenuImage
func (b *ShapeRenderer) GenerateMenuImage(width, height int) *ebiten.Image {
	// Menu background doesn't change color based on hover/press, just theme color
	bgColor := b.theme.MenuColor

	dc := gg.NewContext(width, height)

	// Draw menu background (e.g., a simple rectangle) with rounded corners
	cornerRadius := float64(10)
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.Fill()

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateCheckboxImage implements UiRenderer.GenerateCheckboxImage
func (b *ShapeRenderer) GenerateCheckboxImage(
	width, height int,
	label string,
	componentState ButtonState,
	isChecked bool,
) *ebiten.Image {
	componentBgColor := b.theme.BackgroundColor
	labelColor := b.theme.OnPrimaryColor
	checkmarkColor := b.theme.OnPrimaryColor

	var boxOutlineColor color.RGBA
	switch componentState {
	case ButtonIdle, ButtonDisabled:
		boxOutlineColor = b.theme.PrimaryColor
	case ButtonPressed, ButtonHover:
		boxOutlineColor = b.theme.AccentColor
	}

	dc := gg.NewContext(width, height)

	// Define fixed left padding for the checkbox square
	const checkboxLeftPadding = 5.0
	checkboxSize := float64(min(width, height)) * 0.8
	checkboxPaddingY := (float64(height) - checkboxSize) / 2

	// Text offset starts after the checkbox square and its right padding
	textOffset := checkboxLeftPadding + checkboxSize + 5.0

	// Draw the checkbox component's overall background
	dc.SetRGBA255(int(componentBgColor.R), int(componentBgColor.G), int(componentBgColor.B), int(componentBgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// Draw the checkbox square outline
	dc.SetRGBA255(int(boxOutlineColor.R), int(boxOutlineColor.G), int(boxOutlineColor.B), int(boxOutlineColor.A))
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(checkboxLeftPadding, checkboxPaddingY, checkboxSize, checkboxSize, 3)
	dc.Stroke()

	// Draw the checkmark if checked
	if isChecked {
		dc.SetRGBA255(int(checkmarkColor.R), int(checkmarkColor.G), int(checkmarkColor.B), int(checkmarkColor.A))
		p1x := checkboxLeftPadding + checkboxSize*0.15
		p1y := checkboxPaddingY + checkboxSize*0.5
		p2x := checkboxLeftPadding + checkboxSize*0.5
		p2y := checkboxPaddingY + checkboxSize*0.85
		p3x := checkboxLeftPadding + checkboxSize*0.85
		p3y := checkboxPaddingY + checkboxSize*0.15

		dc.MoveTo(p1x, p1y)
		dc.LineTo(p2x, p2y)
		dc.LineTo(p3x, p3y)
		dc.SetLineWidth(3)
		dc.Stroke()
	}

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw the label text, positioned after the checkbox
	dc.SetRGBA255(int(labelColor.R), int(labelColor.G), int(labelColor.B), int(labelColor.A))
	dc.DrawStringAnchored(label, textOffset, float64(height)/2, 0.0, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateTextFieldImage renders a text field.
func (b *ShapeRenderer) GenerateTextFieldImage(
	width, height int,
	text string,
	componentState ButtonState,
	isFocused bool,
	cursorPos int,
	showCursor bool,
) *ebiten.Image {
	var bgColor, textColor color.RGBA
	switch componentState {
	case ButtonIdle, ButtonDisabled:
		bgColor = b.theme.PrimaryColor
		textColor = b.theme.OnPrimaryColor
	case ButtonPressed, ButtonHover:
		bgColor = b.theme.AccentColor
		textColor = b.theme.OnPrimaryColor
	}

	dc := gg.NewContext(width, height)

	// Draw text field background
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// Draw a thin border
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.SetLineWidth(1)
	dc.Stroke()

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	textX := 5.0 // Padding from left edge
	textY := float64(height) / 2

	// Draw the text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	dc.DrawStringAnchored(text, textX, textY, 0.0, 0.5) // Anchor left-center

	// Draw blinking cursor if focused and showCursor is true
	if isFocused && showCursor {
		// Calculate cursor X position based on text width up to cursorPos
		textRunes := []rune(text)
		textBeforeCursor := string(textRunes[:min(cursorPos, len(textRunes))])

		cursorXOffset, _ := dc.MeasureString(textBeforeCursor)
		cursorX := textX + cursorXOffset

		dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
		dc.SetLineWidth(1)
		dc.DrawLine(cursorX, textY-float64(b.theme.Face.Metrics().Height)/2, cursorX, textY+float64(b.theme.Face.Metrics().Height)/2)
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateLabelImage implements UiRenderer.GenerateLabelImage
func (b *ShapeRenderer) GenerateLabelImage(
	width, height int,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw the label text
	dc.SetRGBA255(int(b.theme.OnPrimaryColor.R), int(b.theme.OnPrimaryColor.G), int(b.theme.OnPrimaryColor.B), int(b.theme.OnPrimaryColor.A)) // Use OnPrimaryColor for label text
	textX := 5.0                                                                                                                              // Small padding from left edge
	textY := float64(height) / 2
	dc.DrawStringAnchored(text, textX, textY, 0.0, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateContainerImage fills a flat background color for the container.
func (b *ShapeRenderer) GenerateContainerImage(width, height int) *ebiten.Image {
	dc := gg.NewContext(width, height)

	dc.SetRGBA255(int(b.theme.BackgroundColor.R), int(b.theme.BackgroundColor.G), int(b.theme.BackgroundColor.B), int(b.theme.BackgroundColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	return ebiten.NewImageFromImage(dc.Image())
}

// Finds the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
