<h1>defimpl</h1>

My background is as a Lisp programmer.  Objects in my programs model
objects in the real world.  Real world objects have identity and so do
the software objects which model them.  If I dent a car, a specific car
gets dented.  To achieve this requires reference rather than value/copy
semantics.

Since I find the interface versus struct reference distinction
confusing and would rather not think about it, I've chosen to define
library interfaces that operate on interface objects.  Also, I find
that I often refactor to pass interfaces rather than structs, so why
not just always start that way.

Each concrete interface requires at least one implementation though,
and those implementations tend to be boilerplate code that is boring
and tedious to implement.  I'd like to automatically generate a
canonical implementation of each interface based on some direction in
the interface's type definition.

Though Go allows fields in struct definitions to have a tag which
various libraies can choose to interpret in some way, this is
apparently not allowed for interface fields (method declarations).

The defimpl binary is a go code preprocessor which reads a source
file, and produces a new source file which contains a struct
definition to serve as an implementation of each interface definition
in the input file, and the methods to support that implementation.

Each field (method declaration) in an input interface definition can
have comments to inform defimpl of the struct methods to be generated
and the struct fields on which they operate.  These are comments
because interface type definitions do not support tags.  These
comments employ the same canonical syntax as struct tags so that we
can leverage the struct tag parsing code.

For any interface method which is meant to read or modify some slot,
that method sould have a signature appropriate to its intended use and
a comment of the form

<pre>
	// defimpl:"verb fieldname"
</pre>

where "verb" is one of the supported verb types, e.g.: read, set,
append, iterate; and fieldname is the name of the struct field on
which the method operates (performs the verb).
