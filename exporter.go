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
	for _, s := range art.SubStats {
		subs = append(subs, GOODSubstat{
			Stat:  s.Stat.goodKey(),
			Value: s.RoundedValue(),
		})
	}
	return GOODArtifact{
		Set:      string(art.Set),
		Rarity:   5, // TODO: Change when you implement a 4* generator!
		Level:    20,
		Slot:     art.Slot.goodKey(),
		MainStat: art.MainStat.goodKey(),
		Subs:     subs,
		Location: "",
		Lock:     false,
	}
}
