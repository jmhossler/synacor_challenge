#ifndef STACK_SYNACOR_JH
#define STACK_SYNACOR_JH
#include <stdint.h>

typedef struct stack Stack;

Stack *new_Stack();
void stack_push(Stack *, uint16_t);
uint16_t stack_pop(Stack *);
int is_empty(Stack *);

#endif
