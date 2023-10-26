package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"github.com/tufin/oasdiff/report"
	"gopkg.in/yaml.v3"
)

var base, revision, filter, filterExtension, format, lang, warnIgnoreFile, errIgnoreFile string
var breakingOnly bool

const (
	formatYAML = "yaml"
	formatText = "text"
	formatHTML = "html"
)

func init() {
	flag.StringVar(&base, "base", "", "path or URL of original OpenAPI spec in YAML or JSON format")
	flag.StringVar(&revision, "revision", "", "path or URL of revised OpenAPI spec in YAML or JSON format")
	flag.BoolVar(&breakingOnly, "breaking-only", true, "display breaking changes only (old method)")
	flag.StringVar(&format, "format", formatYAML, "output format: yaml, text or html")
}

func validateFlags() bool {
	supportedFormats := map[string]bool{"": true, "yaml": true, "text": true, "html": true}
	if !supportedFormats[format] {
		fmt.Fprintf(os.Stderr, "invalid format. Should be yaml, text or html\n")
		return false
	}
	return true
}

func main() {
	flag.Parse()

	if !validateFlags() {
		os.Exit(101)
	}

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	config := diff.NewConfig()
	config.BreakingOnly = breakingOnly

	actual_version, err := load.From(loader, "assets/version1.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load base spec from %q with %v\n", base, err)
		os.Exit(102)
	}

	revised_version, err := load.From(loader, "assets/version2.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load revision spec from %q with %v\n", revision, err)
		os.Exit(103)
	}

	diffReport, err := diff.Get(config, actual_version, revised_version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "diff failed with %v\n", err)
		os.Exit(104)
	}

	switch {
	case format == formatYAML:
		if err = printYAML(diffReport); err != nil {
			fmt.Fprintf(os.Stderr, "failed to print diff YAML with %v\n", err)
			os.Exit(106)
		}
	case format == formatText:
		fmt.Printf("%s", report.GetTextReportAsString(diffReport))
	case format == formatHTML:
		html, err := report.GetHTMLReportAsString(diffReport)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate HTML diff report with %v\n", err)
			os.Exit(107)
		}
		fmt.Printf("%s", html)
	default:
		fmt.Fprintf(os.Stderr, "unknown output format %q\n", format)
		os.Exit(108)
	}

	exitNormally(diffReport.Empty())
}

func exitNormally(diffEmpty bool) {
	if !diffEmpty {
		os.Exit(1)
	}
	os.Exit(0)
}

func printYAML(output interface{}) error {
	if reflect.ValueOf(output).IsNil() {
		return nil
	}

	bytes, err := yaml.Marshal(output)
	if err != nil {
		return err
	}
	fmt.Printf("%s", bytes)
	return nil
}
