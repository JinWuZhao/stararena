package command

import (
	"fmt"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		Context
		Text string
		Want string
	}{
		{
			Text: "j r",
			Want: "cmd-add-player player 3",
		},
		{
			Text: "j b",
			Want: "cmd-add-player player 4",
		},
		{
			Text: "g 10 180",
			Want: "cmd-move-toward player 180 10",
		},
		{
			Text: "w 10 45",
			Want: "cmd-move-toward player 135 10",
		},
		{
			Text: "s 10 45",
			Want: "cmd-move-toward player 315 10",
		},
		{
			Text: "a 10 45",
			Want: "cmd-move-toward player 225 10",
		},
		{
			Text: "d 10 45",
			Want: "cmd-move-toward player 45 10",
		},
		{
			Text: "m a",
			Want: "cmd-set-behavior-mode player attack",
		},
		{
			Text: "m d",
			Want: "cmd-set-behavior-mode player defence",
		},
		{
			Text: "m r",
			Want: "cmd-set-behavior-mode player retreat",
		},
		{
			Text: "u 1",
			Want: "cmd-create-hellion player",
		},
		{
			Text: "u 2",
			Want: "cmd-create-siege-tank player",
		},
		{
			Text: "u 3",
			Want: "cmd-create-thor player",
		},
		{
			Text: "u 4",
			Want: "cmd-create-battlecruiser player",
		},
		{
			Context: Context{
				Unit: "siege-tank",
			},
			Text: "k 1",
			Want: "cmd-issue-ability-siege-tank player siege-mode",
		},
		{
			Context: Context{
				Unit: "siege-tank",
			},
			Text: "k 2",
			Want: "cmd-issue-ability-siege-tank player tank-mode",
		},
		{
			Context: Context{
				Unit: "battlecruiser",
			},
			Text: "k 1",
			Want: "cmd-issue-ability-battlecruiser player yamato-cannon",
		},
		{
			Context: Context{
				Unit: "battlecruiser",
			},
			Text: "k 2 20 180",
			Want: "cmd-issue-ability-battlecruiser player tactical-jump 20 180",
		},
		{
			Context: Context{
				Unit: "thor",
			},
			Text: "k 1",
			Want: "cmd-issue-ability-thor player explosive-payload",
		},
		{
			Context: Context{
				Unit: "thor",
			},
			Text: "k 2",
			Want: "cmd-issue-ability-thor player high-impact-payload",
		},
	}
	for _, p := range tests {
		p.Player = "player"
		p.SC2RedPlayer = 3
		p.SC2BluePlayer = 4
		t.Run(p.Text, func(t *testing.T) {
			command, err := ParseCommand(p.Context, p.Text)
			if err != nil {
				t.Error("MakeCommand() error:", err)
				return
			}
			result := command.String()
			if p.Want != result {
				t.Error("incorrect command string:", result, "want:", p.Want)
				return
			}
			fmt.Println(command.String())
		})
	}
}
