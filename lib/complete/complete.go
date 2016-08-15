// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"log"
	"unsafe"

	clang "github.com/go-clang/v3.7/clang"
)

/*
#cgo pkg-config: python3
#include <Python.h>

static inline int PyArg_ParseTuple_OLL(PyObject *args, char *file, long long *line, long long *col) {
    return PyArg_ParseTuple(args, "sLL", file, line, col);
}
*/
import "C"

//export complete
func complete(self, args *C.PyObject) *C.PyObject {
	var (
		file      C.char
		line, col C.longlong
	)
	C.PyArg_ParseTuple_OLL(args, &file, &line, &col)
	// filepy := C.PyUnicode_AsUTF8(&file)
	// log.Printf("file: %+v\n", filepy)
	log.Printf("file: %+v", C.GoString(&file))
	log.Printf("line %+v", line)
	log.Printf("col: %+v", &col)
	// return C.PyUnicode_FromObject(&file)

	index := clang.NewIndex(0, 0)
	defer index.Dispose()

	gofile := "../testdata/boost-asio_server.cpp"
	tu := index.ParseTranslationUnit(gofile, []string{"-x", "c++", "-std=c++0x", "-stdlib=libc++"}, nil, 15)
	defer tu.Dispose()

	complete := tu.CodeCompleteAt(string(gofile), uint16(line), uint16(col), nil, clang.DefaultCodeCompleteOptions())
	defer complete.Dispose()

	completeResults := complete.Results()

	var buf bytes.Buffer

	for _, r := range completeResults {
		cs := r.CompletionString()

		for i := uint16(0); i < cs.NumChunks(); i++ {
			switch cs.ChunkKind(i) {
			case clang.CompletionChunk_ResultType:
				continue
			}
			buf.WriteString(cs.ChunkText(i))
		}
		buf.Write([]byte("\n"))
	}

	gostr := C.CString(buf.String())
	defer C.free(unsafe.Pointer(gostr))

	return C.PyUnicode_DecodeUTF8(gostr, C.Py_ssize_t(buf.Len()), nil)
}

func main() {}
