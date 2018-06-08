package color

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/mattn/go-colorable"
)

// Testing colors is kinda different. First we test for given colors and their
// escaped formatted results. Next we create some visual tests to be tested.
// Each visual test includes the color name to be compared.
func TestColor(t *testing.T) {
	rb := new(bytes.Buffer)
	DefaultOutput = rb

	NoColor = false

	testColors := []struct {
		text string
		code Attribute
	}{
		{text: "black", code: FgBlack},
		{text: "red", code: FgRed},
		{text: "green", code: FgGreen},
		{text: "yellow", code: FgYellow},
		{text: "blue", code: FgBlue},
		{text: "magent", code: FgMagenta},
		{text: "cyan", code: FgCyan},
		{text: "white", code: FgWhite},
		{text: "hblack", code: FgHiBlack},
		{text: "hred", code: FgHiRed},
		{text: "hgreen", code: FgHiGreen},
		{text: "hyellow", code: FgHiYellow},
		{text: "hblue", code: FgHiBlue},
		{text: "hmagent", code: FgHiMagenta},
		{text: "hcyan", code: FgHiCyan},
		{text: "hwhite", code: FgHiWhite},
	}

	for _, c := range testColors {
		Default(c.code).Print(c.text)

		line, _ := rb.ReadString('\n')
		scannedLine := fmt.Sprintf("%q", line)
		colored := fmt.Sprintf("\x1b[%dm%s\x1b[0m", c.code, c.text)
		escapedForm := fmt.Sprintf("%q", colored)

		fmt.Printf("%s\t: %s\n", c.text, line)

		if scannedLine != escapedForm {
			t.Errorf("Expecting %s, got '%s'\n", escapedForm, scannedLine)
		}
	}
}

func TestColorEquals(t *testing.T) {
	fgblack1 := Default(FgBlack)
	fgblack2 := Default(FgBlack)
	bgblack := Default(BgBlack)
	fgbgblack := Default(FgBlack, BgBlack)
	fgblackbgred := Default(FgBlack, BgRed)
	fgred := Default(FgRed)
	bgred := Default(BgRed)

	if !fgblack1.Equals(fgblack2) {
		t.Error("Two black colors are not equal")
	}

	if fgblack1.Equals(bgblack) {
		t.Error("Fg and bg black colors are equal")
	}

	if fgblack1.Equals(fgbgblack) {
		t.Error("Fg black equals fg/bg black color")
	}

	if fgblack1.Equals(fgred) {
		t.Error("Fg black equals Fg red")
	}

	if fgblack1.Equals(bgred) {
		t.Error("Fg black equals Bg red")
	}

	if fgblack1.Equals(fgblackbgred) {
		t.Error("Fg black equals fg black bg red")
	}
}

func TestNoColor(t *testing.T) {
	rb := new(bytes.Buffer)
	DefaultOutput = rb

	testColors := []struct {
		text string
		code Attribute
	}{
		{text: "black", code: FgBlack},
		{text: "red", code: FgRed},
		{text: "green", code: FgGreen},
		{text: "yellow", code: FgYellow},
		{text: "blue", code: FgBlue},
		{text: "magent", code: FgMagenta},
		{text: "cyan", code: FgCyan},
		{text: "white", code: FgWhite},
		{text: "hblack", code: FgHiBlack},
		{text: "hred", code: FgHiRed},
		{text: "hgreen", code: FgHiGreen},
		{text: "hyellow", code: FgHiYellow},
		{text: "hblue", code: FgHiBlue},
		{text: "hmagent", code: FgHiMagenta},
		{text: "hcyan", code: FgHiCyan},
		{text: "hwhite", code: FgHiWhite},
	}

	for _, c := range testColors {
		p := Default(c.code)
		p.DisableColor()
		p.Print(c.text)

		line, _ := rb.ReadString('\n')
		if line != c.text {
			t.Errorf("Expecting %s, got '%s'\n", c.text, line)
		}
	}

	// global check
	NoColor = true
	defer func() {
		NoColor = false
	}()
	for _, c := range testColors {
		p := Default(c.code)
		p.Print(c.text)

		line, _ := rb.ReadString('\n')
		if line != c.text {
			t.Errorf("Expecting %s, got '%s'\n", c.text, line)
		}
	}

}

func TestColorVisual(t *testing.T) {
	// First Visual Test
	DefaultOutput = colorable.NewColorableStdout()

	Default(FgRed).Printf("red\t")
	Default(BgRed).Print("         ")
	Default(FgRed, Bold).Println(" red")

	Default(FgGreen).Printf("green\t")
	Default(BgGreen).Print("         ")
	Default(FgGreen, Bold).Println(" green")

	Default(FgYellow).Printf("yellow\t")
	Default(BgYellow).Print("         ")
	Default(FgYellow, Bold).Println(" yellow")

	Default(FgBlue).Printf("blue\t")
	Default(BgBlue).Print("         ")
	Default(FgBlue, Bold).Println(" blue")

	Default(FgMagenta).Printf("magenta\t")
	Default(BgMagenta).Print("         ")
	Default(FgMagenta, Bold).Println(" magenta")

	Default(FgCyan).Printf("cyan\t")
	Default(BgCyan).Print("         ")
	Default(FgCyan, Bold).Println(" cyan")

	Default(FgWhite).Printf("white\t")
	Default(BgWhite).Print("         ")
	Default(FgWhite, Bold).Println(" white")
	fmt.Println("")

	// Second Visual test
	Black("black")
	Red("red")
	Green("green")
	Yellow("yellow")
	Blue("blue")
	Magenta("magenta")
	Cyan("cyan")
	White("white")

	// Third visual test
	fmt.Println()
	Set(FgBlue)
	fmt.Println("is this blue?")
	Unset()

	Set(FgMagenta)
	fmt.Println("and this magenta?")
	Unset()

	// Fourth Visual test
	fmt.Println()
	blue := Default(FgBlue).PrintlnFunc()
	blue("blue text with custom print func")

	red := Default(FgRed).PrintfFunc()
	red("red text with a printf func: %d\n", 123)

	put := Default(FgYellow).SprintFunc()
	warn := Default(FgRed).SprintFunc()

	fmt.Fprintf(DefaultOutput, "this is a %s and this is %s.\n", put("warning"), warn("error"))

	info := Default(FgWhite, BgGreen).SprintFunc()
	fmt.Fprintf(DefaultOutput, "this %s rocks!\n", info("package"))

	// Fifth Visual Test
	fmt.Println()

	fmt.Fprintln(DefaultOutput, BlackString("black"))
	fmt.Fprintln(DefaultOutput, RedString("red"))
	fmt.Fprintln(DefaultOutput, GreenString("green"))
	fmt.Fprintln(DefaultOutput, YellowString("yellow"))
	fmt.Fprintln(DefaultOutput, BlueString("blue"))
	fmt.Fprintln(DefaultOutput, MagentaString("magenta"))
	fmt.Fprintln(DefaultOutput, CyanString("cyan"))
	fmt.Fprintln(DefaultOutput, WhiteString("white"))
}
