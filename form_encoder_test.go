package willie_test

import (
	"log"
	"testing"
	"time"

	"github.com/markbates/willie"
	"github.com/stretchr/testify/require"
)

func Test_FormEncoder(t *testing.T) {
	r := require.New(t)

	type Alias struct {
		Name string
		Type string
	}

	val := struct {
		Name             string
		Kids             int
		DesiredGolangLvl float32
		KidNames         []string
		Notes            map[string]string
		Notes2           map[string][]string
		Alias            Alias
		DateCreated      time.Time
	}{
		"Antonio",
		3,
		99.1,
		[]string{
			"Leopoldo",
			"Marco-polo",
			"Pancracia",
		},
		map[string]string{
			"A": "B",
			"C": "D",
			"D": "G",
		},
		map[string][]string{
			"H": []string{"I", "J", "K"},
		},
		Alias{
			"Tony",
			"Friendly",
		},
		time.Now(),
	}

	values, _ := willie.EncodeToURLValues(val)

	log.Println(values)

	r.NotNil(values["Name"])
	r.Equal("Antonio", values.Get("Name"))
	r.Equal("3", values.Get("Kids"))
	r.Equal("99.1", values.Get("DesiredGolangLvl"))

	r.Equal("Leopoldo", values.Get("KidNames[0]"))
	r.Equal("Marco-polo", values.Get("KidNames[1]"))
	r.Equal("Pancracia", values.Get("KidNames[2]"))

	r.Equal("B", values.Get("Notes[A]"))
	r.Equal("Tony", values.Get("Alias.Name"))
	r.Equal("Friendly", values.Get("Alias.Type"))
}
