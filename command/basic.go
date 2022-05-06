package command

import (
	"fmt"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

type JoinGameCmd struct {
	player      string
	sc2PlayerId uint32
}

func (c *JoinGameCmd) Player() string {
	return c.player
}

func (c *JoinGameCmd) SC2PlayerId() uint32 {
	return c.sc2PlayerId
}

func (*JoinGameCmd) New(...any) Command {
	return new(JoinGameCmd)
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
	return fmt.Sprintf("fake-cmd-join-game %s %d", c.player, c.sc2PlayerId)
}

type MoveCmd struct {
	player   string
	distance int32
	angle    int32
}

func (*MoveCmd) New(...any) Command {
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

func (*MoveXCmd) New(...any) Command {
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

func (*ChangeModeCmd) New(...any) Command {
	return new(ChangeModeCmd)
}

func (*ChangeModeCmd) Name() string {
	return "CHANGE_MODE"
}

func (c *ChangeModeCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`m`, "OP"),
		ast.OrdChoice("MODE", nil,
			parsec.Atom(`a`, "attack"),
			parsec.Atom(`d`, "defence"),
			parsec.Atom(`r`, "retreat")))
}

func (c *ChangeModeCmd) Init(ctx Context, query parsec.Queryable) error {
	c.mode = query.GetChildren()[1].GetName()
	c.player = ctx.Player
	return nil
}

func (c *ChangeModeCmd) String() string {
	return fmt.Sprintf("cmd-set-behavior-mode %s %s", c.player, c.mode)
}

type ChangeUnitCmd struct {
	sc2PlayerId uint32
	player      string
	unit        string
}

func (*ChangeUnitCmd) New(...any) Command {
	return new(ChangeUnitCmd)
}

func (*ChangeUnitCmd) Name() string {
	return "CHANGE_UNIT"
}

func (c *ChangeUnitCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`u`, "OP"),
		parsec.Token(`[a-zA-Z0-9-]{1,64}`, "UNIT"))
}

func (c *ChangeUnitCmd) Init(ctx Context, query parsec.Queryable) error {
	unitName := query.GetChildren()[1].GetValue()
	unit, ok := unitDataAcc[unitName]
	if !ok {
		return fmt.Errorf("invalid unit name: %s", unitName)
	}
	if ctx.SC2PlayerId != ctx.SC2RedPlayer && ctx.SC2PlayerId != ctx.SC2BluePlayer {
		return fmt.Errorf("invalid sc2 player: %d", ctx.SC2PlayerId)
	}
	if ctx.Points-unit.Points < 0 {
		return fmt.Errorf("not enough score: %d, need: %d", ctx.Points, unit.Points)
	}
	c.unit = unit.Name
	c.player = ctx.Player
	c.sc2PlayerId = ctx.SC2PlayerId
	return nil
}

func (c *ChangeUnitCmd) String() string {
	return fmt.Sprintf("cmd-create-%s %d %s", c.unit, c.sc2PlayerId, c.player)
}

type IssueSkillCmd struct {
	player string
	unit   string
	skill  Command
}

func (*IssueSkillCmd) New(...any) Command {
	return new(IssueSkillCmd)
}

func (*IssueSkillCmd) Name() string {
	return "ISSUE_SKILL"
}

func (c *IssueSkillCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`k`, "OP"),
		parsec.Token(`[a-zA-Z0-9-\s]{1,128}`, "SKILL"))
}

func (c *IssueSkillCmd) Init(ctx Context, query parsec.Queryable) error {
	if _, ok := unitSkillCtors[ctx.Unit]; !ok {
		return fmt.Errorf("invalid unit name: %s", ctx.Unit)
	}
	skill, err := parseUnitSkill(ctx, query.GetChildren()[1].GetValue())
	if err != nil {
		return fmt.Errorf("parseUnitSkill() error: %w", err)
	}
	c.player = ctx.Player
	c.unit = ctx.Unit
	c.skill = skill
	return nil
}

func (c *IssueSkillCmd) String() string {
	if c.skill == nil {
		return ""
	}
	return fmt.Sprintf("cmd-issue-ability-%s %s %s", c.unit, c.player, c.skill.String())
}
