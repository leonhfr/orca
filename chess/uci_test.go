package chess

// compile time check that UCI implements MoveNotation.
var _ MoveNotation = UCI{}
