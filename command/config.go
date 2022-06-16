package command

var cmdCtors = []Constructor{
	MakeCmdCtor[*SetUnitCmd](),
	MakeCmdCtor[*SetServantsCmd](),
	MakeCmdCtor[*MoveCmd](),
	MakeCmdCtor[*MoveXCmd](),
	MakeCmdCtor[*ChangeModeCmd](),
	MakeCmdCtor[*JoinGameCmd](),
	MakeCmdCtor[*ShowPointsCmd](),
}
