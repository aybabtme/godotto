package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/aybabtme/godotto"
	"github.com/aybabtme/godotto/internal/ottoutil/jsvendor/corejs"
	"github.com/aybabtme/godotto/internal/repl"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/oauth2"

	_ "github.com/robertkrimen/otto/underscore"
)

var prelude = `
Welcome to the DigitalOcean REPL, where all your dreams come true!
`

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

func main() {
	log.SetFlags(0)
	log.SetPrefix("dorepl: ")

	if *apiToken == "" {
		flag.PrintDefaults()
		log.Fatalf("At this time, the REPL requires you to provide an API token")
	}

	gc := godo.NewClient(oauth2.NewClient(oauth2.NoContext,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *apiToken}),
	))
	acc, _, err := gc.Account.Get()
	if err != nil {
		log.Fatalf("can't query DigitalOcean account, is your token valid?\n%v", err)
	}

	vm := otto.New()
	if err := corejs.Load(vm); err != nil {
		log.Fatal(err)
	}
	pkg, err := godotto.Apply(vm, gc)
	if err != nil {
		log.Fatal(err)
	}
	vm.Set("cloud", pkg)

	if len(os.Args[1:]) == 0 {
		// run REPL
		if !terminal.IsTerminal(0) {
			prelude = ""
		} else {
			log.Printf("logged in as %s", acc.Email)
		}

		if err := repl.Run(vm, ">", prelude); err != nil && err != io.EOF {
			log.Fatal(err)
		}
	} else {

		// run scripts

		enc := json.NewEncoder(os.Stdout)
		for _, filename := range os.Args[1:] {
			raw, err := ioutil.ReadFile(filename)
			if err != nil {
				log.Fatal(err)
			}

			script := string(raw[bytes.IndexRune(raw, '\n'):])

			v, err := vm.Run(script)
			if err != nil {
				log.Fatal(err)
			}
			gov, err := v.Export()
			if err != nil {
				log.Fatal(err)
			}
			if v.IsDefined() {
				if err := enc.Encode(gov); err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}
