package command

var cmdCtors = []Constructor{
	MakeCmdCtor[*MoveCmd](),
	MakeCmdCtor[*ChangeModeCmd](),
	MakeCmdCtor[*JoinGameCmd](),
	MakeCmdCtor[*SetTemplateCmd](),
	MakeCmdCtor[*SetWeaponCmd](),
	MakeCmdCtor[*SetAbilityCmd](),
	MakeCmdCtor[*AssignPointsCmd](),
	MakeCmdCtor[*ShowInfoCmd](),
	MakeCmdCtor[*UpvoteCmd](),
	MakeCmdCtor[*SetNoticeCmd](),
}
