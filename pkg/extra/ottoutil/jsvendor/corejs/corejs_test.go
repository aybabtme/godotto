package corejs

import (
	"testing"
	"time"

	"github.com/aybabtme/godotto/pkg/extra/vmtest"
	"github.com/robertkrimen/otto"
)

func TestLoad(t *testing.T) {
	opt := func(vm *otto.Otto) error {
		start := time.Now()
		err := Load(vm)
		t.Logf("loaded in %v", time.Since(start))
		return err
	}
	vmtest.Run(t, nil, `
var msg = "hello world";
assert(msg.startsWith != null, "startsWith should be defined on strings");
assert(msg.startsWith("hello"), "message should start with 'hello'");
    `, opt)
}
