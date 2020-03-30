# Ligujo - A place that holds binds
A server for [Ligi](https://github.com/JohnathanFL/ligi) written in Go that recieves and 
stores information 
about binds in a Ligi program in a 
SQLite database and provides an interface for programs to query information about the program's 
binds.


## Protocol
Ligujo accepts 2 types of requests: requests which add information to the server, and requests 
which, well, request information from the server. Requests which add are *either* **PUT** *or* 
**POST**, and requests are always **GET**. This document will always say put for
compiler -> typeserver, and get for editor -> typeserver.

Request methods are broken down into two parts: location and member. Location specifies which 
symbol or expression we are accessing, and member specifies which part of that symbol or 
expression we wish to access.

For example, to access the type of the variable which begins at line `20` column `4`:
```
GET /at/20/4/type
```
Here `/at/20/4` is the location and `/type` is the member we wish to get.

### Methods
Things to be filled by the requester are in `code format`.

The following is the list of *location specifiers*:
* /bind/`bindID`: Fast way to get to a bind, rather than doing the complete /scope/... or similar.
* /type/`typeID`: Fast way to get info about a type if you already know its ID.
* /scope/`line`/`col`: Select the scope *as of beginning to parse* that line and column.
  * Note that this scope *includes all scopes that occured before it*. For example, the scope of 
    `main` will include the scope of its file.
  * /bind/`bindName`: Look up the binding *as of that scope*.
    * Note that this could be slightly unintuitive in the case of out-of-order `let`s and the like.
  * /binds: A list of all binds that are available in this scope.
    * GET only
  * /parent: The object which owns this scope. That could be a Block or a File.
* /at/`line`/`col`: Select the bind which *begins* at `line:col`.
  * Note this is not the symbol which *runs through* `line:col`, it is that which **begins at** 
    `line:col`.
  * Thus to select `x` in the below file, you would do `/at/2/5`
    ```
    let y = 51
    let x = 10
    ```

The following is the list of *member specifiers*:
* /type: Get the type of that location.
  * GET returns a JSON object structured as follows:
    * A series of objects where the key is a typeID taken either from `primary`
      or from a type found inside 
  * PUT's input is a JSON array of type objects.
    * TypeIDs in PUTs are local to that PUT. (i.e have no relation to the TypeIDs in the DB).
      They are only there for the compiler to indicate to the compiler which is which when
      the types are recursive.
* /comment: Get doc comment (if any) on that location.
  * GET returns plain text. This text may be empty if there is no comment.
  * PUT's input is plain text and may be empty if there is no comment.

The bodies of most gets and puts are structured like so:
* `primary`: The typeID of the object explicitly requested.
  * Only found in GETs.
* A series of objects named for the typeID found either in `primary` or an inner `typeid` field.
  * `name`: The name of the type as specified in the program
    * For anonymous structs, the name is of the form `@struct_line_col`
  * `statics`: An object of the static variables and functions in this type. Each is:
    * `pub`: Either `true` or `false`
    * `name`: The name of that member
    * `typeid`: The type of that member
      This style was chosen to solve recursive structures in a single request.
    * `loc`:
      * `line`: The line at which this was bound
      * `col`: The column at which this was bound
  * `kind`: One of the primitive types or `struct`, `enum`, or `tuple`.
  * If `kind` is not primitive, then one of the following 3 members is also added:
    * `struct`: An array where each field is the same type as `statics`
      * Note that this is an array, while `statics` is an object. This is because order is only 
        important for fields.
    * `enum`:
      * `tagType`: Either `i*` or `u*`. The primitive type that stores the tag.
      * `tags`: An array of objects, one for each discriminator
        * `name`: The name of this tag
        * `contains`: Either empty or a typeID to denote what the union contains here.
        * `val`: The *numeric* value of the tag itself. (`1`, `2`, `-5`, etc)
    * `tuple`: An array of typeIDs.


### Examples
The compiler wants to inform Ligujo of the type of `x` in the following file
