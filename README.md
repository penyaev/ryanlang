Ryanlang is a toy language written in Go. It mainly pursues educational and entertainment purposes and is not meant to be used in any kind of a "serious" context. I used it to solve [2022 Advent of Code](https://adventofcode.com/2022). It's started as a way of killing some time onboard a Ryanair flight, hence the name. It is based on [Thorsten Ball's "Writing An Interpreter In Go" book](https://interpreterbook.com/).

It is a simple tree-walking interpreter and of course it's not optimized for speed or memory or anything. It has some basic type and error reporting system.

To get a taste of what code looks like in ryanlang, check out `src/algo.txt` or `src/std.txt`.

### Language features
Language is a mix of what you would typically see in Javascript and Go with a few differences.

```
// line and block comments
/* are supported */
```
#### defining variables
```
let x = 10;                            // number
let y = "string";                      // string
let sum = func (a, b) { return a+b; }; // func
```

#### functions
```
let f1 = func { return 42; };         // no arguments
let f2 = func (a, b) { return a+b; }; // with arguments
let f3 = func => 42;                  // arrow functions return the expression coming after the arrow "=>". Equivalent to "func {return 42;}"
```

#### inline functions
```
let c = (func(a, b)=>a*b)(10, 5); // 50
```

#### functions can return other functions
```
let add = func(a) => func(b) => a+b;
let add5 = add(5);
add5(10);   // 15
add(5)(10); // also 15
```

#### tuples and multi-variables contexts
```
let a, b = 10, 20;

let f = func => 10, 20; // function returns two values
let xx, yy = f();       // xx=10, yy=20
```

#### conditions
```
if true { } else { };

// can be used with the arrow expression:
let zz = if f3() > 10 => "over10" else => "below10";

// can be chained
if true {
    // ...
} else if true {
    // ...
} else {
    // ...
};
```

#### structs
```
let s = struct {
    v: 10;
    inc: func => this.v++; // "this" refers to the struct
};
```

#### maps
```
let m = map{
    123: 456;
    "a": "b";
    [1, 2, 3]: "c"; // any "hashable" type can be used as a key
};
m.a;           // "b"
m.([1, 2, 3]); // "c"
```

#### dot expressions / accessing map entries, array items, struct fields etc.
```
s.v;         // "v" is evaluated to a string, i.e. this is equivalent to s.("v")
[1, 2, 3].0; // 1: zeroth element of the array

// use brackets when you need to use an expression
m.([1, 2, 3]); // "c"
m.(100+20+3);  // 456
```

#### loops
```
for x in [1, 2, 3] { /*...*/ };       // x=1, 2, 3...
for i, v in [10, 20, 30] { /*...*/ }; // i=0, v=10; i=1, v=20; i=2, v=30;
while x > 10 { x--; };

// you can also iterate over maps
for k, v in m { /*...*/ }; // k iterates over keys, v iterates over values

// arrow expressions in for loops
let squares = for x in [1, 2, 3] => x*x; // a=[1, 4, 9]
```

#### import and exports
```
let std = import("src/std.txt");
std.max(10, 20); // 20

// exports must come first in a file
/*
exports {
    min;
    max;
};

let min = func(a, b) => if a < b => a else => b;
let max = func(a, b) => if a > b => a else => b;
*/
```

#### error reporting
if something goes wrong, ryanlang will do its best to give you an idea of what exactly was not right
```
let zzz = call();

// (interpreter output)
demo.txt:89:11: undefined identifier: call
        demo.txt:89:11: evaluating call
        demo.txt:89:11: evaluating call()
        demo.txt:89:1: evaluating let zzz = call()

```

#### unicode support
```
let すし = "🍣";
let さけ = "🍶";
println(すし + さけ); // "🍣🍶"
```
