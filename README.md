![linux](https://github.com/nikonru/NiLang/actions/workflows/linux.yml/badge.svg)
![windows](https://github.com/nikonru/NiLang/actions/workflows/windows.yml/badge.svg)
# Why NiLang?
The part **NiLa** in the name stands for *Larry Niven* beloved author 
of the famous sci-fi book series *Ringworld*. 
The last two letters **ng** are reference to ISO 639 language code of *Ndonga* - one of the vibrant languages spoken in the Sub-Saharan Africa.

# What's this?
NiLang (Russian: НиЛанг) is a high level language for programming a bot from [TorLand](https://github.com/Slava2001/TorLand).
## File extension
We recommend use `.nil` as a file extension for files containing **NiLang** source code.
# Syntax
```
Using bot

Bool hungry = True
While hungry:
    Int maxEnergy = 1500
    ConsumeSunlight
    If GetEnergy > maxEnergy:
        Fork$ world::Forward
        hungry = False
Dir dir = GetDir
# you may remove `bot::`, since we have already written 'using bot'
bot::Move$ RotateClockwise$ dir  
```

```
Alias Numbers::Int:
    one = 1
    two = 2
    four = 4

Numbers four = numbers::four

Scope names:
    max = 1000

Fun F::Bool$max Int, default Bool:
    Using bot
    ConsumeSunlight
    If GetEnergy > max:
        Return True
    Elif GetEnergy < max:
        Return False
    Else:
        Return default

Fun G::Int:
    Return 5

Fun W:
    ConsumeSunlight

W
Bool flag = F$ 5, False
flag = F$ G, True
```

# Ideas for the future improvements
Here is the list of ideas to implement in the future versions of NiLang. 
The Syntax might be rough and not really compatible with the current version of language.
## Multiple return values
Sometimes it might be really useful to return some value and error or success code, which describes validity of this value.
```
Fun F::Int, Int:
    Return 2, 3

Int x = 0
Int y = 1

x, y = F # x = 2, y = 3
```
## Builtin type conversion function
Writing type conversion functions might be cumbersome, so having such utility built in the language is nice.
```
Bool x = Bool$ 1   # x = True
Int y = Int$ False # y = 0
```
Though for Aliases we have to somehow check whether conversion was successful or not.
```
Alias Error::Int:
    forbidden = 403
    notFound = 404

Error x = error::forbidden

Bool ok = False
x, ok = Error$ 405 # ok = False
x, ok = Error$ 404 # ok = True, x = error::notFound
```
## Lambdas and functions as first-class citizens
Passing function as an argument to another function gives nice functional programming vibes.
```
Fun F$y Dir, Do Fun$x Dir:
    Do$ y

F$ dir::front, Lambda$ z Dir: Move$ z # Move$ dir::front
```
## Simple aliases
Types can have very long names, which may make your code to exceed character limit per line. 
Simple aliases give you an option to shorten your lines of code. 
```
Alias Direction = Dir
Alias Integer = Int

Scope x:
    Scope y:
       Alias Status::Int:
          ok = 1
          bad = 2

Alias Status = x::y::Status
```
## New builtin types
Currently **NiLang** is very boring language, more useful builtin types can solve it.
```
Char x = 'a'
Float y = 10.1
Uint z = 10
```
## Arrays and strings
You can't have a proper programming language without a way to describe a continuous piece of memory. 
Though it opens an interesting question regarding how we should pass it to the function or during assignment (copy or reference).
```
Array::Int a = 1, 2, 3, 4            # initialized array os size 4
Array::Int b = Array$ 4              # uninitialized array os size 4
Array::Char c = "Hello, world!"      # string
Array::Char h = 'C', 'h', 'a', 'r'   # string
String s = "Hello, world!"           # string
Array::Bool d = False, True, True    # initialized array os size 3

Int x = a!0 # x = 1
x = a!1     # x = 2
x = a!10    # undefined behaviour
```
## Generics
At some point of its life any statically typed language needs a way to optimize repeating code. 
One possible solution is an introduction of generics similar to templates in *C++*.
```
Fun$T: F::Bool$x T:
    If x > 10:
        Return True
    Return False

Bool x = F$ 10
x = F$ 10.1
```
## Short assignment
Just a syntax sugar for the economy of characters in a source code.
```
Int x = 0
x += 1 # x = x + 1
x -= 1 # x = x - 1
x *= 2 # x = x * 2
x /= 4 # x = x / 4
```
## For loops
Another way to save few lines of code, while going through an array.
```
For Int i = 0, i < 10, i += 1:
    Move$ dir::front

Array::Dir a = dir::front, dir::left, dir::back 
For direction$a:
    Move$ direction
```
## Objects
Structures, classes, custom data types are beloved in OOP programming paradigm, 
which should be supported if **NiLang** wants to be multi-paradigmatic language.
```
Object Car:
    Int wheels
    String brand

    Fun Car::Car$ numOfWheels Int:
        self.wheels = numOfWheels
        self.brand = "default"
    Fun GetWheels::Int:
        Return self.wheels

Car myCar = Car$ 4

myCar.brand = "Toyota"
Int x = myCar.GetWheels
Bool y = x == myCar.wheels # True
```
## Imports
Writing the whole program in one file can be quite cumbersome. 
Thus having some way to spread code among different files is must have.

Another useful feature can be maintaining import structure, which is parallel 
to the arrangement of files in file system.
```
#src/helpers.nil
Domain helpers

Int x = 10
```
```
#src/main.nil
Domain main

Using helpers 

Int y = helpers.x
```
Aliases for imports can be also useful for importing some code with similar names.
```
#src/main2.nil
Domain main2

Using helpers = hp 

Int y = hp.x
```
Targeted import also helps to maintain clarity of expression.
```
#src/main3.nil
Domain main3

Using helpers::x

Int y = x
```
```
#src/main4.nil
Domain main4

Using helpers::x = z

Int y = z
```
