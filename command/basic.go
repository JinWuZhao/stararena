package command

import (
	"fmt"
	"strconv"
	"strings"

	parsec "github.com/prataprc/goparsec"
)

type JoinGameCmd struct {
	player      string
	playerUID   uint32
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
		parsec.AtomExact(`j`, "OP"))
}

func (c *JoinGameCmd) Init(ctx Context, query parsec.Queryable) error {
	c.player = ctx.Player
	c.playerUID = ctx.PlayerUID
	c.sc2PlayerId = ctx.State.FindAvailablePlayerId(true)
	return nil
}

func (c *JoinGameCmd) String() string {
	if c.isBot {
		return fmt.Sprintf("cmd-add-player %s %d bot", c.player, c.sc2PlayerId)
	}
	return fmt.Sprintf("cmd-add-player %s %d %d", c.player, c.sc2PlayerId, c.playerUID)
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
		ast.OrdChoice("MODE", nil, parsec.Token(`[0-6]`, "INDEX")))
}

func (c *ChangeModeCmd) Init(ctx Context, query parsec.Queryable) error {
	c.mode = query.GetChildren()[1].GetValue()
	c.player = ctx.Player
	return nil
}

func (c *ChangeModeCmd) String() string {
	return fmt.Sprintf("cmd-set-aimode %s %s", c.player, c.mode)
}

type SetWeaponCmd struct {
	player string
	weapon int32
}

func (c *SetWeaponCmd) New() Command {
	return new(SetWeaponCmd)
}

func (c *SetWeaponCmd) Name() string {
	return "SET_WEAPON"
}

func (c *SetWeaponCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`l`, "OP"), parsecUint("WEAPON", 2))
}

func (c *SetWeaponCmd) Init(ctx Context, query parsec.Queryable) error {
	weapon, err := strconv.ParseInt(query.GetChildren()[1].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(weapon): %w", err)
	}
	c.weapon = int32(weapon)
	c.player = ctx.Player
	return nil
}

func (c *SetWeaponCmd) String() string {
	return fmt.Sprintf("cmd-set-weapon %s %d", c.player, c.weapon)
}

type SetAbilityCmd struct {
	player    string
	abilities []string
}

func (c *SetAbilityCmd) New() Command {
	return new(SetAbilityCmd)
}

func (c *SetAbilityCmd) Name() string {
	return "SET_ABILITY"
}

func (c *SetAbilityCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.ManyUntil(c.Name(), nil,
		ast.And("ABILITY", nil,
			parsec.TokenExact(`x?k`, "OP"),
			parsecUint("INDEX", 2)),
		parsec.Atom("", "SEP"),
		ast.End("EOF"))
}

const exAbilityFlag int64 = 10000

func (c *SetAbilityCmd) Init(ctx Context, query parsec.Queryable) error {
	var abilities []string
	for index, node := range query.GetChildren() {
		if index >= 6 {
			break
		}
		abilFlag := node.GetChildren()[0].GetValue()
		ability, err := strconv.ParseInt(node.GetChildren()[1].GetValue(), 10, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt(ability): %w", err)
		}
		if strings.HasPrefix(abilFlag, "x") {
			ability += exAbilityFlag
		}
		abilities = append(abilities, strconv.FormatInt(ability, 10))
	}
	c.abilities = abilities
	c.player = ctx.Player
	return nil
}

func (c *SetAbilityCmd) String() string {
	return fmt.Sprintf("cmd-set-ability %s %s", c.player, strings.Join(c.abilities, ","))
}

type AssignPointsCmd struct {
	player string
	prop   int32
	points int32
}

func (c *AssignPointsCmd) New() Command {
	return new(AssignPointsCmd)
}

func (c *AssignPointsCmd) Name() string {
	return "ASSIGN_POINTS"
}

func (c *AssignPointsCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.TokenExact(`[A-Fa-f]`, "PROP"),
		ast.Maybe("", nil, parsecUint("POINTS", 2)))
}

func (c *AssignPointsCmd) Init(ctx Context, query parsec.Queryable) error {
	var prop int
	switch query.GetChildren()[0].GetValue() {
	case "A", "a":
		prop = 0
	case "B", "b":
		prop = 1
	case "C", "c":
		prop = 2
	case "D", "d":
		prop = 3
	case "E", "e":
		prop = 4
	case "F", "f":
		prop = 5
	default:
		return fmt.Errorf("unkown prop %s", query.GetChildren()[0].GetValue())
	}
	var points int64
	var err error
	if query.GetChildren()[1].GetValue() != "" {
		points, err = strconv.ParseInt(query.GetChildren()[1].GetValue(), 10, 32)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt(points): %w", err)
		}
	} else {
		points = 1
	}
	c.prop = int32(prop)
	c.points = int32(points)
	c.player = ctx.Player
	return nil
}

