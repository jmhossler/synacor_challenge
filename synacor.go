package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var function_map = map[string]func(*VM, uint16) (uint16, error){
	"halt": halt,
	"set":  set,
	"push": push,
	"pop":  pop,
	"eq":   eq,
	"gt":   gt,
	"jmp":  jmp,
	"jt":   jt,
	"jf":   jf,
	"add":  add,
	"mult": mult,
	"mod":  mod,
	"and":  and,
	"or":   or,
	"not":  not,
	"rmem": rmem,
	"wmem": wmem,
	"call": call,
	"ret":  ret,
	"out":  out,
	"in":   in,
	"noop": noop,
}

var op_code_map = map[int]string{
	0:  "halt",
	1:  "set",
	2:  "push",
	3:  "pop",
	4:  "eq",
	5:  "gt",
	6:  "jmp",
	7:  "jt",
	8:  "jf",
	9:  "add",
	10: "mult",
	11: "mod",
	12: "and",
	13: "or",
	14: "not",
	15: "rmem",
	16: "wmem",
	17: "call",
	18: "ret",
	19: "out",
	20: "in",
	21: "noop",
}

type Stack []uint16

type VM struct {
	memory   [32768]uint16
	register [8]uint16
	stack    Stack
	output   []byte
	input    *os.File
}

func main() {
	if len(os.Args[1:]) < 1 {
		panic(errors.New("No arguments given"))
	}

	input_file := process_args(os.Args[1:])

	var byte_array []byte
	var vm VM
	byte_array, err := ioutil.ReadFile(input_file)
	check(err)

	for i := 0; i < len(byte_array); i += 2 {
		vm.memory[int(i/2)] = read_uint16(byte_array[i : i+2])
	}

	tmpl := "index.html.template"

	normal_handler := func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(tmpl)
		check(err)
		t.Execute(w, vm)
	}

	var curr_address uint16 = 0
	step_handler := func(w http.ResponseWriter, r *http.Request) {
		var err error
		curr_address, err = step(&vm, curr_address)
		check(err)

		t, err := template.ParseFiles(tmpl)
		check(err)
		t.Execute(w, vm)
	}

	execute_handler := func(w http.ResponseWriter, r *http.Request) {
		go execute(&vm, curr_address)
		t, err := template.ParseFiles(tmpl)
		check(err)

		t.Execute(w, vm)
	}

	http.HandleFunc("/step", step_handler)
	http.HandleFunc("/", normal_handler)
	http.HandleFunc("/execute", execute_handler)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	//output_memory(os.Stdout, vm)
	//execute(&vm, 0)
	//fmt.Printf("%s\n", vm.output)
}

func execute(vm *VM, address uint16) {
	var err error
	for {
		address, err = step(vm, address)
		if err != nil {
			return
		}
	}
}

func step(vm *VM, address uint16) (uint16, error) {
	op_code := get_val(vm, address)
	f := function_map[op_code_map[int(op_code)]]
	return f(vm, address)
}

func output_memory(f *os.File, vm VM) {
	for index, value := range vm.memory {
		fmt.Fprintf(f, "Address 0x%04x: 0x%04x\n", index, value)
	}
}

func read_uint16(data []byte) (ret uint16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func process_args(array []string) string {
	return array[0]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (s Stack) Push(v uint16) Stack {
	return append(s, v)
}

func (s Stack) Pop() (Stack, uint16) {
	if len(s) == 0 {
		return s, 0
	} else {
		return s[:len(s)-1], s[len(s)-1]
	}
}

func halt(vm *VM, address uint16) (uint16, error) {
	return address, errors.New(fmt.Sprintf("Halted program at address 0x%04x", address))
}

func set(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("set not implemented yet")
}

func push(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("push not implemented yet")
}

func pop(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("pop not implemented yet")
}

func eq(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("eq not implemented yet")
}

func gt(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("gt not implemented yet")
}

func jmp(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("jmp not implemented yet")
}

func jt(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("jt not implemented yet")
}

func jf(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("jf not implemented yet")
}

func add(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("add not implemented yet")
}

func mult(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("mult not implemented yet")
}

func mod(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("mod not implemented yet")
}

func and(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("and not implemented yet")
}

func or(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("or not implemented yet")
}

func not(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("not not implemented yet")
}

func rmem(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("rmem not implemented yet")
}

func wmem(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("wmem not implemented yet")
}

func call(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("call not implemented yet")
}

func ret(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("ret not implemented yet")
}

func out(vm *VM, address uint16) (uint16, error) {
	vm.output = append(vm.output, byte(get_val(vm, address+1)))
	return address + 2, nil
}

func in(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("in not implemented yet")
}

func noop(vm *VM, address uint16) (uint16, error) {
	return address + 1, nil
}

func get_val(vm *VM, address uint16) uint16 {
	if int(address) < len(vm.memory) {
		return vm.memory[int(address)]
	} else if int(address)%len(vm.memory) < len(vm.register) {
		return vm.register[int(address)%len(vm.memory)]
	} else {
		return 0
	}
}
