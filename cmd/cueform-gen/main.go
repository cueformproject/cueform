package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/format"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/jcrqr/cueform/cmd/cueform-gen/internal/config"
	"github.com/jcrqr/cueform/cmd/cueform-gen/internal/generate"
	"github.com/spf13/cobra"
)

type options struct {
	debug          bool
	configFilename string
}

func main() {
	opts := &options{}
	cmd := &cobra.Command{
		Use:   "cueform-gen",
		Short: "Generates Terraform provider schemas in CUE format",
		RunE: func(cmd *cobra.Command, args []string) error {
			slogOpts := &slog.HandlerOptions{}
			if opts.debug {
				slogOpts.Level = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, slogOpts)))

			return doGenerate(cmd, args, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.configFilename, "config", "c", config.DefaultFilename, "path to the configuration file")
	cmd.Flags().BoolVarP(&opts.debug, "debug", "d", false, "enable debug level logging")

	cmd.Execute()
}

func doGenerate(cmd *cobra.Command, _ []string, opts *options) error {
	ctx := cmd.Context()

	cfg, err := config.Read(opts.configFilename)
	if err != nil {
		return err
	}

	slog.Debug("Read configuration",
		slog.String("command", cfg.Command),
		slog.String("output", cfg.Output),
		slog.Int("providers", len(cfg.Providers)))

	workdir, err := os.MkdirTemp(os.TempDir(), "cueform-gen-*")
	if err != nil {
		return err
	}

	err = writeTerraformConfig(filepath.Join(workdir, "cueform.tf.json"), cfg)
	if err != nil {
		return err
	}

	tf, err := tfexec.NewTerraform(workdir, cfg.Command)
	if err != nil {
		return err
	}

	slog.Debug("Initializing Terraform module", slog.String("workdir", workdir))
	if err := tf.Init(ctx); err != nil {
		return err
	}

	slog.Debug("Fetching Terraform provider schemas", slog.String("workdir", workdir))
	schemas, err := tf.ProvidersSchema(ctx)
	if err != nil {
		return err
	}

	for _, p := range cfg.Providers {
		var schema *tfjson.ProviderSchema
		for k, v := range schemas.Schemas {
			if strings.HasSuffix(k, p.Source) {
				schema = v
				break
			}
		}

		if schema == nil {
			slog.Warn("Provider schema not found", slog.String("provider", p.Source))
			continue
		}

		f, err := generate.Transform(p, schema)
		if err != nil {
			return err
		}

		if err := writeSchema(f); err != nil {
			return err
		}

		slog.Debug("Generated schema", slog.String("provider", p.Source), slog.String("filename", f.Filename))
	}

	return nil
}

func writeTerraformConfig(filename string, cfg *config.Config) error {
	tfcfg := map[string]map[string]map[string]map[string]string{
		"terraform": {
			"required_providers": {},
		},
	}
	for _, p := range cfg.Providers {
		tfcfg["terraform"]["required_providers"][p.Name] = map[string]string{
			"source":  p.Source,
			"version": p.Version,
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(tfcfg)
}

func writeSchema(f *ast.File) error {
	data, err := format.Node(f, format.Simplify())
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(f.Filename), 0777); err != nil {
		return err
	}

	osf, err := os.Create(f.Filename)
	if err != nil {
		return err
	}
	defer osf.Close()

	if _, err := osf.Write(data); err != nil {
		return err
	}
	return nil
}
