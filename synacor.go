package main

import (
  "fmt"
  "io/ioutil"
  "os"
  "bytes"
  "encoding/binary"
)

type Stack []uint16

type VM struct {
  memory [32768]uint16
  register [8]uint16
  stack Stack
}

func main() {
  input_file := process_args(os.Args[1:])
  var vm VM
  var byte_array []byte
  byte_array, err := ioutil.ReadFile(input_file)
  check(err)

  for i := 0; i < len(byte_array); i += 2 {
    vm.memory[int(i / 2)] = read_uint16(byte_array[i:i+2])
  }

  output_memory(os.Stdout, vm)
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

func process_args(array []string) (string){
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
    return s[:len(s) - 1], s[len(s) - 1]
  }
}

func halt(w http.ResponseWriter, address uint16) uint16 {

  return address
}

func out(address uint16) uint16 {

