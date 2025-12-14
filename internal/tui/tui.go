package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TuiPosition represents a position in the grid layout
type TuiPosition struct {
	Row, Col int
}

// TuiManager manages the TUI application state and navigation
type TuiManager struct {
	app        *tview.Application
	focusables [][]tview.Primitive
	currentPos TuiPosition
	allPanels  []tview.Primitive
}

// NewTuiManager creates a new TUI manager instance
func NewTuiManager() *TuiManager {
	return &TuiManager{
		app:        tview.NewApplication(),
		currentPos: TuiPosition{Row: 0, Col: 0},
		allPanels:  make([]tview.Primitive, 0),
	}
}

// SetFocusables defines the grid layout of focusable panels
func (t *TuiManager) SetFocusables(focusables [][]tview.Primitive) {
	t.focusables = focusables

	// Flatten to get all unique panels
	panelMap := make(map[tview.Primitive]bool)
	for _, row := range focusables {
		for _, panel := range row {
			panelMap[panel] = true
		}
	}

	t.allPanels = make([]tview.Primitive, 0, len(panelMap))
	for panel := range panelMap {
		t.allPanels = append(t.allPanels, panel)
	}
}

// UpdateBorders updates all panel borders based on current focus
func (t *TuiManager) UpdateBorders() {
	// Reset all panels to unfocused state
	for _, panel := range t.allPanels {
		if box, ok := panel.(interface {
			SetBorderColor(tcell.Color) *tview.Box
			SetBorderAttributes(tcell.AttrMask) *tview.Box
		}); ok {
			box.SetBorderColor(tcell.ColorWhite)
			box.SetBorderAttributes(tcell.AttrNone)
		}
	}

	// Set current panel to focused state
	if len(t.focusables) > 0 && t.currentPos.Row < len(t.focusables) {
		currentPanel := t.focusables[t.currentPos.Row][t.currentPos.Col]
		if box, ok := currentPanel.(interface {
			SetBorderColor(tcell.Color) *tview.Box
			SetBorderAttributes(tcell.AttrMask) *tview.Box
		}); ok {
			box.SetBorderColor(tcell.ColorGreen)
			box.SetBorderAttributes(tcell.AttrBold)
		}

		t.app.SetFocus(currentPanel)
	}
}

// Navigate moves focus in the specified direction
func (t *TuiManager) Navigate(dr, dc int) {
	newRow := t.currentPos.Row + dr
	newCol := t.currentPos.Col + dc

	// Boundary checks
	if newRow < 0 || newRow >= len(t.focusables) {
		return
	}
	if newCol < 0 || newCol >= len(t.focusables[newRow]) {
		return
	}

	t.currentPos.Row = newRow
	t.currentPos.Col = newCol
	t.UpdateBorders()
}

// SetRoot sets the root primitive for the application
func (t *TuiManager) SetRoot(root tview.Primitive, fullscreen bool) *TuiManager {
	t.app.SetRoot(root, fullscreen)
	return t
}

// SetupNavigation configures keyboard navigation (vim keys + arrows)
func (t *TuiManager) SetupNavigation() *TuiManager {
	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Vim keybindings
		switch event.Rune() {
		case 'h':
			t.Navigate(0, -1)
			return nil
		case 'l':
			t.Navigate(0, 1)
			return nil
		case 'k':
			t.Navigate(-1, 0)
			return nil
		case 'j':
			t.Navigate(1, 0)
			return nil
		}

		// Arrow keys
		switch event.Key() {
		case tcell.KeyLeft:
			t.Navigate(0, -1)
			return nil
		case tcell.KeyRight:
			t.Navigate(0, 1)
			return nil
		case tcell.KeyUp:
			t.Navigate(-1, 0)
			return nil
		case tcell.KeyDown:
			t.Navigate(1, 0)
			return nil
		case tcell.KeyEsc:
			t.app.Stop()
			return nil
		}

		return event
	})
	return t
}

// SetInputCapture allows custom input handling
func (t *TuiManager) SetInputCapture(
	capture func(event *tcell.EventKey) *tcell.EventKey,
) *TuiManager {
	t.app.SetInputCapture(capture)
	return t
}

// Initialize sets up initial focus and borders
func (t *TuiManager) Initialize() *TuiManager {
	t.UpdateBorders()
	return t
}

// Run starts the TUI application
func (t *TuiManager) Run() error {
	return t.app.Run()
}

// Stop stops the TUI application
func (t *TuiManager) Stop() {
	t.app.Stop()
}

// GetApp returns the underlying tview application
func (t *TuiManager) GetApp() *tview.Application {
	return t.app
}

// GetCurrentPosition returns the current focus position
func (t *TuiManager) GetCurrentPosition() TuiPosition {
	return t.currentPos
}

// GetCurrentPanel returns the currently focused panel
func (t *TuiManager) GetCurrentPanel() tview.Primitive {
	if len(t.focusables) > 0 &&
		t.currentPos.Row < len(t.focusables) &&
		t.currentPos.Col < len(t.focusables[t.currentPos.Row]) {
		return t.focusables[t.currentPos.Row][t.currentPos.Col]
	}
	return nil
}
