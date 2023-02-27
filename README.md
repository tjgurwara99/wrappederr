# wrappederr is a simple program to integrate with `go vet` checks

It statically analyses all functions which include error in their return type and check whether these errors are wrapped using the `errors.Wrap` function.

To use it, install the program and run `go vet` like so:

```bash
go vet --vettool=$(which wrappederr) ./...
```

NOTE: This is a POC where there is a possibility that the error is indeed returned wrapped from another
function but it unwrapped in the current function, causing it to display false positives.