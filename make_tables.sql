-- Holds info about types themselves
DROP TABLE IF EXISTS Types;
CREATE TABLE Types (
  id INT PRIMARY KEY, -- The typeid, as determined by the compiler
  name TEXT, -- The name assigned to this type.
  pos TEXT, -- The position this type was made on. NULL for builtin
  -- What type of type this type is
  -- 0: Void
  -- 1-1024: u1-1024
  -- 1025-2048: i1-1024
  -- 2049: usize
  -- 2050: isize
  -- 2051: bool
  -- The following 3 make use of the "contained" field to specify what's inside them
  -- 2052: tuple: `len` is the number of items in the tuple. `Fields` has fields named 0..len
  -- 2053: struct: `Fields` has a field for each field
  -- 2054: enum: We have an entry in `Fields`
  -- 2055: slice: `contained` is the type inside the slice
  -- 2056: array: `contained` is a single type ID. len is the length of the array
  -- 2057: const: `contained` is a single type ID
  -- 2058: pure: contains typeid of a function
  -- 2059: comptime: contains typeid
  -- 2060: fn: Fields has a list of its arguments. 
  type INT,
  contained INT,
  len INT,
);

DROP TABLE IF EXISTS Statics;
CREATE TABLE Statics (
  -- The fact that these 2 are prim means we can only have one static of a name in a type
  -- The type the static is
  type INT PRIMARY KEY,
  name TEXT PRIMARY KEY,
  contains INT, -- The type this static holds
);


-- Holds fields, enums, tuple members, function params/rets, etc
DROP TABLE IF EXISTS Fields;
CREATE TABLE Fields (
  -- Note that you may only have 1 field of a name in a type
  type INT PRIMARY KEY, -- The type which holds this field
  name TEXT PRIMARY KEY, -- The name of the field. 0..len for tuples. May be `_` on a ret
  contains INT, -- The type we hold in the field
  isRet INT, -- Special for fn: 1 for if this arg is the return
);

-- Finally, the table that holds info about the binds themselves.
-- Note this is intended to exclude binds that are statics
-- We don't care if it's a let/var/cvar because that's baked into its type.
-- -- John made bind specs, the type baker made em equal
DROP TABLE IF EXISTS BINDS;
CREATE TABLE Binds (
  pos TEXT PRIMARY KEY, -- The position of this bind, in line:col format
  name TEXT, -- Name assigned to this bind
  type INT, -- The type this bind holds
);
