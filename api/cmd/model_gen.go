package cmd

type Replacable struct {
	Token string
	Value string
}

type Field struct {
	Name     string
	Type     string
	Unique   bool
	Required bool
}

type ModelConfig struct {
	Name   string
	Fields []*Field
}

type ModelGenerator struct {
	name   string
	fields []*Field
}

func NewModelGenerator(mc *ModelConfig) *ModelGenerator {
	return &ModelGenerator{mc.Name, mc.Fields}
}

func (mg *ModelGenerator) GetReplacables() []*Replacable {
	return []*Replacable{
		{Token: "ModelName", Value: ""},
	}
}

func (mg *ModelGenerator) GetPackagePath() string {
	return "internal/models"
}

func (mg *ModelGenerator) GetStubPath() string {
	return "./api/stubs/model.txt"
}

func (mg *ModelGenerator) Generate() error {
	return nil
}

// How I want the api to look like
// modelGen := cmd.NewModelGenerator(&cmd.ModelConfig{
// 	Fields: []*Field{{Type: "uint", Name: "ID", Unique: true, Required: true}}
// })
