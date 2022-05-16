package data

type Unit struct {
	Name   string
	Cost   int64
	Reward int64
}

var unitByName = map[string]Unit{
	UnitHellion.Name:       UnitHellion,
	UnitSiegeTank.Name:     UnitSiegeTank,
	UnitThor.Name:          UnitThor,
	UnitBattlecruiser.Name: UnitBattlecruiser,
}

func GetUnitByName(name string) (unit Unit, ok bool) {
	unit, ok = unitByName[name]
	return
}

var (
	UnitHellion = Unit{
		Name:   "hellion",
		Cost:   0,
		Reward: 100,
	}
	UnitSiegeTank = Unit{
		Name:   "siege-tank",
		Cost:   100,
		Reward: 150,
	}
	UnitThor = Unit{
		Name:   "thor",
		Cost:   300,
		Reward: 200,
	}
	UnitBattlecruiser = Unit{
		Name:   "battlecruiser",
		Cost:   500,
		Reward: 300,
	}
)
