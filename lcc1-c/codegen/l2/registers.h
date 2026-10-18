typedef struct {
    char name[4];
    bool taken;
} _register;

void regs_init();
char* take_register();
void free_register(char* name);
