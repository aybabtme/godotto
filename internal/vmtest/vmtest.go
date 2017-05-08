package vmtest

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/net/context"

	"github.com/aybabtme/godotto"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/mockcloud"
	"github.com/robertkrimen/otto"
)

// A RunOption is applied on the otto VM before the test begins.
type RunOption func(vm *otto.Otto) error

// Run the JS source against godotto.
func Run(t testing.TB, cloud cloud.Client, src string, opts ...RunOption) {

	if cloud == nil {
		cloud = mockcloud.Client(nil)
	}

	vm := otto.New()

	pkg, err := godotto.Apply(context.Background(), vm, cloud)
	if err != nil {
		t.Fatal(err)
	}
	vm.Set("cloud", pkg)
	vm.Set("equals", func(call otto.FunctionCall) otto.Value {
		vm := call.Otto
		got, err := call.Argument(0).Export()
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		want, err := call.Argument(1).Export()
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		ok, cause := deepEqual(got, want)
		if ok {
			return otto.UndefinedValue()
		}
		msg := "assertion failed!\n" + cause

		if len(call.ArgumentList) > 2 {
			format, err := call.ArgumentList[2].ToString()
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			msg += "\n" + format
		}
		ottoutil.Throw(vm, msg)
		return otto.UndefinedValue()
	})
	vm.Set("assert", func(call otto.FunctionCall) otto.Value {
		vm := call.Otto
		v, err := call.Argument(0).ToBoolean()
		if err != nil {
			ottoutil.Throw(vm, err.Error())
		}
		if v {
			return otto.UndefinedValue()
		}
		msg := "assertion failed!"
		if len(call.ArgumentList) > 1 {
			format, err := call.ArgumentList[1].ToString()
			if err != nil {
				ottoutil.Throw(vm, err.Error())
			}
			msg += "\n" + format
		}
		ottoutil.Throw(vm, msg)
		return otto.UndefinedValue()
	})
	script, err := vm.Compile("", src)
	if err != nil {
		t.Fatalf("invalid code: %v", err)
	}

	for _, opt := range opts {
		if err := opt(vm); err != nil {
			t.Fatalf("can't apply option: %v", err)
		}
	}

	if _, err := vm.Run(script); err != nil {
		if oe, ok := err.(*otto.Error); ok {
			t.Fatal(oe.String())
		} else {
			t.Fatal(err)
		}
	}
}

func deepEqual(lhs, rhs interface{}) (bool, string) {
	switch lhsv := lhs.(type) {
	case map[string]interface{}:
		rhsv, ok := rhs.(map[string]interface{})
		if !ok {
			return false, fmt.Sprintf("type %T != %T", lhs, rhs)
		}

		for k, v := range lhsv {
			if ok, cause := deepEqual(v, rhsv[k]); !ok {
				return false, fmt.Sprintf("lhs-key %q: %s", k, cause)
			}
		}
		for k, v := range rhsv {
			if ok, cause := deepEqual(v, lhsv[k]); !ok {
				return false, fmt.Sprintf("rhs-key %q: %s", k, cause)
			}
		}
		return true, ""
	case []map[string]interface{}:
		rhsv, ok := rhs.([]map[string]interface{})
		if !ok {
			return false, fmt.Sprintf("type %T != %T", lhs, rhs)
		}

		if len(lhsv) != len(rhsv) {
			return false, fmt.Sprintf("len(lhs) = %v, len(rhs) = %v", len(lhsv), len(rhsv))
		}

		for k, v := range lhsv {
			if ok, cause := deepEqual(v, rhsv[k]); !ok {
				return false, fmt.Sprintf("lhs-key %v: %s", k, cause)
			}
		}
		for k, v := range rhsv {
			if ok, cause := deepEqual(v, lhsv[k]); !ok {
				return false, fmt.Sprintf("rhs-key %v: %s", k, cause)
			}
		}
		return true, ""
	default:
		return reflect.DeepEqual(lhs, rhs), fmt.Sprintf("lhs=%#v (%T)\nrhs=%#v  (%T)", lhs, lhs, rhs, rhs)
	}
}