func (c *AssignPointsCmd) String() string {
	return fmt.Sprintf("cmd-assign-points %s %d %d", c.player, c.prop, c.points)
}

type ShowInfoCmd struct {
	kind  string
	index int32
}

func (c *ShowInfoCmd) New() Command {
	return new(ShowInfoCmd)
}

func (c *ShowInfoCmd) Name() string {
	return "SHOW_INFO"
}

func (c *ShowInfoCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`i`, "OP"),
		ast.OrdChoice("KIND", nil,
			parsec.Atom(`l`, "WEAPON"),
			parsec.Atom(`k`, "ABILITY")),
		parsecUint("INDEX", 2))
}

func (c *ShowInfoCmd) Init(ctx Context, query parsec.Queryable) error {
	index, err := strconv.ParseInt(query.GetChildren()[2].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(index): %w", err)
	}
	c.kind = query.GetChildren()[1].GetValue()
	c.index = int32(index)
	return nil
}

func (c *ShowInfoCmd) String() string {
	if c.kind == "l" {
		return fmt.Sprintf("cmd-show-weapon %d", c.index)
	}
	return fmt.Sprintf("cmd-show-ability %d", c.index)
}

type GiftItemCmd struct {
	player string
	kind   uint32
	number uint32
}

func (c *GiftItemCmd) New() Command {
	return new(GiftItemCmd)
}

type GiftItemOpts func(*GiftItemCmd)

func GiftItemOptsGift(player string, gift string, number uint32) GiftItemOpts {
	return func(cmd *GiftItemCmd) {
		cmd.player = player
		if gift == "辣条" {
			cmd.kind = 1
		} else if gift == "小花花" {
			cmd.kind = 2
		} else if gift == "粉丝团灯牌" {
			cmd.kind = 3
		}
		cmd.number = number
	}
}

func (*GiftItemCmd) NewWithOpts(opts ...GiftItemOpts) Command {
	cmd := new(GiftItemCmd)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func (c *GiftItemCmd) Name() string {
	return "GIFT_ITEM"
}

func (c *GiftItemCmd) Parser(_ *parsec.AST) parsec.Parser {
	panic("not support")
}

func (c *GiftItemCmd) Init(_ Context, _ parsec.Queryable) error {
	panic("not support")
}

func (c *GiftItemCmd) String() string {
	return fmt.Sprintf("cmd-apply-gift %s %d %d", c.player, c.kind, c.number)
}

type UpvoteCmd struct {
	player string
}

func (c *UpvoteCmd) New() Command {
	return new(UpvoteCmd)
}

func (c *UpvoteCmd) Name() string {
	return "UPVOTE"
}

func (c *UpvoteCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`赞`, "UPVOTE"))
}

func (c *UpvoteCmd) Init(ctx Context, query parsec.Queryable) error {
	c.player = ctx.Player
	return nil
}

func (c *UpvoteCmd) String() string {
	return fmt.Sprintf("cmd-apply-gift %s %d %d", c.player, 0, 1)
}

type SetTemplateCmd struct {
	player   string
	template int32
}

func (c *SetTemplateCmd) New() Command {
	return new(SetTemplateCmd)
}

func (c *SetTemplateCmd) Name() string {
	return "SET_TEMPLATE"
}

func (c *SetTemplateCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`t`, "OP"),
		parsecUint("TEMPLATE", 2))
}

func (c *SetTemplateCmd) Init(ctx Context, query parsec.Queryable) error {
	tpl, err := strconv.ParseInt(query.GetChildren()[1].GetValue(), 10, 32)
	if err != nil {
		return fmt.Errorf("strconv.ParseInt(tpl): %w", err)
	}
	c.template = int32(tpl)
	c.player = ctx.Player
	return nil
}

func (c *SetTemplateCmd) String() string {
	return fmt.Sprintf("cmd-set-template %s %d", c.player, c.template)
}

type SetNoticeCmd struct {
	notice string
}

func (c *SetNoticeCmd) New() Command {
	return new(SetNoticeCmd)
}

func (c *SetNoticeCmd) Name() string {
	return "SET_NOTICE"
}

func (c *SetNoticeCmd) Parser(ast *parsec.AST) parsec.Parser {
	return ast.And(c.Name(), nil,
		parsec.AtomExact(`n:`, "OP"),
		ast.Maybe("VALUE", nil, parsec.Token(`.+`, "NOTICE")))
}

func (c *SetNoticeCmd) Init(ctx Context, query parsec.Queryable) error {
	if ctx.Player != ctx.Streamer {
		return fmt.Errorf("%s has no permission", ctx.Player)
	}
	c.notice = query.GetChildren()[1].GetValue()
	return nil
}

func (c *SetNoticeCmd) String() string {
	return fmt.Sprintf("cmd-set-notice %s", c.notice)
}
