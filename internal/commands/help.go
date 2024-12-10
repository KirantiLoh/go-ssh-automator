package commands

import "fmt"

func ShowHelp() {
	fmt.Println(`
Usage: ssh-automator --config <filename> 

Options:
  --config          Path to config file
    `)
}
