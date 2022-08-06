package genshinartis

type artifactSlot int

const (
	SlotFlower artifactSlot = iota
	SlotPlume
	SlotSands
	SlotGoblet
	SlotCirclet
)

func (t artifactSlot) String() string {
	switch t {
	case SlotFlower:
		return "Flower of Life"
	case SlotPlume:
		return "Plume of Death"
	case SlotSands:
		return "Sands of Eon"
	case SlotGoblet:
		return "Goblet of Eonothem"
	case SlotCirclet:
		return "Circlet of Logos"
	}
	return "Unknown"
}

func (t artifactSlot) goodKey() string {
	switch t {
	case SlotFlower:
		return "flower"
	case SlotPlume:
		return "plume"
	case SlotSands:
		return "sands"
	case SlotGoblet:
		return "goblet"
	case SlotCirclet:
		return "circlet"
	}
	return "Unknown"
}
