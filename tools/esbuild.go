//go:build dev

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

func main() {
	bundle := flag.Bool("bundle", false, "Bundle all dependencies into the output files")
	outFile := flag.String("outfile", "", "The output file (for one entry point)")
	outDir := flag.String("outdir", "", "The output directory (for multiple entry points)")
	platform := flag.String("platform", "", "Platform target (browser | node | neutral)")
	minify := flag.Bool("minify", false, "Minify the output (sets all --minify-* flags)")
	sourceMap := flag.Bool("sourcemap", false, "Emit a source map")
	format := flag.String("format", "", "Output format (iife | cjs | esm, no default when not bundling, otherwise default is iife when platform is browser and cjs when platform is node)")
	target := flag.String("target", "", "Environment target (e.g. es2017, chrome58, firefox57, safari11, edge16, node10, ie9, opera45, default esnext)")
	watch := flag.Bool("watch", false, "Watch mode: rebuild on file system changes")

	flag.Parse()

	var p api.Platform
	switch *platform {
	case "browser":
		p = api.PlatformBrowser
	case "node":
		p = api.PlatformNode
	case "neutral":
		p = api.PlatformNeutral
	default:
		p = api.PlatformBrowser
	}

	var f api.Format
	switch *format {
	case "iife":
		f = api.FormatIIFE
	case "cjs":
		f = api.FormatCommonJS
	case "esm":
		f = api.FormatESModule
	default:
		f = api.FormatDefault
	}

	var t api.Target

	switch *target {
	case "es2015":
		t = api.ES2015
	case "es2016":
		t = api.ES2016
	case "es2017":
		t = api.ES2017
	case "es2018":
		t = api.ES2018
	case "es2019":
		t = api.ES2019
	case "es2020":
		t = api.ES2020
	case "es2021":
		t = api.ES2021
	case "es2022":
		t = api.ES2022
	case "es2023":
		t = api.ES2023
	case "es2024":
		t = api.ES2024
	default:
		t = api.ESNext
	}

	opts := api.BuildOptions{
		EntryPoints:       flag.Args(),
		Bundle:            *bundle,
		Outfile:           *outFile,
		Outdir:            *outDir,
		Format:            f,
		Platform:          p,
		Target:            t,
		MinifySyntax:      *minify,
		MinifyWhitespace:  *minify,
		MinifyIdentifiers: *minify,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
	}

	if *sourceMap {
		opts.Sourcemap = api.SourceMapLinked
	}

	if *watch {
		ctx, err := api.Context(opts)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := ctx.Watch(api.WatchOptions{}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		<-make(chan struct{})
		return
	}

	result := api.Build(opts)

	if len(result.Errors) != 0 {
		for _, err := range result.Errors {
			fmt.Println(err)
		}

		os.Exit(1)
	}
}
