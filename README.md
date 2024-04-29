# Go-email

Craft MIME comformant emails in Go

## Relevant Documents

Package makes best effort to conform to following standards:

- [RFC 5322](https://datatracker.ietf.org/doc/html/rfc5322) — Internet Message Format. Specification for Email.
- [RFC 2045](https://datatracker.ietf.org/doc/html/rfc2045) — MIME Part 1: Message Body Formatting. Extends email support for non-ASCII body data including attachments.
- [RFC 2046](https://datatracker.ietf.org/doc/html/rfc2046) — MIME Part 2: Media Types. Defines MIME media types eg text, image, audio, multipart etc.
- [RFC 2047](https://datatracker.ietf.org/doc/html/rfc2047) — MIME Part 3: Message Header Formatting. Extends email support for non-ASCII header data.
- [RFC 2049](https://datatracker.ietf.org/doc/html/rfc2049) — MIME Part 5: MIME conformance. Minimal standard to be considered MIME comformant.
- [RFC 2231](https://datatracker.ietf.org/doc/html/rfc2231) — MIME Parameter Value and Encoded Word Extensions. Defines extensions that address limitations in previous MIME RFCs. These mechanisms should be used only when necessary.
- [RFC 5321](https://datatracker.ietf.org/doc/html/rfc5321) — Simple Mail Transfer Protocol

Notes
Emails are considered within the context of SMTP, specfically:

- [`Local-part`](https://datatracker.ietf.org/doc/html/rfc5321#section-4.5.3.1.1) part of address maximum length is 64 octets
- [`Domain`](https://datatracker.ietf.org/doc/html/rfc5321#section-4.5.3.1.1) part of address maximum length is 255 octets

Message validation will return an error if these limits are exceeded

## Scope

This package is primarily used for crafting new email messages in the context of an end user. It is not meant for use by Mail User Agents (MUA) or Mail Transfer Agents (MTA). As such, conformity to standards is prioritised towards common usage.
