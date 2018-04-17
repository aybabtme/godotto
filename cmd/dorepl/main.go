package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"github.com/aybabtme/godotto"
	"github.com/aybabtme/godotto/pkg/extra/ottoutil/jsvendor/corejs"
	"github.com/aybabtme/godotto/pkg/extra/repl"
	"github.com/aybabtme/godotto/pkg/extra/do/cloud"
	"github.com/aybabtme/godotto/pkg/extra/do/spycloud"
	"github.com/aybabtme/godotto/pkg/extra/godoos"
	jsssh "github.com/aybabtme/godotto/pkg/extra/ssh"

	"github.com/digitalocean/godo"
	"github.com/robertkrimen/otto"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
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

func main() {
	apiToken := flag.String("api.token", defaultToken, "token to use to communicate with the DO API")
	flag.Parse()

	log.SetFlags(0)
	log.SetPrefix("dorepl: ")

	if *apiToken == "" {
		flag.PrintDefaults()
		log.Fatalf("At this time, the REPL requires you to provide an API token")
	}

	gc := godo.NewClient(oauth2.NewClient(oauth2.NoContext,
		oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *apiToken}),
	))
	acc, _, err := gc.Account.Get(context.TODO())
	if err != nil {
		log.Fatalf("can't query DigitalOcean account, is your token valid?\n%v", err)
	}

	vm := otto.New()
	if err := corejs.Load(vm); err != nil {
		log.Fatal(err)
	}

	cloud, spy := spycloud.Client(cloud.New(cloud.UseGodo(gc)))
	defer enumerateLeftover(spy)

	ctx := context.Background()
	pkg, err := godotto.Apply(ctx, vm, cloud)
	if err != nil {
		log.Fatal(err)
	}
	vm.Set("cloud", pkg)

	ospkg, err := godoos.Apply(vm)
	if err != nil {
		log.Fatal(err)
	}
	vm.Set("os", ospkg)

	auth, done := sshAgent()
	defer done()
	if s, cleanup, err := jsssh.Apply(ctx, vm, auth); err != nil {
		log.Fatal(err)
	} else {
		defer cleanup()
		vm.Set("ssh", s)
	}

	if len(flag.Args()) == 0 {
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

func sshAgent() (ssh.AuthMethod, func()) {
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil, func() {}
	}
	return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers), func() {
		_ = sshAgent.Close()
	}
}

func enumerateLeftover(spy func(...spycloud.Spy)) {
	var once sync.Once
	print := func() {
		log.Print("quitting! the following resources were created")
	}
	spy(
		spycloud.Droplets(func(v *godo.Droplet) {
			once.Do(print)
			log.Printf("- Droplet: %d", v.ID)
		}),
		spycloud.Volumes(func(v *godo.Volume) {
			once.Do(print)
			log.Printf("- Volume: %q", v.ID)
		}),
		spycloud.Snapshots(func(v *godo.Snapshot) {
			once.Do(print)
			log.Printf("- Snapshot: %q", v.ID)
		}),
		spycloud.Domains(func(v *godo.Domain) {
			once.Do(print)
			log.Printf("- Domain: %q", v.Name)
		}),
		spycloud.Records(func(v *godo.DomainRecord) {
			once.Do(print)
			log.Printf("- DomainRecord: %q", v.ID)
		}),
		spycloud.FloatingIPs(func(v *godo.FloatingIP) {
			once.Do(print)
			log.Printf("- FloatingIP: %q", v.IP)
		}),
		spycloud.Keys(func(v *godo.Key) {
			once.Do(print)
			log.Printf("- Key: %q", v.ID)
		}),
		spycloud.Tags(func(v *godo.Tag) {
			once.Do(print)
			log.Printf("- Tag: %q", v.Name)
		}),
		spycloud.LoadBalancers(func(v *godo.LoadBalancer) {
			once.Do(print)
			log.Printf("- Load Balancer: %q", v.Name)
		}),
		spycloud.Snapshots(func(v *godo.Snapshot) {
			once.Do(print)
			log.Printf("- Snapshot: %q", v.ID)
		}),
		spycloud.Firewalls(func(v *godo.Firewall) {
			once.Do(print)
			log.Printf("- Firewall: %q", v.Name)
		}),
	)
}
