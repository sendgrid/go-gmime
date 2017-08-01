#include <stdlib.h>
#include <strings.h>
#include <gmime/gmime.h>

GMimeMessage *gmime_parse (const char *buffer, size_t len);
char* gmime_get_content_string (GMimeObject *object);
char* gmime_get_content_type_string (GMimeObject *object);
gboolean gmime_is_text_part (GMimeObject *object);
GByteArray *gmime_get_bytes (GMimeObject *object);
