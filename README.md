# Go-email

Craft MIME comformant emails in Go

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

## Scope

This package is primarily used for crafting new email messages in the context of an end user. It is not meant for use by Mail User Agents (MUA) or Mail Transfer Agents (MTA). As such, conformity to standards is prioritised towards common usage.
