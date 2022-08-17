package command

import (
	"fmt"
	"testing"
)

type FakeGameState struct{}

func (f *FakeGameState) FindAvailablePlayerId(humanOnly bool) uint32 {
	return 3
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		Context
		Text string
		Want string
	}{
		{
			Text: "j",
			Want: "cmd-add-player player 3",
		},
		{
			Text: "g 10 180",
			Want: "cmd-move-toward player 180 10",
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
			Text: "l 12",
			Want: "cmd-set-weapon player 12",
		},
		{
			Text: "k 5",
			Want: "cmd-set-ability player 5",
		},
		{
			Text: "xk 5",
			Want: "cmd-set-ability player 10005",
		},
		{
			Text: "k 5 k 6 k 7 k 8 k 9",
			Want: "cmd-set-ability player 5,6,7,8,9",
		},
		{
			Text: "k 5 xk 6 k 7 xk 8 k 9",
			Want: "cmd-set-ability player 5,10006,7,10008,9",
		},
		{
			Text: "a 4",
			Want: "cmd-assign-points player 0 4",
		},
		{
			Text: "b 4",
			Want: "cmd-assign-points player 1 4",
		},
		{
			Text: "c 4",
			Want: "cmd-assign-points player 2 4",
		},
		{
			Text: "d 4",
			Want: "cmd-assign-points player 3 4",
		},
		{
			Text: "e 4",
			Want: "cmd-assign-points player 4 4",
		},
		{
			Text: "f 4",
			Want: "cmd-assign-points player 5 4",
		},
		{
			Text: "a",
			Want: "cmd-assign-points player 0 1",
		},
		{
			Text: "b",
			Want: "cmd-assign-points player 1 1",
		},
		{
			Text: "c",
			Want: "cmd-assign-points player 2 1",
		},
		{
			Text: "d",
			Want: "cmd-assign-points player 3 1",
		},
		{
			Text: "e",
			Want: "cmd-assign-points player 4 1",
		},
		{
			Text: "i l 1",
			Want: "cmd-show-weapon 1",
		},
		{
			Text: "i k 1",
			Want: "cmd-show-ability 1",
		},
		{
			Text: "赞",
			Want: "cmd-apply-gift player 0 1",
		},
		{
			Text: "t 1",
			Want: "cmd-set-template player 1",
		},
		{
			Context: Context{
				Player: "星际竞技场",
			},
			Text: "n: 大家好",
			Want: "cmd-set-notice 大家好",
		},
	}
	for _, p := range tests {
		if p.Player == "" {
			p.Player = "player"
		}
		p.SC2RedPlayer = 3
		p.SC2BluePlayer = 4
		p.State = new(FakeGameState)
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
