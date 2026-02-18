package internal

import (
	"flag"
	"log"
)

type Inputs struct {
	Branch string
	Repos  []string
	Probe  bool
	Merge  bool
	Init   bool
}

func GetInputs() Inputs {
	configPath := flag.String("config", "", "Path to a JSON config file (relative or absolute)")
	branch := flag.String("branch", "", "The target branch you want to search for")
	probe := flag.Bool("probe", false, "If set to `true` no approvals will be made but all the outputs for testing will be made")
	autoMerge := flag.Bool("merge", false, "If set to `true` any approved prs will be merged automatically")
	init := flag.Bool("init", false, "Run the interactive config wizard to create a new config file")
	flag.Parse()

	explicitFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		explicitFlags[f.Name] = true
	})

	inputs := Inputs{}

	if explicitFlags["config"] {
		cfg, err := LoadConfig(*configPath)
		if err != nil {
			log.Fatalf("Failed to load config file %q: %v", *configPath, err)
		}
		inputs.Branch = cfg.Branch
		inputs.Repos = cfg.Repos
		inputs.Probe = cfg.Probe
		inputs.Merge = cfg.Merge
	}

	if explicitFlags["branch"] {
		inputs.Branch = *branch
	}
	if explicitFlags["probe"] {
		inputs.Probe = *probe
	}
	if explicitFlags["merge"] {
		inputs.Merge = *autoMerge
	}

	if len(flag.Args()) > 0 {
		inputs.Repos = flag.Args()
	}

	inputs.Init = *init

	return inputs
}
