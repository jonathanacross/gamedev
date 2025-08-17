package main

import (
	"image"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	// For font metrics
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
// It acts as a factory for creating all themed UI elements.
type BareBonesUiGenerator struct {
	theme BareBonesTheme
}

// NewMenu creates a new Menu instance. Moved here for consistency.
func (b *BareBonesUiGenerator) NewMenu(x, y, width int, parentUi *Ui) *Menu {
	m := &Menu{
		component: component{
			Bounds: image.Rectangle{
				Min: image.Point{X: x, Y: y},
				Max: image.Point{X: x + width, Y: y}, // Max Y will be adjusted later
			},
		},
		items:       []*MenuItem{},
		isVisible:   false,
		theme:       b.theme, // Use the generator's theme
		uiGenerator: b,       // Pass reference to this generator
		parentUi:    parentUi,
		justOpened:  false,
	}
	return m
}

// NewButton creates a new Button instance with the specified dimensions, label, and theme.
// It also passes a reference to itself for dynamic updates.
func (b *BareBonesUiGenerator) NewButton(x, y, width, height int, label string) *Button {
	idle := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label)
	pressed := b.generateButtonImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, label)
	hover := b.generateButtonImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, label)
	disabled := b.generateButtonImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, label) // Example: disabled state image

	return &Button{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idle, pressed, hover, disabled),
		Label:                label,
		onClick:              nil, // Click handler set separately
		uiGenerator:          b,   // Pass reference to this generator
	}
}

// NewDropDown creates a new DropDown instance, generating its state-specific images.
// The `parentUi` parameter is removed as it's not directly used by DropDown.
func (b *BareBonesUiGenerator) NewDropDown(x, y, width, height int, initialLabel string, menu *Menu) *DropDown {
	// Generate specific dropdown images using this generator's methods
	idleImg := b.generateDropdownImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, initialLabel)
	hoverImg := b.generateDropdownImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, initialLabel)
	pressedImg := b.generateDropdownImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, initialLabel)   // Often same as hover for dropdown pressed
	disabledImg := b.generateDropdownImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, initialLabel) // Example: darker version

	return &DropDown{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg),
		Label:                initialLabel,
		SelectedOption:       initialLabel,
		menu:                 menu,
		uiGenerator:          b, // Pass generator reference
	}
}

// NewMenuItem creates a new MenuItem instance, generating its state-specific images.
// It also passes a reference to itself for dynamic updates.
func (b *BareBonesUiGenerator) NewMenuItem(x, y, width, height int, label string) *MenuItem {
	idleImg := b.generateMenuItemImage(width, height, b.theme.MenuColor, b.theme.OnPrimaryColor, label)
	hoverImg := b.generateMenuItemImage(width, height, b.theme.MenuItemHoverColor, b.theme.OnPrimaryColor, label)
	pressedImg := b.generateMenuItemImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, label)
	disabledImg := idleImg // Default disabled to idle, can be customized

	return &MenuItem{
		interactiveComponent: NewInteractiveComponent(x, y, width, height, idleImg, pressedImg, hoverImg, disabledImg),
		Label:                label,
		onClick:              nil, // Click handler set separately
		uiGenerator:          b,   // Pass reference to this generator
	}
}

