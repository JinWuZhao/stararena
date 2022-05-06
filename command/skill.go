package command

import (
	"fmt"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

type SiegeMode struct {
	operator string
}

func (*SiegeMode) New(params ...any) Command {
	return &SiegeMode{
		operator: params[0].(string),
	}
}

func (*SiegeMode) Name() string {
	return "SIEGE_MODE"
}

func (c *SiegeMode) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil, parsec.AtomExact(c.operator, "OP"))
}

func (c *SiegeMode) Init(Context, parsec.Queryable) error {
	return nil
}

func (c *SiegeMode) String() string {
	return "siege-mode"
}

type TankMode struct {
	operator string
}

func (*TankMode) New(params ...any) Command {
	return &TankMode{
		operator: params[0].(string),
	}
}

func (*TankMode) Name() string {
	return "TANK_MODE"
}

func (c *TankMode) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(c.operator, "OP"))
}

func (c *TankMode) Init(Context, parsec.Queryable) error {
	return nil
}

func (c *TankMode) String() string {
	return "tank-mode"
}

type ExplosivePayload struct {
	operator string
}

func (*ExplosivePayload) New(params ...any) Command {
	return &ExplosivePayload{
		operator: params[0].(string),
	}
}

func (*ExplosivePayload) Name() string {
	return "EXPLOSIVE_PAYLOAD"
}

func (c *ExplosivePayload) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(c.operator, "OP"))
}

func (c *ExplosivePayload) Init(Context, parsec.Queryable) error {
	return nil
}

func (c *ExplosivePayload) String() string {
	return "explosive-payload"
}

type HighImpactPayload struct {
	operator string
}

func (*HighImpactPayload) New(params ...any) Command {
	return &HighImpactPayload{
		operator: params[0].(string),
	}
}

func (*HighImpactPayload) Name() string {
	return "HIGH_IMPACT_PAYLOAD"
}

func (c *HighImpactPayload) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(c.operator, "OP"))
}

func (c *HighImpactPayload) Init(Context, parsec.Queryable) error {
	return nil
}

func (c *HighImpactPayload) String() string {
	return "high-impact-payload"
}

type YamatoCannon struct {
	operator string
}

func (*YamatoCannon) New(params ...any) Command {
	return &YamatoCannon{
		operator: params[0].(string),
	}
}

func (*YamatoCannon) Name() string {
	return "YAMATO_CANNON"
}

func (c *YamatoCannon) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(c.operator, "OP"))
}

func (c *YamatoCannon) Init(Context, parsec.Queryable) error {
	return nil
}

func (c *YamatoCannon) String() string {
	return "yamato-cannon"
}

type TacticalJump struct {
	operator string
	distance int32
	angle    int32
}

func (*TacticalJump) New(params ...any) Command {
	return &TacticalJump{
		operator: params[0].(string),
	}
}

func (*TacticalJump) Name() string {
	return "TACTICAL_JUMP"
}

func (c *TacticalJump) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(c.operator, "OP"),
		parsecInt("DISTANCE", 4),
		parsecInt("ANGLE", 3))
}

func (c *TacticalJump) Init(_ Context, query parsec.Queryable) error {
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
	return nil
}

func (c *TacticalJump) String() string {
	return fmt.Sprintf("tactical-jump %d %d", c.distance, c.angle)
}
