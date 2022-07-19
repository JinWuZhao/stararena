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
			Text: "m 0",
			Want: "cmd-set-aimode player 0",
		},
		{
			Text: "m 1",
			Want: "cmd-set-aimode player 1",
		},
		{
			Text: "m 2",
			Want: "cmd-set-aimode player 2",
		},
		{
			Text: "m 3",
			Want: "cmd-set-aimode player 3",
		},
		{
			Text: "m 4",
			Want: "cmd-set-aimode player 4",
		},
		{
			Text: "i 12",
			Want: "cmd-set-weapon player 12",
		},
		{
			Text: "p 1 4",
			Want: "cmd-assign-points player 1 4",
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
