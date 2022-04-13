package main

import "gozar/cmd"

var Version string

func main() {
	cmd.Version = Version
	cmd.Execute()
}
