package chess

// compile time check that UCI implements MoveNotation.
var _ MoveNotation = UCI{}

// compile time check that UCIChess960 implements MoveNotation.
var _ MoveNotation = UCIChess960{}
