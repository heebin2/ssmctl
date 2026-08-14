package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/heebin2/ssmctl/internal/ssm"
)

const zshCompletionScript = `#compdef ssmctl

if (( ! $+functions[compdef] )); then
  autoload -Uz compinit
  compinit -i -D
fi

_ssmctl() {
  local -a commands positional instances
  local config_path=""
  local i=2

  while (( i < CURRENT )); do
    if [[ "${words[i]}" == "-config" ]]; then
      (( i++ ))
      if (( i < CURRENT )); then
        config_path="${words[i]}"
      fi
    else
      positional+=("${words[i]}")
    fi
    (( i++ ))
  done

  if (( CURRENT > 2 )) && [[ "${words[CURRENT-1]}" == "-config" ]]; then
    _files
    return
  fi

  if (( ${#positional} == 0 )); then
    commands=(
      'init:initialize config'
      'list:list configured instances'
      'completion:generate shell completion'
    )
    _describe 'ssmctl command' commands

    if [[ -n "$config_path" ]]; then
      instances=("${(@f)$(command ssmctl -config "$config_path" completion __instances 2>/dev/null)}")
    else
      instances=("${(@f)$(command ssmctl completion __instances 2>/dev/null)}")
    fi
    (( ${#instances} )) && compadd -M 'l:|=* r:|=*' -a instances
    return
  fi

  if [[ "${positional[1]}" == "completion" ]] && (( ${#positional} == 1 )); then
    compadd zsh
  fi
}

compdef _ssmctl ssmctl
`

func writeCompletion(args []string, configPath string, out io.Writer) error {
	if len(args) != 1 {
		return errors.New("usage: ssmctl completion <zsh>")
	}

	switch args[0] {
	case "zsh":
		_, err := io.WriteString(out, zshCompletionScript)
		return err
	case "__instances":
		cfg, err := ssm.LoadConfig(configPath)
		if err != nil {
			return nil
		}
		for _, name := range cfg.InstanceNames() {
			if _, err := fmt.Fprintln(out, name); err != nil {
				return err
			}
		}
		return nil
	default:
		return errors.New("usage: ssmctl completion <zsh>")
	}
}
