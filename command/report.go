package command

import (
	"fmt"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

type DamageUnitReport struct {
	Player string
	Target string
	Damage int32
}

func (*DamageUnitReport) New(...any) Command {
	return new(DamageUnitReport)
}

func (*DamageUnitReport) Name() string {
	return "ATTACK_UNIT"
}

func (c *DamageUnitReport) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`report-damage-unit`, "OP"),
		parsecWS(),
		parsec.TokenExact(`[^\s]+`, "PLAYER"),
		parsecWS(),
		parsec.TokenExact(`[^\s]+`, "TARGET"),
		parsecWS(),
		parsec.TokenExact(`[0-9]+`, "DAMAGE"))
}

func (c *DamageUnitReport) Init(_ Context, query parsec.Queryable) error {
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

func (c *DamageUnitReport) String() string {
	return fmt.Sprintf(`report-damage-unit %s %s %d`, c.Player, c.Target, c.Damage)
}

type KillUnitReport struct {
	Player string
	Target string
	Unit   string
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
		parsecWS(),
		parsec.TokenExact(`[^\s]+`, "PLAYER"),
		parsecWS(),
		parsec.TokenExact(`[^\s]+`, "TARGET"),
		parsecWS(),
		parsec.TokenExact(`[a-z0-9-]+`, "UNIT"))
}

func (c *KillUnitReport) Init(_ Context, query parsec.Queryable) error {
	c.Player = query.GetChildren()[2].GetValue()
	c.Target = query.GetChildren()[4].GetValue()
	c.Unit = query.GetChildren()[6].GetValue()
	return nil
}

func (c *KillUnitReport) String() string {
	return fmt.Sprintf(`report-kill-unit %s %s %s`, c.Player, c.Target, c.Unit)
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
