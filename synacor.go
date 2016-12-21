package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
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
	Output      *bufio.Writer
	MemoryTrace *bufio.Writer
	Input       *bufio.Reader
}

type state struct {
	x, y   int
	path   string
	weight int
	lastOp string
}

const (
	post = "POST"
	get  = "GET"
)

var finalMap = [4][4]string{
	[4]string{"22", "-", "9", "*"},
	[4]string{"+", "4", "-", "18"},
	[4]string{"4", "*", "11", "*"},
	[4]string{"*", "8", "-", "1"},
}

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

	minPath := getMinPath()
	fmt.Println(minPath)

	vm.Output = bufio.NewWriter(os.Stdout)
	vm.Input = bufio.NewReader(os.Stdin)

	/* Commented out because extracted algorithm
	memoryOut, err := os.Create("MemoryTrace.txt")
	check(err)
	vm.MemoryTrace = bufio.NewWriter(memoryOut)
	*/
	execute(&vm, &currAddress)
	//fmt.Printf("%s\n", vm.Output)
}

func newAck(i uint16) func(uint16, uint16) uint16 {
	var memo = make(map[uint16]map[uint16]uint16)
	var f func(uint16, uint16) uint16
	f = func(x, y uint16) uint16 {
		if v, ok := memo[x]; ok {
			if n_v, ok := v[y]; ok {
				return n_v
			}
		} else {
			memo[x] = make(map[uint16]uint16)
		}
		if x == 0 {
			val := (y + 1) & 0x7FFF
			memo[x][y] = val
			return val
		} else if y == 0 {
			val := f(x-1, i)
			memo[x][y] = val
			return val
		} else {
			val := f(x-1, f(x, y-1))
			memo[x][y] = val
			return val
		}
	}
	return f
}

func findReg() uint16 {
	for i := uint16(0x7FFF); i > 0; i-- {
		if newAck(i)(4, 1) == 6 {
			return i
		}
	}
	return 0
}

func applyOp(op string, a, b int) int {
	switch op {
	case "*":
		return a * b
	case "-":
		return a - b
	case "+":
		return a + b
	}
	return 0
}

func getMinPath() string {
	visited := []state{state{weight: 22}}

	for len(visited) > 0 {
		curr := visited[0]
		visited = visited[1:]

		if curr.x == 3 && curr.y == 3 && curr.weight == 30 {
			return curr.path
		} else if curr.x == 3 && curr.y == 3 {
			continue
		}

		var newState state
		if curr.x-1 >= 0 {
			newState = state{x: curr.x - 1, y: curr.y, path: curr.path + "s"}
			if v, err := strconv.Atoi(finalMap[newState.x][newState.y]); err == nil {
				newState.weight = applyOp(curr.lastOp, curr.weight, v)
			} else {
				newState.weight = curr.weight
				newState.lastOp = finalMap[newState.x][newState.y]
			}
			if newState.weight > 0 && !(newState.x == 0 && newState.y == 0) {
				visited = append(visited, newState)
			}

		}
		if curr.x+1 < 4 {
			newState = state{x: curr.x + 1, y: curr.y, path: curr.path + "n"}
			if v, err := strconv.Atoi(finalMap[newState.x][newState.y]); err == nil {
				newState.weight = applyOp(curr.lastOp, curr.weight, v)
			} else {
				newState.weight = curr.weight
				newState.lastOp = finalMap[newState.x][newState.y]
			}
			if newState.weight > 0 && !(newState.x == 0 && newState.y == 0) {
				visited = append(visited, newState)
			}
		}
		if curr.y-1 >= 0 {
			newState = state{x: curr.x, y: curr.y - 1, path: curr.path + "w"}
			if v, err := strconv.Atoi(finalMap[newState.x][newState.y]); err == nil {
				newState.weight = applyOp(curr.lastOp, curr.weight, v)
			} else {
				newState.weight = curr.weight
				newState.lastOp = finalMap[newState.x][newState.y]
			}
			if newState.weight > 0 && !(newState.x == 0 && newState.y == 0) {
				visited = append(visited, newState)
			}
		}
		if curr.y+1 < 4 {
			newState = state{x: curr.x, y: curr.y + 1, path: curr.path + "e"}
			if v, err := strconv.Atoi(finalMap[newState.x][newState.y]); err == nil {
				newState.weight = applyOp(curr.lastOp, curr.weight, v)
			} else {
				newState.weight = curr.weight
				newState.lastOp = finalMap[newState.x][newState.y]
			}
			if newState.weight > 0 && !(newState.x == 0 && newState.y == 0) {
				visited = append(visited, newState)
			}
		}
	}
	return ""
}

func replaceNewLines(s bytes.Buffer) template.HTML {
	str := s.String()
	return template.HTML(strings.Replace(str, "\n", "<br>", -1))
}

/*
func format_mem(vm *VM) string {
	var str string
	for i := 0; i < len(vm.Memory); i++ {
		str += fmt.Sprintf("%05d: %05d\n", uint16(i), vm.Memory[i])
	}
	return str
}
*/

