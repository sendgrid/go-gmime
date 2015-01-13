package gmime

/*
#cgo pkg-config: gmime-2.6
#include <stdlib.h>
#include <gmime/gmime.h>

// We do this to determine the type of address returned since they are
// contained in macros on internet-address.h and CGO can't access macros.
static gboolean is_internet_address_mailbox(GTypeInstance *obj) {
    return INTERNET_ADDRESS_IS_MAILBOX(obj);
}
static gboolean is_internet_address_group(GTypeInstance *obj) {
    return INTERNET_ADDRESS_IS_GROUP(obj);
}

*/
import "C"
import (
	"unsafe"
)

// XXX The underlying type-unsafe C implementation is modeled here with a type-safe
// Object-Oriented design, which is a modified form of the Composite Design Pattern.
// Specifically, Address is the Component, Mailbox is the Leaf, and Group & AddressList
// co-operate to act as the Composite. SEE TestCompositeAddress.
type Address interface {
	Janitor
	Name() string
	SetName(string)
}

type Mailbox interface {
	Address
	Email() string
	SetEmail(string)
}

type Group interface {
	Address
	Members() AddressList
	SetMembers(AddressList)
	AddMember(Address) int
}

type AddressList interface {
	Janitor
	GetLength() int
	Clear()
	Add(Address) int
	Prepend(AddressList)
	Append(AddressList)
	Insert(Address, int)
	Remove(Address) bool
	RemoveAt(int) bool
	Contains(Address) bool
	IndexOf(Address) int
	GetAddress(int) Address
	SetAddress(Address, int)
	ToString(bool) string
}

type internetAddress struct {
	*PointerMixin
}

type rawAddress interface {
	Address
	rawAddress() *C.InternetAddress
}

type rawInternetAddressList interface {
	AddressList
	rawList() *C.InternetAddressList
}

func isMailbox(address *C.InternetAddress) bool {
	return (address != nil) && gobool(C.is_internet_address_mailbox((*C.GTypeInstance)(unsafe.Pointer(address))))
}

func isGroup(address *C.InternetAddress) bool {
	return (address != nil) && gobool(C.is_internet_address_group((*C.GTypeInstance)(unsafe.Pointer(address))))
}

func (a *internetAddress) SetName(name string) {
	var cName *C.char = C.CString(name)
	C.internet_address_set_name(a.rawAddress(), cName)
	C.free(unsafe.Pointer(cName))
}

func (a *internetAddress) Name() string {
	name := C.internet_address_get_name(a.rawAddress())
	return C.GoString(name)
}

func (a *internetAddress) rawAddress() *C.InternetAddress {
	return (*C.InternetAddress)(a.pointer())
}

func CastAddress(ca *C.InternetAddress) Address {
	if isMailbox(ca) {
		return CastMailboxAddress((*C.InternetAddressMailbox)(unsafe.Pointer(ca)))
	} else if isGroup(ca) {
		return CastGroupAddress((*C.InternetAddressGroup)(unsafe.Pointer(ca)))
	} else {
		panic("bad cast: not MailboxAddress nor GroupAddress")
	}
}

type MailboxAddress struct {
	*internetAddress
}

func CastMailboxAddress(ma *C.InternetAddressMailbox) *MailboxAddress {
	ptr := CastPointer(C.gpointer(ma))
	ia := &internetAddress{ptr}
	return &MailboxAddress{ia}
}

func NewMailboxAddress(name string, email string) *MailboxAddress {
	var cName *C.char = C.CString(name)
	var cEmail *C.char = C.CString(email)
	defer C.free(unsafe.Pointer(cName))
	defer C.free(unsafe.Pointer(cEmail))
	ai := C.internet_address_mailbox_new(cName, cEmail)
	defer unref(C.gpointer(ai))
	return CastMailboxAddress((*C.InternetAddressMailbox)(C.gpointer(ai)))
}

func (a *MailboxAddress) SetEmail(email string) {
	var cEmail *C.char = C.CString(email)
	C.internet_address_mailbox_set_addr((*C.InternetAddressMailbox)(unsafe.Pointer(a.rawAddress())), cEmail)
	C.free(unsafe.Pointer(cEmail))
}

func (a *MailboxAddress) Email() string {
	email := C.internet_address_mailbox_get_addr((*C.InternetAddressMailbox)(unsafe.Pointer(a.rawAddress())))
	return C.GoString(email)
}

type GroupAddress struct {
	*internetAddress
}

