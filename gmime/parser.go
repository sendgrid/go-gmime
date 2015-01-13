package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

type Parser interface {
	Janitor
	ConstructMessage() Message
	ConstructPart() Object
	Tell() int64
	Eos() bool
	SetScanFrom(bool)
}

type aParser struct {
	*PointerMixin
}

func CastParser(p *C.GMimeParser) Parser {
	return &aParser{CastPointer(C.gpointer(p))}
}

func NewParserWithStream(stream Stream) Parser {
	rawStream := stream.(rawStream)

	parser := C.g_mime_parser_new_with_stream(rawStream.rawStream())
	defer unref(C.gpointer(parser))
	return CastParser(parser)
}

func (p *aParser) ConstructMessage() Message {
	message := C.g_mime_parser_construct_message(p.rawParser())
	defer unref(C.gpointer(message))
	return CastMessage(message)
}

func (p *aParser) ConstructPart() Object {
	object := C.g_mime_parser_construct_part(p.rawParser())
	defer unref(C.gpointer(object))
	return objectAsSubclass(object)
}

func (p *aParser) Tell() int64 {
	cint := C.g_mime_parser_tell(p.rawParser())
	return int64(cint)
}

func (p *aParser) Eos() bool {
	cbool := C.g_mime_parser_eos(p.rawParser())
	return gobool(cbool)
}

func (p *aParser) SetScanFrom(scanFrom bool) {
	C.g_mime_parser_set_scan_from(p.rawParser(), gbool(scanFrom))
}

func (p *aParser) rawParser() *C.GMimeParser {
	return (*C.GMimeParser)(p.pointer())
}
