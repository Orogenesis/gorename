package replacer

import (
	"bytes"
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

var defaultOptions = options{
	rootDir: ".",
	logf:    log.Printf,
	ignoreFolders: []string{
		".git",
		"vendor",
	},
}

type Replace struct {
	prevPath string
	newPath  string
	opts     options
}

type options struct {
	ignoreFolders []string
	useModules    bool
	printResult   bool
	rootDir       string
	logf          func(string, ...interface{})
}

type ClientOption func(*options)

func IgnoreFolder(name string) ClientOption {
	return func(o *options) { o.ignoreFolders = append(o.ignoreFolders, name) }
}

func UseModules(useModules bool) ClientOption {
	return func(o *options) { o.useModules = useModules }
}

func PrintResult(printResult bool) ClientOption {
	return func(o *options) { o.printResult = printResult }
}

func RootDir(rootDir string) ClientOption {
	return func(o *options) { o.rootDir = rootDir }
}

func Logf(logf func(string, ...interface{})) ClientOption {
	return func(o *options) { o.logf = logf }
}

func New(prevPath, newPath string, opt ...ClientOption) *Replace {
	opts := defaultOptions
	for _, o := range opt {
		o(&opts)
	}

	return &Replace{
		prevPath: prevPath,
		newPath:  newPath,
		opts:     opts,
	}
}

func (r *Replace) Run() error {
	fset := token.NewFileSet()
	err := filepath.Walk(r.opts.rootDir, r.walkFn(fset))
	if err != nil {
		return fmt.Errorf("can't replace imports: %w", err)
	}

	if !r.opts.useModules {
		return nil
	}

	err = r.renameModulePath()
	if err != nil {
		return fmt.Errorf("can't rename module path: %w", err)
	}

	return nil
}

func (r *Replace) renameModulePath() (err error) {
	filePath := filepath.Join(r.opts.rootDir, "go.mod")
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	modFile, err := modfile.Parse(fi.Name(), buf.Bytes(), nil)
	if err != nil {
		return err
	}

	modulePath := modFile.Module.Mod.Path
	if modulePath == r.newPath {
		return nil
	}

	err = modFile.AddModuleStmt(r.newPath)
	if err != nil {
		return err
	}

	outBytes, err := modFile.Format()
	err = ioutil.WriteFile(filePath, outBytes, fi.Mode())
	return nil
}

func (r *Replace) walkFn(fset *token.FileSet) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		for _, f := range r.opts.ignoreFolders {
			if !info.IsDir() {
				continue
			}

			if info.Name() != f {
				continue
			}

			return filepath.SkipDir
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		f, err := NewFile(fset, path)
		if err != nil {
			return err
		}

		f.RewriteImport(r.prevPath, r.newPath)
		var buf bytes.Buffer
		err = f.Fprint(&buf)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(path, buf.Bytes(), info.Mode())
		if err != nil {
			return err
		}

		if r.opts.printResult {
			r.opts.logf(path)
		}

		return nil
	}
}
