package ottoutil

import "github.com/robertkrimen/otto"

func GetObject(vm *otto.Otto, obj *otto.Object, name string) otto.Value {
	v, err := obj.Get(name)
	if err != nil {
		Throw(vm, err.Error())
	}
	return v
}

func String(vm *otto.Otto, v otto.Value) string {
	if !v.IsDefined() {
		return ""
	}
	s, err := v.ToString()
	if err != nil {
		Throw(vm, err.Error())
	}
	return s
}

func Int(vm *otto.Otto, v otto.Value) int {
	i, err := v.ToInteger()
	if err != nil {
		Throw(vm, err.Error())
	}
	return int(i)
}

func Float64(vm *otto.Otto, v otto.Value) float64 {
	f, err := v.ToFloat()
	if err != nil {
		Throw(vm, err.Error())
	}
	return f
}

func Bool(vm *otto.Otto, v otto.Value) bool {
	b, err := v.ToBoolean()
	if err != nil {
		Throw(vm, err.Error())
	}
	return b
}

func Throw(vm *otto.Otto, str string) {
	value, _ := vm.Call("new Error", nil, str)
	panic(value)
}
