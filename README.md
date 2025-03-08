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
## Usage of compiler
Simply pass your source code file as the first argument to the compiler and it will generate 
file with extension `.tor` and `botlang` in it.
```
$./nilang bot.nil
```
If you wish to use more options use `--help` flag to get their names and descriptions.
```
$./nilang ---help
```
# The hitchhiker's guide to NiLang
## Rule №1
**No brackets are allowed.**
## Comments
Character `#` is used to define the beginning of the comment, the part of code which is 
completely ignored by the compiler. Once comment has begun it ends only at the end of line
```
Int variable = 0 #this commentary goes until the end of line
Bool anotherVariable = False
# this comment covers the whole line 
```
## Variable declaration
**NiLang** is **statically typed** language, hence all variables must be declared and initialized first.

Currently the language supports 3 builtin types:
* `Int` - representing integer numbers;
* `Bool` - representing `True` and `False`;
* `Dir` - representing possible directions, which can be used to control a bot and accessed via `dir` alias:
   * `front`;
   * `frontRight`;
   * `frontLeft`;
   * `right`;
   * `left`;
   * `back`;
   * `backRight`;
   * `backLeft`.
```
Int integer = 10
Bool truth = True
Bool lie = False
Dir direction = dir::front
```
Did you notice something? 
Variables in **NiLang** are very humble and **must have a name, which begins with a lower case letter**.
On the opposite **types are very proud and always begins with an upper case letter**.

The following code won't compile.
```
Int Bad = 10
```
But in the rest, naming of variables is pretty similar to other programming languages.
All names down below are valid in **NiLang**.
```
Int _val = 0
Int val_1 = 1
Int aBcD1234 = 2
Int _a_B_c_D_1_2_3__4_ = 3
```


## Variable assignment
After declaration you are free to use assign variables as you wish
```
Int x = 10
Int y = 16
x = 12
y = x
```
But keep in mind that variable and value you would like to assign to it must have the same type.

The following code won't compile.
```
Bool val = True
Int x = 1
val = x
```
## Operators
But simple assignment is boring, thus you can use different operators in your code 
to do logic and calculation.
```
Int x = 1
Int y = 2

Bool isEqual = x == y
Bool isGreater = x > y
Bool isGreaterOrEqual = isEqual Or isGreater
```
Down below is the complete list of all supported operators in **NiLang**.
### Logical
* `And` - the operator returns `True` when both the conditions in consideration are satisfied. Otherwise it returns `False`. For example, `x And y` returns `True` when both `x` and `y` are `True`;
* `Or` - the operator returns `True` when one (or both) of the conditions in consideration is satisfied. Otherwise it returns `False`. 
For example, `x Or y` returns `True` if one of `x` or `y` is `True`. Of course, it returns `True` when both a and b are `True`.
* `Not` - the operator returns `True` the condition in consideration is not satisfied. Otherwise it returns `False`. For example, `Not x` returns `True` if `x` is `False`.
### Comparison
* `==` - (Equal To) operator checks whether the two given operands are equal or not,
If so, it returns `True`. Otherwise, it returns `False`. For example, `5==5` will return `True`.
* `!=` - (Not Equal To) operator checks whether the two given operands are equal or not. If not, it returns `True`. Otherwise, it returns `False`. It is the exact boolean complement of the `==` operator. For example, `5!=5` will return `False`.
* `>` - (Greater Than) operator checks whether the first operand is greater than the second operand. If so, it returns `True`. Otherwise, it returns `False`. For example, `6>5` will return `True`.
* `<` - (Less Than) operator checks whether the first operand is lesser than the second operand. If so, it returns `True`. Otherwise, it returns `False`. For example, `6<5` will return `False`.
* `>=` - (Greater Than Equal To) operator checks whether the first operand is greater than or equal to the second operand. If so, it returns `True`. Otherwise, it returns `False`. For example, `5>=5` will return `True`.
* `<=` - (Less Than Equal To) operator checks whether the first operand is lesser than or equal to the second operand. If so, it returns `True`. Otherwise, it returns `False`. For example, `5<=5` will also return `True`.
### Arithmetic 
* `+` - (Addition) operator adds up two numbers. For example, `x = 5 + 2` will write value `7` to the variable `x`.
* Infix `-` - (Subtraction) operator subtracts the second number from the first one. For example, `x = 5 - 2` will write value `3` to the variable `x`.
* Prefix `-` - (Negation) operator negates a number. For example, `x = - 5` will write value `-5` to the variable `x`.
* `*` - (Multiplication) operator multiplies two numbers. For example, `x = 5 * 2` will write value `10` to the variable `x`.
* `/` - (integer Division) operator divides the first number from the second one dropping the reminder. For example, `x = 5 / 2` will write value `2` to the variable `x`. Another example, `y = 9 / 3` will write value `3`
to the variable `y`.
* `**` - (Power) operator raises the first number to a power equal to the second one. For example, `x = 5 ** 3` will write value `125` to the variable `x`.
### Precedence
Operators are applied in the following order, starting from the highest:
* `::` - Scope resolution (see Scopes)
* Function call (see Functions)
* `**`
* `Not`, unary `-`
* `*`, `/`
* `+`, `-`
* `>=`, `>`, `<`, `<=`
* `==`
* `And`, `Or`
## Conditional statements
What **if** you have to make a decision based on some condition? Exactly for such case 
there is *if-statement* in **NiLang**.
```
If 10 > 8:
    x = 10
```
Instead of cumbersome brackets **NiLang** uses elegant indentations to define code, which would be executed 
in case of fulfilling condition after keyword `If`.

