package internal

import (
	"flag"
)

type Inputs struct {
	Branch string
	Repos  []string
	Probe  bool
}

func GetInputs() Inputs {
	branch := flag.String("branch", "", "The target branch you want to search for")
	probe := flag.Bool("probe", false, "If set to `true` no approvals will be made but all the outputs for testing will be made")
	flag.Parse()

	return Inputs{*branch, flag.Args(), *probe}
}
