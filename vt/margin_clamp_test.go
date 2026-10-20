package vt

import "testing"

// Regression tests for https://github.com/charmbracelet/x/issues/974.
//
// DECSTBM and DECSLRM must clamp their exclusive margins to the current
// screen. An application that still draws for the height it had a moment
// ago can send "CSI 1;29 r" after a shrink to 27 rows; without clamping,
// the next reverse index panics inside InsertLineArea.
func TestDECSTBMClampsBottomPastScreen(t *testing.T) {
	t.Parallel()
	term := newTestTerminal(t, 40, 27)
	if _, err := term.Write([]byte("\x1b[1;29r")); err != nil {
		t.Fatal(err)
	}
	scroll := term.scr.ScrollRegion()
	if scroll.Min.Y != 0 || scroll.Max.Y != 27 {
		t.Fatalf("scroll region Y = [%d,%d), want [0,27)", scroll.Min.Y, scroll.Max.Y)
	}
	// Home, then reverse index: must not panic.
	if _, err := term.Write([]byte("top\x1b[H\x1bM")); err != nil {
		t.Fatal(err)
	}
}

func TestDECSTBMTopPastScreenLeavesPriorRegion(t *testing.T) {
	t.Parallel()
	term := newTestTerminal(t, 40, 27)
	if _, err := term.Write([]byte("\x1b[1;2r")); err != nil {
		t.Fatal(err)
	}
	// Top and bottom both past the screen: the sequence is ignored and the
	// prior region survives.
	if _, err := term.Write([]byte("\x1b[28;29r")); err != nil {
		t.Fatal(err)
	}
	scroll := term.scr.ScrollRegion()
	if scroll.Min.Y != 0 || scroll.Max.Y != 2 {
		t.Fatalf("scroll region Y = [%d,%d), want [0,2)", scroll.Min.Y, scroll.Max.Y)
	}
}

func TestDECSTBMLastLineRegionStillValid(t *testing.T) {
	t.Parallel()
	term := newTestTerminal(t, 40, 2)
	// "CSI 2;3 r" on a 2-row screen: the exclusive end clamps to the height,
	// leaving a valid [1,2) region instead of rejecting the sequence.
	if _, err := term.Write([]byte("\x1b[2;3r")); err != nil {
		t.Fatal(err)
	}
	scroll := term.scr.ScrollRegion()
	if scroll.Min.Y != 1 || scroll.Max.Y != 2 {
		t.Fatalf("scroll region Y = [%d,%d), want [1,2)", scroll.Min.Y, scroll.Max.Y)
	}
}

func TestDECSLRMClampsRightPastScreen(t *testing.T) {
	t.Parallel()
	term := newTestTerminal(t, 40, 27)
	if _, err := term.Write([]byte("\x1b[?69h\x1b[1;60s")); err != nil {
		t.Fatal(err)
	}
	scroll := term.scr.ScrollRegion()
	if scroll.Min.X != 0 || scroll.Max.X != 40 {
		t.Fatalf("scroll region X = [%d,%d), want [0,40)", scroll.Min.X, scroll.Max.X)
	}
	// Insert line in the clamped region: must not panic.
	if _, err := term.Write([]byte("\x1b[H\x1b[L")); err != nil {
		t.Fatal(err)
	}
}

func TestDECSLRMLeftPastScreenLeavesPriorRegion(t *testing.T) {
	t.Parallel()
	term := newTestTerminal(t, 40, 27)
	if _, err := term.Write([]byte("\x1b[?69h\x1b[1;10s")); err != nil {
		t.Fatal(err)
	}
	// Left and right both past the screen: the sequence is ignored and the
	// prior region survives.
	if _, err := term.Write([]byte("\x1b[50;60s")); err != nil {
		t.Fatal(err)
	}
	scroll := term.scr.ScrollRegion()
	if scroll.Min.X != 0 || scroll.Max.X != 10 {
		t.Fatalf("scroll region X = [%d,%d), want [0,10)", scroll.Min.X, scroll.Max.X)
	}
}
