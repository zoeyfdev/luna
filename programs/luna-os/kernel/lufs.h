#include "stdbool.h"

typedef struct {
    long int* Address;
    char* Name;
} File;

extern void fcreate(char* name, long int size);
extern void fwrite(char* name, char* buffer);
extern File* fopen(char* filename, bool complain_on_not_found);
extern long int fgetsize(char* filename);
extern long int* ffnt(char* filename);
extern void flist();
extern long int* fntf(char* name);
extern void fstrap();

