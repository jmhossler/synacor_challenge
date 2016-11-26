#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "virtual_machine.h"
#include "memory.h"
#include "stack.h"

typedef uint16_t (*ptr_to_func)(VM *, uint16_t);

ptr_to_func get_func(uint16_t);

uint16_t halt(VM *, uint16_t);
uint16_t out(VM *, uint16_t);
uint16_t noop(VM *, uint16_t);
uint16_t set(VM *, uint16_t);
uint16_t push(VM *, uint16_t);
uint16_t pop(VM *, uint16_t);
uint16_t eq(VM *, uint16_t);
uint16_t gt(VM *, uint16_t);
uint16_t jmp(VM *, uint16_t);
uint16_t jt(VM *, uint16_t);
uint16_t jf(VM *, uint16_t);
uint16_t add(VM *, uint16_t);
uint16_t mult(VM *, uint16_t);
uint16_t mod(VM *, uint16_t);
uint16_t and(VM *, uint16_t);
uint16_t or(VM *, uint16_t);
uint16_t not(VM *, uint16_t);
uint16_t rmem(VM *, uint16_t);
uint16_t wmem(VM *, uint16_t);
uint16_t call(VM *, uint16_t);
uint16_t ret(VM *, uint16_t);
uint16_t in(VM *, uint16_t);

ptr_to_func function_ptrs[22] = { *halt, //0
                                  *set,  //1
                                  *push, //2
                                  *pop,  //3
                                  *eq,   //4
                                  *gt,   //5
                                  *jmp,  //6
                                  *jt,   //7
                                  *jf,   //8
                                  *add,  //9
                                  *mult, //10
                                  *mod,  //11
                                  *and,  //12
                                  *or,   //13
                                  *not,  //14
                                  *rmem, //15
                                  *wmem, //16
                                  *call, //17
                                  *ret,  //18
                                  *out,  //19
                                  *in,   //20
                                  *noop, //21
                                  };

struct vm {
  Memory *memory;
  Stack *stack;
  uint16_t current_address;
};

VM *
new_VM(FILE *fp)
{
  VM *vm = malloc(sizeof(VM));
  vm->memory = new_Memory(fp, 0x8000);
  vm->stack = new_Stack();
  vm->current_address = 0;

  return vm;
}

void
execute(VM *vm, uint16_t address)
{
  if(address > 0x7FFF)
  {
    printf("Finished iterating over memory\n");
    return;
  }
  //display_reg(vm->memory);
  uint16_t op_code = get_val(vm->memory, address);
  //printf("Function %u\n", (unsigned int) op_code);
  ptr_to_func func = get_func(op_code);
  if(func == NULL)
  {
    address += 1;
  } else {
    address = (*func)(vm, address);
  }
  execute(vm, address);
  return;
}

ptr_to_func
get_func(uint16_t op_code)
{
  if(op_code > 21)
  {
    fprintf(stderr, "Out of range index into function_ptrs, %u\n", ((unsigned int) op_code));
    exit(1);
  }
  return function_ptrs[op_code];
}

uint16_t
halt(VM *vm, uint16_t address)
{
  printf("Halted at %u\n", (unsigned int) address);
  exit(0);
}

uint16_t
out(VM *vm, uint16_t address)
{
  printf("%c", (char) get_val(vm->memory, address + 1));

  return address + 2;
}

uint16_t
noop(VM *vm, uint16_t address)
{
  return address + 1;
}

uint16_t
set(VM *vm, uint16_t address)
{
  set_val(vm->memory, get_val(vm->memory, address + 1) + 32768, get_val(vm->memory, address + 2));
  return address + 3;
}

uint16_t
push(VM *vm, uint16_t address)
{
  stack_push(vm->stack, get_val(vm->memory, address + 1));
  return address + 2;
}

uint16_t
pop(VM *vm, uint16_t address)
{
  set_val(vm->memory, address + 1, stack_pop(vm->stack));
  return address + 2;
}

uint16_t
eq(VM *vm, uint16_t address)
{
  if(get_val(vm->memory, address + 2) == get_val(vm->memory, address + 3))
  {
    set_val(vm->memory, address + 1, 1);
  } else
  {
    set_val(vm->memory, address + 1, 0);
  }
  return address + 4;
}

uint16_t
gt(VM *vm, uint16_t address)
{
  if(get_val(vm->memory, address + 2) > get_val(vm->memory, address + 3))
  {
    set_val(vm->memory, address + 1, 1);
  } else
  {
    set_val(vm->memory, address + 1, 0);
  }
  return address + 4;
}

uint16_t
jmp(VM *vm, uint16_t address)
{
  return get_val(vm->memory, address + 1);
}

uint16_t
jt(VM *vm, uint16_t address)
{
  if(get_val(vm->memory, address + 1) != 0)
  {
    return get_val(vm->memory, address + 2);
  } else
  {
    return address + 3;
  }
}

uint16_t
jf(VM *vm, uint16_t address)
{
  if(get_val(vm->memory, address + 1) == 0)
  {
    return get_val(vm->memory, address + 2);
  } else
  {
    return address + 3;
  }
}

uint16_t
add(VM *vm, uint16_t address)
{
  uint16_t result = get_val(vm->memory, address + 2) + get_val(vm->memory, address + 3);
  set_val(vm->memory, get_val(vm->memory, address + 1), result % 32768);
  return address + 4;
}

uint16_t
mult(VM *vm, uint16_t address)
{
  uint16_t result = get_val(vm->memory, address + 2) + get_val(vm->memory, address + 3);
  set_val(vm->memory, get_val(vm->memory, address + 1), result % 32768);
  return address + 4;
}

uint16_t
mod(VM *vm, uint16_t address)
{
  uint16_t result = get_val(vm->memory, address + 2) % get_val(vm->memory, address + 3);
  set_val(vm->memory, get_val(vm->memory, address + 1), result % 32768);
  return address + 4;
}

uint16_t
and(VM *vm, uint16_t address)
{
  set_val(vm->memory, get_val(vm->memory, address + 1), get_val(vm->memory, address + 2) & get_val(vm->memory, address + 3));
  return address + 4;
}

uint16_t
or(VM *vm, uint16_t address)
{
  set_val(vm->memory, get_val(vm->memory, address + 1), get_val(vm->memory, address + 2) | get_val(vm->memory, address + 3));
  return address + 4;
}

uint16_t
not(VM *vm, uint16_t address)
{
  set_val(vm->memory, get_val(vm->memory, address + 1), ~(get_val(vm->memory, address + 2)) & 0x7FFF);
  return address + 3;
}

uint16_t
rmem(VM *vm, uint16_t address)
{
  set_val(vm->memory, get_val(vm->memory, address + 1), get_val(vm->memory, get_val(vm->memory, address + 2)));
  return address + 3;
}

uint16_t
wmem(VM *vm, uint16_t address)
{
  set_val(vm->memory, get_val(vm->memory, address + 1), get_val(vm->memory, address + 2));
  return address + 3;
}

uint16_t
call(VM *vm, uint16_t address)
{
  stack_push(vm->stack, address + 2);
  return get_val(vm->memory, address + 1);
}

uint16_t
ret(VM *vm, uint16_t address)
{
  if(is_empty(vm->stack))
  {
    return halt(vm, address);
  } else
  {
    return stack_pop(vm->stack);
  }
}

uint16_t
in(VM *vm, uint16_t address)
{
  char val;
  fscanf(stdin, "%c", &val);
  set_val(vm->memory, address + 1, (uint16_t) val);
  return address + 2;
}
