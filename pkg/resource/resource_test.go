package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
)

func TestResourceAttributes_SetDefault(t *testing.T) {
	a := &ResourceAttributes{
		Attrs: map[string]interface{}{
			"setString":  nil,
			"setBool":    nil,
			"setNumber":  nil,
			"listString": nil,
			"listBool":   nil,
			"listNumber": nil,
			"mapString":  nil,
			"mapBool":    nil,
			"mapNumber":  nil,

			"setStringEmpty":  []string{},
			"setBoolEmpty":    []bool{},
			"setNumberEmpty":  []int{},
			"listStringEmpty": []string{},
			"listBoolEmpty":   []bool{},
			"listNumberEmpty": []int{},
			"mapStringEmpty":  map[string]string{},
			"mapBoolEmpty":    map[string]bool{},
			"mapNumberEmpty":  map[string]int{},

			"restrictions": []map[string][]interface{}{
				{
					"geo_restriction": []interface{}{
						map[string]interface{}{
							"locations": []string{},
						},
					},
				},
			},
		},
		Type: cty.Object(map[string]cty.Type{
			"setString":  cty.Set(cty.String),
			"setBool":    cty.Set(cty.Bool),
			"setNumber":  cty.Set(cty.Number),
			"listString": cty.List(cty.String),
			"listBool":   cty.List(cty.Bool),
			"listNumber": cty.List(cty.Number),
			"mapString":  cty.Map(cty.String),
			"mapBool":    cty.Map(cty.Bool),
			"mapNumber":  cty.Map(cty.Number),

			"setStringEmpty":  cty.Set(cty.String),
			"setBoolEmpty":    cty.Set(cty.Bool),
			"setNumberEmpty":  cty.Set(cty.Number),
			"listStringEmpty": cty.List(cty.String),
			"listBoolEmpty":   cty.List(cty.Bool),
			"listNumberEmpty": cty.List(cty.Number),
			"mapStringEmpty":  cty.Map(cty.String),
			"mapBoolEmpty":    cty.Map(cty.Bool),
			"mapNumberEmpty":  cty.Map(cty.Number),

			"restrictions": cty.List(cty.Object(map[string]cty.Type{
				"geo_restriction": cty.List(cty.Object(map[string]cty.Type{
					"locations": cty.Set(cty.String),
				})),
			})),
		}),
	}
	cases := map[string]struct {
		value interface{}
		path  []string
		exist bool
	}{
		"setString":  {path: []string{"setString"}, exist: false},
		"setBool":    {path: []string{"setBool"}, exist: false},
		"setNumber":  {path: []string{"setNumber"}, exist: false},
		"listString": {path: []string{"listString"}, exist: false},
		"listBool":   {path: []string{"listBool"}, exist: false},
		"listNumber": {path: []string{"listNumber"}, exist: false},
		"mapString":  {path: []string{"mapString"}, exist: false},
		"mapBool":    {path: []string{"mapBool"}, exist: false},
		"mapNumber":  {path: []string{"mapNumber"}, exist: false},

		"setStringEmpty":  {path: []string{"setStringEmpty"}, exist: false},
		"setBoolEmpty":    {path: []string{"setBoolEmpty"}, exist: false},
		"setNumberEmpty":  {path: []string{"setNumberEmpty"}, exist: false},
		"listStringEmpty": {path: []string{"listStringEmpty"}, exist: false},
		"listBoolEmpty":   {path: []string{"listBoolEmpty"}, exist: false},
		"listNumberEmpty": {path: []string{"listNumberEmpty"}, exist: false},
		"mapStringEmpty":  {path: []string{"mapStringEmpty"}, exist: false},
		"mapBoolEmpty":    {path: []string{"mapBoolEmpty"}, exist: false},
		"mapNumberEmpty":  {path: []string{"mapNumberEmpty"}, exist: false},

		"restrictions": {path: []string{"restrictions", "0", "geo_restriction", "0", "locations"}, exist: false},
	}

	for k, v := range cases {
		t.Run(k, func(t *testing.T) {
			if err := a.SetDefault(v.path); err != nil {
				t.Errorf("SetDefault() error = %v", err)
			}

			val, exist := a.Get(k) // TODO
			assert.Equal(t, v.exist, exist)
			if exist {
				assert.Equal(t, v.value, val)
			}
		})
	}
}
