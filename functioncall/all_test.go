package functioncall

import (
	"encoding/json"
	"reflect"
	"testing"

	. "github.com/otiai10/mint"
)

func TestFunctions(t *testing.T) {
	funcs := Funcs{}
	Expect(t, funcs).TypeOf("functioncall.Funcs")
}

func TestFunctions_MarshalJSON(t *testing.T) {
	repeat := func(word string, count int) (r string) {
		for i := 0; i < count; i++ {
			r += word
		}
		return r
	}
	funcs := Funcs{
		"repeat": Func{Value: repeat, Description: "Repeat given string N times", Parameters: Params{
			{Name: "word", Type: "string", Description: "String to be repeated", Required: true},
			{Name: "count", Type: "number", Description: "How many times to repeat", Required: true},
		}},
	}
	b, err := funcs.MarshalJSON()
	Expect(t, err).ToBe(nil)

	v := []map[string]any{}
	err = json.Unmarshal(b, &v)
	Expect(t, err).ToBe(nil)

	Expect(t, v).Query("0.name").ToBe("repeat")
	Expect(t, v).Query("0.description").ToBe("Repeat given string N times")
	Expect(t, v).Query("0.parameters.type").ToBe("object")
	Expect(t, v).Query("0.parameters.properties.word.type").ToBe("string")
	Expect(t, v).Query("0.parameters.required.1").ToBe("count")
}

func TestAs(t *testing.T) {
	repeat := func(word string, count int) (r string) {
		for i := 0; i < count; i++ {
			r += word
		}
		return r
	}
	funcs := Funcs{
		"repeat": Func{Value: repeat, Description: "Repeat given string N times", Parameters: Params{
			{Name: "word", Type: "string", Description: "String to be repeated", Required: true},
			{Name: "count", Type: "number", Description: "How many times to repeat", Required: true},
		}},
	}
	a := As[[]map[string]any](funcs)
	Expect(t, a).TypeOf("[]map[string]interface {}")
	Expect(t, a).Query("0.name").ToBe("repeat")
	Expect(t, a).Query("0.parameters.type").ToBe("object")
}

func TestParams_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		params  Params
		want    []byte
		wantErr bool
	}{
		{
			name: "nested_example1",
			params: []Param{
				{
					Name:        "quality",
					Type:        "object",
					Description: "",
					Required:    true,
					Items: []Param{
						{
							Name:        "pros",
							Type:        "array",
							Description: "Write 3 points why this text is well written",
							Required:    true,
							Items: []Param{
								{Type: "string"},
							},
						},
					},
				},
			},
			want:    []byte(`{"properties":{"quality":{"properties":{"pros":{"description":"Write 3 points why this text is well written","items":{"type":"string"},"type":"array"}},"required":["pros"],"type":"object"}},"required":["quality"],"type":"object"}`),
			wantErr: false,
		},
		{
			name: "nested_example2",
			params: []Param{
				{
					Name:     "ingredients",
					Type:     "array",
					Required: true,
					Items: []Param{
						{
							Type: "object",
							Items: []Param{
								{
									Name:     "name",
									Type:     "string",
									Required: true,
								},
								{
									Name: "unit",
									Type: "string",
									// Enum: []any{"grams", "ml", "cups", "pieces", "teaspoons"},
									Required: true,
								},
								{
									Name:     "amount",
									Type:     "number",
									Required: true,
								},
							},
						},
					},
				},
				{
					Name:     "instructions",
					Type:     "array",
					Required: true,
					Items: []Param{
						{
							Type: "string",
						},
					},
					Description: "Steps to prepare the recipe (no numbering)",
				},
				{
					Name:        "time_to_cook",
					Type:        "number",
					Description: "Total time to prepare the recipe in minutes",
					Required:    true,
				},
			},
			want:    []byte(`{"type":"object","required":["ingredients","instructions","time_to_cook"],"properties":{"ingredients":{"type":"array","items":{"type":"object","required":["name","unit","amount"],"properties":{"name":{"type":"string"},"unit":{"type":"string"},"amount":{"type":"number"}}}},"instructions":{"type":"array","items":{"type":"string"},"description":"Steps to prepare the recipe (no numbering)"},"time_to_cook":{"type":"number","description":"Total time to prepare the recipe in minutes"}}}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotMap := map[string]interface{}{}
			err = json.Unmarshal(got, &gotMap)
			Expect(t, err).ToBe(nil)
			wantMap := map[string]interface{}{}
			err = json.Unmarshal(tt.want, &wantMap)
			Expect(t, err).ToBe(nil)
			if !reflect.DeepEqual(gotMap, wantMap) {
				t.Errorf("MarshalJSON() got = %s, want %s", gotMap, wantMap)
			}

		})
	}
}
