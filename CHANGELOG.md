New Features:
* Arithmetic operators

Fixes:
* Stack overflow error provoked by scope resolution operator
* Compilation with invalid indentation
* Print of `If` statement in AST form
* Putting a comment in the middle don't compile
* Line with only white spaces don't compile
* Comparison of non-integers
* Compiling function with all `Return` statements inside of an `If` statements
* `While` statement at the end of a function, which leads to the compiler escaping function earlier than it should
