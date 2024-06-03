package header

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/jimtsao/go-email/folder"
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
// note: Domain literals and groups are not supported
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
//	phrase          =   1*(word / encoded-word)
//	word            =   atom / quoted-string
//	atom            =   [CFWS] 1*atext [CFWS]
//	dot-atom        =   [CFWS] dot-atom-text [CFWS]
//	dot-atom-text   =   1*atext *("." 1*atext)
//	quoted-string   =   [CFWS] DQUOTE *([FWS] qcontent) [FWS] DQUOTE [CFWS]
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
	// parse addresses
	addrs, err := mail.ParseAddressList(a.Value)
	if err != nil {
		return fmt.Errorf("%s: %w", a.Field, err)
	}

	// check sender only 1 single address
	if a.Field == AddressSender && len(addrs) > 1 {
		return fmt.Errorf("%s: %s", a.Field, "must not contain more than 1 address")
	}

	// smtp restriction: local-part max 64 octets, domain max 255 octets
	for _, addr := range addrs {
		local, domain, _ := strings.Cut(addr.Address, "@")
		if len(local) > 64 {
			return fmt.Errorf("%s: address part exceeds max length 64 bytes (%q)", a.Name(), local)
		} else if len(domain) > 255 {
			return fmt.Errorf("%s: address part exceeds max length 255 bytes (%q)", a.Name(), domain)
		}
	}

	return nil
}

func (a Address) String() string {
	sb := &strings.Builder{}
	f := folder.New(sb)
	f.Write(a.Name() + ": ")
	var fallback string

	switch a.Field {
	case AddressSender:
		// single address
		if addr, err := mail.ParseAddress(a.Value); err != nil {
			fallback = a.Value
		} else {
			a.writeAddress(addr, f)
		}
	case AddressFrom, AddressReplyTo, AddressTo, AddressCc, AddressBcc:
		// multiple address
		if addrs, err := mail.ParseAddressList(a.Value); err != nil {
			fallback = a.Value
		} else {
			for i := 0; i < len(addrs); i++ {
				if i > 0 {
					f.Write(",")
				}
				a.writeAddress(addrs[i], f)
			}
		}
	default:
		fallback = a.Value
	}

	if fallback != "" || f.Err != nil {
		return fmt.Sprintf("%s: %s\r\n", a.Field, fallback)
	}

	f.Close()
	return sb.String()
}

func (a Address) writeAddress(addr *mail.Address, f *folder.Folder) {
	// net/mail address ouputs 'quoted-string angle-addr' or 'encoded-word angle-addr' format
	var q, e, d string

	// local part
	if addr.Name != "" {
		local := (&mail.Address{Name: addr.Name}).String()
		local = local[:len(local)-len(" <@>")]
		if len(local) > 8 && local[:8] == "=?utf-8?" {
			e = local
		} else {
			q = local
		}
	}

	// domain part
	if addr.Address != "" {
		d = (&mail.Address{Address: addr.Address}).String()
	}

	// write to encoder
	if q != "" {
		// quoted string: [CFWS] DQUOTE *([FWS] qcontent) [FWS] DQUOTE [CFWS]
		// format: [1]quoted[3][space]string[2][space]angle-addr[1]
		f.Write(1)
		for idx, qp := range strings.Split(q, " ") {
			if idx > 0 {
				f.Write(3, " ")
			}
			f.Write(qp)
		}
		f.Write(2, " ", d, 1)
	} else if e != "" {
		// encoded-word: [CFWS] between encoded words (whereby upon parsing space is ignored)
		// format: encoded-word[2][space]angle-addr[1]
		f.Write(e, 2, " ", d, 1)
	} else {
		// angle-addr: [CFWS] "<" local @ domain ">" [CFWS]
		// format: [1]<addr-spec>[1]
		f.Write(1, d, 1)
	}
}
