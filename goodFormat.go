package genshinartis

const goodFormatKey = "GOOD"
const goodVersion = 1
const goodExportSource = "Jarv ArtifactGEN"

type GOODArtifact struct {
	Set      string        `json:"setKey"`
	Rarity   int           `json:"rarity"`
	Level    int           `json:"level"`
	Slot     string        `json:"slotKey"`
	MainStat string        `json:"mainStatKey"`
	Subs     []GOODSubstat `json:"substats"`
	Location string        `json:"location"`
	Lock     bool          `json:"lock"`
}

type GOODSubstat struct {
	Stat  string  `json:"key"`
	Value float32 `json:"value"`
}

type GOODExport struct {
	Artifacts []GOODArtifact `json:"artifacts"`
	Format    string         `json:"format"`
	Version   int            `json:"version"`
	Source    string         `json:"source"`
}

func ExportToGOOD(arts []*Artifact) GOODExport {
	goodArts := []GOODArtifact{}
	for _, a := range arts {
		goodArts = append(goodArts, artifactToGOOD(a))
	}
	return GOODExport{
		Artifacts: goodArts,
		Format:    goodFormatKey,
		Version:   goodVersion,
		Source:    goodExportSource,
	}
}

func artifactToGOOD(art *Artifact) GOODArtifact {
	subs := []GOODSubstat{}
	for _, ss := range art.SubStats {
		subs = append(subs, GOODSubstat{
			Stat:  goodStatKey(ss.Stat),
			Value: ss.Value,
		})
	}
	return GOODArtifact{
		Set:      string(art.Set),
		Rarity:   5, // TODO: Change when you implement a 4* generator!
		Level:    20,
		Slot:     goodSlotKey(art.Slot),
		MainStat: goodStatKey(art.MainStat),
		Subs:     subs,
		Location: "",
		Lock:     false,
	}
}

func goodSlotKey(s artifactSlot) string {
	switch s {
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
	return "unknown"
}

func goodStatKey(s artifactStat) string {
	switch s {
	case HP:
		return "hp"
	case ATK:
		return "atk"
	case DEF:
		return "def"
	case HPP:
		return "hp_"
	case ATKP:
		return "atk_"
	case DEFP:
		return "def_"
	case EnergyRecharge:
		return "enerRech_"
	case ElementalMastery:
		return "eleMas"
	case CritRate:
		return "critRate_"
	case CritDmg:
		return "critDMG_"
	case PyroDMG:
		return "pyro_dmg_"
	case ElectroDMG:
		return "electro_dmg_"
	case CryoDMG:
		return "cryo_dmg_"
	case HydroDMG:
		return "hydro_dmg_"
	case AnemoDMG:
		return "anemo_dmg_"
	case GeoDMG:
		return "geo_dmg_"
	case DendroDMG:
		return "dendro_dmg_"
	case PhysDMG:
		return "physical_dmg_"
	case HealingBonus:
		return "heal_"
	}
	return "unknown"
}
