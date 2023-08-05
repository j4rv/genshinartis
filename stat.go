package genshinartis

import (
	"math/rand"
)

type stat int

const (
	HP stat = iota
	ATK
	DEF
	HPP
	ATKP
	DEFP
	EnergyRecharge
	ElementalMastery
	CritRate
	CritDmg
	PyroDMG
	ElectroDMG
	CryoDMG
	HydroDMG
	AnemoDMG
	GeoDMG
	DendroDMG
	PhysDMG
	HealingBonus

	GlobalDMGBonus
	BaseDMGIncrease
)

var substatValues map[stat][4]float32 = map[stat][4]float32{
	HP:               {209.13, 239.00, 268.88, 298.75},
	ATK:              {13.62, 15.56, 17.51, 19.45},
	DEF:              {16.20, 18.52, 20.83, 23.15},
	HPP:              {4.08, 4.66, 5.25, 5.83},
	ATKP:             {4.08, 4.66, 5.25, 5.83},
	DEFP:             {5.10, 5.83, 6.56, 7.29},
	ElementalMastery: {16.32, 18.65, 20.98, 23.31},
	EnergyRecharge:   {4.53, 5.18, 5.83, 6.48},
	CritRate:         {2.72, 3.11, 3.50, 3.89},
	CritDmg:          {5.44, 6.22, 6.99, 7.77},
}

var mainStatValues map[stat]float32 = map[stat]float32{
	HP:               4780,
	ATK:              311,
	HPP:              46.6,
	ATKP:             46.6,
	DEFP:             58.3,
	ElementalMastery: 186.5,
	EnergyRecharge:   51.8,
	PyroDMG:          46.6,
	ElectroDMG:       46.6,
	CryoDMG:          46.6,
	HydroDMG:         46.6,
	AnemoDMG:         46.6,
	GeoDMG:           46.6,
	DendroDMG:        46.6,
	PhysDMG:          58.3,
	CritRate:         31.1,
	CritDmg:          62.2,
	HealingBonus:     35.9,
}

func (s stat) String() string {
	switch s {
	case HP:
		return "HP"
	case ATK:
		return "ATK"
	case DEF:
		return "DEF"
	case HPP:
		return "HP%"
	case ATKP:
		return "ATK%"
	case DEFP:
		return "DEF%"
	case EnergyRecharge:
		return "Energy Recharge%"
	case ElementalMastery:
		return "Elemental Mastery"
	case CritRate:
		return "CRIT Rate%"
	case CritDmg:
		return "CRIT DMG%"
	case PyroDMG:
		return "Pyro DMG%"
	case ElectroDMG:
		return "Electro DMG%"
	case CryoDMG:
		return "Cryo DMG%"
	case HydroDMG:
		return "Hydro DMG%"
	case AnemoDMG:
		return "Anemo DMG%"
	case GeoDMG:
		return "Geo DMG%"
	case DendroDMG:
		return "Dendro DMG%"
	case PhysDMG:
		return "Physical DMG%"
	case HealingBonus:
		return "Healing Bonus%"
	}
	return "Unknown"
}

func (s stat) RandomRollValue() float32 {
	return substatValues[s][rand.Intn(4)]
}

// Weights from https://genshin-impact.fandom.com/wiki/Artifacts/Distribution
// And https://genshin-impact.fandom.com/wiki/Artifacts/Stats

var sandsWeightedStats = map[stat]int{
	HPP:              26_680,
	ATKP:             26_660,
	DEFP:             26_660,
	EnergyRecharge:   10_000,
	ElementalMastery: 10_000,
}

var gobletWeightedStats = map[stat]int{
	HPP:              19_175,
	ATKP:             19_175,
	DEFP:             19_150,
	PyroDMG:          5_000,
	ElectroDMG:       5_000,
	CryoDMG:          5_000,
	HydroDMG:         5_000,
	AnemoDMG:         5_000,
	GeoDMG:           5_000,
	DendroDMG:        5_000,
	PhysDMG:          5_000,
	ElementalMastery: 2_500,
}

var circletWeightedStats = map[stat]int{
	HPP:              22_000,
	ATKP:             22_000,
	DEFP:             22_000,
	CritRate:         10_000,
	CritDmg:          10_000,
	HealingBonus:     10_000,
	ElementalMastery: 4_000,
}

const (
	flatSubstatWeight   = 150
	commonSubstatWeight = 100
	critSubstatWeight   = 75
)

func weightedSubstats(mainStat stat) map[stat]int {
	weightedSubs := map[stat]int{
		HP:               flatSubstatWeight,
		ATK:              flatSubstatWeight,
		DEF:              flatSubstatWeight,
		HPP:              commonSubstatWeight,
		ATKP:             commonSubstatWeight,
		DEFP:             commonSubstatWeight,
		EnergyRecharge:   commonSubstatWeight,
		ElementalMastery: commonSubstatWeight,
		CritRate:         critSubstatWeight,
		CritDmg:          critSubstatWeight,
	}
	delete(weightedSubs, mainStat)
	return weightedSubs
}
