#pragma bits 16

int foo() {
    int* a = 1;
    int** b = 2;

    a = (int) *b;
}
