package command

import (
	"fmt"
	"log"
	"strconv"

	parsec "github.com/prataprc/goparsec"
)

var (
	cmdAST, cmdParser       = makeParser()
	skillASTs, skillParsers = makeSkillParsers()
)

type Context struct {
	SC2RedPlayer  uint32
	SC2BluePlayer uint32
	Player        string
	Unit          string
}

type Command interface {
	New(params ...any) Command
	Name() string
	Parser(ast *parsec.AST) parsec.Parser
	Init(ctx Context, query parsec.Queryable) error
	String() string
}

type Constructor func() Command

func MakeCmdCtor[T Command](args ...any) Constructor {
	return func() Command {
		var skill T
		return skill.New(args...)
	}
}

func makeCmdParsers(commands []Constructor, ast *parsec.AST) []interface{} {
	parsers := make([]interface{}, 0, len(commands))
	for _, cmd := range commands {
		parsers = append(parsers, cmd().Parser(ast))
	}
	return parsers
}

func makeParser() (*parsec.AST, parsec.Parser) {
	ast := parsec.NewAST("COMMAND", 20)
	return ast,
		ast.Kleene("TEXT", nil,
			ast.OrdChoice("COMMAND", nil,
				makeCmdParsers(cmdCtors, ast)...),
			ast.End("END"))
}

func parseText(ast *parsec.AST, parser parsec.Parser, text string) parsec.Queryable {
	scanner := parsec.NewScanner([]byte(text))
	query, _ := ast.Parsewith(parser, scanner)
	log.Printf("parseText('%s'):\n", text)
	ast.Prettyprint()
	return query
}

func makeCommand(ctors []Constructor, ctx Context, query parsec.Queryable) (Command, error) {
	if len(query.GetChildren()) != 1 {
		return nil, fmt.Errorf("invalid command")
	}
	var cmd Command
	var err error
	for _, ctor := range ctors {
		newCmd := ctor()
		if query.GetChildren()[0].GetName() == newCmd.Name() {
			err = newCmd.Init(ctx, query.GetChildren()[0])
			if err != nil {
				err = fmt.Errorf("init command %s error: %w", newCmd.Name(), err)
			} else {
				cmd = newCmd
			}
			break
		}
	}
	return cmd, err
}

func parseCommand(ctors []Constructor, ast *parsec.AST, parser parsec.Parser, ctx Context, text string) (Command, error) {
	query := parseText(ast, parser, text)
	cmd, err := makeCommand(ctors, ctx, query)
	if err != nil {
		return nil, fmt.Errorf("makeCommand() error: %w", err)
	}
	return cmd, nil
}

func ParseCommand(ctx Context, text string) (Command, error) {
	return parseCommand(cmdCtors, cmdAST, cmdParser, ctx, text)
}

func makeSkillCtor[T Command](op string) Constructor {
	return func() Command {
		var skill T
		return skill.New(op)
	}
}

func makeUnitSkillParsers(commands []Constructor, ast *parsec.AST) []interface{} {
	parsers := make([]interface{}, 0, len(commands))
	for _, cmd := range commands {
		parsers = append(parsers, cmd().Parser(ast))
	}
	return parsers
}

func makeSkillParsers() (map[string]*parsec.AST, map[string]parsec.Parser) {
	asts := make(map[string]*parsec.AST)
	parsers := make(map[string]parsec.Parser)
	for unit, ctors := range unitSkillCtors {
		ast := parsec.NewAST(unit+"_SKILL", 20)
		asts[unit] = ast
		parsers[unit] = ast.Kleene("TEXT", nil,
			ast.OrdChoice("SKILL", nil,
				makeUnitSkillParsers(ctors, ast)...),
			ast.End("END"))
	}
	return asts, parsers
}

func makeSkillCommand(ctx Context, query parsec.Queryable) (Command, error) {
	if len(query.GetChildren()) != 1 {
		return nil, fmt.Errorf("invalid skill")
	}
	var skill Command
	var err error
	for _, skillCtor := range unitSkillCtors[ctx.Unit] {
		newSkill := skillCtor()
		if query.GetChildren()[0].GetName() == newSkill.Name() {
			err = newSkill.Init(ctx, query.GetChildren()[0])
			if err != nil {
				err = fmt.Errorf("init skill %s error: %w", newSkill.Name(), err)
			} else {
				skill = newSkill
			}
			break
		}
	}
	return skill, err
}

func parseUnitSkill(ctx Context, text string) (Command, error) {
	ast, ok := skillASTs[ctx.Unit]
	if !ok {
		return nil, nil
	}
	parser := skillParsers[ctx.Unit]
	query := parseText(ast, parser, text)
	return makeSkillCommand(ctx, query)
}

func parsecInt(name string, maxLen int) parsec.Parser {
	return parsec.Token(`-?[0-9]{1,`+strconv.Itoa(maxLen)+"}", name)
}

func parsecUint(name string, maxLen int) parsec.Parser {
	return parsec.Token(`[0-9]{1,`+strconv.Itoa(maxLen)+"}", name)
}

func parsecWS() parsec.Parser {
	return parsec.TokenExact(`\s+`, "WS")
}
