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
	"strings"
)

var functionMap = map[string]func(*VM, uint16) (uint16, error){
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

var opCodeMap = map[int]string{
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

// Stack simple stack for memory values
type Stack []uint16

// VM the 'vm'
type VM struct {
	Memory      [32768]uint16
	Register    [8]uint16
	Stack       Stack
	Output      bytes.Buffer
	MemoryTrace bytes.Buffer
	Input       *os.File
}

const (
	post = "POST"
	get  = "GET"
)

func main() {
	if len(os.Args[1:]) < 1 {
		panic(errors.New("No arguments given"))
	}

	inputFile := processArgs(os.Args[1:])

	var byteArray []byte
	var vm VM
	byteArray, err := ioutil.ReadFile(inputFile)
	check(err)

	for i := 0; i < len(byteArray); i += 2 {
		vm.Memory[int(i/2)] = readUint16(byteArray[i : i+2])
	}
	var currAddress uint16

	tmpl := "index.html.template"

	normalHandler := func(w http.ResponseWriter, r *http.Request) {
		report, err := template.New("report").Funcs(template.FuncMap{"replaceNewLines": replaceNewLines}).ParseFiles(tmpl)
		check(err)

		report.Execute(w, vm)
	}

	stepHandler := func(w http.ResponseWriter, r *http.Request) {
		var err error
		if r.Method == post {
			currAddress, err = step(&vm, currAddress)
			if err != nil {
				vm.Output.WriteString(err.Error() + "\n")
			}
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	executeHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == post {
			go execute(&vm, &currAddress)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	resetHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == post {
			currAddress = 0
			vm.Output.WriteString("\n--------\n")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	http.HandleFunc("/", normalHandler)
	http.HandleFunc("/execute", executeHandler)
	http.HandleFunc("/step", stepHandler)
	http.HandleFunc("/reset", resetHandler)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))

	//outputMemory(os.Stdout, vm)
	//execute(&vm, 0)
	//fmt.Printf("%s\n", vm.Output)
}

func replaceNewLines(s bytes.Buffer) template.HTML {
	str := s.String()
	return template.HTML(strings.Replace(str, "\n", "<br>", -1))
}

/*
func format_mem(vm *VM) string {
	var str string
>>>>>>> origin/master
	for i := 0; i < len(vm.Memory); i++ {
		str += fmt.Sprintf("%05d: %05d\n", uint16(i), vm.Memory[i])
	}
	return str
}
*/

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
	opCode := *getVal(vm, address)
	f := functionMap[opCodeMap[int(opCode)]]
	return f(vm, address)
}

func outputMemory(f *os.File, vm VM) {
	for index, value := range vm.Memory {
		fmt.Fprintf(f, "Address 0x%04x: 0x%04x\n", index, value)
	}
}

func readUint16(data []byte) (ret uint16) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func processArgs(array []string) string {
	return array[0]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Push pushes stuff on
func (s Stack) Push(v uint16) Stack {
	return append(s, v)
}

func (s Stack) isEmpty() bool {
	return len(s) == 0
}

// Pop pops stuff off
func (s Stack) Pop() (Stack, uint16) {
	if s.isEmpty() {
		return s, 0
	}
	return s[:len(s)-1], s[len(s)-1]
}

func halt(vm *VM, address uint16) (uint16, error) {
	return address, fmt.Errorf("Halted program at address 0x%04x", address)
}

func set(vm *VM, address uint16) (uint16, error) {
	ptr := getVal(vm, address+1)
	*ptr = *getVal(vm, address+2)
	return address + 3, nil
}

func push(vm *VM, address uint16) (uint16, error) {
	vm.Stack.Push(*getVal(vm, address+1))
	return address + 2, nil
}

func pop(vm *VM, address uint16) (uint16, error) {
	var value uint16
	vm.Stack, value = vm.Stack.Pop()
	*getVal(vm, address+1) = value
	return address + 2, nil
}

func eq(vm *VM, address uint16) (uint16, error) {
	a := getVal(vm, address+1)
	b := getVal(vm, address+2)
	c := getVal(vm, address+3)
	if *b == *c {
		*a = 1
	} else {
		*a = 0
	}
	return address + 4, nil
}

func gt(vm *VM, address uint16) (uint16, error) {
	a := getVal(vm, address+1)
	b := getVal(vm, address+2)
	c := getVal(vm, address+3)
	if *b > *c {
		*a = 1
	} else {
		*a = 0
	}
	return address + 4, nil
}

func jmp(vm *VM, address uint16) (uint16, error) {
	return *getVal(vm, address+1), nil
}

func jt(vm *VM, address uint16) (uint16, error) {
	if *getVal(vm, address+1) != 0 {
		return *getVal(vm, address+2), nil
	}
	return address + 3, nil
}

func jf(vm *VM, address uint16) (uint16, error) {
	if *getVal(vm, address+1) == 0 {
		return *getVal(vm, address+2), nil
	}
	return address + 3, nil
}

func add(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = (*getVal(vm, address+2) + *getVal(vm, address+3)) & 0x7FFF
	return address + 4, nil
}

func mult(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = (*getVal(vm, address+2) * *getVal(vm, address+3)) & 0x7FFF
	return address + 4, nil
}

func mod(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = uint16(int(*getVal(vm, address+2)) % int(*getVal(vm, address+3)))
	return address + 4, nil
}

func and(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = *getVal(vm, address+2) & (*getVal(vm, address+3))
	return address + 4, nil
}

func or(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = *getVal(vm, address+2) | *getVal(vm, address+3)
	return address + 4, nil
}

func not(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = ^(*getVal(vm, address+2)) & 0x7FFF
	return address + 3, nil
}

func rmem(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = *getVal(vm, *getVal(vm, address+2))
	return address + 3, nil
}

func wmem(vm *VM, address uint16) (uint16, error) {
	*getVal(vm, address+1) = *getVal(vm, address+2)
	return address + 3, nil
}

func call(vm *VM, address uint16) (uint16, error) {
	vm.Stack.Push(address + 2)
	return *getVal(vm, address+1), nil
}

func ret(vm *VM, address uint16) (uint16, error) {
	if vm.Stack.isEmpty() {
		return halt(vm, address)
	}
	var value uint16
	vm.Stack, value = vm.Stack.Pop()
	return value, nil

}

func out(vm *VM, address uint16) (uint16, error) {
	vm.Output.WriteRune(rune(*getVal(vm, address+1)))
	return address + 2, nil
}

func in(vm *VM, address uint16) (uint16, error) {
	return address, fmt.Errorf("in not implemented yet")
}

func noop(vm *VM, address uint16) (uint16, error) {
	return address + 1, nil
}

func getVal(vm *VM, address uint16) *uint16 {
	val := vm.Memory[address]

	vm.MemoryTrace.WriteString(fmt.Sprintf("Address %d: %d\n", int(address), int(val)))
	if val > 0x7FFF {
		return &vm.Register[val&0x7FFF]
	}
	return &val
}

func setVal(vm *VM, address uint16, value uint16) {
	if int(address) < len(vm.Memory) {
		vm.Memory[address] = value
	} else if int(address)%len(vm.Memory) < len(vm.Register) {
		vm.Register[int(address)%len(vm.Memory)] = value
	}
}