func execute(vm *VM, address *uint16) {
	var err error
	beenSet := false
	for {
		if *address == 0x1571 || *address == 0x1572 {
			vm.Memory[*address] = 21
			vm.Register[0] = 6
			// doing this because I know the answer
			vm.Register[7] = 0x6486
			// vm.Register[7] = findReg()
		}
		*address, err = step(vm, *address)
		//vm.MemoryTrace.Flush()
		if *address > 522 && !beenSet {
			vm.Register[7] = 9
			beenSet = true
		}
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

func outputMemory(f *bufio.Writer, vm VM) {
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
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: halt", address)))
	return address, fmt.Errorf("Halted program at address 0x%04x", address)
}

func set(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: set %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2))))
	ptr := getVal(vm, address+1)
	*ptr = *getVal(vm, address+2)
	return address + 3, nil
}

func push(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: push %04x\n", address, *getVal(vm, address+1))))
	vm.Stack = vm.Stack.Push(*getVal(vm, address+1))
	return address + 2, nil
}

func pop(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: pop %04x\n", address, *getVal(vm, address+1))))
	var value uint16
	if vm.Stack.isEmpty() {
		return address + 1, fmt.Errorf("Address %d: Cannot pop from empty stack", address)
	}
	vm.Stack, value = vm.Stack.Pop()
	*getVal(vm, address+1) = value
	return address + 2, nil
}

func eq(vm *VM, address uint16) (uint16, error) {
	a := getVal(vm, address+1)
	b := getVal(vm, address+2)
	c := getVal(vm, address+3)
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: eq %04x %04x %04x\n", address, *a, *b, *c)))
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
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: gt %04x %04x %04x\n", address, *a, *b, *c)))
	if *b > *c {
		*a = 1
	} else {
		*a = 0
	}
	return address + 4, nil
}

func jmp(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: jmp %04x\n", address, *getVal(vm, address+1))))
	return *getVal(vm, address+1), nil
}

func jt(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: jt %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2))))
	if *getVal(vm, address+1) != 0 {
		return *getVal(vm, address+2), nil
	}
	return address + 3, nil
}

func jf(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: jf %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2))))
	if *getVal(vm, address+1) == 0 {
		return *getVal(vm, address+2), nil
	}
	return address + 3, nil
}

func add(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: add %04x %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2), *getVal(vm, address+3))))
	*getVal(vm, address+1) = (*getVal(vm, address+2) + *getVal(vm, address+3)) & 0x7FFF
	return address + 4, nil
}

func mult(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: mult %04x %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2), *getVal(vm, address+3))))
	*getVal(vm, address+1) = (*getVal(vm, address+2) * *getVal(vm, address+3)) & 0x7FFF
	return address + 4, nil
}

func mod(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: mod %04x %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2), *getVal(vm, address+3))))
	*getVal(vm, address+1) = uint16(int(*getVal(vm, address+2)) % int(*getVal(vm, address+3)))
	return address + 4, nil
}

func and(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: and %04x %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2), *getVal(vm, address+3))))
	*getVal(vm, address+1) = *getVal(vm, address+2) & (*getVal(vm, address+3))
	return address + 4, nil
}

func or(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: or %04x %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2), *getVal(vm, address+3))))
	*getVal(vm, address+1) = *getVal(vm, address+2) | *getVal(vm, address+3)
	return address + 4, nil
}

func not(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: not %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2))))
	*getVal(vm, address+1) = ^(*getVal(vm, address+2)) & 0x7FFF
	return address + 3, nil
}

func rmem(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: rmem %04x %04x\n", address, *getVal(vm, address+1), *getVal(vm, address+2))))
	*getVal(vm, address+1) = *getVal(vm, *getVal(vm, address+2))
	return address + 3, nil
}

func wmem(vm *VM, address uint16) (uint16, error) {
	b := *getVal(vm, address+2)
	a := *getVal(vm, address+1)
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: wmem %04x %04x\n", address, a, b)))
	*getVal(vm, a) = b
	return address + 3, nil
}

func call(vm *VM, address uint16) (uint16, error) {
	vm.Stack = vm.Stack.Push(address + 2)
	val := *getVal(vm, address+1)
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: call %04x\n", address, val)))
	return val, nil
}

func ret(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: ret\n", address)))
	if vm.Stack.isEmpty() {
		return address, fmt.Errorf("Address %04x: pop from empty stack\n", address)
	}
	var value uint16
	vm.Stack, value = vm.Stack.Pop()
	return value, nil

}

func out(vm *VM, address uint16) (uint16, error) {
	val := byte(*getVal(vm, address+1))
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: out %c\n", address, val)))
	vm.Output.Write([]byte{val})
	vm.Output.Flush()
	return address + 2, nil
}

func in(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: in %04x\n", address, address+1)))
	var val uint16
	v, err := vm.Input.ReadByte()
	if err != nil {
		for err == nil {
			v, err = vm.Input.ReadByte()
		}
	}
	val = uint16(v)
	*getVal(vm, address+1) = val
	return address + 2, nil
}

func noop(vm *VM, address uint16) (uint16, error) {
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %04x: noop\n", address)))
	return address + 1, nil
}

func getVal(vm *VM, address uint16) *uint16 {
	val := &vm.Memory[address]

	if *val > 0x7FFF {
		val = &vm.Register[(*val)&0x7FFF]
	}
	//vm.MemoryTrace.Write([]byte(fmt.Sprintf("Address %d: %d\n", int(address), int(*val))))
	//vm.MemoryTrace.Flush()
	return val
}
