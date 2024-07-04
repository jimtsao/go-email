# Go-email

Craft MIME comformant emails in Go

## Install

`go get github.com/jimtsao/go-email`

## Usage

General usage

```go
m := goemail.New()
m.From = "alice@example.com"
m.To = `Bob <bob@example.com>`
m.Bcc = `"Eve the eavesdropper" <eve@example.com>`
m.Subject = "Secret Plans"
m.Body = "<b>Attack at Dawn!</b>"
m.AddHeader(header.MessageID("<local@host.com>"))

// header validation
err := m.Validate()
if len(err) > 0 {
    // handle errors
}

raw := m.Raw()
// use in gmail api, aws ses etc
```

Custom email

```go
headers := []header.Header{
    header.MIMEVersion{},
    header.Address{Field: header.AddressFrom, Value: "alice@example.com"},
    header.Address{Field: header.AddressTo, Value: "Bob <bob@example.com>"},
    header.Subject("Secret Plan"),
    header.MessageID("<local@host.com>"),
    // insert own header that satisfies header.Header interface
}

text := mime.NewEntity(
    []header.Header{header.NewContentType(
        "text/plain",
        header.NewMIMEParams("charset", "us-ascii"),
    ),
    }, "Attack at Dawn!")

html := mime.NewEntity(
    []header.Header{header.NewContentType(
        "text/html",
        header.NewMIMEParams("charset", "utf-8"),
    ),
    }, "<b>Attack at Dawn!</b>")

alt := mime.NewMultipartAlternative(headers, []*mime.Entity{text, html})

// header validation
for _, h := range alt.Headers {
    if err := h.Validate(); err != nil {
        // handle error
    }
}

// generate raw email output, use in gmail api, aws ses etc.
raw := alt.String()
```

## Features

General

- [x] email header validation
- [x] non us-ascii support for header and mime parameter values
- [x] non us-ascii support for email body

Folding

- [x] header (78 octet limit)
- [x] RFC 2045 base64 (76 octet limit)
- [x] support for folding priority

Highly customisable and extensible

- [x] header.Header interface
- [x] folder.Foldable interface
- [x] standalone syntax checking library
- [x] standalone folding library

## Relevant Documents

Package makes best effort to conform to following standards:

- [RFC 5322](https://datatracker.ietf.org/doc/html/rfc5322) — Internet Message Format. Specification for emails.
- [RFC 2045](https://datatracker.ietf.org/doc/html/rfc2045) — MIME Part 1: Message Body Formatting. Supports non-ascii body including attachments.
- [RFC 2046](https://datatracker.ietf.org/doc/html/rfc2046) — MIME Part 2: Media Types. Text, image, video, audio, application, multipart and message.
- [RFC 2047](https://datatracker.ietf.org/doc/html/rfc2047) — MIME Part 3: Message Header Formatting. Supports non-ascii headers.
- [RFC 2049](https://datatracker.ietf.org/doc/html/rfc2049) — MIME Part 5: MIME conformance.
- [RFC 2231](https://datatracker.ietf.org/doc/html/rfc2231) — MIME Parameter Value and Encoded Word Extensions. Supports non-ascii header parameters.
- [RFC 2183](https://datatracker.ietf.org/doc/html/rfc2183) — Communicating Presentation Information in Internet Messages: The Content-Disposition Header Field.
- [RFC 5321](https://datatracker.ietf.org/doc/html/rfc5321) — Simple Mail Transfer Protocol. Imposes some length limits on various parts of message.