// NewCheckbox creates a new Checkbox instance, generating all its state-specific images.
// Updated signature to take width and height explicitly.
func (b *BareBonesUiGenerator) NewCheckbox(x, y, width, height int, label string, initialChecked bool) *Checkbox {
	// Colors that stay constant for the overall component background and label text
	componentBgColor := b.theme.BackgroundColor
	labelColor := b.theme.OnPrimaryColor
	checkmarkColor := b.theme.OnPrimaryColor // Checkmark color remains constant

	// Generate images for unchecked states
	uncheckedIdle := b.generateCheckboxImage(width, height, componentBgColor, b.theme.PrimaryColor, checkmarkColor, labelColor, false, label)
	uncheckedPressed := b.generateCheckboxImage(width, height, componentBgColor, b.theme.AccentColor, checkmarkColor, labelColor, false, label) // Box outline changes to AccentColor
	uncheckedHover := b.generateCheckboxImage(width, height, componentBgColor, b.theme.AccentColor, checkmarkColor, labelColor, false, label)   // Box outline changes to AccentColor
	uncheckedDisabled := b.generateCheckboxImage(width, height, componentBgColor, b.theme.PrimaryColor, checkmarkColor, labelColor, false, label)

	// Generate images for checked states
	checkedIdle := b.generateCheckboxImage(width, height, componentBgColor, b.theme.PrimaryColor, checkmarkColor, labelColor, true, label)
	checkedPressed := b.generateCheckboxImage(width, height, componentBgColor, b.theme.AccentColor, checkmarkColor, labelColor, true, label) // Box outline changes to AccentColor
	checkedHover := b.generateCheckboxImage(width, height, componentBgColor, b.theme.AccentColor, checkmarkColor, labelColor, true, label)   // Box outline changes to AccentColor
	checkedDisabled := b.generateCheckboxImage(width, height, componentBgColor, b.theme.PrimaryColor, checkmarkColor, labelColor, true, label)

	stateImages := CheckboxStateImages{
		UncheckedIdle:     uncheckedIdle,
		UncheckedPressed:  uncheckedPressed,
		UncheckedHover:    uncheckedHover,
		UncheckedDisabled: uncheckedDisabled,
		CheckedIdle:       checkedIdle,
		CheckedPressed:    checkedPressed,
		CheckedHover:      checkedHover,
		CheckedDisabled:   checkedDisabled,
	}

	// Use the NewCheckbox constructor
	return NewCheckbox(x, y, width, height, label, initialChecked, b, stateImages)
}

// NewTextField creates a new TextField instance.
func (b *BareBonesUiGenerator) NewTextField(x, y, width, height int, initialText string) *TextField {
	// Generate images for different states of the text field
	idle := b.generateTextFieldImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, initialText, false, 0, 0)
	hover := b.generateTextFieldImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, initialText, false, 0, 0)
	pressed := b.generateTextFieldImage(width, height, b.theme.AccentColor, b.theme.OnPrimaryColor, initialText, true, len(initialText), 0) // Focused state on press
	disabled := b.generateTextFieldImage(width, height, b.theme.PrimaryColor, b.theme.OnPrimaryColor, initialText, false, 0, 0)

	return NewTextField(x, y, width, height, initialText, b, idle, pressed, hover, disabled)
}

