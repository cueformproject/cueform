package config

import (
	"os"

	_ "embed"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
)

const DefaultFilename = "cueform.gen.cue"

var (
	schema cue.Value

	//go:embed config.cue
	schemaContents []byte
)

func init() {
	ctx := cuecontext.New()
	inst := ctx.CompileBytes(schemaContents)
	if err := inst.Err(); err != nil {
		panic(err)
	}
	schema = inst.LookupPath(cue.MakePath(cue.Def("Config")))
}

type Config struct {
	Command   string               `json:"command"`
	Output    string               `json:"output"`
	Providers map[string]*Provider `json:"providers"`
}

type Provider struct {
	Name     string `json:"name"`
	Source   string `json:"source"`
	Version  string `json:"version"`
	Filename string `json:"filename"`
}

func Read(filename string) (*Config, error) {
	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	ctx := cuecontext.New()

	cv := ctx.CompileBytes(contents, cue.Filename(filename))
	if err := cv.Err(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := cv.Unify(schema).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
