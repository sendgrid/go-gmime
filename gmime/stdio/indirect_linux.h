#include <stdio.h>
extern ssize_t (*c_reader)(void *, char *, size_t);
extern ssize_t (*c_writer)(void *, const char *, size_t);
extern int (*c_seeker)(void *, off64_t *, int);
extern int (*c_closer)(void *);
