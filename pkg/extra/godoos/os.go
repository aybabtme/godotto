package godoos

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/robertkrimen/otto"
)

var q = otto.Value{}

func Apply(vm *otto.Otto) (v otto.Value, err error) {
	root, err := vm.Object(`({})`)
	if err != nil {
		return q, err
	}

	for _, applier := range []struct {
		Name   string
		Method func(otto.FunctionCall) otto.Value
	}{
		{"read_file", readFile},
		{"write_file", writeFile},
		{"sleep", sleep},
		{"args", args},
	} {
		if err := root.Set(applier.Name, applier.Method); err != nil {
			return q, fmt.Errorf("preparing method %q, %v", applier.Name, err)
		}
	}

	return root.Value(), nil
}

func args(all otto.FunctionCall) otto.Value {
	return ottoutil.ToValue(all.Otto, os.Args)
}

func readFile(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	filename := ottoutil.String(vm, all.Argument(0))
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return ottoutil.ToValue(vm, string(data))
}

func writeFile(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	filename := ottoutil.String(vm, all.Argument(0))
	content := ottoutil.String(vm, all.Argument(1))
	perm := ottoutil.String(vm, all.Argument(2))

	permuint, err := strconv.ParseUint(perm, 8, 32)
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}

	err = ioutil.WriteFile(filename, []byte(content), os.FileMode(permuint))
	if err != nil {
		ottoutil.Throw(vm, err.Error())
	}
	return q
}

func sleep(all otto.FunctionCall) otto.Value {
	vm := all.Otto
	dur := ottoutil.Duration(vm, all.Argument(0))
	time.Sleep(dur)
	return q
}
