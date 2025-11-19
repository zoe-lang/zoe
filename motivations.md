
Like for many languages, zoe was born from the sentiment that no language offered this set of requirements :

 - compile time null-access checking
 - trait-based inheritance/type composing
 - sum types and their associated pattern matching/clean variable destructuring
 - functional patterns
 - closures
 - garbage collection (mostly) without stop-the-world pauses
 - errors as values instead of catch/throw semantics with enforced explicit error handling that tries to stay mostly out of the way
 - compile time execution without macros or type definition changes

It tries to be useful for all manner of programming and does not try to focus on any given world problem. In a sense, it tries to fit a spot in between Rust, Golang, and C#.

# Compile-time null-access checking

No modern language should exist without it. Runtime null-pointer exception simply should not happen ; this is too common a mistake and it falls within the compiler's responsability to prevent its user from shooting themself in the foot.

# Trait-based inheritance

This choice is esthetical in nature. I just like them.

# Garbage collection (mostly) without stop-the-world pauses

# Compared to other languages :

  - Rust : 
  - Go :
  - Zig :
  - Typescript : 
  - C/C++ :

