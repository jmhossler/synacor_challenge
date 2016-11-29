#ifndef MEMORY_SYNACOR_JH_2016
#define MEMORY_SYNACOR_JH_2016
#include <stdio.h>
#include <stdint.h>

typedef struct memory Memory;

Memory *new_Memory(FILE *, uint16_t);
uint16_t *get_val(Memory *, uint16_t);
void set_val(Memory *, uint16_t, uint16_t);
void display_reg(Memory *);

#endif
