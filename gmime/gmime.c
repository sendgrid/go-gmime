#include "gmime.h"


GMimeMessage *gmime_parse (const char *buffer, size_t len) {
	GMimeStream *stream = g_mime_stream_mem_new_with_buffer (buffer, len);
	GMimeParser *parser = g_mime_parser_new_with_stream (stream);
	g_object_unref (stream);
	GMimeMessage *message = g_mime_parser_construct_message (parser, NULL);
	g_object_unref (parser);
	if (!message) {
		return NULL; 
	}

	return message;
}

char* gmime_get_content_string (GMimeObject *object) {
	if (!GMIME_IS_TEXT_PART (object)) {
		return NULL;
	}
	return g_mime_text_part_get_text ((GMimeTextPart *) object);
}

char* gmime_get_content_type_string (GMimeObject *object) {
	GMimeContentType *ctype = g_mime_object_get_content_type (object);
	return g_mime_content_type_get_mime_type (ctype);
}

gboolean gmime_is_text_part (GMimeObject *object) {
	return GMIME_IS_TEXT_PART (object);
}

gboolean gmime_is_multi_part (GMimeObject *object) {
	return GMIME_IS_MULTIPART (object);
}

void gmime_type_name(GMimeObject *object){
	printf("Name: %s\n", G_OBJECT_TYPE_NAME (object));
}

GByteArray *gmime_get_bytes (GMimeObject *object) {
	GMimeStream *stream;
	GMimeDataWrapper *content;
	GByteArray *buf;

	if (!(content = g_mime_part_get_content ((GMimePart *) object)))
		return NULL;
	stream = g_mime_stream_mem_new ();
	ssize_t size = g_mime_data_wrapper_write_to_stream (content, stream);
	g_mime_stream_flush (stream);

	buf = g_mime_stream_mem_get_byte_array ((GMimeStreamMem *) stream);
	g_mime_stream_mem_set_owner ((GMimeStreamMem *) stream, FALSE);

	g_object_unref (stream);
	return buf;
}
