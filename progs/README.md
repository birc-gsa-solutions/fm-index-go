# Executables

Put the source for your executables in this directory, in one sub-directory per executable. To build executables, and putting the result at the root of the repository, you can run

```sh
> GOBIN=$PWD go install ./...
```

in the root of the repository. This will build all the Go code in the repository and put the binary executables at the root, ready to be tested by the test framework.
