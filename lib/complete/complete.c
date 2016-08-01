#include <Python.h>

PyObject *complete(PyObject *, PyObject *);

static PyMethodDef GoLibclangMethods[]
    = {{"complete", complete, METH_VARARGS, "completion engine for CPython."},
       {NULL, NULL, 0, NULL}};

static struct PyModuleDef golibclang_module
    = {PyModuleDef_HEAD_INIT, "golibclang", NULL, -1, GoLibclangMethods};

PyMODINIT_FUNC PyInit_golibclang(void)
{
  return PyModule_Create(&golibclang_module);
}
