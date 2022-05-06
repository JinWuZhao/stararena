package data

type Unit struct {
	Name   string
	Points int64
}

var (
	UnitHellion = Unit{
		Name:   "hellion",
		Points: 0,
	}
	UnitSiegeTank = Unit{
		Name:   "siege-tank",
		Points: 100,
	}
	UnitThor = Unit{
		Name:   "thor",
		Points: 300,
	}
	UnitBattlecruiser = Unit{
		Name:   "battlecruiser",
		Points: 500,
	}
)
