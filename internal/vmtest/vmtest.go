package vmtest

import (
	"flag"
	"os"
	"testing"

	"golang.org/x/oauth2"

	"github.com/aybabtme/godotto"
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

var apiToken = flag.String("api.token", defaultToken, "token to use to communicate with the DO API")

// Run the JS source against godotto.
func Run(t testing.TB, src string) {
	client := godo.NewClient(oauth2.NewClient(oauth2.NoContext,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *apiToken}),
	))
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

	if _, err := vm.Run(script); err != nil {
		oe := err.(*otto.Error)
		t.Fatalf(oe.String())
	}
}
