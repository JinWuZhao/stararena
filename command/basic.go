package command

import (
	"fmt"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

type JoinGameCmd struct {
	player      string
	sc2PlayerId uint32
	isBot       bool
}

func (c *JoinGameCmd) Player() string {
	return c.player
}

func (c *JoinGameCmd) SC2PlayerId() uint32 {
	return c.sc2PlayerId
}

func (*JoinGameCmd) New() Command {
	return new(JoinGameCmd)
}

type JoinGameOpts func(*JoinGameCmd)

func JoinGameOptsPlayer(player string, sc2PlayerId uint32) JoinGameOpts {
	return func(cmd *JoinGameCmd) {
		cmd.player = player
		cmd.sc2PlayerId = sc2PlayerId
	}
}

func JoinGameOptsBot() JoinGameOpts {
	return func(cmd *JoinGameCmd) {
		cmd.isBot = true
	}
}

func (*JoinGameCmd) NewWithOpts(opts ...JoinGameOpts) Command {
	cmd := new(JoinGameCmd)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func (*JoinGameCmd) Name() string {
	return "JOIN_GAME"
}

func (c *JoinGameCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`j`, "OP"),
		ast.OrdChoice("SC2PLAYER", nil,
			parsec.Atom(`r`, "RED"),
			parsec.Atom(`b`, "BLUE")))
}

func (c *JoinGameCmd) Init(ctx Context, query parsec.Queryable) error {
	switch query.GetChildren()[1].GetName() {
	case "RED":
		c.sc2PlayerId = ctx.SC2RedPlayer
	case "BLUE":
		c.sc2PlayerId = ctx.SC2BluePlayer
	}
	c.player = ctx.Player
	return nil
}

func (c *JoinGameCmd) String() string {
	if c.isBot {
		return fmt.Sprintf("cmd-add-player %s %d bot", c.player, c.sc2PlayerId)
	}
	return fmt.Sprintf("cmd-add-player %s %d", c.player, c.sc2PlayerId)
}

type LeaveGameCmd struct {
	player string
}

func (c *LeaveGameCmd) New() Command {
	return new(LeaveGameCmd)
}

type LeaveGameOpts func(*LeaveGameCmd)

func LeaveGameOptsPlayer(player string) LeaveGameOpts {
	return func(cmd *LeaveGameCmd) {
		cmd.player = player
	}
}

func (c *LeaveGameCmd) NewWithOpts(opts ...LeaveGameOpts) Command {
	cmd := new(LeaveGameCmd)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func (c *LeaveGameCmd) Name() string {
	return "LEAVE_GAME"
}

func (c *LeaveGameCmd) Parser(_ *parsec.AST) parsec.Parser {
	panic("not support")
}

func (c *LeaveGameCmd) Init(_ Context, _ parsec.Queryable) error {
	panic("not support")
}

func (c *LeaveGameCmd) String() string {
	return fmt.Sprintf("cmd-remove-player %s", c.player)
}

type MoveCmd struct {
	player   string
	distance int32
	angle    int32
}

func (*MoveCmd) New() Command {
	return new(MoveCmd)
}

func (*MoveCmd) Name() string {
	return "MOVE"
}

func (c *MoveCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`g`, "OP"),
		parsecInt("DISTANCE", 4),
		parsecInt("ANGLE", 3))
}

func (c *MoveCmd) Init(ctx Context, query parsec.Queryable) error {
	distance, err := strconv.ParseInt(query.GetChildren()[1].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(distance) error: %w", err)
	}
	angle, err := strconv.ParseInt(query.GetChildren()[2].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(angle) error: %w", err)
	}
	if angle < 0 {
		angle = 0
	} else if angle > 360 {
		angle = 360
	}
	c.distance = int32(distance)
	c.angle = int32(angle)
	c.player = ctx.Player
	return nil
}

func (c *MoveCmd) String() string {
	return fmt.Sprintf("cmd-move-toward %s %d %d", c.player, c.angle, c.distance)
}

type MoveXCmd struct {
	player    string
	direction int32
	distance  int32
	angle     int32
}

func (*MoveXCmd) New() Command {
	return new(MoveXCmd)
}

func (*MoveXCmd) Name() string {
	return "MOVEX"
}

func (c *MoveXCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		ast.OrdChoice("OP", nil,
			parsec.AtomExact(`w`, "UP"),
			parsec.AtomExact(`s`, "DOWN"),
			parsec.AtomExact(`a`, "LEFT"),
			parsec.AtomExact(`d`, "RIGHT")),
		parsecInt("DISTANCE", 4),
		ast.Maybe("OPTION", nil,
			parsecInt("ANGLE", 3)))
}

func (c *MoveXCmd) Init(ctx Context, query parsec.Queryable) error {
	op := query.GetChildren()[0].GetName()
	distance, err := strconv.ParseInt(query.GetChildren()[1].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(distance) error: %w", err)
	}
	var angle int64
	if angleValue := query.GetChildren()[2].GetValue(); angleValue != "" {
		angle, err = strconv.ParseInt(angleValue, 10, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt(angle) error: %w", err)
		}
		if angle < -90 {
			angle = -90
		} else if angle > 90 {
			angle = 90
		}
	}
	switch op {
	case "RIGHT":
		c.direction = 0
	case "UP":
		c.direction = 90
	case "LEFT":
		c.direction = 180
	case "DOWN":
		c.direction = 270
	}
	c.distance = int32(distance)
	c.angle = int32(angle)
	c.player = ctx.Player
	return nil
}

func (c *MoveXCmd) String() string {
	return fmt.Sprintf("cmd-move-toward %s %d %d", c.player, c.direction+c.angle, c.distance)
}

type ChangeModeCmd struct {
	player string
	mode   string
}

func (*ChangeModeCmd) New() Command {
	return new(ChangeModeCmd)
}

func (*ChangeModeCmd) Name() string {
	return "CHANGE_MODE"
}

func (c *ChangeModeCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`m`, "OP"),
		ast.OrdChoice("MODE", nil,
			parsec.Atom(`0`, "manual"),
			parsec.Atom(`1`, "attack"),
			parsec.Atom(`2`, "hunter"),
			parsec.Atom(`3`, "defence"),
			parsec.Atom(`4`, "retreat")))
}

func (c *ChangeModeCmd) Init(ctx Context, query parsec.Queryable) error {
	c.mode = query.GetChildren()[1].GetValue()
	c.player = ctx.Player
	return nil
}

func (c *ChangeModeCmd) String() string {
	return fmt.Sprintf("cmd-set-aimode %s %s", c.player, c.mode)
}

type ShowPointsCmd struct {
	player string
}

func (c *ShowPointsCmd) New() Command {
	return new(ShowPointsCmd)
}

func (c *ShowPointsCmd) Name() string {
	return "SHOW_POINTS"
}

func (c *ShowPointsCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`pt`, "OP"))
}

func (c *ShowPointsCmd) Init(ctx Context, query parsec.Queryable) error {
	c.player = ctx.Player
	return nil
}

func (c *ShowPointsCmd) String() string {
	return fmt.Sprintf("cmd-show-points %s", c.player)
}
