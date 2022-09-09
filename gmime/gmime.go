package gmime

import "C"
import (
	"net/mail"
	"unsafe"
)

// #cgo pkg-config: gmime-3.0
// #include "gmime.h"
import "C"

var (
	cStringEmpty       = C.CString("")
	cStringAlternative = C.CString("alternative")
	cStringMixed       = C.CString("mixed")
	cStringRelated     = C.CString("related")
	cStringCharset     = C.CString("charset")
	cStringCharsetUTF8 = C.CString("utf-8")

	cStringText   = C.CString("text")
	cStringPlain  = C.CString("plain")
	cStringHTML   = C.CString("html")
	cStringBase64 = C.CString("base64")

	cStringContentID               = C.CString("Content-Id")
	cStringHeaderFormat            = C.CString("%s: %s\n")
	cStringContentTransferEncoding = C.CString("Content-Transfer-Encoding")
)

// This function call automatically by runtime
func init() {
	C.g_mime_init()
	format := C.g_mime_format_options_get_default()
	C.g_mime_format_options_set_newline_format(format, C.GMIME_NEWLINE_FORMAT_DOS)
}

// Shutdown is really needed only for valgrind
func Shutdown() {
	C.g_mime_shutdown()
}

// convert from Go bool to C gboolean
func gbool(b bool) C.gboolean {
	if b {
		return C.gboolean(1)
	}
	return C.gboolean(0)
}

// convert from C gboolean to Go bool
func gobool(b C.gboolean) bool {
	return b != C.gboolean(0)
}

// free up memory
func unref(referee C.gpointer) {
	C.g_object_unref(referee)
}

// ParseAddressList parses and returns address list
func ParseAddressList(addrs string) []*mail.Address {
	cAddrs := C.CString(addrs)
	defer C.free(unsafe.Pointer(cAddrs))
	parsedAddrs := C.internet_address_list_parse(C.g_mime_parser_options_get_default(), cAddrs)
	if parsedAddrs == nil {
		return nil
	}
	// dont move this up next to instatiation. we dont want to free this if parseAddrs is nil
	// gmime will free it
	defer C.g_object_unref((C.gpointer)(unsafe.Pointer(parsedAddrs)))
	nAddrs := C.internet_address_list_length(parsedAddrs)
	if nAddrs <= 0 {
		return nil
	}

	var i C.int
	goAddrs := make([]*mail.Address, nAddrs)
	for i = 0; i < nAddrs; i++ {
		address := C.internet_address_list_get_address(parsedAddrs, i)
		gAddr := convertToGoAddress(address)
		goAddrs[i] = gAddr
	}
	return goAddrs
}

func convertToGoAddress(addr *C.InternetAddress) *mail.Address {
	var gAddr mail.Address
	name := C.internet_address_get_name(addr)
	address := C.internet_address_mailbox_get_addr((*C.InternetAddressMailbox)(unsafe.Pointer(addr)))
	if name != nil {
		gAddr.Name = C.GoString(name)
	}
	gAddr.Address = C.GoString(address)
	return &gAddr
}
