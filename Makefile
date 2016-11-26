all: vm

vm: main.o memory.o virtual_machine.o stack.o
	gcc main.o memory.o virtual_machine.o stack.o -o vm

main.o: src/main.c virtual_machine.o memory.o stack.o
	gcc -c src/main.c

virtual_machine.o: src/virtual_machine.c src/virtual_machine.h memory.o stack.o
	gcc -c src/virtual_machine.c

memory.o: src/memory.c src/memory.h
	gcc -c src/memory.c

stack.o: src/stack.c src/stack.h
	gcc -c src/stack.c

clean:
	rm *.o &> /dev/null
