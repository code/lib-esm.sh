package server

import (
	"encoding/json"
	"errors"
	"path"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type ImportMap struct {
	Imports map[string]string            `json:"imports"`
	Scopes  map[string]map[string]string `json:"scopes"`
}

type TransformInput struct {
	Code      string          `json:"code"`
	ImportMap json.RawMessage `json:"importMap"`
	Filename  string          `json:"filename"`
	Target    string          `json:"target"`
	SourceMap bool            `json:"sourceMap"`
}

type TransformOutput struct {
	Code string `json:"code"`
	Map  string `json:"map,omitempty"`
}

func transform(input TransformInput) (out TransformOutput, err error) {
	target := api.ESNext
	if input.Target != "" {
		if t, ok := targets[input.Target]; ok {
			target = t
		} else {
			err = errors.New("<400> invalid target")
			return
		}
	}

	loader := api.LoaderJS
	extname := path.Ext(input.Filename)
	switch extname {
	case ".jsx":
		loader = api.LoaderJSX
	case ".ts":
		loader = api.LoaderTS
	case ".tsx":
		loader = api.LoaderTSX
	}

	var importMap ImportMap
	err = json.Unmarshal(input.ImportMap, &importMap)
	if err != nil {
		return
	}

	imports := map[string]string{}
	trailingSlashImports := map[string]string{}
	jsxImportSource := ""

	// todo: use importMap.Scopes
	if len(importMap.Imports) > 0 {
		for key, value := range importMap.Imports {
			if value != "" {
				if strings.HasSuffix(key, "/") {
					trailingSlashImports[key] = value
				} else {
					if key == "@jsxImportSource" {
						jsxImportSource = value
					}
					imports[key] = value
				}
			}
		}
	}

	onResolver := func(args api.OnResolveArgs) (api.OnResolveResult, error) {
		path := args.Path
		if value, ok := imports[path]; ok {
			path = value
		} else {
			for key, value := range trailingSlashImports {
				if strings.HasPrefix(path, key) {
					path = value + path[len(key):]
					break
				}
			}
		}
		return api.OnResolveResult{
			Path:     path,
			External: true,
		}, nil
	}
	stdin := &api.StdinOptions{
		Contents:   input.Code,
		ResolveDir: "/",
		Sourcefile: input.Filename,
		Loader:     loader,
	}
	jsx := api.JSXTransform
	if jsxImportSource != "" {
		jsx = api.JSXAutomatic
	}
	opts := api.BuildOptions{
		Outdir:           "/esbuild",
		Stdin:            stdin,
		Platform:         api.PlatformBrowser,
		Format:           api.FormatESModule,
		Target:           target,
		JSX:              jsx,
		JSXImportSource:  jsxImportSource,
		Bundle:           true,
		TreeShaking:      api.TreeShakingFalse,
		MinifyWhitespace: true,
		MinifySyntax:     true,
		Write:            false,
		Plugins: []api.Plugin{
			{
				Name: "resolver",
				Setup: func(build api.PluginBuild) {
					build.OnResolve(api.OnResolveOptions{Filter: ".*"}, onResolver)
				},
			},
		},
	}
	if input.SourceMap {
		opts.Sourcemap = api.SourceMapExternal
	}
	ret := api.Build(opts)
	if len(ret.Errors) > 0 {
		err = errors.New("<400> failed to validate code: " + ret.Errors[0].Text)
		return
	}
	if len(ret.OutputFiles) == 0 {
		err = errors.New("<400> failed to validate code: no output files")
		return
	}
	for _, file := range ret.OutputFiles {
		if strings.HasSuffix(file.Path, ".js") {
			out.Code = string(file.Contents)
		} else if strings.HasSuffix(file.Path, ".map") {
			out.Map = string(file.Contents)
		}
	}
	return
}
