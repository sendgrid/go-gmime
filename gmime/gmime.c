#include <string.h>
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

/**
 * gmime_text_part_set_text:
 * @mime_part: a #GMimeTextPart
 * @text: the text in utf-8
 *
 * Sets the specified text as the content and updates the charset parameter on the Content-Type header.
 **/
void
gmime_text_part_set_text (GMimeTextPart *mime_part, const char *text)
{
	GMimeContentType *content_type;
	GMimeStream *filtered, *stream;
	GMimeContentEncoding encoding;
	GMimeDataWrapper *content;
	GMimeFilter *filter;
	const char *charset;
	GMimeCharset mask;
	size_t len;

	g_return_if_fail (GMIME_IS_TEXT_PART (mime_part));
	g_return_if_fail (text != NULL);

	len = strlen (text);


  charset = "utf-8";

	content_type = g_mime_object_get_content_type ((GMimeObject *) mime_part);
	g_mime_content_type_set_parameter (content_type, "charset", charset);

	stream = g_mime_stream_mem_new_with_buffer (text, len);


	content = g_mime_data_wrapper_new_with_stream (stream, GMIME_CONTENT_ENCODING_DEFAULT);
	g_object_unref (stream);

	g_mime_part_set_content ((GMimePart *) mime_part, content);
	g_object_unref (content);

	encoding = g_mime_part_get_content_encoding ((GMimePart *) mime_part);

	/* if the user has already specified encoding the content with base64/qp/uu, don't change it */
	if (encoding > GMIME_CONTENT_ENCODING_BINARY)
		return;


  g_mime_part_set_content_encoding ((GMimePart *) mime_part, GMIME_CONTENT_ENCODING_8BIT);

}
