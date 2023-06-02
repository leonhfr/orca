package search

// compile time check that noPawnTable implements transpositionPawnTable.
var _ transpositionPawnTable = noPawnTable{}
