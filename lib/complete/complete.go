// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"sync"
	"unsafe"

	"github.com/davecgh/go-spew/spew"
	clang "github.com/go-clang/v3.8/clang"
)

/*
#cgo pkg-config: python3
#include <Python.h>

static inline int PyArg_ParseTuple_LL(PyObject *args, long long *a, long long *b) {
    return PyArg_ParseTuple(args, "LL", a, b);
}
*/
import "C"

//export complete
func complete(self, args *C.PyObject) *C.PyObject {
	var line, col C.longlong
	C.PyArg_ParseTuple_LL(args, &line, &col)

	index := clang.NewIndex(0, 0)
	defer index.Dispose()

	file := "../testdata/boost-asio_server.cpp"
	tu := index.ParseTranslationUnit(file, []string{"-x", "c++", "-std=c++0x", "-stdlib=libc++"}, nil, 15)
	defer tu.Dispose()

	complete := tu.CodeCompleteAt(string(file), uint32(line), uint32(col), nil, clang.DefaultCodeCompleteOptions())
	defer complete.Dispose()

	completeResults := complete.Results()

	var buf bytes.Buffer

	var (
		ch = make(chan []byte, int(complete.NumResults()))
		wg sync.WaitGroup
	)

	wg.Add(int(complete.NumResults()))
	for _, r := range completeResults {
		cs := r.CompletionString()

		go func() {
			defer wg.Done()

			var bu bytes.Buffer
			for i := uint32(0); i < cs.NumChunks(); i++ {
				switch cs.ChunkKind(i) {
				case clang.CompletionChunk_ResultType:
					continue
				}
				bu.WriteString(cs.ChunkText(i))
			}

			bu.Write([]byte("\n"))
			ch <- bu.Bytes()
		}()
	}

	wg.Wait()
	log.Printf("\x1b[34;40mch\x1b[0m:\n%+v\n", spew.Sdump(ch))
	var i int
	for b := range ch {
		buf.Write(b)
		i++

		if i > len(ch) {
			break
		}
	}

	gostr := C.CString(buf.String())
	defer C.free(unsafe.Pointer(gostr))

	return C.PyUnicode_FromString(gostr)
}

func main() {}
