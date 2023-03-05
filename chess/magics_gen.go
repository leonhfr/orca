//go:build ignore

package main

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"

	"github.com/leonhfr/orca/chess"
)

//go:embed magics.go.tmpl
var magicsTemplate string

const magicsPath = "magics.go"

func main() {
	var magicLiterals []MagicLiteral
	for _, magicType := range []struct {
		pieceType chess.PieceType
		name      string
	}{{chess.Rook, "rook"}, {chess.Bishop, "bishop"}} {
		fmt.Printf("Generating %s magics...\n", magicType.name)

		magicLiterals = append(magicLiterals, MagicLiteral{
			Name:   magicType.name,
			Magics: chess.FindMagics(magicType.pieceType),
		})
	}

	code, err := formatCode(magicLiterals)
	if err != nil {
		log.Fatalf("could not generate source code: %s", err)
	}

	err = writeCode(code, magicsPath)
	if err != nil {
		log.Fatalf("could not write source code: %s", err)
	}
}

// MagicLiteral contains the data to pass magic literals to the template.
type MagicLiteral struct {
	Name   string
	Magics [64]chess.Magic
}

// formatCode formats the magic literals to valid Go code using the template.
func formatCode(magicLiterals []MagicLiteral) ([]byte, error) {
	t := template.Must(template.New("").Parse(magicsTemplate))

	var code bytes.Buffer
	err := t.ExecuteTemplate(&code, "", struct {
		MagicLiterals []MagicLiteral
	}{
		MagicLiterals: magicLiterals,
	})
	if err != nil {
		return nil, err
	}

	formatted, err := format.Source(code.Bytes())
	if err != nil {
		return nil, err
	}

	return formatted, nil
}

// writeCode writes the code to the file.
func writeCode(code []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	_, err = w.Write(code)
	if err != nil {
		return err
	}

	return w.Flush()
}
