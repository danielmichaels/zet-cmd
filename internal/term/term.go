// Credit to https://github.com/rwxrob/term which much of this is taken from

package term

import (
	"bufio"
	"fmt"
	"github.com/danielmichaels/zet-cmd/internal/esc"
	"log"
	"os"
	"os/exec"
)

func init() {
	SetInteractive(DetectInteractive())
}

var (
	Reset      string
	Bright     string
	Bold       string
	Dim        string
	Italic     string
	Under      string
	Blink      string
	BlinkF     string
	Reverse    string
	Hidden     string
	Strike     string
	BoldItalic string
	Black      string
	Red        string
	Green      string
	Yellow     string
	Blue       string
	Magenta    string
	Cyan       string
	White      string
	BBlack     string
	BRed       string
	BGreen     string
	BYellow    string
	BBlue      string
	BMagenta   string
	BCyan      string
	BWhite     string
	HBlack     string
	HRed       string
	HGreen     string
	HYellow    string
	HBlue      string
	HMagenta   string
	HCyan      string
	HWhite     string
	BHBlack    string
	BHRed      string
	BHGreen    string
	BHYellow   string
	BHBlue     string
	BHMagenta  string
	BHCyan     string
	BHWhite    string
	X          string
	B          string
	I          string
	U          string
	BI         string
)

// Prompt prints the given message if the terminal IsInteractive and
// reads the string by calling Read. The argument signature is identical
// as that passed to fmt.Printf().
func Prompt(form string, args ...any) string {
	if IsInteractive() {
		fmt.Printf(form, args...)
	}
	return Read()
}

var interactive bool

func IsInteractive() bool { return interactive }

// DetectInteractive returns true if the output is to an interactive
// terminal (not piped in any way).
func DetectInteractive() bool {
	if f, _ := os.Stdout.Stat(); (f.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

// Read reads a single line of input and chomps the \r?\n. Also see
// ReadHidden.
func Read() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// SetInteractive forces the interactive internal state affecting output
// including calling AttrOn (true) or AttrOff (false).
func SetInteractive(to bool) {
	interactive = to
	if to {
		AttrOn()
	} else {
		AttrOff()
	}
}

var attron bool

// AttrOff sets all the terminal attributes to zero values (empty strings).
// Note that this does not affect anything in the esc subpackage (which
// contains the constants from the VT100 specification). Sets the
// AttrAreOn bool to false.
func AttrOff() {
	attron = false
	Reset = ""
	Bright = ""
	Bold = ""
	Dim = ""
	Italic = ""
	Under = ""
	Blink = ""
	BlinkF = ""
	Reverse = ""
	Hidden = ""
	Strike = ""
	BoldItalic = ""
	Black = ""
	Red = ""
	Green = ""
	Yellow = ""
	Blue = ""
	Magenta = ""
	Cyan = ""
	White = ""
	BBlack = ""
	BRed = ""
	BGreen = ""
	BYellow = ""
	BBlue = ""
	BMagenta = ""
	BCyan = ""
	BWhite = ""
	HBlack = ""
	HRed = ""
	HGreen = ""
	HYellow = ""
	HBlue = ""
	HMagenta = ""
	HCyan = ""
	HWhite = ""
	BHBlack = ""
	BHRed = ""
	BHGreen = ""
	BHYellow = ""
	BHBlue = ""
	BHMagenta = ""
	BHCyan = ""
	BHWhite = ""
	X = ""
	B = ""
	I = ""
	U = ""
	BI = ""
}

// AttrOn sets all the terminal attributes to zero values (empty strings).
// Note that this does not affect anything in the esc subpackage (which
// contains the constants from the VT100 specification). Sets the
// AttrAreOn bool to true.
func AttrOn() {
	attron = true
	Reset = esc.Reset
	Bright = esc.Bright
	Bold = esc.Bold
	Dim = esc.Dim
	Italic = esc.Italic
	Under = esc.Under
	Blink = esc.Blink
	BlinkF = esc.BlinkF
	Reverse = esc.Reverse
	Hidden = esc.Hidden
	Strike = esc.Strike
	Black = esc.Black
	Red = esc.Red
	Green = esc.Green
	Yellow = esc.Yellow
	Blue = esc.Blue
	Magenta = esc.Magenta
	Cyan = esc.Cyan
	White = esc.White
	BBlack = esc.BBlack
	BRed = esc.BRed
	BGreen = esc.BGreen
	BYellow = esc.BYellow
	BBlue = esc.BBlue
	BMagenta = esc.BMagenta
	BCyan = esc.BCyan
	BWhite = esc.BWhite
	HBlack = esc.HBlack
	HRed = esc.HRed
	HGreen = esc.HGreen
	HYellow = esc.HYellow
	HBlue = esc.HBlue
	HMagenta = esc.HMagenta
	HCyan = esc.HCyan
	HWhite = esc.HWhite
	BHBlack = esc.BHBlack
	BHRed = esc.BHRed
	BHGreen = esc.BHGreen
	BHYellow = esc.BHYellow
	BHBlue = esc.BHBlue
	BHMagenta = esc.BHMagenta
	BHCyan = esc.BHCyan
	BHWhite = esc.BHWhite
	X = esc.Reset
	B = esc.Bold
	I = esc.Italic
	U = esc.Under
	BI = esc.BoldItalic
}

// Exec checks for existence of first argument as an executable on the
// system and then runs it without exiting in a way that is supported
// across all architectures that Go supports. The stdin, stdout, and stderr are
// connected directly to that of the calling program. Sometimes this is
// insufficient and the UNIX-specific SysExec is preferred. For example,
// when handing over control to a terminal editor such as Vim.
func Exec(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing name of executable")
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Out returns the standard output of the executed command as
// a string. Errors are logged but not returned.
func Out(args ...string) string {
	if len(args) == 0 {
		log.Println("missing name of executable")
		return ""
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Println(err)
		return ""
	}
	out, err := exec.Command(path, args[1:]...).Output()
	if err != nil {
		log.Println(err)
	}
	return string(out)
}
