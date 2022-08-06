package genshinartis

import (
	"encoding/json"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestRandomArtifactFromDomain(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var set1Count, set2Count int
	set1, set2 := "Emblem", "Shimenawa"

	// Generate 1000 artifacts from two sets
	for i := 0; i < 1000; i++ {
		art := RandomArtifactFromDomain(set1, set2)
		if art.Set == artifactSet(set1) {
			set1Count++
		} else if art.Set == artifactSet(set2) {
			set2Count++
		} else {
			t.Error("Unexpected artifact set: " + art.Set)
		}
	}

	// Then check that the chances of getting an artifact from either set is ~50%
	if set1Count < 450 || set1Count > 550 {
		t.Error("Too many or too few artifacts from set 1: " + strconv.Itoa(set1Count))
	}
	if set2Count < 450 || set2Count > 550 {
		t.Error("Too many or too few artifacts from set 1: " + strconv.Itoa(set2Count))
	}
}

func TestRemoveTrashArtifacts(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	var artis []*Artifact
	set1, set2 := "Emblem", "Shimenawa"

	// Generate 1000 artifacts from two sets
	for i := 0; i < 10000; i++ {
		artis = append(artis, RandomArtifactFromDomain(set1, set2))
	}

	subs := map[artifactStat]float32{
		ATKP:           1,
		CritRate:       1,
		CritDmg:        1,
		EnergyRecharge: 0.5,
		ATK:            0.25,
	}
	filtered := RemoveTrashArtifacts(artis, subs, 5)
	for _, a := range filtered {
		t.Log(*a)
	}
}

func TestExportToGOOD(t *testing.T) {
	var artis []*Artifact
	for i := 0; i < 1_000_000; i++ {
		artis = append(artis, RandomArtifact(DomainBase4Chance))
	}
	subs := map[artifactStat]float32{
		ATKP:             1,
		CritRate:         1,
		CritDmg:          1,
		ElementalMastery: 0.33,
		EnergyRecharge:   0.25,
		ATK:              0.25,
	}
	artis = RemoveTrashArtifacts(artis, subs, 10)
	export := ExportToGOOD(artis)
	b, err := json.Marshal(export)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("goodExportTest.json", b, 0755)
}

func TestExportToGOODEmblemHell(t *testing.T) {
	// 6 months of Emblem -> 1620 runs
	// 6 months of Emblem -> 4320 runs if max refreshing
	var artis []*Artifact
	for i := 0; i < 4320; i++ {
		artis = append(artis, RandomArtifactFromDomain("EmblemOfSeveredFate", "ShimenawasReminiscence"))
	}
	subs := map[artifactStat]float32{
		CritRate:         1,
		CritDmg:          1,
		ATKP:             0.8,
		EnergyRecharge:   0.8,
		ElementalMastery: 0.5,
		ATK:              0.25,
	}
	artis = RemoveTrashArtifacts(artis, subs, 10)
	export := ExportToGOOD(artis)
	b, err := json.Marshal(export)
	if err != nil {
		t.Error(err)
	}
	os.WriteFile("goodExportEmblemTest.json", b, 0755)
}