func CastGroupAddress(ga *C.InternetAddressGroup) *GroupAddress {
	ia := &internetAddress{CastPointer(C.gpointer(ga))}
	return &GroupAddress{ia}
}

func NewGroupAddress(name string) *GroupAddress {
	var cName *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cName))
	ga := C.internet_address_group_new(cName)
	defer unref(C.gpointer(ga))
	return CastGroupAddress((*C.InternetAddressGroup)(unsafe.Pointer(ga)))
}

func (group *GroupAddress) SetMembers(list AddressList) {
	ial := list.(rawInternetAddressList)
	C.internet_address_group_set_members((*C.InternetAddressGroup)(group.pointer()), ial.rawList())
}

func (group *GroupAddress) Members() AddressList {
	cList := C.internet_address_group_get_members((*C.InternetAddressGroup)(group.pointer()))
	if cList != nil {
		return CastInternetAddressList(cList)
	}
	return nil
}

func (group *GroupAddress) AddMember(newAddress Address) int {
	address := newAddress.(rawAddress)
	offset := C.internet_address_group_add_member((*C.InternetAddressGroup)(group.pointer()), address.rawAddress())
	return int(offset)
}

type InternetAddressList struct {
	*PointerMixin
}

type rawAddressList interface {
	AddressList
	rawList() *C.InternetAddressList
}

func CastInternetAddressList(cList *C.InternetAddressList) *InternetAddressList {
	return &InternetAddressList{CastPointer(C.gpointer(cList))}
}

func NewAddressList() *InternetAddressList {
	cList := C.internet_address_list_new()
	defer unref(C.gpointer(cList))
	return CastInternetAddressList(cList)
}

func (l *InternetAddressList) GetLength() int {
	return int(C.internet_address_list_length(l.rawList()))
}

func (l *InternetAddressList) Clear() {
	C.internet_address_list_clear(l.rawList())
}

func (l *InternetAddressList) Add(address Address) int {
	rawAddress := address.(rawAddress)
	offset := C.internet_address_list_add(l.rawList(), rawAddress.rawAddress())
	return int(offset)
}

func (l *InternetAddressList) Prepend(newList AddressList) {
	C.internet_address_list_prepend(l.rawList(), newList.(rawInternetAddressList).rawList())
}

func (l *InternetAddressList) Append(newList AddressList) {
	C.internet_address_list_append(l.rawList(), newList.(rawInternetAddressList).rawList())
}

func (l *InternetAddressList) Insert(address Address, index int) {
	C.internet_address_list_insert(l.rawList(), (C.int)(index), address.(rawAddress).rawAddress())
}

func (l *InternetAddressList) Remove(address Address) bool {
	ret := C.internet_address_list_remove(l.rawList(), address.(rawAddress).rawAddress())
	return gobool(ret)
}

func (l *InternetAddressList) RemoveAt(index int) bool {
	ret := C.internet_address_list_remove_at(l.rawList(), (C.int)(index))
	return gobool(ret)
}

func (l *InternetAddressList) Contains(address Address) bool {
	rawAddress := address.(rawAddress)
	ret := C.internet_address_list_contains(l.rawList(), rawAddress.rawAddress())
	return gobool(ret)
}

func (l *InternetAddressList) IndexOf(address Address) int {
	rawAddress := address.(rawAddress)
	ret := C.internet_address_list_index_of(l.rawList(), rawAddress.rawAddress())
	return int(ret)
}

func (l *InternetAddressList) GetAddress(index int) Address {
	cAddress := C.internet_address_list_get_address(l.rawList(), (C.int)(index))
	return CastAddress(cAddress)
}

func (l *InternetAddressList) SetAddress(address Address, index int) {
	rawAddress := address.(rawAddress)
	C.internet_address_list_set_address(l.rawList(), (C.int)(index), rawAddress.rawAddress())
}

func (l *InternetAddressList) ToString(encode bool) string {
	addresses := C.internet_address_list_to_string(l.rawList(), gbool(encode))
	if addresses != nil {
		addressString := C.GoString(addresses)
		defer C.free(unsafe.Pointer(addresses))
		return addressString
	}
	return ""
}

func (l *InternetAddressList) rawList() *C.InternetAddressList {
	return (*C.InternetAddressList)(l.pointer())
}

func ParseString(str string) AddressList {
	var cStr *C.char = C.CString(str)
	defer C.free(unsafe.Pointer(cStr))

	cList := C.internet_address_list_parse_string(cStr)
	defer unref(C.gpointer(cList))

	if cList != nil {
		return CastInternetAddressList(cList)
	}
	return nil
}
