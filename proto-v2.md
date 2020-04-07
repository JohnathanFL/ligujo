# Protocol - Second attempt

The protocol is split into two parts: Compiler methods and Editor methods. Each is meant to be used
only by compilers and editors to update information and retrieve information, respectively.

Ligujo shall not distinguish between PUT and POST, but this documentation shall always say PUT.

The following two sections are formatted as follows:
* *(GET|PUT) URL*
  * *description and argument specification*

When a url must be followed by something else, it is formatted as follows:
* *(GET|PUT) /something/...*
  * *else*/...
    * *more*/...

## Type Object
* `name`: The name of the type as specified in the program
  * For anonymous structs, the name is of the form `@struct_line_col`
  * Note that this is not unique
* `pos`: The position this type was bound at, in `line:col` format
* `statics`: An object of the static variables and functions in this type. Each is:
  * `pub`: Either `true` or `false`
  * `typeid`: The type of that member
    This style was chosen to solve recursive structures in a single request.
  * `loc`:
    * `line`: The line at which this was bound
    * `col`: The column at which this was bound
* If the type is not a primitive, then one of the following 3 members is also added:
  * `struct`: An array where each field is the same type as `statics`
    * Note that this is an array, while `statics` is an object. This is because order is only 
      important for fields.
  * `enum`:
    * `tagType`: Either `i*` or `u*`. The primitive type that stores the tag.
    * `tags`: An array of objects, one for each discriminator
      * `name`: The name of this tag
      * `contains`: Either non-existent or a typeID to denote what the union contains here.
      * `val`: The *numeric* value of the tag itself. (`1`, `2`, `-5`, etc)
  * `tuple`: An array of typeIDs.

## Type List
* `primary`: The typeID of the bind found at *line*:*col*
* A series of keys->type objects where each key is a typeID.

## Compiler Methods
* PUT /mktype
  * Registers one or more new types with Ligujo
  * Argument is a JSON object where each key is the typeID for that type and each value is
    a valid type object as outlined in the [Type Object](#Type object) section
    * TypeIDs are managed by the compiler, *not* Ligujo.
* PUT /at/*line*:*col*
  * Registers a binding at *line*:*col*
  * Ligujo takes care of figuring out which scope this is in and registering it there
  * Argument is a JSON object as follows:
    * `name`: The name of the bind itself
    * `pub`: `true` or `false` for whether the variable is public
    * `type`: The typeID for the type of that bind


## Editor Methods
* GET /at/*line*:*col*
  * Returns a JSON object formatted as follows
    * `primary`: The typeID of the bind found at *line*:*col*
    * A series of keys->type objects where each key is a typeID.
