package genshinartis

var weaponPJWSFullStacks = weapon{
	baseAtk: 674,
	stats:   map[stat]float32{CritRate: 22.1, ATKP: 3.2 * 7, GlobalDMGBonus: 12},
}

var weaponHomaPassiveOff = weapon{
	baseAtk: 608,
	stats:   map[stat]float32{CritDmg: 66.2, HPP: 20},
	passive: func(s map[stat]float32) map[stat]float32 {
		return map[stat]float32{ATK: s[HP] * 0.008}
	},
}

var weaponHomaPassiveOn = weapon{
	baseAtk: 608,
	stats:   map[stat]float32{CritDmg: 66.2, HPP: 20},
	passive: func(s map[stat]float32) map[stat]float32 {
		return map[stat]float32{ATK: s[HP] * 1.8}
	},
}
