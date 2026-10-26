#pragma once

void lcc_fatal(char* label, int line, ...);
void lcc_error(char* label, int line, bool* marker, ...);
void lcc_info(char* label, char* str);
void lcc_ice(char* label, int line, ...);
