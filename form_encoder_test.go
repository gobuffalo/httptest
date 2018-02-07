package willie_test

import (
	"testing"

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
		Alias            Alias
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
		},
		Alias{
			"Tony",
			"Friendly",
		},
	}

	encoded := willie.EncodeToFormValues(val)

	r.NotNil(encoded["Name"])
	r.Equal("Antonio", encoded.Get("Name"))
	r.Equal("3", encoded.Get("Kids"))
	r.Equal("99.1", encoded.Get("DesiredGolangLvl"))

	r.Equal("Leopoldo", encoded.Get("KidNames[0]"))
	r.Equal("Marco-polo", encoded.Get("KidNames[1]"))
	r.Equal("Pancracia", encoded.Get("KidNames[2]"))

	r.Equal("B", encoded.Get("Notes.A"))
	r.Equal("Tony", encoded.Get("Alias.Name"))
	r.Equal("Friendly", encoded.Get("Alias.Type"))
}
