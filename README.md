# unipdf-issue-insert-svg

Reproducible Go program that demonstrates [`unipdf`](https://github.com/unidoc/unipdf) issues with SVG images.

We tested with Go version `1.23.0` and UniPDF versions `v3.62.0` and `v3.63.0` (see the `go.mod` file).

We can demonstrate three different outcomes when attempting to add an SVG image to an existing PDF document:

1. Adding the SVG image succeeds without an error and the image is actually rendered in the document.
2. Adding the SVG image does not result in an error, but the image is not rendered in the document.
3. Adding the SVG image crashes the program.

## Reproduce

- set the `UNIPDF_LICENSE_KEY` and `UNIPDF_CUSTOMER_NAME` environment variables

- `$ go get ./...`
- `$ go run main.go ./input.pdf`

- see success/error log messages written to stdout
  - note that some test cases result in a panic

- see generated PDF files in `./out` for reference
  - note that some of the generated PDF files that
  do not print an error to stdout are actually blank

## Analysis

We have reason to believe that outcome (2) is caused whenever a `<g>` element with a `transform` attribute is present in the SVG file.

Regarding outcome (3), we have observed the following runtime errors:

- `runtime error: index out of range [0] with length 0`
- `runtime error: invalid memory address or nil pointer dereference`