Beware, **NiLang** allows **usage only of 4-space long indentations**, any tabulations or shorter/longer indentations in code would lead to compilation error.

By using the keyword `Else` you can specify what to do in the alternative case.
```
If 8 > 10:
    x = 8
Else:
    x = 10
```
And as a way to avoid cumbersome nested *if-statements* the language provides keyword `Elif`.

The following code
```
If 8 > 10:
    x = 8
Elif 11 > 10:
    x = 11
Else:
    x = 10
```
is equivalent to 
```
If 8 > 10:
    x = 8
Else:
    If 11 > 10:
        x = 11
    Else:
        x = 10
```
Keep in mind that **NiLang** is **strongly typed** language and won't allow you to use anything but boolean 
value as a condition.

The code below won't compile.
```
Int x = 1
If x:
    x = 10
```
But, we can fix it by simply using *"equal to"* `==` operator.
```
Int x = 1
If x == 1:
    x = 10
```
## Functions
Functions are reusable bits of code. Functions must be fun, 
that's why their declaration starts with `Fun` keyword.
```
Fun DoNothing:
    Int useless = 0
```
We have declared quite useless function, it does nothing and returns nothing as well.
Let's make it return some `Int` value.
```
Fun UselessNumber::Int:
    Return 5
```
Looks better, but kinda useless anyway. What if our function could check some value and return result 
of its examination.
```
Fun IsGreaterThan5::Bool$ x Int:
    Return x > 5 
```
Since right now we've written only one line functions, let us make a little bit more complex one.
```
Fun ComplexFunction::Int$ x Int, z Int, y Bool:
    If y:
        If x > z:
            Return 0
    If x == z:
        Return 1
    Elif x > z:
        Return 2
    Else:
        Return 3
    Return 0
```
By the way, functions are also very proud members of **NiLang**, hence **function names must begin with an upper case letter**.

Now we have plenty of functions, but how can we use them? Don't worry, call syntax is pretty simple in **NiLang**.

To call a function just write it name.
```
DoNothing
Int x = UselessNumber
```
Character `$` is used to signalize beginning of the function arguments, which are separated by commas. 

```
Bool truth = IsGreaterThan5$ 6
Int complexResult = ComplexFunction$ 10, UselessNumber, IsGreaterThan5$ 4
# in ComplexFunction
# x is 10
# z is UselessNumber, which returns 5
# y is IsGreaterThan5$ 4, which returns False for 4
complexResult = ComplexFunction$ 10, UselessNumber, True
complexResult = ComplexFunction$ 10, 12, False
```

You may wonder what signalizes the end? And the answer is simple - **end of line** or `:`. 

