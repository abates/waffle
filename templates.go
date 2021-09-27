package waffle

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed internal/templates
var internal embed.FS
var templates fs.FS

func init() {
	// create filesystem rooted at templates
	var err error
	templates, err = fs.Sub(internal, "internal/templates")
	if err != nil {
		panic("Failed to build subdirectory FS for internal/templates")
	}
}

func getTmplName(filename string) string {
	dir, tmplName := path.Split("/" + filename)
	if strings.HasPrefix(tmplName, "!") {
		tmplName = "." + strings.TrimPrefix(tmplName, "!")
	}

	return fmt.Sprintf("%s%s", dir, strings.TrimSuffix(tmplName, ".tmpl"))
}

type templateBuilder struct {
	input     fs.FS
	dest      string
	root      *template.Template
	templates []string
}

func newTemplateBuilder(dest string) (*templateBuilder, error) {
	tb := &templateBuilder{
		input:     templates,
		dest:      dest,
		root:      template.New("root"),
		templates: []string{},
	}

	return tb, tb.init()
}

func (tb *templateBuilder) init() error {
	return fs.WalkDir(tb.input, ".", func(filename string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		if strings.HasSuffix(filename, ".tmpl") {
			tmplName := getTmplName(filename)
			subTmpl := tb.root.New(tmplName)
			if !strings.HasPrefix(d.Name(), "$") {
				tb.templates = append(tb.templates, tmplName)
			}

			var tmplContent []byte
			tmplContent, err = fs.ReadFile(tb.input, filename)
			if err == nil {
				_, err = subTmpl.Parse(string(tmplContent))
			}
		}
		return err
	})
}

func (tb *templateBuilder) executeTemplate(buf *bytes.Buffer, name string, config Config) error {
	defer buf.Reset()
	err := tb.root.ExecuteTemplate(buf, name, config)
	if err == nil {
		b := buf.Bytes()
		if err == nil && strings.HasSuffix(name, ".go") {
			if b, err = format.Source(buf.Bytes()); err != nil {
				// still write the file contents, but report the
				// formatting error
				println("Failed to format " + name + ": " + err.Error())
				b = buf.Bytes()
				// don't stop processing
			}
		}

		// TODO: should filemode be configurable?
		dest := filepath.Join(tb.dest, filepath.FromSlash(name))
		err = os.MkdirAll(filepath.Dir(dest), 0755)
		if err == nil {
			err = ioutil.WriteFile(dest, b, 0644)
		}
	}
	return err
}

func (tb *templateBuilder) execute(config Config) error {
	buf := bytes.NewBuffer(make([]byte, 8192))
	for _, tmplName := range tb.templates {
		err := tb.executeTemplate(buf, tmplName, config)

		if err != nil {
			return err
		}
	}
	return nil
}

func ExecuteTemplates(destDir string, config Config) error {
	tb, err := newTemplateBuilder(destDir)
	if err == nil {
		err = tb.execute(config)
	}
	return err
}
