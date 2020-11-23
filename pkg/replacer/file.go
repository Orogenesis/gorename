package replacer

import (
	"go/token"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/dstutil"
)

type File struct {
	dstFile *dst.File
}

func NewFile(fset *token.FileSet, filePath string) (*File, error) {
	dstFile, err := decorator.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	f := &File{dstFile}
	return f, nil
}

func (f *File) RewriteImport(prevPath, newPath string) (rewrote bool) {
	f.Apply(func(cursor *dstutil.Cursor) bool {
		node := cursor.Node()
		if typedNode, ok := node.(*dst.ImportSpec); ok {
			if mustMatchPath(typedNode.Path.Value, prevPath) {
				rewrote = true
				typedNode.Path.Value = strings.Replace(typedNode.Path.Value, prevPath, newPath, 1)
			}
		}

		return true
	})

	return rewrote
}

func (f *File) Fprint(w io.Writer) error {
	err := decorator.Fprint(w, f.dstFile)
	return err
}

func (f *File) Apply(fn func(cursor *dstutil.Cursor) bool) dst.Node {
	return dstutil.Apply(f.dstFile, fn, nil)
}

func mustMatchPath(currentPath, newPath string) bool {
	unquotedPath, err := strconv.Unquote(currentPath)
	if err != nil {
		panic(err)
	}

	path, err := filepath.Rel(unquotedPath, newPath)
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(path, "..") {
		return false
	}

	return true
}
