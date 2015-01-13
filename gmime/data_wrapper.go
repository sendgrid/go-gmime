package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>
*/
import "C"

type DataWrapper interface {
	Janitor
	Encoding() ContentEncoding
	WriteToStream(stream Stream) uintptr
	Stream() Stream
}

type rawDataWrapper interface {
	DataWrapper
	rawDataWrapper() *C.GMimeDataWrapper
}

type aDataWrapper struct {
	*PointerMixin
}

func CastDataWrapper(dw *C.GMimeDataWrapper) *aDataWrapper {
	return &aDataWrapper{CastPointer(C.gpointer(dw))}
}

func NewDataWrapper() DataWrapper {
	dw := C.g_mime_data_wrapper_new()
	defer unref(C.gpointer(dw))
	return CastDataWrapper(dw)
}

func NewDataWrapperWithStream(stream Stream, encoding ContentEncoding) DataWrapper {
	rawStream := stream.(rawStream)
	rawEncoding := encoding.(rawContentEncoding)

	dw := C.g_mime_data_wrapper_new_with_stream(rawStream.rawStream(), rawEncoding.rawContentEncoding())
	defer unref(C.gpointer(dw))
	return CastDataWrapper(dw)
}

func (d *aDataWrapper) Stream() Stream {
	return CastStream(C.g_mime_data_wrapper_get_stream(d.rawDataWrapper()))
}

func (d *aDataWrapper) Encoding() ContentEncoding {
	return CastContentEncoding(C.g_mime_data_wrapper_get_encoding(d.rawDataWrapper()))
}

func (d *aDataWrapper) WriteToStream(stream Stream) uintptr {
	rawStream := stream.(rawStream)
	return uintptr(C.g_mime_data_wrapper_write_to_stream(d.rawDataWrapper(), rawStream.rawStream()))
}

func (d *aDataWrapper) rawDataWrapper() *C.GMimeDataWrapper {
	return (*C.GMimeDataWrapper)(d.pointer())
}
