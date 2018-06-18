package zoe

// Symbols are stored in a way that allows to know when they should be recomputed
// if a file changes.
// Also, they store a pointer to whoever is using them, which has to keep being updated.
