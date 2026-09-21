void copy(unsigned char* dst, int ds, int de, unsigned char* src, int ss, int se) {
    int i = ds;
    int j = ss;

    while (j < se && i < de) {
        dst[i] = src[j];
        i++;
        j++;
    }
}
