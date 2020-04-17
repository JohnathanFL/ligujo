-- Holds info about types themselves
DROP TABLE IF EXISTS Types;
CREATE TABLE Types (
  -- ids 0..=2051 are reserved for the builtins, and will match the type field below
  id INT PRIMARY KEY, -- The typeid, as determined by the compiler
  name TEXT, -- The human name assigned to this type.
  pos TEXT, -- The position this type was made on. NULL for builtin
  -- What type of type this type is
  -- 0: void
  -- 1-1024: u1-1024
  -- 1025-2048: i1-1024
  -- 2049: usize
  -- 2050: isize
  -- 2051: bool
  -- 2052: f16
  -- 2053: f32
  -- 2054: f64
  -- 2055: f128
  -- 2056: type
  -- 2057: tuple: `len` is the number of items in the tuple. `Fields` has fields named 0..len
  -- 2058: struct: `Fields` has a field for each field
  -- 2059: enum: We have an entry in `Fields`. `contained` is the tagType
  -- 2060: slice: `contained` is the type inside the slice
  -- 2061: array: `contained` is a single type ID. len is the length of the array
  -- 2062: const: `contained` is a single type ID
  -- 2063: pure: contains typeid of a function
  -- 2064: comptime: contains typeid
  -- 2065: fn: Fields has a list of its arguments.
  type INT,
  -- Value is irrelevant(NULL) unless specified in type's comment
  contained INT,
  -- Value is irrelevant(NULL) unless specified in type's comment
  len INT
);

-- Holds statics, fields, enums, tuple members, function params/rets, etc
DROP TABLE IF EXISTS Fields;
CREATE TABLE Fields (
  -- Note that you may only have 1 field and 1 static of a name in a particular type
  -- Emphasis on the above: you can have a field x AND a static x
  -- This is because you have to access the field through instance.x, and the static through Type.x
  static INT, -- Whether it's a static (1) or a field/enum (0)
  type INT, -- The type which holds this field in it
  name TEXT, -- The name of the field. 0..len for tuples. May be `_` on a ret
  access INT, -- 0, 1, or 2 for private, readonly, or public
  contains INT, -- The type we hold in the field
  isRet INT, -- Special for fn: 1 for if this arg is the return -- Ignored in seniorSEM for now
  PRIMARY KEY (static, type, name)
);

-- Holds links between types to share the *same* things
-- If a type has an entry that conflicts with what's in here, the first has precedence
-- Thus, "MyVec2 + Vec3Ext" can 'inherit' things from MyVec2, while supersceding others
-- Note that this link denotes that toType.static1 and fromType.static1 refer to the same 
-- memory location
-- Currently unused
DROP TABLE IF EXISTS Links;
CREATE TABLE Links (
  fromType INT,
  toType INT,
  PRIMARY KEY(fromType, toType)
);

-- Finally, the table that holds info about the binds themselves.
-- Note this is intended to exclude binds that are statics
-- We don't care if it's a let/var/cvar because that's baked into its type.
-- -- John made bind specs, the type baker made em equal
DROP TABLE IF EXISTS Binds;
CREATE TABLE Binds (
  pos TEXT PRIMARY KEY, -- The position of this bind, in line:col format
  name TEXT, -- Name assigned to this bind
  type INT -- The type this bind holds (T in let x: T)
);
