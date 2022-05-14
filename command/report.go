package command

import (
	"fmt"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

type AttackUnitReport struct {
	Player string
	Target string
	Damage int32
}

func (*AttackUnitReport) New(...any) Command {
	return new(AttackUnitReport)
}

func (*AttackUnitReport) Name() string {
	return "ATTACK_UNIT"
}

func (c *AttackUnitReport) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`report-attack-unit`, "OP"),
		parsec.TokenExact(`\s+`, "WS"),
		parsec.TokenExact(`[^\s]+`, "PLAYER"),
		parsec.TokenExact(`\s+`, "WS"),
		parsec.TokenExact(`[^\s]+`, "PLAYER"),
		parsec.TokenExact(`\s+`, "WS"),
		parsec.TokenExact(`[0-9]{1,10}`, "DAMAGE"))
}

func (c *AttackUnitReport) Init(_ Context, query parsec.Queryable) error {
	player := query.GetChildren()[2].GetValue()
	target := query.GetChildren()[4].GetValue()
	damage, err := strconv.ParseInt(query.GetChildren()[6].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(damage) error: %w", err)
	}
	c.Player = player
	c.Target = target
	c.Damage = int32(damage)
	return nil
}

func (c *AttackUnitReport) String() string {
	return fmt.Sprintf(`report-attack-unit %s %s %d`, c.Player, c.Target, c.Damage)
}

type KillUnitReport struct {
	Player string
	Target string
}

func (*KillUnitReport) New(...any) Command {
	return new(KillUnitReport)
}

func (*KillUnitReport) Name() string {
	return "KILL_UNIT"
}

func (c *KillUnitReport) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`report-kill-unit`, "OP"),
		parsec.TokenExact(`\s+`, "WS"),
		parsec.TokenExact(`[^\s]+`, "PLAYER"),
		parsec.TokenExact(`\s+`, "WS"),
		parsec.TokenExact(`[^\s]+`, "PLAYER"))
}

func (c *KillUnitReport) Init(_ Context, query parsec.Queryable) error {
	c.Player = query.GetChildren()[2].GetValue()
	c.Target = query.GetChildren()[4].GetValue()
	return nil
}

func (c *KillUnitReport) String() string {
	return fmt.Sprintf(`report-kill-unit %s %s`, c.Player, c.Target)
}

type VictoryReport struct {
	SC2PlayerId uint32
}

func (c *VictoryReport) New(...any) Command {
	return new(VictoryReport)
}

func (c *VictoryReport) Name() string {
	return "VICTORY"
}

func (c *VictoryReport) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`report-victory`, "OP"),
		parsec.TokenExact(`\s+`, "WS"),
		parsecInt("SC2PLAYER", 2))
}

func (c *VictoryReport) Init(_ Context, query parsec.Queryable) error {
	sc2PlayerId, err := strconv.ParseUint(query.GetChildren()[2].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseUint(sc2PlayerId) error: %w", err)
	}
	c.SC2PlayerId = uint32(sc2PlayerId)
	return nil
}

func (c *VictoryReport) String() string {
	return fmt.Sprintf("report-victory %d", c.SC2PlayerId)
}
