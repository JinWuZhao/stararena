package command

var cmdCtors = []Constructor{
	MakeCmdCtor[*MoveCmd](),
	MakeCmdCtor[*MoveXCmd](),
	MakeCmdCtor[*ChangeModeCmd](),
	MakeCmdCtor[*JoinGameCmd](),
	MakeCmdCtor[*SetWeaponCmd](),
	MakeCmdCtor[*SetAbilityCmd](),
	MakeCmdCtor[*AssignPointsCmd](),
	MakeCmdCtor[*ShowInfoCmd](),
}
