package compiler

// Streaming lexer for the compiler.

// Lexer catches the following errors:
// 1. decorators parentheses correctly start and end, with arguments
// correctly formatted and using allowed characters
// 2. the decorator name after @ is a valid decorator
// 3. the resource type value is one of the valid resource types
// 4. the logical name value uses allowed characters
// 5. Resources are separated by 2+ newlines (as determined by Resource type
// 		and Logical name occurrence within the block)
// 6. Expressions between { and } are variables in an appropriate format.
