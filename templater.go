package migo

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type TemplateType string

const (
	TemplateTypeSQLUp   TemplateType = "sql-up"
	TemplateTypeSQLDown              = "sql-down"
	TemplateTypeGo                   = "go"
	TemplateTypeVersion              = "version"
)

type Template struct {
	Type TemplateType
	file string
}

func (t *Template) Build(v *Version) ([]byte, error) {
	tmpl, err := template.ParseFiles(t.file)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]interface{}{
		"version": map[string]interface{}{
			"name":    v.Name,
			"verSafe": strings.Replace(v.v.String(), ".", "_", -1),
			"ver":     v.v.String(),
		},
	})
	return buf.Bytes(), nil
}

func (t *Template) BuildVersion() (string, error) {
	tmpl, err := template.ParseFiles(t.file)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]interface{}{
		"timestamp": time.Now().Unix(),
	})

	return strings.Replace(buf.String(), "\n", "", -1), nil
}

type Templater struct {
	templates map[TemplateType]Template
}

func (t *Templater) LoadTemplates(path string) error {
	t.templates = make(map[TemplateType]Template)
	err := t.loadSQLTemplates(path)
	if err != nil {
		return err
	}

	err = t.loadGoTemplates(path)
	if err != nil {
		return err
	}

	err = t.loadVersionTemplate(path)
	if err != nil {
		return err
	}

	return nil
}

func (t *Templater) ContentForTemplateType(tmplType TemplateType, v *Version) ([]byte, error) {
	tmpl := t.templates[tmplType]
	return tmpl.Build(v)
}

func (t *Templater) TampleteWithType(tmplType TemplateType) *Template {
	tmpl, found := t.templates[tmplType]
	if !found {
		return nil
	}

	return &tmpl
}

func (t *Templater) loadSQLTemplates(path string) error {
	upPath := t.migoPath(path) + "/tmpl/sql/up.sql"
	if _, err := os.Stat(upPath); os.IsNotExist(err) {
		err = os.MkdirAll(t.migoPath(path)+"/tmpl/sql", os.ModePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(upPath, []byte(SQLUpDefaultTemplateData), 0755)
		if err != nil {
			return err
		}
	}

	downPath := t.migoPath(path) + "/tmpl/sql/down.sql"
	if _, err := os.Stat(downPath); os.IsNotExist(err) {
		err = ioutil.WriteFile(downPath, []byte(SQLDownDefaultTemplateData), 0755)
		if err != nil {
			return err
		}
	}

	t.templates[TemplateTypeSQLUp] = Template{
		file: upPath,
		Type: TemplateTypeSQLUp,
	}

	t.templates[TemplateTypeSQLDown] = Template{
		file: downPath,
		Type: TemplateTypeSQLDown,
	}

	return nil
}

func (t *Templater) loadGoTemplates(path string) error {
	goPath := t.migoPath(path) + "/tmpl/go/go.tmpl"
	if _, err := os.Stat(goPath); os.IsNotExist(err) {
		err = os.MkdirAll(t.migoPath(path)+"/tmpl/go", os.ModePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(goPath, []byte(GoDefaultTemplateData), 0755)
		if err != nil {
			return err
		}
	}

	t.templates[TemplateTypeGo] = Template{
		file: goPath,
		Type: TemplateTypeGo,
	}

	return nil
}

func (t *Templater) loadVersionTemplate(path string) error {
	verPath := t.migoPath(path) + "/tmpl/version/version.tmpl"
	if _, err := os.Stat(verPath); os.IsNotExist(err) {
		err = os.MkdirAll(t.migoPath(path)+"/tmpl/version", os.ModePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(verPath, []byte(VersionDefaultTemplateData), 0755)
		if err != nil {
			return err
		}
	}

	t.templates[TemplateTypeVersion] = Template{
		file: verPath,
		Type: TemplateTypeVersion,
	}

	return nil
}

func (t *Templater) migoPath(path string) string {
	return path + "/migo"
}
