package environments

import (
	"bytes"
	"github.com/Mrbruin-ctrl/moon-cli/lib/go-ansi"
	"github.com/Mrbruin-ctrl/moon-cli/src/internal/debug"
	"github.com/Mrbruin-ctrl/moon-cli/src/internal/util"
	"github.com/c-bata/go-prompt"
	"github.com/olekukonko/tablewriter"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Command struct {
	Text            string
	HokKey          prompt.Key
	Environments    []Environment
	Provider        string
	Description     string
	Alias           string
	MatchFn         func(string, string) bool
	Fn              func([]string, io.Writer)
	ExtendArguments []prompt.Suggest
	ExtendOptions   []prompt.Suggest
}

func checkEnvironments(envs []Environment) bool {
	if len(envs) == 0 {
		return true
	} else {
		for _, env := range envs {
			if env == ActiveEnvironment {
				return true
			}
		}
	}
	return false
}

func (cmd *Command) MatchAndExecute(args string, out io.Writer) bool {
	if checkEnvironments(cmd.Environments) && cmd.MatchFn != nil && cmd.MatchFn(cmd.Text, args) {
		cmd.Fn(strings.Split(args, " "), out)
		return true
	}
	return false
}

func (cmd *Command) keyBindFn(buffer *prompt.Buffer) {
	out := &bytes.Buffer{}
	cmd.Fn([]string{}, out)
	buffer.InsertText("\n"+out.String(), false, false)
}

func BuildPromptKeyBinds() []prompt.KeyBind {
	var keybinds []prompt.KeyBind
	for _, cmd := range DefaultCommands {
		x := cmd // for closure
		if cmd.HokKey > 0 {
			keybinds = append(keybinds, prompt.KeyBind{
				Key: x.HokKey,
				Fn:  x.keyBindFn,
			})
		}
	}
	return keybinds
}

func IgnoreCaseMatch(text string, cmd string) bool {
	return strings.EqualFold(cmd, text)
}

func StartWithMatch(text string, cmd string) bool {
	return strings.HasPrefix(cmd, text)
}

func getKeyName(k prompt.Key) string {
	if k >= prompt.F1 && k <= prompt.F12 {
		return "F" + strconv.Itoa(int(k-prompt.F1+1))
	}
	return ""
}

func getEnv(envs []Environment) string {
	if len(envs) != 0 {
		return strings.Join(util.Slice(envs, reflect.TypeOf([]string(nil))).([]string), ",")
	} else {
		return "All"
	}
}

// 打印帮助信息
func HelpInfo(args []string, out io.Writer) {
	Logo()
	table := tablewriter.NewWriter(ansi.NewAnsiStdout())
	table.SetHeader([]string{"Command", "Hotkey", "Description", "Environment"})
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})
	table.SetColWidth(120)
	table.SetCenterSeparator("|")
	cmds := DefaultCommands
	for _, v := range Environments {
		cmds = append(cmds, v.Commands...)
	}
	for _, v := range cmds {
		table.Append([]string{v.Text, getKeyName(v.HokKey), v.Description, getEnv(v.Environments)})
	}

	// keyboard
	for _, v := range [][]string{
		{"", "", "", ""},
		{"", "Ctrl+A", "Go to the beginning of the line (Home)", ""},
		{"", "Ctrl+E", "Go to the end of the line (End)", ""},
		{"", "Ctrl+P", "Previous command (Up arrow)", ""},
		{"", "Ctrl+N", "Next command (Down arrow)", ""},
		{"", "Ctrl+F", "Forward one character", ""},
		{"", "Ctrl+B", "Backward one character", ""},
		{"", "Ctrl+D", "Delete character under the cursor", ""},
		{"", "Ctrl+H", "Delete character before the cursor (Backspace)", ""},
		{"", "Ctrl+W", "Cut the word before the cursor to the clipboard", ""},
		{"", "Ctrl+K", "Cut the line after the cursor to the clipboard", ""},
		{"", "Ctrl+U", "Cut the line before the cursor to the clipboard", ""},
		{"", "Ctrl+L", "Clear the screen", ""},
	} {
		table.Append(v)
	}
	table.Render()
}

var DefaultCommands []Command

func init() {
	DefaultCommands = []Command{
		{
			Text:        "x-help",
			HokKey:      prompt.F1,
			Description: "帮助信息",
			MatchFn:     IgnoreCaseMatch,
			Fn:          HelpInfo,
		},
		{
			Text:        "exit",
			HokKey:      prompt.F12,
			Description: "退出",
			MatchFn:     IgnoreCaseMatch,
			Fn: func(args []string, out io.Writer) {
				_, _ = out.Write([]byte("Bye!\n"))
				os.Exit(0)
			},
		},
		{
			Text:        "clear",
			Description: "清屏",
			MatchFn:     IgnoreCaseMatch,
			Fn: func(args []string, out io.Writer) {
				consoleWriter := prompt.GetConsoleWriter()
				consoleWriter.EraseScreen()
				consoleWriter.CursorGoTo(0, 0)
				debug.AssertNoError(consoleWriter.Flush())
			},
		},
	}
}
