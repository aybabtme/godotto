package ssh

import (
	"fmt"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/aybabtme/godotto/internal/vmtest"
	"github.com/robertkrimen/otto"
)

func TestApply(t *testing.T) {

	host, user, port, auth, done := server(t)
	defer done()

	src := fmt.Sprintf(`
assert(ssh != null, "package should be loaded");
assert(ssh.shell != null, "shell function should be defined");

var host = %[1]q;
var user = %[1]q;
var port = %[1]q;

var shell = ssh.shell(host, {"user":user, "port": port});
var out = shell("echo hello");

assert(out == 'you sent "echo hello"', "should have read response from server");

`, host, user, port)

	vmtest.Run(t, `

    `, func(vm *otto.Otto) error {
		pkg, err := Apply(vm, auth)
		if err != nil {
			return err
		}
		return vm.Set("ssh", pkg)
	})
}

func server(t testing.TB) (host, user, port string, auth ssh.AuthMethod, close func()) {
	user := "testuser"
	password := "tiger"
	auth = ssh.Password(password)

	// An SSH server is represented by a ServerConfig, which holds
	// certificate details and handles authentication of ServerConns.
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			// Should use constant-time compare (or better, salt+hash) in
			// a production setting.
			if c.User() == user && string(pass) == password {
				return nil, nil
			}
			return nil, fmt.Errorf("password rejected for %q", c.User())
		},
	}

	// Once a ServerConfig has been configured, connections can be
	// accepted.
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("failed to listen for connection")
	}

	go func() {
		nConn, err := listener.Accept()
		if err != nil {
			t.Fatal("failed to accept incoming connection")
		}

		// Before use, a handshake must be performed on the incoming
		// net.Conn.
		_, chans, reqs, err := ssh.NewServerConn(nConn, config)
		if err != nil {
			t.Fatal("failed to handshake")
		}
		// The incoming Request channel must be serviced.
		go ssh.DiscardRequests(reqs)

		// Service the incoming Channel channel.
		for newChannel := range chans {
			// Channels have a type, depending on the application level
			// protocol intended. In the case of a shell, the type is
			// "session" and ServerShell may be used to present a simple
			// terminal interface.
			if newChannel.ChannelType() != "session" {
				newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
				continue
			}
			channel, requests, err := newChannel.Accept()
			if err != nil {
				t.Fatal("could not accept channel.")
			}

			// Sessions have out-of-band requests such as "shell",
			// "pty-req" and "env".  Here we handle only the
			// "shell" request.
			go func(in <-chan *ssh.Request) {
				for req := range in {
					ok := false
					switch req.Type {
					case "shell":
						ok = true
						if len(req.Payload) > 0 {
							// We don't accept any
							// commands, only the
							// default shell.
							ok = false
						}
					}
					req.Reply(ok, nil)
				}
			}(requests)

			term := terminal.NewTerminal(channel, "> ")

			go func() {
				defer channel.Close()
				for {
					line, err := term.ReadLine()
					if err != nil {
						break
					}
					fmt.Println(line)
					fmt.Fsprintf(channel, "you send %q", line)
				}
			}()
		}
	}()

	host, port, _ = net.SplitHostPort(listener.Addr().String())
	return host, user, port, auth, listener.Close
}
