package command

var cmdCtors = []Constructor{
	MakeCmdCtor[*MoveCmd](),
	MakeCmdCtor[*MoveXCmd](),
	MakeCmdCtor[*ChangeModeCmd](),
	MakeCmdCtor[*JoinGameCmd](),
	MakeCmdCtor[*ShowPointsCmd](),
}
