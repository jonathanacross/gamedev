package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// ShapeTheme holds the color scheme and font face for the UI components.
type ShapeTheme struct {
	PrimaryAccentColor color.Color
	BackgroundColor    color.Color
	SurfaceColor       color.Color
	TextColor          color.Color
	BorderColor        color.Color
	Face               font.Face
}

// Helper function to adjust the brightness of a color.
func adjustBrightness(c color.Color, factor float64) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(float64(r>>8) * factor),
		G: uint8(float64(g>>8) * factor),
		B: uint8(float64(b>>8) * factor),
		A: uint8(a >> 8),
	}
}

// Helper function to desaturate a color.
func desaturateColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	gray := uint8(float64(r>>8)*0.299 + float64(g>>8)*0.587 + float64(b>>8)*0.114)
	return color.RGBA{R: gray, G: gray, B: gray, A: uint8(a >> 8)}
}

// ShapeRenderer implements the UiRenderer interface.
type ShapeRenderer struct {
	theme ShapeTheme
}

// GenerateButtonImage creates an image for a button in a specific state.
func (r *ShapeRenderer) GenerateButtonImage(width, height int, textContent string, state ButtonState) *ebiten.Image {
	dc := gg.NewContext(width, height)
	var bgColor color.Color
	var textColor color.Color

	switch state {
	case ButtonIdle:
		bgColor = r.theme.SurfaceColor
		textColor = r.theme.TextColor
	case ButtonHover:
		bgColor = adjustBrightness(r.theme.SurfaceColor, 1.2)
		textColor = r.theme.TextColor
	case ButtonPressed:
		bgColor = adjustBrightness(r.theme.SurfaceColor, 0.8)
		textColor = r.theme.TextColor
	case ButtonDisabled:
		bgColor = desaturateColor(r.theme.SurfaceColor)
		textColor = desaturateColor(r.theme.TextColor)
	}

	// Draw button background
	cornerRadius := 9.0
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.FillPreserve()

	// Draw border
	if state != ButtonDisabled {
		dc.SetColor(r.theme.BorderColor)
	} else {
		dc.SetColor(desaturateColor(r.theme.BorderColor))
	}
	dc.SetLineWidth(1)
	dc.Stroke()

	// Draw the text
	dc.SetFontFace(r.theme.Face)
	dc.SetColor(textColor)
	dc.DrawStringAnchored(textContent, float64(width)/2, float64(height)/2, 0.5, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateDropdownImage creates an image for a dropdown button.
func (r *ShapeRenderer) GenerateDropdownImage(width, height int, textContent string, state ButtonState) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw the button background and border
	var bgColor color.Color
	var textColor color.Color
	switch state {
	case ButtonIdle:
		bgColor = r.theme.SurfaceColor
		textColor = r.theme.TextColor
	case ButtonHover:
		bgColor = adjustBrightness(r.theme.SurfaceColor, 1.2)
		textColor = r.theme.TextColor
	case ButtonPressed:
		bgColor = adjustBrightness(r.theme.SurfaceColor, 0.8)
		textColor = r.theme.TextColor
	case ButtonDisabled:
		bgColor = desaturateColor(r.theme.SurfaceColor)
		textColor = desaturateColor(r.theme.TextColor)
	}

	cornerRadius := 9.0
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.FillPreserve()

	if state != ButtonDisabled {
		dc.SetColor(r.theme.BorderColor)
	} else {
		dc.SetColor(desaturateColor(r.theme.BorderColor))
	}
	dc.SetLineWidth(1)
	dc.Stroke()

	// Draw the text
	dc.SetFontFace(r.theme.Face)
	dc.SetColor(textColor)
	dc.DrawStringAnchored(textContent, float64(width)/2-10, float64(height)/2, 0.5, 0.5)

	// Draw the dropdown arrow
	arrowHeight := float64(height) / 7
	arrowWidth := 2 * arrowHeight
	padding := float64(width) * 0.05
	arrowX := float64(width) - arrowWidth - padding
	arrowY := float64(height)/2 - arrowHeight/2

	dc.MoveTo(arrowX, arrowY)
	dc.LineTo(arrowX+arrowWidth/2, arrowY+arrowHeight)
	dc.LineTo(arrowX+arrowWidth, arrowY)
	dc.SetColor(textColor)
	dc.SetLineWidth(2)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateMenuItemImage creates an image for a menu item.
func (r *ShapeRenderer) GenerateMenuItemImage(width, height int, textContent string, state ButtonState) *ebiten.Image {
	dc := gg.NewContext(width, height)
	var bgColor, textColor color.Color

	switch state {
	case ButtonIdle:
		bgColor = r.theme.SurfaceColor
		textColor = r.theme.TextColor
	case ButtonHover:
		bgColor = r.theme.PrimaryAccentColor
		textColor = color.White // High contrast for accent color
	case ButtonPressed:
		bgColor = adjustBrightness(r.theme.PrimaryAccentColor, 0.8)
		textColor = color.White
	case ButtonDisabled:
		bgColor = desaturateColor(r.theme.SurfaceColor)
		textColor = desaturateColor(r.theme.TextColor)
	}

	cornerRadius := 5.0
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.SetColor(bgColor)
	dc.Fill()

	dc.SetFontFace(r.theme.Face)
	dc.SetColor(textColor)
	dc.DrawStringAnchored(textContent, float64(10), float64(height)/2, 0, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateMenuImage creates an image for the menu's background.
func (r *ShapeRenderer) GenerateMenuImage(width, height int) *ebiten.Image {
	dc := gg.NewContext(width, height)

	cornerRadius := 5.0
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.SetColor(r.theme.SurfaceColor)
	dc.FillPreserve()

	dc.SetColor(r.theme.BorderColor)
	dc.SetLineWidth(1)
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateCheckboxImage creates an image for a checkbox.
func (r *ShapeRenderer) GenerateCheckboxImage(width, height int, label string, componentState ButtonState, isChecked bool) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw the checkbox square
	checkboxSize := 15
	checkboxX, checkboxY := 5.0, (float64(height)-float64(checkboxSize))/2

	boxColor := r.theme.SurfaceColor
	borderColor := r.theme.BorderColor
	if componentState == ButtonDisabled {
		boxColor = desaturateColor(boxColor)
		borderColor = desaturateColor(borderColor)
	} else if isChecked {
		boxColor = r.theme.PrimaryAccentColor
		borderColor = r.theme.PrimaryAccentColor
	}

	dc.DrawRectangle(checkboxX, checkboxY, float64(checkboxSize), float64(checkboxSize))
	dc.SetColor(boxColor)
	dc.Fill()

	dc.SetColor(borderColor)
	dc.SetLineWidth(1)
	dc.DrawRectangle(checkboxX, checkboxY, float64(checkboxSize), float64(checkboxSize))
	dc.Stroke()

	// Draw the checkmark
	if isChecked {
		checkColor := r.theme.TextColor
		if componentState == ButtonDisabled {
			checkColor = desaturateColor(checkColor)
		}
		dc.SetColor(checkColor)
		dc.SetLineWidth(2)
		dc.DrawLine(checkboxX+3, checkboxY+float64(checkboxSize)/2, checkboxX+float64(checkboxSize)/2, checkboxY+float64(checkboxSize)-3)
		dc.DrawLine(checkboxX+float64(checkboxSize)/2, checkboxY+float64(checkboxSize)-3, checkboxX+float64(checkboxSize)-3, checkboxY+3)
		dc.Stroke()
	}

	// Draw the label
	dc.SetFontFace(r.theme.Face)
	dc.SetColor(r.theme.TextColor)
	dc.DrawStringAnchored(label, checkboxX+float64(checkboxSize)+10, float64(height)/2, 0, 0.5)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateTextFieldImage creates an image for a text field.
func (r *ShapeRenderer) GenerateTextFieldImage(width, height int, textContent string, componentState ButtonState, isFocused bool, cursorPos int, showCursor bool) *ebiten.Image {
	dc := gg.NewContext(width, height)

	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.SetColor(r.theme.SurfaceColor)
	dc.Fill()

	borderColor := r.theme.BorderColor
	if isFocused {
		borderColor = r.theme.PrimaryAccentColor
	}
	if componentState == ButtonDisabled {
		borderColor = desaturateColor(borderColor)
	}
	dc.SetColor(borderColor)
	dc.SetLineWidth(2)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Stroke()

	dc.SetFontFace(r.theme.Face)
	dc.SetColor(r.theme.TextColor)
	dc.DrawString(textContent, 5, float64(height)/2)

	if isFocused && showCursor {
		cursorX := 5.0
		dc.SetFontFace(r.theme.Face)
		w, _ := dc.MeasureString(textContent[:cursorPos])
		cursorX += w
		dc.SetColor(r.theme.TextColor)
		dc.SetLineWidth(1)
		dc.DrawLine(cursorX, 5, cursorX, float64(height-5))
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateLabelImage creates an image for a static text label.
func (r *ShapeRenderer) GenerateLabelImage(width, height int, textContent string) *ebiten.Image {
	dc := gg.NewContext(width, height)
	dc.SetColor(color.Transparent)
	dc.Clear()

	dc.SetFontFace(r.theme.Face)
	dc.SetColor(r.theme.TextColor)
	dc.DrawString(textContent, 0, float64(height)/2)

	return ebiten.NewImageFromImage(dc.Image())
}

// GenerateContainerImage creates an image for a container's background.
func (r *ShapeRenderer) GenerateContainerImage(width, height int) *ebiten.Image {
	dc := gg.NewContext(width, height)

	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.SetColor(r.theme.BackgroundColor)
	dc.Fill()

	dc.SetColor(r.theme.BorderColor)
	dc.SetLineWidth(2)
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Stroke()

	return ebiten.NewImageFromImage(dc.Image())
}
