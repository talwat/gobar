// A super minimalistic library for progressbars in golang.
package gobar

import (
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

// The bar class, do not initialize directly, instead, use NewBar.
type Bar struct {
	Value         int64
	Total         int64
	Description   string
	Done          string
	Throttle      int64
	lastDisplayed time.Time
}

// Makes a new bar and initializes it.
// This will call bar.Display(), so be careful not to call it to early in your program.
//
//nolint:gomnd // Magic numbers for configuration.
func NewBar(value int64, total int64, description string, done string) *Bar {
	bar := Bar{Value: value, Total: total, Description: description, Done: done, Throttle: 65}
	bar.Init()

	return &bar
}

// Gets the terminal size.
func getTermSize() int {
	// Fallback size in case of error.
	const fallback = 80

	width, _, err := term.GetSize(0)
	if err != nil {
		return fallback
	}

	return width
}

// Initializes a new bar and renders it in a 0% progress state.
func (bar *Bar) Init() {
	bar.Display(true)
}

// Makes some newlines and finishes the bar.
// This is usually called automatically once the bar is full.
func (bar *Bar) Finish() {
	bar.Value = bar.Total
	bar.Display(true)
	fmt.Fprintf(os.Stderr, "\n%s\n", bar.Done)
}

// Increments the bar by num and displays it.
func (bar *Bar) Increment(num int64) {
	bar.Add(num)
	bar.Display(false)

	if bar.Value >= bar.Total {
		bar.Finish()
	}
}

// Adds an int to the value of the bar.
// This also adds a rate aswell with the same number.
// Rate defines whether to add a rate.
func (bar *Bar) Add(num int64) {
	num64 := num
	bar.Value += num64
}

// Display the bar.
// This does not change the bar, or increment the value.
func (bar *Bar) Display(ignoreThrottle bool) {
	// Are we running in a terminal? If not, exit.
	if !term.IsTerminal(0) {
		return
	}

	// If we have displayed in the last <throttle> seconds, return.
	if time.Since(bar.lastDisplayed).Milliseconds() < bar.Throttle && !ignoreThrottle {
		return
	}

	// The percentage of the bar that is filled.
	// This, contradictory to the name, is a decimal, not a percentage.
	percentage := float64(bar.Value) / float64(bar.Total)

	//nolint:gomnd // You need to multiply by 100 to turn a decimal into a percent.
	percentDisplay := fmt.Sprintf("%.0f%%", percentage*100)

	// Using the value -6 because
	//   -1 for the spaces.
	//   -2 for the brackets.
	//   -2 Padding for cursor.
	//   -2 for the saucer.
	const reduce = 7

	// This is the width of the inside of the bar, not the entire bar including the brackets.
	// [======     ]
	//  ^^^^^^^^^^^
	width := getTermSize() - len(bar.Description) - len(percentDisplay) - reduce

	// Spamming float64 because golang division rounds to "help" you when using ints.
	// The amount of the bar to fill.
	// [======     ]
	//  ^^^^^^
	fill := int(percentage * float64(width))

	// The amount of empty space.
	// [======     ]
	//        ^^^^^
	empty := width - fill

	filledDisplay := strings.Repeat("=", fill)
	emptyDisplay := strings.Repeat(" ", empty)
	output := fmt.Sprintf("\r%s %s [%s>%s] ", bar.Description, percentDisplay, filledDisplay, emptyDisplay)

	fmt.Fprint(os.Stderr, output)

	bar.lastDisplayed = time.Now()
}

// Implementation for io.Writer.
func (bar *Bar) Write(data []byte) (int, error) {
	num := len(data)
	bar.Increment(int64(num))

	return num, nil
}
