
# Stuff that may hold names

Two way to add names ; they are either member of something (that will be reached through '.'), or in
a scope. So we should handle a Context that will know where the object trying to register itself
should go.

Context.RegisterName() (with some kind of flag ?)

Name registration should be right at the point where the name is found to avoid redefinitions.