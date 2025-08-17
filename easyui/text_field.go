package main

import (
	"log"
	// For cursor blinking
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	// Required for image.Rectangle, image.Point
)

// TextField represents a single-line text input field.
type TextField struct {
	interactiveComponent                       // Embed for common state handling (hover, press)
	Text                 string                // Current text content
	isFocused            bool                  // True if this text field is currently focused and receiving input
	cursorPos            int                   // Current cursor position within the text
	blinkTimer           float64               // Timer for cursor blinking
	uiGenerator          *BareBonesUiGenerator // Reference to generator for image regeneration
}

// NewTextField creates a new TextField instance.
// It is now a standalone function.
func NewTextField(x, y, width, height int, initialText string, uiGen *BareBonesUiGenerator) *TextField {
	// Generate images for different states of the text field
	idle := uiGen.generateTextFieldImage(width, height, uiGen.theme.PrimaryColor, uiGen.theme.OnPrimaryColor, initialText, false, 0, 0)
	hover := uiGen.generateTextFieldImage(width, height, uiGen.theme.AccentColor, uiGen.theme.OnPrimaryColor, initialText, false, 0, 0)
	pressed := uiGen.generateTextFieldImage(width, height, uiGen.theme.AccentColor, uiGen.theme.OnPrimaryColor, initialText, true, len(initialText), 0) // Focused state on press
	disabled := uiGen.generateTextFieldImage(width, height, uiGen.theme.PrimaryColor, uiGen.theme.OnPrimaryColor, initialText, false, 0, 0)

	tf := &TextField{
		interactiveComponent: NewInteractiveComponent(x, y, width, height,
			idle, pressed, hover, disabled),
		Text:        initialText,
		cursorPos:   len(initialText), // Cursor at end initially
		uiGenerator: uiGen,
	}
	return tf
}

// Update handles text input, cursor blinking, and general interactive component updates.
func (tf *TextField) Update() {
	tf.interactiveComponent.Update() // Handle hover/pressed states

	// Handle cursor blinking
	tf.blinkTimer += ebiten.ActualTPS() / 60.0 // Assuming 60 TPS for consistent blink rate
	if tf.blinkTimer >= 120 {                  // Roughly blink every 2 seconds (120 frames at 60 TPS)
		tf.blinkTimer = 0
	}

	if tf.isFocused {
		// Handle character input
		for _, char := range ebiten.InputChars() {
			if char >= ' ' && char <= '~' { // Basic printable ASCII characters
				tf.Text = tf.Text[:tf.cursorPos] + string(char) + tf.Text[tf.cursorPos:]
				tf.cursorPos++
			}
		}

		// Handle backspace
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			if tf.cursorPos > 0 {
				tf.Text = tf.Text[:tf.cursorPos-1] + tf.Text[tf.cursorPos:]
				tf.cursorPos--
			}
		}

		// Handle Delete (Forward Delete) - Ebiten doesn't have a direct "delete" key code easily,
		// but we can simulate it by moving cursor right and then backspace.
		// For simplicity, let's omit explicit "Delete" key handling for now, focusing on Backspace.

		// Handle left/right arrow keys for cursor movement
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if tf.cursorPos > 0 {
				tf.cursorPos--
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if tf.cursorPos < len(tf.Text) {
				tf.cursorPos++
			}
		}

		// Regenerate the image to reflect new text/cursor (always use idle color for base text field)
		tf.idleImg = tf.uiGenerator.generateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.uiGenerator.theme.PrimaryColor,   // Background color of the field
			tf.uiGenerator.theme.OnPrimaryColor, // Text color
			tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))

		// Update other states too, so hover/pressed effects apply correctly even with new text
		tf.hoverImg = tf.uiGenerator.generateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.uiGenerator.theme.AccentColor, // Background color changes on hover
			tf.uiGenerator.theme.OnPrimaryColor,
			tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))

		tf.pressedImg = tf.uiGenerator.generateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.uiGenerator.theme.AccentColor, // Background color changes on pressed
			tf.uiGenerator.theme.OnPrimaryColor,
			tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))

		tf.disabledImg = tf.uiGenerator.generateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.uiGenerator.theme.PrimaryColor, // Stays primary for disabled
			tf.uiGenerator.theme.OnPrimaryColor,
			tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))
	}
}

// Draw draws the text field's current state image to the screen.
func (tf *TextField) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(tf.Bounds.Min.X), float64(tf.Bounds.Min.Y))
	screen.DrawImage(tf.GetCurrentStateImage(), op)
}

// HandlePress calls the embedded interactiveComponent's HandlePress method.
func (tf *TextField) HandlePress() {
	tf.interactiveComponent.HandlePress()
}

// HandleRelease calls the embedded interactiveComponent's HandleRelease method.
func (tf *TextField) HandleRelease() {
	tf.interactiveComponent.HandleRelease()
}

// HandleClick toggles the focus state of the text field.
func (tf *TextField) HandleClick() {
	// If the click happened on THIS text field, focus it.
	cx, cy := ebiten.CursorPosition()
	if ContainsPoint(tf.Bounds, cx, cy) {
		if !tf.isFocused {
			log.Println("TextField focused!")
			tf.isFocused = true
			tf.cursorPos = len(tf.Text) // Move cursor to end on initial click
		}
	} else {
		// If click was outside, lose focus
		if tf.isFocused {
			log.Println("TextField unfocused.")
			tf.isFocused = false
		}
	}

	// Re-render image to reflect focus change
	tf.idleImg = tf.uiGenerator.generateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.uiGenerator.theme.PrimaryColor,
		tf.uiGenerator.theme.OnPrimaryColor,
		tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))
	tf.hoverImg = tf.uiGenerator.generateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.uiGenerator.theme.AccentColor, // Use accent for hover
		tf.uiGenerator.theme.OnPrimaryColor,
		tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))
	tf.pressedImg = tf.uiGenerator.generateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.uiGenerator.theme.AccentColor, // Use accent for pressed
		tf.uiGenerator.theme.OnPrimaryColor,
		tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))
	tf.disabledImg = tf.uiGenerator.generateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.uiGenerator.theme.PrimaryColor, // Stays primary for disabled
		tf.uiGenerator.theme.OnPrimaryColor,
		tf.Text, tf.isFocused, tf.cursorPos, int(tf.blinkTimer))
}
