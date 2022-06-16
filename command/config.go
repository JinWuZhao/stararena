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

var unitDataAcc = map[string]string{
	"aa": "MengskHellion",
	"ba": "SiegeTank",
	"ca": "ThorAP",
	"da": "Battlecruiser",
}
