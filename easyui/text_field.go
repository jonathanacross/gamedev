package main

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// keyRepeatTracker manages the state for key repeat functionality.
type keyRepeatTracker struct {
	lastKeyEvent   time.Time
	key            ebiten.Key
	repeatCount    int
	initialDelay   time.Duration
	repeatInterval time.Duration
}

// newKeyRepeatTracker initializes a new tracker with default settings.
func newKeyRepeatTracker() *keyRepeatTracker {
	return &keyRepeatTracker{
		initialDelay:   400 * time.Millisecond,
		repeatInterval: 50 * time.Millisecond,
	}
}

// isReadyToRepeat checks if the key is ready for a new repeat event.
func (krt *keyRepeatTracker) isReadyToRepeat() bool {
	if krt.key == ebiten.Key(0) {
		return false // No key is being held
	}

	elapsed := time.Since(krt.lastKeyEvent)
	if krt.repeatCount == 0 {
		return elapsed >= krt.initialDelay
	} else {
		return elapsed >= krt.repeatInterval
	}
}

// TextField represents a single-line text input field.
type TextField struct {
	interactiveComponent
	Text       string
	isFocused  bool
	cursorPos  int
	blinkTimer float64
	renderer   UiRenderer
	keyRepeat  *keyRepeatTracker
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
		keyRepeat: newKeyRepeatTracker(),
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
		tf.handleTextInput()

		// Handle left/right arrow keys for cursor movement (without repeat)
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
		tf.regenerateImages(tf.isFocused, showCursor)
	}
}

// handleTextInput handles all character and backspace input, including key repeat.
func (tf *TextField) handleTextInput() {
	// Handle character input from keyboard
	for _, char := range ebiten.InputChars() {
		tf.Text = tf.Text[:tf.cursorPos] + string(char) + tf.Text[tf.cursorPos:]
		tf.cursorPos++
	}

	// Handle backspace with key repeat
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
		if tf.cursorPos > 0 {
			tf.Text = tf.Text[:tf.cursorPos-1] + tf.Text[tf.cursorPos:]
			tf.cursorPos--
			tf.keyRepeat.lastKeyEvent = time.Now()
			tf.keyRepeat.key = ebiten.KeyBackspace
			tf.keyRepeat.repeatCount = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		if tf.keyRepeat.key == ebiten.KeyBackspace && tf.keyRepeat.isReadyToRepeat() {
			if tf.cursorPos > 0 {
				tf.Text = tf.Text[:tf.cursorPos-1] + tf.Text[tf.cursorPos:]
				tf.cursorPos--
				tf.keyRepeat.lastKeyEvent = time.Now()
				tf.keyRepeat.repeatCount++
			}
		}
	} else if tf.keyRepeat.key == ebiten.KeyBackspace && !ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		// Reset state if key is released
		tf.keyRepeat.key = ebiten.Key(0)
		tf.keyRepeat.repeatCount = 0
	}
}

// regenerateImages updates the cached images for the text field's states.
func (tf *TextField) regenerateImages(isFocused bool, showCursor bool) {
	tf.idleImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonIdle, isFocused, tf.cursorPos, showCursor)

	tf.hoverImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonHover, isFocused, tf.cursorPos, showCursor)

	tf.pressedImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonPressed, isFocused, tf.cursorPos, showCursor)

	tf.disabledImg = tf.renderer.GenerateTextFieldImage(
		tf.Bounds.Dx(), tf.Bounds.Dy(),
		tf.Text, ButtonDisabled, isFocused, tf.cursorPos, showCursor)
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
	// Use a centralized approach to manage focus
	rootUi := tf.GetRootUi()
	if rootUi != nil {
		cx, cy := ebiten.CursorPosition()
		clickedInside := ContainsPoint(tf, cx, cy)
		if clickedInside {
			rootUi.SetFocusedComponent(tf.self)
		} else {
			rootUi.SetFocusedComponent(nil)
		}
	}
}

// Focus is called by the Ui when this component gains focus.
func (tf *TextField) Focus() {
	if !tf.isFocused {
		log.Println("TextField focused!")
		tf.isFocused = true
		tf.cursorPos = len(tf.Text)
		tf.blinkTimer = 0
		tf.regenerateImages(tf.isFocused, tf.isFocused)
	}
}

// Unfocus is called by the Ui when this component loses focus.
func (tf *TextField) Unfocus() {
	if tf.isFocused {
		log.Println("TextField unfocused.")
		tf.isFocused = false
		tf.regenerateImages(tf.isFocused, tf.isFocused)
	}
}
