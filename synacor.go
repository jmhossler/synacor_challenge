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
	Memory   [32768]uint16
	Register [8]uint16
	Stack    Stack
	Output   bytes.Buffer
	Input    *os.File
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
	var curr_address uint16 = 0

	tmpl := "index.html.template"
	memory_tmpl := "memory.html"
	mem_vals_tmpl := "mem_vals.html"

	normal_handler := func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(tmpl)
		t.Execute(w, vm)
	}

	step_handler := func(w http.ResponseWriter, r *http.Request) {
		var err error
		if r.Method == "POST" {
			curr_address, err = step(&vm, curr_address)
			if err != nil {
				vm.output += err.Error() + "\n"
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	execute_handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			go execute(&vm, &curr_address)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	reset_handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			curr_address = 0
			vm.output += "\n--------\n"
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	http.HandleFunc("/", normal_handler)
	http.HandleFunc("/execute", execute_handler)
	http.HandleFunc("/step", step_handler)
	http.HandleFunc("/reset", reset_handler)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	//output_memory(os.Stdout, vm)
	//execute(&vm, 0)
	//fmt.Printf("%s\n", vm.output)
}

func format_mem(vm *VM) bytes.Buffer {
	var str bytes.Buffer
	for i := 0; i < len(vm.memory); i++ {
		str += fmt.Sprintf("%05d: %05d\n", uint16(i), vm.memory[i])
	}
	return str
}

func execute(vm *VM, address *uint16) {
	var err error
	for {
		*address, err = step(vm, *address)
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

func (s Stack) isEmpty() bool {
	return len(s) == 0
}

func (s Stack) Pop() (Stack, uint16) {
	if s.isEmpty() {
		return s, 0
	} else {
		return s[:len(s)-1], s[len(s)-1]
	}
}

func halt(vm *VM, address uint16) (uint16, error) {
	return address, errors.New(fmt.Sprintf("Halted program at address 0x%04x", address))
}

func set(vm *VM, address uint16) (uint16, error) {
	vm.register[get_val(vm, address+1)] = get_val(vm, address+2)
	return address + 3, nil
}

func push(vm *VM, address uint16) (uint16, error) {
	vm.stack.Push(get_val(vm, address+1))
	return address + 2, nil
}

func pop(vm *VM, address uint16) (uint16, error) {
	var value uint16
	vm.stack, value = vm.stack.Pop()
	set_val(vm, get_val(vm, address+1), value)
	return address + 2, nil
}

func eq(vm *VM, address uint16) (uint16, error) {
	a := get_val(vm, address+1)
	b := get_val(vm, address+2)
	c := get_val(vm, address+3)
	if b == c {
		set_val(vm, a, 1)
	} else {
		set_val(vm, a, 0)
	}
	return address + 4, nil
}

func gt(vm *VM, address uint16) (uint16, error) {
	a := get_val(vm, address+1)
	b := get_val(vm, address+2)
	c := get_val(vm, address+3)
	if b > c {
		set_val(vm, a, 1)
	} else {
		set_val(vm, a, 0)
	}
	return address + 4, nil
}

func jmp(vm *VM, address uint16) (uint16, error) {
	return get_val(vm, address+1), nil
}

func jt(vm *VM, address uint16) (uint16, error) {
	if get_val(vm, address+1) != 0 {
		return get_val(vm, address+2), nil
	} else {
		return address + 3, nil
	}
}

func jf(vm *VM, address uint16) (uint16, error) {
	if get_val(vm, address+1) == 0 {
		return get_val(vm, address+2), nil
	} else {
		return address + 3, nil
	}
}

func add(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), (get_val(vm, address+2)+get_val(vm, address+3))&0x7FFF)
	return address + 4, nil
}

func mult(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), (get_val(vm, address+2)*get_val(vm, address+3))&0x7FFF)
	return address + 4, nil
}

func mod(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), uint16(int(get_val(vm, address+2))%int(get_val(vm, address+3))))
	return address + 4, nil
}

func and(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), get_val(vm, address+2)&get_val(vm, address+3))
	return address + 4, nil
}

func or(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), get_val(vm, address+2)|get_val(vm, address+3))
	return address + 4, nil
}

func not(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), ^(get_val(vm, address+2))&0x7FFF)
	return address + 3, nil
}

func rmem(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), get_val(vm, get_val(vm, address+2)))
	return address + 3, nil
}

func wmem(vm *VM, address uint16) (uint16, error) {
	set_val(vm, get_val(vm, address+1), get_val(vm, address+2))
	return address + 3, nil
}

func call(vm *VM, address uint16) (uint16, error) {
	vm.stack.Push(address + 2)
	return get_val(vm, address+1), nil
}

func ret(vm *VM, address uint16) (uint16, error) {
	if vm.stack.isEmpty() {
		return halt(vm, address)
	} else {
		var value uint16
		vm.stack, value = vm.stack.Pop()
		return value, nil
	}
}

func out(vm *VM, address uint16) (uint16, error) {
	vm.output += string(rune(get_val(vm, address+1)))
	return address + 2, nil
}

func in(vm *VM, address uint16) (uint16, error) {
	return address, errors.New("in not implemented yet")
}

func noop(vm *VM, address uint16) (uint16, error) {
	return address + 1, nil
}

func get_val(vm *VM, address uint16) uint16 {
	var val uint16
	if int(address) < len(vm.memory) {
		val = vm.memory[int(address)]
	} else if int(address)%len(vm.memory) < len(vm.register) {
		val = vm.register[int(address)%len(vm.memory)]
	} else {
		val = 0
	}
	fmt.Printf("Address %d: %d\n", int(address), int(val))
	if val > 0x7FFF {
		return vm.register[val-0x8000]
	} else {
		return val
	}
}

func set_val(vm *VM, address uint16, value uint16) {
	if int(address) < len(vm.memory) {
		vm.memory[address] = value
	} else if int(address)%len(vm.memory) < len(vm.register) {
		vm.register[int(address)%len(vm.memory)] = value
	}
}