Instead of building insufferable
nested constructions with a lot of brackets **NiLang** facilitates more comprehensive approach to the
way of writing code without any brackets.
## Scopes
To keep number of name collisions low **NiLang** utilizes the concept of named scopes, which helps you
to isolate similarly named entities in the different blocks of code. Scopes are also humble and thus 
**scopes names begin with lower case letter**. 
```
Scope first:
    Fun GetNum::Int:
        Return 1  

Scope second:
    Fun GetNum::Int:
        Return 2
```
To access entities in the scopes use the *scope resolution* - `::` operator.
```
Int x = first::GetNum # x = 1
Int y = second::GetNum # y = 2
```
If you wish to avoid writing scope name each time you want to access entity out of there, 
you may use `Using` keyword to add these entities to the current scope.
```
Using first

Int x = GetNum # x = 1
Int y = second::GetNum # y = 2
Int z = first::GetNum # z = 1
```
Of course, scopes can be nested and you may use *scope resolution* - `::` operator after `Using` keyword.
```
Scope food:
    Scope berry:
        Scope blackberry:
            Bool isTasty = True

Bool isTasty = food::berry::blackberry::isTasty

Using food::berry
isTasty = blackberry::isTasty
```
## Aliases
Sometimes there is a need to give some numerical values comprehensive names, exactly for such
reason exists keyword `Alias`. Since aliases are similar to types, **an alias name also must begin with 
an upper case letter**.
```
Alias Code::Int:
    ok = 100
    bad = 101
    notFound = 404
```
Alternatively you may use `Bool` as a hidden type. 
```
Alias BooleanCode::Bool:
    ok = True
    bad = False
    notFound = False
```
To access values of alias you have to use its name, **but with the first letter in lower case**.
```
Code myCode = code::ok
BooleanCode myBooleanCode = booleanCode::bad
```
You can use `Using` keyword with aliases in the similar fashion as with scopes.
```
Alias State::Int:
    healthy = 1
    ill = 2

Using state
State myState = healthy
```

Keep in mind that you can use only *primitive* types `Bool` or `Int` for defining aliases and
all possible values must be given as *literal expressions*. 

The following code won't compile.
```
Int x = 100
Alias Codes::Int:
    ok = x # not literal expression
    bad = 101
    notFound = 404

Alias Right::Dir: # not primitive type
    right = dir::frontRight # not literal expression
```
## Bot control functions
Currently the following functions are built in the language and are located int the `bot` scope. 
You can use them "out of the box" to control the bot's behaviour:
* `Split$Dir` - bot makes its own copy in the given direction if it has enough energy, a new
bot faces the same direction as its "parent";
* `Fork$Dir` - similar to `Split`, but a new bot creates a new colony and mutates with some probability;
* `Bite$Dir` - bot attacks in the given direction stealing part of an energy from the victim;
* `ConsumeSunlight` - bot consumes the light of the nearest star increasing its energy level;
* `AbsorbMinerals` - bot absorbs minerals under it increasing its energy level;
* `IsEmpty::Bool$Dir` - checks whether cell in the given direction empty or not, returns `True` if it is;
* `IsSibling::Bool$Dir` - checks whether cell in the given direction occupied by bot of the same specie or not, returns `True` if it is;
* `IsFriend::Bool$Dir` - checks whether cell in the given direction occupied by bot of the same colony, returns `True` if it is;
* `GetLuminosity::Int$Dir` - returns the difference in the luminosity between the cell in the given direction
and current position of the bot;
* `GetMineralization::Int$Dir` - returns the difference in the mineralization between the cell in the given direction
and current position of the bot;
* `Sleep` - bot skips one world cycle;
* `Move$Dir` - move the bot on one cell towards the given direction;
* `Face$Dir` - face the bot towards the given direction;
## Loop
Keyword `While` is used to describe block of code, which repeats multiple times **while** condition 
after it is satisfied. The code down below calls `Move` function exactly 10 times.
```
Int x = 0

While x < 10:
    bot::Move$ dir::front
    x = x + 1
```
Sometimes it can be useful to escape *while-loop* earlier.
You can do that with the help of `Break` keyword.
```
Int x = 0

While x < 10:
    x = x + 1
    If bot::IsEmpty$ dir::front:
        Break
```
Or you might want to skip the current iteration. The keyword `Continue` is used for that.
```
Int x = 0

While x < 10:
    x = x + 1
    If Not bot::IsEmpty$ dir::front:
        Continue
    bot::Move$ dir::front
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
## Short assignments
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

Maintaining import structure, which is parallel 
to the arrangement of files in file system, might be also
quite handful.
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
Aliases for imports can be also useful for importing some 
code with similarly named entities.
```
#src/main2.nil
Domain main2

Using helpers = hp 

Int y = hp.x
```
Targeted import might also help to avoid repeating use of resolution operator,
while keeping the problem of name collisions away.
```
#src/main3.nil
Domain main3

Using helpers::x

Int y = x
```
In all other matters, targeted import must be similar to the simple one.
```
#src/main4.nil
Domain main4

Using helpers::x = z

Int y = z
```
