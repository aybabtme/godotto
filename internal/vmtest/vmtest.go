package vmtest

import (
	"flag"
	"os"
	"testing"

	"golang.org/x/oauth2"

	"github.com/aybabtme/godotto"
	"github.com/aybabtme/godotto/internal/do"
	"github.com/aybabtme/godotto/internal/do/cloud"
	"github.com/aybabtme/godotto/internal/ottoutil"
	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
)

var defaultToken = func() string {
	for _, env := range []string{
		"DIGITALOCEAN_ACCESS_TOKEN",
		"DIGITALOCEAN_TOKEN",
		"DIGITAL_OCEAN_TOKEN",
		"DIGITAL_OCEAN_ACCESS_TOKEN",
		"DO_TOKEN",
	} {
		if s := os.Getenv(env); s != "" {
			return s
		}
	}
	return ""
}()

var (
	apiToken = flag.String("api.token", defaultToken, "token to use to communicate with the DO API")
)

// A RunOption is applied on the otto VM before the test begins.
type RunOption func(vm *otto.Otto) error

// Run the JS source against godotto.
func Run(t testing.TB, src string, opts ...RunOption) {

	u, done := do.Stub()
	defer done()

	var gc *godo.Client
	if *apiToken != "" {
		gc = godo.NewClient(oauth2.NewClient(oauth2.NoContext,
			oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *apiToken}),
		))
	} else {
		gc = godo.NewClient(nil)
		gc.BaseURL = u
	}

	client := cloud.New(cloud.UseGodo(gc))

	vm := otto.New()

	pkg, err := godotto.Apply(vm, client)
	if err != nil {
		t.Fatal(err)
	}
	vm.Set("cloud", pkg)
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
			msg += "\n" + call.CallerLocation() + " | " + format
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
		oe := err.(*otto.Error)
		t.Fatalf(oe.String())
	}
}
