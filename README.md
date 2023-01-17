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
