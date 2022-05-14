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
			Want: "fake-cmd-join-game player 3",
		},
		{
			Text: "j b",
			Want: "fake-cmd-join-game player 4",
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
			Context: Context{
				Points: 0,
			},
			Text: "u 1",
			Want: "cmd-create-hellion 3 player",
		},
		{
			Context: Context{
				Points: 100,
			},
			Text: "u 2",
			Want: "cmd-create-siege-tank 3 player",
		},
		{
			Context: Context{
				Points: 300,
			},
			Text: "u 3",
			Want: "cmd-create-thor 3 player",
		},
		{
			Context: Context{
				Points: 500,
			},
			Text: "u 4",
			Want: "cmd-create-battlecruiser 3 player",
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
		p.SC2PlayerId = 3
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

func TestParseReport(t *testing.T) {
	tests := []struct {
		Text string
		Want string
	}{
		{
			Text: `report-attack-unit attacker target 100`,
			Want: `report-attack-unit attacker target 100`,
		},
		{
			Text: `report-kill-unit attacker target`,
			Want: `report-kill-unit attacker target`,
		},
		{
			Text: "report-victory 3",
			Want: "report-victory 3",
		},
	}
	for _, p := range tests {
		t.Run(p.Text, func(t *testing.T) {
			report, err := ParseReport(Context{}, p.Text)
			if err != nil {
				t.Error("ParseReport() error:", err)
				return
			}
			result := report.String()
			if p.Want != result {
				t.Error("incorrect report string:", result, "want:", p.Want)
				return
			}
			fmt.Println(report.String())
		})
	}
}
