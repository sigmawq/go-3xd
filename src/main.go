package main

import (
	"fmt"
	"os"
	"bytes"
	"bufio"
	"encoding/json"
	"gopkg.in/yaml.v2"
	"github.com/alexflint/go-arg"
)

func RunImport(verbose bool, all bool) {
	var prefix string
	if all {
		prefix = "[All]"
	}

	fmt.Printf("%v Import begins...\n", prefix)
	for i := 1; i < 100; i++ {
		if i % 10 == 0 && verbose {
			fmt.Printf("%v Outputting %v row..\n", prefix, i)
		}
	}

	fmt.Printf("%v 100 rows have been imported\n", prefix)
}	

type ImportCmd struct {}

type SyncCmd struct {
	All bool`arg:"-a,--all" help:"Display the [All] prefix"`
}

type Config struct {
	All bool `json: "all" yaml: "all"`
	Verbose bool `json: "verbose" yaml:"verbose"`
}

func loadJson() Config {
	file, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func loadYaml() Config {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}	

	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func initialize() (Config, bool) {
	// Default behaviour: Read YAML or JSON file. If such file exists the provided arguments will be overridden.
	// Note: config.JSON will be prioritised over config.yaml. But if config.yaml exists alongside with config.json there will be a warning
	var jsonExists bool
	var yamlExists bool
	if _, err := os.Stat("config.json"); err == nil {
		jsonExists = true
	}
	if _, err := os.Stat("config.yaml"); err == nil {
		yamlExists = true
	}

	var config Config
	if jsonExists && yamlExists {
		fmt.Println("WARNING: both config.json and config.yaml are present in the root folder. In such a case JSON file will be prioritised.")
	}

	if jsonExists {
		fmt.Println("Using JSON config.")
		config = loadJson()
		return config, true
	} else if yamlExists {
		fmt.Println("Using YAML config.")
		config = loadYaml()
		return config, true
	}

	return config, false
}

func main() {
	var args struct {
		Import *ImportCmd `arg:"subcommand:import" help:"Import 100 rows"`
		Sync   *SyncCmd   `arg:"subcommand:sync" help:"Import 100 rows with an optional sync"`
		Verbose bool `arg:"-v,--verbose" help:"Display verbose ouput"`
	}
	parser := arg.MustParse(&args)

	var all bool
	if args.Import != nil {
		all = false
	} else if args.Sync != nil {
		all = args.Sync.All
	} else {
		var buf bytes.Buffer
		bufW := bufio.NewWriter(&buf)
		parser.WriteHelp(bufW)
		bufW.Flush()
		fmt.Println(buf.String())
		return
	}

	config, hasConfig := initialize()
	if hasConfig {
		RunImport(config.Verbose, config.All)
		return
	}

	RunImport(args.Verbose, all)
}
