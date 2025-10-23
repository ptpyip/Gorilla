# The Gorilla Programming Language
An interpreted language developed to learn the making of interpreter, following the book: "Writing an interpreter in Go".

## Ch.1 Lexing

## Ch.2 Parsing

### 2 strategies
1. Top-down
    Start from root nodeof the AST
    Example:
    - **Recursive descent parsing**
    - Early parsing
    - Predictive parsing


2. Bottom-up
    Start from leaf nodes

### Trad-offs
- Speed
- formal proff of correctness
- error-recovery
- error detection

### Pratt parser
In ch.2, we are going to write a **top-down operator precedence** parser, one kind of *recursionve decent parser.

