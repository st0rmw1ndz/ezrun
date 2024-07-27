package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/google/shlex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Flags struct {
	ConfigFile string
}

type Program struct {
	Name    string            `toml:"name"`
	Command string            `toml:"command"`
	Env     map[string]string `toml:"env"`
}

type Config struct {
	ChoiceCommand string    `toml:"choice_command"`
	Programs      []Program `toml:"programs"`
}

func parseFlags() Flags {
	var flags Flags
	flag.StringVar(&flags.ConfigFile, "c", "config.toml", "configuration file")

	flag.Parse()

	return flags
}

func getConfigFilePath(input string) (string, error) {
	if filepath.IsAbs(input) {
		return input, nil
	}

	if strings.ContainsAny(input, string(filepath.Separator)) {
		absPath, err := filepath.Abs(input)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	configDirectory := os.Getenv("XDG_CONFIG_HOME")
	if configDirectory == "" {
		homeDir, _ := os.UserHomeDir()
		configDirectory = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDirectory, "ezrun", input), nil
}

func chooseProgram(choiceCommand string, programs []Program) (Program, error) {
	var programNames string
	for _, program := range programs {
		programNames += program.Name + "\n"
	}

	cmdArgs, err := shlex.Split(choiceCommand)
	if err != nil {
		return Program{}, err
	}
	var out bytes.Buffer

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Stdin = strings.NewReader(programNames)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return Program{}, err
	}

	progName := strings.ReplaceAll(out.String(), "\n", "")
	for _, program := range programs {
		if program.Name == progName {
			return program, nil
		}
	}
	return Program{}, nil
}

func runSelectedProgram(program Program) error {
	cmdArgs, err := shlex.Split(program.Command)
	if err != nil {
		return err
	}

	env := os.Environ()
	for key, value := range program.Env {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func main() {
	flags := parseFlags()

	configFilePath := os.ExpandEnv(flags.ConfigFile)
	configFilePath, err := getConfigFilePath(configFilePath)
	if err != nil {
		fmt.Printf("error parsing config file path: %v\n", err)
		os.Exit(1)
	}

	var config Config
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		fmt.Printf("error parsing config file: %v\n", err)
		os.Exit(1)
	}

	selectedProgram, err := chooseProgram(config.ChoiceCommand, config.Programs)
	if err != nil {
		fmt.Printf("error choosing program: %v\n", err)
		os.Exit(1)
	}
	selectedProgram.Command = strings.ReplaceAll(selectedProgram.Command, "~/", "$HOME/")
	selectedProgram.Command = os.ExpandEnv(selectedProgram.Command)

	runSelectedProgram(selectedProgram)
}
