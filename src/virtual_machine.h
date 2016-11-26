#ifndef VM_SYNACOR_JH_2016
#define VM_SYNACOR_JH_2016
#include <stdio.h>
#include <stdint.h>

typedef struct vm VM;

VM *new_VM(FILE *);
void execute(VM *, uint16_t);

#endif
