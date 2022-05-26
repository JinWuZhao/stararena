package data

type Unit struct {
	Name string
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
		Name: "hellion",
	}
	UnitSiegeTank = Unit{
		Name: "siege-tank",
	}
	UnitThor = Unit{
		Name: "thor",
	}
	UnitBattlecruiser = Unit{
		Name: "battlecruiser",
	}
)
