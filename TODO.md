### general
- disallow break/continue outside of their scopes 
  - i.e. now it's possible to do a break or continue inside a function
    even though there's no enclosing loop cycle
    e.g. func { break; }
- use `StaticNull` and `StaticBool` everywhere
- arrow expressions in while loop and in for loop work differently
- not sure I understand why values are passed as copies and their outer scopes are not affected if theyre modified inside the function? thats how it should work, I just dont know why it works
  - `let x = 10; decreasetozero(x); println(x);`
- `break ("hello");` parses into wtf?
- support `return` without arguments
- `iferr()` and error handling in compiler is ugly
- implement something like resolve? i.e. break from a group, or like make something that usually resolves to nothing (null) resolve to something, and also stop execution
- `func=>return this.x++;` is executed successfully although it does not make sense
- `return return x` is possible although it does not make sense
- cyclic imports
- do not expand `+=` to asignment expr, keep it, so that it's possible to optimize `a += [i]` by running `append(a, i)` 
- ```
  this refers to wtf
  let s = struct{
    x: struct{
      a: this.b+1;
      b: 10;
    };
  };
  debugger();
  println(s.x.a);
  ```

### VM
#### language support
- [x] support break
- [ ] support continue
- [ ] support maps
- [x] support structs
- [ ] support exports for modules
- [ ] support tuples and multi-identifier initializations
- [x] support field assign
- [ ] type checks: it should not be possible to put a value of a different type into a var
- [ ] support arrow functions in loops
#### dev tools
- locations
- foreign vars name for debug
- smth like a source map? maybe via annotations
#### refactoring
- break/continue via return is probably ugly
- structs use string keys to index values, replace with numeric symbols? treat structs just as a collection of vars in a closure and thats it?
- uint8/uint16 limits: no more than uint8 args to a function etc
- object table deduplication
- heap profiling
- pre-compile modules so that we don't have to compile them on-the-fly which is expensive


### Eval