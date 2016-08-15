go-libclang
===========

C completion library written in Go.  
Currentpy not works.

## Required

On GNU/Linux

```sh
apt-get install llvm-dev clang libclang-dev
```

## Build
go-libclang depends [go-clang/v3.7](https://github.com/go-clang/v3.7)

```sh
CGO_LDFLAGS="-L`llvm-config --libdir`" go get -u -v -x github.com/go-clang/v3.7/...
```

fetch the go-libclang.

```sh
go get -u -v -x github.com/zchee/go-libclang/...
```

building.

```sh
make
```

test use boost-asio example code.

```sh
make run
```
