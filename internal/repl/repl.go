// Copyright (c) 2012 Robert Krimen

// Modified from https://github.com/robertkrimen/otto

// Package repl implements a REPL (read-eval-print loop) for otto.
package repl

import (
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/robertkrimen/otto"
	"gopkg.in/readline.v1"
)

// Run runs a REPL with the given prompt and prelude.
func Run(vm *otto.Otto, prompt, prelude string) error {
	if prompt == "" {
		prompt = ">"
	}

	prompt = strings.Trim(prompt, " ")
	prompt += " "

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       prompt,
		AutoComplete: &autoCompleter{vm},
	})
	if err != nil {
		return err
	}

	if prelude != "" {
		if _, err := io.Copy(rl.Stderr(), strings.NewReader(prelude+"\n")); err != nil {
			return err
		}

		rl.Refresh()
	}

	var d []string

	for {
		l, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				if d != nil {
					d = nil

					rl.SetPrompt(prompt)
					rl.Refresh()

					continue
				}

				break
			}

			return err
		}

		if l == "" {
			continue
		}

		d = append(d, l)

		s, err := vm.Compile("repl", strings.Join(d, "\n"))
		if err != nil {
			rl.SetPrompt(strings.Repeat(" ", len(prompt)))
		} else {
			rl.SetPrompt(prompt)

			d = nil

			v, err := vm.Eval(s)
			if err != nil {
				if oerr, ok := err.(*otto.Error); ok {
					io.Copy(rl.Stdout(), strings.NewReader(oerr.String()))
				} else {
					io.Copy(rl.Stdout(), strings.NewReader(err.Error()))
				}
			} else {

				if !v.IsDefined() {
				} else {

					gov, err := toGo(v)
					if err != nil {
						io.Copy(rl.Stdout(), strings.NewReader(err.Error()))
					} else {
						data, _ := json.MarshalIndent(gov, "", "  ")
						log.Printf("%#v", gov)
						rl.Stdout().Write(append(data, "\n"...))
					}
				}

			}
		}

		rl.Refresh()
	}

	return rl.Close()
}

func toGo(v otto.Value) (interface{}, error) {
	gov, err := v.Export()
	if err != nil {
		return nil, err
	}

	var recurseType func(gov interface{}) (interface{}, error)
	recurseType = func(gov interface{}) (interface{}, error) {
		switch tgov := gov.(type) {

		case map[string]interface{}:
			out := make(map[string]interface{}, len(tgov))
			for k, v := range tgov {
				outel, err := recurseType(v)
				if err != nil {
					return nil, err
				}
				out[k] = outel
			}
			return out, nil

		case []otto.Value:
			var out []interface{}
			for _, el := range tgov {
				outel, err := toGo(el)
				if err != nil {
					return nil, err
				}
				out = append(out, outel)
			}
			return out, nil

		default:
			return tgov, nil
		}
	}
	return recurseType(gov)
}
