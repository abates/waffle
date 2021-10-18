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

func newTemplateBuilder(input fs.FS, dest string) (*templateBuilder, error) {
	tb := &templateBuilder{
		input:     input,
		dest:      dest,
		root:      template.New("root"),
		templates: []string{},
	}

	return tb, tb.init()
}

func (tb *templateBuilder) init() error {
	return fs.WalkDir(tb.input, ".", func(filename string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			if err != nil {
				Logger.Logf("Template Builder <fail>failed</fail> loading %s: %v", filename, err)
			}
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
				Logger.Logf("Parsing template %v", filename)
				_, err = subTmpl.Parse(string(tmplContent))
			}

			if err != nil {
				Logger.Logf("Template Builder <fail>failed</fail> to load template %s: %v", filename, err)
			}
		}
		return err
	})
}

func (tb *templateBuilder) executeTemplate(buf *bytes.Buffer, name string, config Config) error {
	defer buf.Reset()
	err := tb.root.ExecuteTemplate(buf, name, config)
	if err == nil {
		dest := filepath.Join(tb.dest, filepath.FromSlash(name))
		err = os.MkdirAll(filepath.Dir(dest), 0755)

		if err == nil && strings.HasSuffix(name, ".go") {
			b := buf.Bytes()
			if b, err = format.Source(buf.Bytes()); err != nil {
				// still write the file contents, but report the
				// formatting error
				err = fmt.Errorf("failed to format go source file %s: %w", dest, err)
				b = buf.Bytes()
			}

			err1 := ioutil.WriteFile(dest, b, 0644)
			if err == nil && err1 != nil {
				err = fmt.Errorf("failed to write file %s: %w", dest, err)
			}
		} else if err != nil {
			err = fmt.Errorf("failed to create directory %s: %w", filepath.Dir(dest), err)
		}

	} else {
		err = fmt.Errorf("failed to execute template %s: %v", name, err)
	}
	return err
}

func (tb *templateBuilder) execute(config Config) error {
	buf := bytes.NewBuffer(make([]byte, 8192))
	for _, tmplName := range tb.templates {
		err := tb.executeTemplate(buf, tmplName, config)

		if err == nil {
			Logger.Logf("<success>%s</success>", tmplName)
		} else {
			Logger.Logf("<fail>%s</fail>: %v", tmplName, err)
			return err
		}
	}
	return nil
}

func ExecuteTemplates(srcDir string, destDir string, config Config) error {
	templates, err := fs.Sub(internal, fmt.Sprintf("internal/templates/%s", srcDir))
	if err == nil {
		tb, err := newTemplateBuilder(templates, destDir)
		if err == nil {
			err = tb.execute(config)
		}
	} else {
		err = fmt.Errorf("Failed to mount %q: %w", srcDir, err)
	}
	return err
}
