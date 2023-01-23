# Smart-Calculator

## Stage 1/7 & 2/7: 2+2+

Write a program that reads two integer numbers from the same line and prints their
sum in the standard output. Numbers can be positive, negative, or zero.
Add a loop to calculate the sum of two numbers continuously.

## Stage 3/7: Count them all

Teach the calculator to read an unlimited sequence of numbers.

- Add to the calculator the ability to read an unlimited sequence of numbers.
- Add a /help command to print some information about the program.
- If you encounter an empty line, do not output anything.

### Examples
```
> 4 5 -2 3
10
> 4 7
11
> 6
6
> /help
The program calculates the sum of numbers
> /exit
Bye!
```

## Stage 4/7: Add subtractions

The program must receive the addition + and subtraction - operators as an input
to distinguish operations from each other. It must support both unary and binary 
minus operators. Moreover, If the user has entered several same operators following
each other, the program still should work.

### Examples
```
> 8
8
> -2 + 4 - 5 + 6
3
> 9 +++ 10 -- 8
27
> 3 --- 5
-2
> 14       -   12
2
> /exit
Bye!
```

## Stage 5/7: Error!

The program should print Invalid expression in cases when the given expression
has an invalid format. If a user enters an invalid command, the program must 
print Unknown command. All messages must be printed without quotes. 

The program must never throw an exception.

### Examples
```
> 8 + 7 - 4
11
> abc
Invalid expression
> 123+
Invalid expression
> +15
15
> 18 22
Invalid expression
>
> -22
-22
> 22-
Invalid expression
> /go
Unknown command
> /exit
Bye!
```

## Stage 6/7: Variables

The program should support variables. Use a `map[string]int` to store them.

Rules for variables:

   - We suppose that the name of a variable (identifier) can contain only Latin letters.
   - A variable can have a name consisting of more than one letter.
   - The case is also important; your program should be able to take both lowercase and uppercase variables.
   - The value can be an integer number or a value of another variable.
   - It should be possible to set a new value to an existing variable.
   - To print the value of a variable, you should just type its name.

Incorrect spelling or declaration of variables should also throw an exception with the corresponding message to the user:
- First, the variable is checked for correctness. If the user inputs an invalid variable name, then the output should be 
"Invalid identifier".
```
> a2a
Invalid identifier
> n22
Invalid identifier
```
- If a variable is valid but not declared yet, the program should print "Unknown variable".
```
> a = 8
> b = c
Unknown variable
> e
Unknown variable
```
- If an identifier or value of a variable is invalid during variable declaration, the program must print a message like 
the one below.
```
> a1 = 8
Invalid identifier
> n1 = a2a
Invalid identifier
> n = a2a
Invalid assignment
> a = 7 = 8
Invalid assignment
```
Handle as many incorrect inputs as possible. The program must never throw an exception of any kind.
### Examples
```
> a  =  3
> b= 4
> c =5
> a + b - c
2
> b - c + 4 - a
0
> X = 800
> Y=4
> Z= 5
> X + Y + Z
809
> BIG = 9000
> BIG
9000
> big
Unknown variable
> /exit
Bye!
```
## Stage 7/7: I've got the power
In the final stage, it remains to add operations: multiplication *, integer division / and parentheses (...). 
They have a higher priority than addition + and subtraction -.

## Objectives
- Your program should support multiplication *, integer division / and parentheses (...).
- Modify the result of the /help command to explain all possible operators.
- The program should not stop until the user enters the /exit command.
- Note that a sequence of + (like +++ or +++++) is an admissible operator that should be interpreted as a single plus.
A sequence of - (like -- or ---) is also an admissible operator and its meaning depends on the length. If a user enters
a sequence of * or /, the program must print a message that the expression is invalid.

### Examples
```
> 8 * 3 + 12 * (4 - 2)
48
> 2 - 2 + 3
3
> 4 * (2 + 3
Invalid expression
> -10
-10
> a=4
> b=5
> c=6
> a*2+b*3+c*(2+3)
53
> 1 +++ 2 * 3 -- 4
11
> 3 *** 5
Invalid expression
> 4+3)
Invalid expression
> /command
Unknown command
> /exit
Bye!
```