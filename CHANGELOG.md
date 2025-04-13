New Features:
* You can get a version of a compiler from the WebAssembly with `getVersion` function in js. 

Fixes:
* Phantom statement at the end of AST due to incorrect parsing EOF token;
* Incorrect compilation of the `Not` operator;
* Compilation of function returning value without return in all branches;
* Crashes on the wrong indentation.
