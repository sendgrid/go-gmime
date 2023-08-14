#include <stdlib.h>
#include <strings.h>
#include <gmime/gmime.h>

GMimeMessage *gmime_parse (const char *buffer, size_t len);
char* gmime_get_content_string (GMimeObject *object);
char* gmime_get_content_type_string (GMimeObject *object);
char* gmime_get_content_disposition(GMimeObject *object);
gboolean gmime_is_multi_part (GMimeObject *object);
gboolean gmime_is_part (GMimeObject *object);
gboolean gmime_is_text_part (GMimeObject *object);
gboolean gmime_is_content_type (GMimeObject *object);
void gmime_type_name(GMimeObject *object);
GByteArray *gmime_get_bytes (GMimeObject *object);
char* gmime_get_content_string_full (GMimeObject *object);
