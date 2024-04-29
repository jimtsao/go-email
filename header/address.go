package header

import (
	"fmt"
	"net/mail"
	"strings"
)

type AddressField string

const (
	AddressFrom    AddressField = "From"
	AddressSender  AddressField = "Sender"
	AddressReplyTo AddressField = "Reply-To"
	AddressTo      AddressField = "To"
	AddressCc      AddressField = "Cc"
	AddressBcc     AddressField = "Bcc"
)

// Address represents an Originator or Destination Address header field
//
// usage:
//
//	addr := Address{Field: AddressFrom, Value: "alice@secret.com"}
//	addr := Address{Field: AddressTo, Value: "alice@secret.com, bob@secret.com"}
//	addr := Address{Field: AddressBcc, Value: "Eavesdrop Eve <eve@secret.com>"}
//
// note: domain literals and groups are not currently supported
//
// Syntax:
//
//	from            =   "From:" mailbox-list CRLF
//	sender          =   "Sender:" mailbox CRLF
//	reply-to        =   "Reply-To:" address-list CRLF
//	to              =   "To:" address-list CRLF
//	cc              =   "Cc:" address-list CRLF
//	bcc             =   "Bcc:" [address-list / CFWS] CRLF
//	address         =   mailbox / group
//	addr-spec       =   local-part "@" domain
//	local-part      =   dot-atom / quoted-string
//	domain          =   dot-atom / domain-literal
//	domain-literal  =   [CFWS] "[" *([FWS] dtext) [FWS] "]" [CFWS]
//	mailbox         =   name-addr / addr-spec
//	name-addr       =   [display-name] angle-addr
//	angle-addr      =   [CFWS] "<" addr-spec ">" [CFWS]
//	group           =   display-name ":" [group-list] ";" [CFWS]
//	display-name    =   phrase
//	mailbox-list    =   (mailbox *("," mailbox))
//	address-list    =   (address *("," address))
//	group-list      =   mailbox-list / CFWS
//	phrase          =   1*word
//	word            =   atom / quoted-string
//	atom            =   [CFWS] 1*atext [CFWS]
//	quoted-string   =   DQUOTE *([FWS] qcontent) [FWS] DQUOTE
//	qcontent        =   qtext / quoted-pair
//	qtext           =   %d33 / %d35-91 / %d93-126
type Address struct {
	Field AddressField
	Value string
}

func (a Address) Name() string {
	return string(a.Field)
}

func (a Address) Validate() error {
	if addrs, err := mail.ParseAddressList(a.Value); err != nil {
		return fmt.Errorf("%s: %w", a.Field, err)
	} else if a.Field == AddressSender && len(addrs) > 1 {
		return fmt.Errorf("%s: %s", a.Field, "must not contain more than 1 address")
	}
	return nil
}

func (a Address) String() string {
	var val string

	switch a.Field {
	case AddressSender:
		// single address
		if addr, err := mail.ParseAddress(a.Value); err != nil {
			val = a.Value
		} else {
			val = addr.String()
		}
	case AddressFrom, AddressReplyTo, AddressTo, AddressCc, AddressBcc:
		// multiple address
		if addrs, err := mail.ParseAddressList(a.Value); err != nil {
			val = a.Value
		} else {
			formatted := []string{}
			for _, addr := range addrs {
				formatted = append(formatted, addr.String())
			}
			val = strings.Join(formatted, ",")
		}
	default:
		val = a.Value
	}

	return fmt.Sprintf("%s: %s\r\n", a.Field, val)
}