// NewLabel creates a new Label instance.
func (b *BareBonesUiGenerator) NewLabel(x, y, width, height int, text string) *Label {
	// Labels are static, so only an idle image is needed.
	labelImage := b.generateLabelImage(width, height, b.theme.OnPrimaryColor, text) // Text color from theme
	return NewLabel(x, y, width, height, text, b, labelImage)
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

// generateDropdownImage creates an Ebiten image for a dropdown's specific state.
// No rounded corners as requested.
func (b *BareBonesUiGenerator) generateDropdownImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
) *ebiten.Image {
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

// generateMenuItemImage creates an Ebiten image for a menu item.
// Rounded corners for menu items.
func (b *BareBonesUiGenerator) generateMenuItemImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw menu item background with rounded corners for consistency
	cornerRadius := float64(height) * 0.1 // Slightly smaller radius for menu items
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
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

// generateMenuImage creates an Ebiten image for the menu background.
// Rounded corners for the overall menu background.
func (b *BareBonesUiGenerator) generateMenuImage(
	width, height int,
	bgColor color.RGBA,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Draw menu background (e.g., a simple rectangle) with rounded corners
	cornerRadius := float64(10) // Small fixed radius for the overall menu background
	dc.SetRGBA255(int(bgColor.R), int(bgColor.G), int(bgColor.B), int(bgColor.A))
	dc.DrawRoundedRectangle(0, 0, float64(width), float64(height), cornerRadius)
	dc.Fill()

	return ebiten.NewImageFromImage(dc.Image())
}

// generateCheckboxImage creates an Ebiten image for a checkbox, including the box, checkmark, and label.
// It now takes specific colors for the component background, box outline, checkmark, and label text.
func (b *BareBonesUiGenerator) generateCheckboxImage(
	width, height int,
	componentBgColor color.RGBA, // The color of the component's full background
	boxOutlineColor color.RGBA, // Color for the square outline
	checkmarkColor color.RGBA, // Color for the checkmark
	labelColor color.RGBA, // Color for the text label
	isChecked bool,
	label string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Define fixed left padding for the checkbox square
	const checkboxLeftPadding = 5.0 // Adjust as needed
	checkboxSize := float64(min(width, height)) * 0.8
	checkboxPaddingY := (float64(height) - checkboxSize) / 2

	// Text offset starts after the checkbox square and its right padding
	textOffset := checkboxLeftPadding + checkboxSize + 5.0 // Small gap between box and text

	// Draw the checkbox component's overall background
	dc.SetRGBA255(int(componentBgColor.R), int(componentBgColor.G), int(componentBgColor.B), int(componentBgColor.A))
	dc.DrawRectangle(0, 0, float64(width), float64(height))
	dc.Fill()

	// Draw the checkbox square outline
	dc.SetRGBA255(int(boxOutlineColor.R), int(boxOutlineColor.G), int(boxOutlineColor.B), int(boxOutlineColor.A))
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(checkboxLeftPadding, checkboxPaddingY, checkboxSize, checkboxSize, 3) // Small rounded corners for the box
	dc.Stroke()

	// Draw the checkmark if checked
	if isChecked {
		dc.SetRGBA255(int(checkmarkColor.R), int(checkmarkColor.G), int(checkmarkColor.B), int(checkmarkColor.A))
		// Points are relative to the checkbox square's top-left (checkboxLeftPadding, checkboxPaddingY)
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
	dc.SetRGBA255(int(labelColor.R), int(labelColor.G), int(labelColor.B), int(labelColor.A)) // Use labelColor for text color
	dc.DrawStringAnchored(label, textOffset, float64(height)/2, 0.0, 0.5)                     // Anchor left-center after checkbox

	return ebiten.NewImageFromImage(dc.Image())
}

// generateTextFieldImage creates an Ebiten image for a text field, including its background, text, and optional cursor.
func (b *BareBonesUiGenerator) generateTextFieldImage(
	width, height int,
	bgColor, textColor color.RGBA,
	text string,
	isFocused bool,
	cursorPos int,
	blinkTimer int, // Used for cursor blinking logic
) *ebiten.Image {
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

	// Draw blinking cursor if focused and visible (blinkTimer determines visibility)
	if isFocused && (blinkTimer < 60) { // Cursor visible for first half of the blink cycle
		// Calculate cursor X position based on text width up to cursorPos
		textRunes := []rune(text)
		textBeforeCursor := string(textRunes[:min(cursorPos, len(textRunes))]) // Handle cursor at end or beyond

		// Corrected: Capture both return values of MeasureString, use only width
		cursorXOffset, _ := dc.MeasureString(textBeforeCursor)
		cursorX := textX + cursorXOffset

		dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
		dc.SetLineWidth(1)
		dc.DrawLine(cursorX, textY-float64(b.theme.Face.Metrics().Height)/2, cursorX, textY+float64(b.theme.Face.Metrics().Height)/2)
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

// generateLabelImage creates an Ebiten image for a static text label.
func (b *BareBonesUiGenerator) generateLabelImage(
	width, height int,
	textColor color.RGBA,
	text string,
) *ebiten.Image {
	dc := gg.NewContext(width, height)

	// Labels typically have a transparent background, or the background of their parent.
	// We just need to draw the text.
	// You could fill with b.theme.BackgroundColor here if you want a solid background for labels.

	if b.theme.Face != nil {
		dc.SetFontFace(b.theme.Face)
	}

	// Draw the label text
	dc.SetRGBA255(int(textColor.R), int(textColor.G), int(textColor.B), int(textColor.A))
	textX := 5.0 // Small padding from left edge
	textY := float64(height) / 2
	dc.DrawStringAnchored(text, textX, textY, 0.0, 0.5) // Anchor left-center

	return ebiten.NewImageFromImage(dc.Image())
}

// Helper function to find the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
