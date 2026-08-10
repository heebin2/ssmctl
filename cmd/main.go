package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/heebin2/ssmctl/internal/ssm"
)

type command int

const (
	cmdInit command = iota
	cmdList
	cmdConnect // connect to instance (ssmctl <instance-name>)
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		exit("error: cannot determine home directory")
	}
	defaultConfig := filepath.Join(homeDir, ".ssm.yml")

	var cfgPath string
	flag.StringVar(&cfgPath, "config", defaultConfig, "path to config file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: ssmctl [-config path] <init|list|instance-name>\n")
	}
	flag.Parse()

	if flag.NArg() < 1 {
		exit("usage: ssmctl [-config path] <init|list|instance-name>")
	}

	arg := flag.Arg(0)
	cmd := parseCommand(arg)

	switch cmd {
	case cmdInit:
		if err := ssm.InitConfig(cfgPath); err != nil {
			exit(fmt.Sprintf("error: %v", err))
		}
		fmt.Printf("config initialized at %s\n", cfgPath)
		return
	}

	cfg, err := ssm.LoadConfig(cfgPath)
	if err != nil {
		handleConfigError(cfgPath, err)
	}

	switch cmd {
	case cmdList:
		ssm.PrintList(cfg)
	case cmdConnect:
		if err := ssm.Start(cfg, arg); err != nil {
			exit(fmt.Sprintf("error: %v", err))
		}
	}
}

func parseCommand(arg string) command {
	switch arg {
	case "init":
		return cmdInit
	case "list":
		return cmdList
	default:
		// any other argument is treated as instance name
		return cmdConnect
	}
}

func handleConfigError(cfgPath string, err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
	fmt.Fprintf(os.Stderr, "\nhelp:\n")
	fmt.Fprintf(os.Stderr, "  1. initialize config from AWS EC2 instances:\n")
	fmt.Fprintf(os.Stderr, "     ssmctl init\n\n")
	fmt.Fprintf(os.Stderr, "  2. or manually create %s with:\n", cfgPath)
	fmt.Fprintf(os.Stderr, "     global:\n")
	fmt.Fprintf(os.Stderr, "       user: ec2-user\n")
	fmt.Fprintf(os.Stderr, "     instances:\n")
	fmt.Fprintf(os.Stderr, "       my-instance:\n")
	fmt.Fprintf(os.Stderr, "         target: i-1234567890abcdef0\n")
	fmt.Fprintf(os.Stderr, "  3. or use: ssmctl -config <path> list\n\n")

	abs, _ := filepath.Abs(cfgPath)
	fmt.Fprintf(os.Stderr, "expected location: %s\n", abs)
	os.Exit(1)
}

func exit(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
	os.Exit(1)
}
