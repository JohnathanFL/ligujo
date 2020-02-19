# Ligujo - A place that holds binds
A server for [Ligi](https://github.com/JohnathanFL/ligi) written in Go that recieves and 
stores information 
about binds in a Ligi program in a 
SQLite database and provides an interface for programs to query information about the program's 
binds.


## Protocol
Ligujo accepts 2 types of requests: requests which add information to the server, and requests 
which, well, request information from the server. Requests which add are *either* **PUT** *or* 
**POST**, and requests are always **GET**.

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
* /scope/`line`/`col`: Select the scope *as of beginning to parse* that line and column.
  * /bind/`bindName`: Look up the binding *as of that scope*.
    * Note that this could be slightly unintuitive in the case of out-of-order `let`s and the like.
  * /binds: A list of all binds that are available in this scope.
  * /parent: The object which owns this scope. That could be a Block or a File.
* /at/`line`/`col`: Select the object which *begins* at `line:col`.
  * Note this is not the symbol which *runs through* `line:col`, it is that which **begins at** 
    `line:col`.
  * This could be an identifier like `x` or `@type`, or an operator like `or` or `+`.
    * If it's an operator, its type is the type given by evaluating that expression.
* /

### **POST**
This interface is expected to be used by the interpreter/compiler to provide information about the program.

This interface shall accept the following methods:

### **GET**
