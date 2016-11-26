#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "virtual_machine.h"

void check_file(FILE *, char *);

int
main(int argc, char **argv)
{
  if(argc != 2)
  {
    fprintf(stderr, "Incorrect number of arguments, 2 expected\n");
    exit(1);
  }

  char *input_file = argv[1];
  FILE *fp = fopen(input_file, "r");
  check_file(fp, input_file);

  VM *virtual_machine = new_VM(fp);
  fclose(fp);

  execute(virtual_machine, 0);
  return 0;
}

void
check_file(FILE *fp, char *filename)
{
  if(!fp)
  {
    fprintf(stderr, "%s could not be opened\n", filename);
    exit(1);
  }
}
