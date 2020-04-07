#!/bin/dash
sqlite3 'type.db' < make_tables.sql

# Now we generate based on the following layout:
#  # 0-4096 are reserved for compiler types
#  # 0: void
#  # 1-1024: u1-1024
#  # 1025-2048: i1-1024
#  # 2049: usize
#  # 2050: isize
#  # 2051: f16
#  # 2052: f32
#  # 2053: f64
#  # 2054: f128
#  # 2055: bool
#  # 2056: type
#    # All of the above are builtin and are in the DB as described above
#    # All of the below are used to make each *new* type
#  # The following 3 make use of the "contained" field to specify what's inside them
#  # 2057: tuple: `len` is the number of items in the tuple. `Fields` has fields named 0..len
#  # 2058: struct: `Fields` has a field for each field
#  # 2059: enum: We have an entry in `Fields`
#  # 2060: slice: `contained` is the type inside the slice
#  # 2061: array: `contained` is a single type ID. len is the length of the array
#  # 2062: const: `contained` is a single type ID
#  # 2063: pure: contains typeid of a function
#  # 2064: comptime: contains typeid
#  # 2065: fn: Fields has a list of its arguments.

addType() {
  # Function arguments are the fields below, in the same order
  echo "INSERT INTO Types (id, name, pos, type, contained, len) VALUES ($1, '$2', $3, $4, $5, $6);" | sqlite3 'type.db'
}

# For when you don't need pos, contained, or len and type==id
addBasicType() {
  # $1: TypeID
  # $2: TypeName
  addType $1 $2 NULL $1 NULL NULL
}

addType 0 "void" "NULL" 0 NULL NULL

for i in $(seq 1 1024); do
  # The unsigned form
  addBasicType $i "u$i"
  si=$((i + 1024))
  # The signed form is just an offset of 1024
  addBasicType $si "i$i"
done

addType 2049 "usize"
addBasicType 2050 "isize"
addBasicType 2051 "bool"
addBasicType 2052 "type"
