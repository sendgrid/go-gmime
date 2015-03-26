package gmime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailboxAddress(t *testing.T) {
	var mailbox *MailboxAddress = NewMailboxAddress("box1", "kickass1@example.com")
	assert.Equal(t, mailbox.Name(), "box1")
	assert.Equal(t, mailbox.Email(), "kickass1@example.com")
	mailbox.SetName("box2")
	assert.Equal(t, mailbox.Name(), "box2")
	mailbox.SetEmail("kickass2@example.com")
	assert.Equal(t, mailbox.Email(), "kickass2@example.com")
}

func TestGroupAddress(t *testing.T) {
	var group *GroupAddress = NewGroupAddress("group1")
	assert.Equal(t, group.Name(), "group1")
	var members AddressList = NewAddressList()

	members.Add(NewMailboxAddress("box1", "kickass1@example.com"))
	members.Add(NewMailboxAddress("box2", "kickass2@example.com"))

	group.SetMembers(members)
	for i := 0; i < members.GetLength(); i++ {
		assert.True(t, group.Members().Contains(members.GetAddress(i)))
	}

	unicode := NewMailboxAddress("日本語", "kickass3@example.com")
	group.AddMember(unicode)
	assert.True(t, group.Members().Contains(unicode))
}

func TestAddressList(t *testing.T) {
	var members AddressList = NewAddressList()

	first := NewMailboxAddress("box1", "kickass1@example.com")
	second := NewMailboxAddress("box2", "kickass2@example.com")
	members.Add(second)
	members.Insert(first, 0)

	mailbox, ok := members.GetAddress(0).(*MailboxAddress)
	assert.True(t, ok)
	assert.Equal(t, mailbox.Name(), "box1")
	assert.Equal(t, mailbox.Email(), "kickass1@example.com")

	mailbox, ok = members.GetAddress(1).(*MailboxAddress)
	assert.True(t, ok)
	assert.Equal(t, mailbox.Name(), "box2")
	assert.Equal(t, mailbox.Email(), "kickass2@example.com")

	unicode := NewMailboxAddress("日本語", "kickass3@example.com")
	var unicodeMembers AddressList = NewAddressList()
	unicodeMembers.Add(unicode)
	members.Append(unicodeMembers)

	raw := "box1 <kickass1@example.com>, box2 <kickass2@example.com>, 日本語 <kickass3@example.com>"
	assert.Equal(t, members.ToString(false), raw)
	encoded := "box1 <kickass1@example.com>, box2 <kickass2@example.com>, =?UTF-8?b?5pel5pys6Kqe?= <kickass3@example.com>"
	assert.Equal(t, members.ToString(true), encoded)

	parsed := ParseString(raw)
	for i := 0; i < members.GetLength(); i++ {
		parsedMailbox, ok1 := parsed.GetAddress(i).(*MailboxAddress)
		assert.True(t, ok1)
		memberMailbox, ok2 := members.GetAddress(i).(*MailboxAddress)
		assert.True(t, ok2)
		assert.Equal(t, parsedMailbox.Name(), memberMailbox.Name())
		assert.Equal(t, parsedMailbox.Email(), memberMailbox.Email())
	}

	parsed.Walk(func(a Address) {
		_, ok1 := a.(*MailboxAddress)
		assert.True(t, ok1)
	})

	parsed = ParseString(encoded)
	for i := 0; i < members.GetLength(); i++ {
		parsedMailbox, ok1 := parsed.GetAddress(i).(*MailboxAddress)
		assert.True(t, ok1)
		memberMailbox, ok2 := members.GetAddress(i).(*MailboxAddress)
		assert.True(t, ok2)
		assert.Equal(t, parsedMailbox.Name(), memberMailbox.Name())
		assert.Equal(t, parsedMailbox.Email(), memberMailbox.Email())
	}

	for i, each := range members.Slice() {
		parsedMailbox, ok1 := parsed.GetAddress(i).(*MailboxAddress)
		assert.True(t, ok1)
		memberMailbox, ok2 := each.(*MailboxAddress)
		assert.True(t, ok2)
		assert.Equal(t, parsedMailbox.Name(), memberMailbox.Name())
		assert.Equal(t, parsedMailbox.Email(), memberMailbox.Email())
	}

	replacement := NewMailboxAddress("box 4", "kickass4@example.com")
	members.SetAddress(replacement, members.IndexOf(unicode))
	members.RemoveAt(members.IndexOf(first))

	assert.False(t, members.Contains(first))
	assert.True(t, members.Contains(second))
	assert.False(t, members.Contains(unicode))
	assert.True(t, members.Contains(replacement))

	assert.False(t, members.Remove(first))
	assert.True(t, members.Remove(second))

	members.Prepend(unicodeMembers)
	assert.Equal(t, members.IndexOf(unicode), 0)
	assert.Equal(t, members.IndexOf(replacement), 1)

	members.Clear()
	assert.Equal(t, members.GetLength(), 0)
}

func TestCompositeAddress(t *testing.T) {

	// XXX This shows that an arbitrarily complex composite structure can be modeled
	var members AddressList = NewAddressList()
	members.Add(NewMailboxAddress("box1", "kickass1@example.com"))
	var subGroup *GroupAddress = NewGroupAddress("subgroup")
	subGroup.SetMembers(members)
	var subMembers AddressList = NewAddressList()
	subMembers.Add(subGroup)
	subMembers.Add(NewMailboxAddress("box2", "kickass2@example.com"))
	var superGroup *GroupAddress = NewGroupAddress("supergroup")
	superGroup.SetMembers(subMembers)
	var superMembers AddressList = NewAddressList()
	superMembers.Add(superGroup)

	arbitrarilyNested := "supergroup: subgroup: box1 <kickass1@example.com>;, box2 <kickass2@example.com>;"
	assert.Equal(t, superMembers.ToString(false), arbitrarilyNested)

	// XXX This demonstrates that type safety is enforced
	var address Address = superGroup.Members().GetAddress(0)
	mailbox, isMailbox := address.(*MailboxAddress)
	assert.False(t, isMailbox)
	assert.Nil(t, mailbox)

	group, isGroup := address.(*GroupAddress)
	assert.True(t, isGroup)
	assert.NotNil(t, group)
}
