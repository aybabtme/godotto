package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/robertkrimen/otto"
)

func TestApply(t *testing.T) {

	host, user, port, auth, done := server(t)
	defer done()

	src := fmt.Sprintf(`
assert(ssh != null, "package should be loaded");
assert(ssh.session != null, "session function should be defined");

var host = %[1]q;
var user = %[2]q;
var port = %[3]q;

var droplet = {"name":"derp", "id":456789, "public_ipv4": host};

var session = ssh.session(droplet, {"user":user, "port": port});
assert(session.exec != null, "session should have exec method");
assert(session.close != null, "session should have close method");

try {
	assert(session.exec("1") == 'you sent "1"', "should have read response from server");
	assert(session.exec("22") == 'you sent "22"', "should have read response from server");
	assert(session.exec("333") == 'you sent "333"', "should have read response from server");
	assert(session.exec("echo hello") == 'you sent "echo hello"', "should have read response from server");
} finally {
	session.close();
}

`, host, user, port)

	vmtest.Run(t, src, func(vm *otto.Otto) error {
		pkg, err := Apply(vm, auth)
		if err != nil {
			return err
		}
		return vm.Set("ssh", pkg)
	})
}

func server(t testing.TB) (host, user, port string, auth ssh.AuthMethod, close func() error) {
	user = "testuser"
	password := "tiger"
	auth = ssh.Password(password)

	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			if c.User() == user && string(pass) == password {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("Failed to generate private key")
	}

	private, err := ssh.NewSignerFromKey(k)
	if err != nil {
		panic("Failed to parse private key")
	}

	config.AddHostKey(private)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("failed to listen for connection")
	}

	go func() {
		for {
			nConn, err := listener.Accept()
			if err != nil {
				return
			}

			_, chans, reqs, err := ssh.NewServerConn(nConn, config)
			if err != nil {
				panic("failed to handshake")
			}

			go ssh.DiscardRequests(reqs)

			for newChannel := range chans {

				if newChannel.ChannelType() != "session" {
					newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
					continue
				}
				channel, requests, err := newChannel.Accept()
				if err != nil {
					panic("could not accept channel.")
				}

				req := <-requests

				var (
					line string
					ok   bool
				)
				switch req.Type {
				case "exec":
					if len(req.Payload) > 4 {
						line = string(req.Payload[4:])
						ok = true
					}
				}
				req.Reply(ok, nil)
				if ok {
					fmt.Fprintf(channel, "you sent %q", line)
					channel.CloseWrite()
					channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
				}
				channel.Close()
			}
		}
	}()

	host, port, _ = net.SplitHostPort(listener.Addr().String())
	return host, user, port, auth, listener.Close
}
