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

func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}
	return x
}

// A simple helper function to adjust the brightness of a color.
func adjustBrightness(c color.Color, factor float64) color.Color {
	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(clamp(float64(r>>8)*factor, 0, 255)),
		G: uint8(clamp(float64(g>>8)*factor, 0, 255)),
		B: uint8(clamp(float64(b>>8)*factor, 0, 255)),
		A: uint8(a >> 8),
	}
}

// A simple helper function to desaturate a color.
func desaturateColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	gray := uint8(clamp(float64(r>>8)*0.299+float64(g>>8)*0.587+float64(b>>8)*0.114, 0, 255))
	return color.RGBA{R: gray, G: gray, B: gray, A: uint8(a >> 8)}
}

// ShapeRenderer implements the UiRenderer interface.
type ShapeRenderer struct {
	theme ShapeTheme
}

// getStateColor returns a modified color based on the component's state.
func (r *ShapeRenderer) getStateColor(baseColor color.Color, state ButtonState) color.Color {
	switch state {
	case ButtonIdle:
		return baseColor
	case ButtonHover:
		return adjustBrightness(baseColor, 1.2)
	case ButtonPressed:
		return adjustBrightness(baseColor, 0.8)
	case ButtonDisabled:
		return desaturateColor(baseColor)
	default:
		return baseColor
	}
}

// getButtonColors returns the background and text color for a standard button or dropdown.
func (r *ShapeRenderer) getButtonColors(state ButtonState) (bgColor, textColor color.Color) {
	bgColor = r.getStateColor(r.theme.SurfaceColor, state)
	textColor = r.getStateColor(r.theme.TextColor, state)
	return
}

// getCheckableColors returns the colors for a checkbox or radio button based on its state and checked status.
func (r *ShapeRenderer) getCheckableColors(componentState ButtonState, isChecked bool) (boxColor, borderColor, checkDotColor color.Color) {
	if isChecked {
		boxColor = r.getStateColor(r.theme.PrimaryAccentColor, componentState)
		borderColor = r.getStateColor(r.theme.PrimaryAccentColor, componentState)
		checkDotColor = r.getStateColor(r.theme.TextColor, componentState)
	} else { // unchecked
		boxColor = r.getStateColor(r.theme.SurfaceColor, componentState)
		borderColor = r.getStateColor(r.theme.BorderColor, componentState)
		checkDotColor = r.getStateColor(r.theme.TextColor, componentState)
	}
	return
}

// GenerateButtonImage creates an image for a button in a specific state.
func (r *ShapeRenderer) GenerateButtonImage(width, height int, textContent string, state ButtonState, isChecked bool) *ebiten.Image {
	dc := gg.NewContext(width, height)
	var bgColor color.Color
	if isChecked {
		bgColor = r.getStateColor(r.theme.PrimaryAccentColor, state)
	} else {
		bgColor = r.getStateColor(r.theme.SurfaceColor, state)
	}
	textColor := r.getStateColor(r.theme.TextColor, state)

	// Draw button background
	cornerRadius := 9.0
	dc.SetColor(bgColor)
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.FillPreserve()

	// Draw border
	borderColor := r.getStateColor(r.theme.BorderColor, state)
	dc.SetColor(borderColor)
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
	bgColor, textColor := r.getButtonColors(state)

	// Draw the button background and border
	cornerRadius := 9.0
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.SetColor(bgColor)
	dc.FillPreserve()

	borderColor := r.getStateColor(r.theme.BorderColor, state)
	dc.SetColor(borderColor)
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
		textColor = r.theme.TextColor
	case ButtonPressed:
		bgColor = adjustBrightness(r.theme.PrimaryAccentColor, 0.8)
		textColor = r.theme.TextColor
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

	boxColor, borderColor, checkColor := r.getCheckableColors(componentState, isChecked)

	dc.DrawRectangle(checkboxX, checkboxY, float64(checkboxSize), float64(checkboxSize))
	dc.SetColor(boxColor)
	dc.Fill()

	dc.SetColor(borderColor)
	dc.SetLineWidth(1)
	dc.DrawRectangle(checkboxX, checkboxY, float64(checkboxSize), float64(checkboxSize))
	dc.Stroke()

	// Draw the checkmark
	if isChecked {
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

// GenerateRadioButtonImage creates an image for a radio button.
func (r *ShapeRenderer) GenerateRadioButtonImage(width, height int, label string, componentState ButtonState, isChecked bool) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw the radio button circle
	circleRadius := 7.0
	circleX, circleY := 5.0+circleRadius, (float64(height)/2)+0.5 // Center the circle vertically in the component, +0.5 to fix pixel alignment

	circleColor, borderColor, dotColor := r.getCheckableColors(componentState, isChecked)

	dc.DrawCircle(circleX, circleY, circleRadius)
	dc.SetColor(circleColor)
	dc.Fill()

	dc.SetColor(borderColor)
	dc.SetLineWidth(1)
	dc.DrawCircle(circleX, circleY, circleRadius)
	dc.Stroke()

	// Draw the inner dot if checked
	if isChecked {
		dotColor = r.getStateColor(r.theme.TextColor, componentState)
		dc.SetColor(dotColor)
		dc.DrawCircle(circleX, circleY, circleRadius*0.4) // Smaller inner circle
		dc.Fill()
	}

	// Draw the label
	dc.SetFontFace(r.theme.Face)
	dc.SetColor(r.theme.TextColor)
	dc.DrawStringAnchored(label, circleX+circleRadius+10, float64(height)/2, 0, 0.5)

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
