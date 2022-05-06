package command

import "github.com/jinwuzhao/stararena/data"

var cmdCtors = []commandConstructor{
	makeCmdCtor[*ChangeUnitCmd](),
	makeCmdCtor[*IssueSkillCmd](),
	makeCmdCtor[*MoveCmd](),
	makeCmdCtor[*MoveXCmd](),
	makeCmdCtor[*ChangeModeCmd](),
	makeCmdCtor[*JoinGameCmd](),
}

var unitDataAcc = map[string]data.Unit{
	"1": data.UnitHellion,
	"2": data.UnitSiegeTank,
	"3": data.UnitThor,
	"4": data.UnitBattlecruiser,
}

var unitSkillCtors = map[string][]commandConstructor{
	data.UnitSiegeTank.Name: {
		makeSkillCtor[*SiegeMode]("1"),
		makeSkillCtor[*TankMode]("2"),
	},
	data.UnitThor.Name: {
		makeSkillCtor[*ExplosivePayload]("1"),
		makeSkillCtor[*HighImpactPayload]("2"),
	},
	data.UnitBattlecruiser.Name: {
		makeSkillCtor[*YamatoCannon]("1"),
		makeSkillCtor[*TacticalJump]("2"),
	},
}
