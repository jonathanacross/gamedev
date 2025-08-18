package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// TextField represents a single-line text input field.
type TextField struct {
	interactiveComponent
	Text       string
	isFocused  bool
	cursorPos  int
	blinkTimer float64
	renderer   UiRenderer
}

// NewTextField creates a new TextField instance.
func NewTextField(x, y, width, height int, initialText string, renderer UiRenderer) *TextField {
	// Initial image generation, assuming not focused and no cursor initially
	idle := renderer.GenerateTextFieldImage(width, height, initialText, ButtonIdle, false, 0, false)
	hover := renderer.GenerateTextFieldImage(width, height, initialText, ButtonHover, false, 0, false)
	// For initial focused/pressed image, assume cursor is visible if focused
	pressed := renderer.GenerateTextFieldImage(width, height, initialText, ButtonPressed, true, len(initialText), true)
	disabled := renderer.GenerateTextFieldImage(width, height, initialText, ButtonDisabled, false, 0, false)

	tf := &TextField{
		Text:      initialText,
		cursorPos: len(initialText), // Cursor at end initially
		renderer:  renderer,
	}
	tf.interactiveComponent = NewInteractiveComponent(x, y, width, height,
		idle, pressed, hover, disabled, tf)

	return tf
}

// Update handles text input, cursor blinking, and general interactive component updates.
func (tf *TextField) Update() {
	tf.interactiveComponent.Update() // Handle hover/pressed states

	// Handle cursor blinking
	tf.blinkTimer += ebiten.ActualTPS() / 60.0
	if tf.blinkTimer >= 120 { // Blink about every 2 seconds (120 frames at 60 TPS)
		tf.blinkTimer = 0
	}

	// Determine if cursor should be shown based on blinkTimer
	showCursor := tf.isFocused && (tf.blinkTimer < 60)

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

		// Regenerate the image to reflect new text/cursor
		tf.idleImg = tf.renderer.GenerateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.Text, ButtonIdle, tf.isFocused, tf.cursorPos, showCursor)

		tf.hoverImg = tf.renderer.GenerateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.Text, ButtonHover, tf.isFocused, tf.cursorPos, showCursor)

		tf.pressedImg = tf.renderer.GenerateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.Text, ButtonPressed, tf.isFocused, tf.cursorPos, showCursor)

		tf.disabledImg = tf.renderer.GenerateTextFieldImage(
			tf.Bounds.Dx(), tf.Bounds.Dy(),
			tf.Text, ButtonDisabled, tf.isFocused, tf.cursorPos, showCursor)
	}
}

// Draw draws the text field's current state image to the screen.
func (tf *TextField) Draw(screen *ebiten.Image) {
	tf.interactiveComponent.Draw(screen)
}

// HandlePress sets the interactive component to the pressed state.
func (tf *TextField) HandlePress() {
	tf.interactiveComponent.HandlePress()
}

// HandleRelease resets the interactive component's state after a mouse release.
func (tf *TextField) HandleRelease() {
	tf.interactiveComponent.HandleRelease()
}

// HandleClick toggles the focus state of the text field.
func (tf *TextField) HandleClick() {
	// If the click happened on this text field, focus it.
	cx, cy := ebiten.CursorPosition()
	clickedInside := ContainsPoint(tf, cx, cy)

	if clickedInside {
		if !tf.isFocused {
			log.Println("TextField focused!")
			tf.isFocused = true
			tf.cursorPos = len(tf.Text)
			tf.blinkTimer = 0
		}
	} else {
		// If click was outside, lose focus
		if tf.isFocused {
			log.Println("TextField unfocused.")
			tf.isFocused = false
		}
	}

	// Re-render image to reflect focus change and current cursor state
	// Determine showCursor based on new focus state (and potentially reset blink timer if focused)
	showCursor := tf.isFocused // Initially assume true if focused, will be handled by update loop's blink logic

	tf.idleImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonIdle, tf.isFocused, tf.cursorPos, showCursor)
	tf.hoverImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonHover, tf.isFocused, tf.cursorPos, showCursor)
	tf.pressedImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonPressed, tf.isFocused, tf.cursorPos, showCursor)
	tf.disabledImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonDisabled, tf.isFocused, tf.cursorPos, showCursor)
}
