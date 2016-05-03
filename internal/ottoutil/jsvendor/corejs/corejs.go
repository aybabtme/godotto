// Package corejs adds core-js functionalities to an otto VM.
package corejs

import "github.com/robertkrimen/otto"

// Load the core-js library in the vm.
func Load(vm *otto.Otto) error {
	_, err := vm.Run(source)
	return err
}
