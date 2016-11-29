#include <stdlib.h>
#include <stdint.h>
#include <stdio.h>
#include "memory.h"

uint16_t little_endian_to_16(uint8_t, uint8_t);

struct memory {
  uint16_t *array;
  uint16_t size;
  uint16_t *registers;
};

Memory *
new_Memory(FILE *fp, uint16_t size)
{
  Memory *memory = malloc(sizeof(Memory));
  memory->size = size;
  memory->array = malloc(sizeof(uint16_t) * (size + 8));

  uint16_t i = 0;
  uint8_t little, big;
  fscanf(fp, "%c%c", &little, &big);
  while(!feof(fp))
  {
    memory->array[i++] = little_endian_to_16(little, big);
    fscanf(fp, "%c%c", &little, &big);
  }
  memory->array[i++] = little_endian_to_16(little, big);

  for(i = 0; i < 8; ++i)
  {
    memory->array[memory->size + i] = 0;
  }

  return memory;
}

uint16_t
little_endian_to_16(uint8_t little, uint8_t big)
{
  return (((uint16_t) big) << 8) + (uint16_t) little;
}

uint16_t
get_val(Memory *memory, uint16_t address)
{
  if(address >= memory->size + 8)
  {
    fprintf(stderr, "Out of bounds index\n");
    exit(1);
  }
  if(memory->array[address] > 0x7FFF) {
    return memory->array[memory->array[address] - 0x8000];
  }
  return memory->array[address];
}

void
set_val(Memory *memory, uint16_t address, uint16_t val)
{
  if(address >= memory->size + 8)
  {
    fprintf(stderr, "Out of bounds index\n");
    exit(1);
  }
  memory->array[address] = val;
}

void
display_reg(Memory *memory)
{
  int i;
  for(i = 0; i < 8; ++i)
  {
    printf("REG: %d\tVAL: %u\n", i, (unsigned int) memory->array[i + memory->size]);
  }
}
