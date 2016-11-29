#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "stack.h"

typedef struct node {
  uint16_t data;
  struct node *next;
} Node;

Node *new_Node(uint16_t);

struct stack {
  Node *head;
};

Stack *
new_Stack()
{
  Stack *stack = malloc(sizeof(Stack));
  stack->head = NULL;
  return stack;
}

void
stack_push(Stack *s, uint16_t v)
{
  Node *node = new_Node(v);
  if(s->head == NULL)
  {
    s->head = node;
  } else
  {
    node->next = s->head;
    s->head = node;
  }
}

uint16_t
stack_pop(Stack *s)
{
  if(is_empty(s))
  {
    fprintf(stderr, "Pop from empty stack\n");
    exit(1);
  } else
  {
    uint16_t temp_val = s->head->data;
    s->head = s->head->next;

    return temp_val;
  }
}

int
is_empty(Stack *s)
{
  return s->head == NULL;
}

Node *
new_Node(uint16_t val)
{
  Node *node = malloc(sizeof(Node));
  node->data = val;
  node->next = NULL;

  return node;
}
