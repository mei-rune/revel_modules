package pongo2

import (
	"io"
	"strings"

	p2 "github.com/flosch/pongo2"
	"github.com/revel/revel"
)

// Adapter for HAML Templates.
type PongoTemplate struct {
	name     string
	template *p2.Template
	engine   *PongoEngine
}

func (tmpl PongoTemplate) Name() string {
	return tmpl.name
}

// return a 'revel.Template' from HAML's template.
func (tmpl PongoTemplate) Render(wr io.Writer, arg interface{}) error {
	err := tmpl.template.ExecuteWriter(p2.Context(arg.(map[string]interface{})), wr)
	if nil != err {
		if e, ok := err.(*p2.Error); ok {
			rerr := &revel.Error{
				Title:       "Template Execution Error",
				Path:        tmpl.name,
				Description: e.ErrorMsg,
				Line:        e.Line,
				//SourceLines: tmpl.Content(),
			}
			if revel.DevMode {
				rerr.SourceLines = tmpl.Content()
			}
			return rerr
		}
	}
	return err
}

func (tmpl PongoTemplate) Content() []string {
	pa, ok := tmpl.engine.loader.TemplatePaths[tmpl.Name()]
	if !ok {
		pa, ok = tmpl.engine.loader.TemplatePaths[strings.ToLower(tmpl.Name())]
	}
	content, _ := revel.ReadLines(pa)
	return content
}

type PongoEngine struct {
	loader                *revel.TemplateLoader
	templateSetBybasePath map[string]*p2.TemplateSet
	templates             map[string]*p2.Template
}

func (engine *PongoEngine) ParseAndAdd(templateName string, templateSource string, basePath string) *revel.Error {
	templateSet := engine.templateSetBybasePath[basePath]
	if nil == templateSet {
		templateSet = p2.NewSet(basePath, p2.MustNewLocalFileSystemLoader(basePath))
		engine.templateSetBybasePath[basePath] = templateSet
	}

	tpl, err := templateSet.FromString(templateSource)
	if nil != err {
		_, line, description := parsePongo2Error(err)
		return &revel.Error{
			Title:       "Template Compilation Error",
			Path:        templateName,
			Description: description,
			Line:        line,
			SourceLines: strings.Split(templateSource, "\n"),
		}
	}

	engine.templates[strings.ToLower(templateName)] = tpl
	return nil
}

func parsePongo2Error(err error) (templateName string, line int, description string) {
	pongoError := err.(*p2.Error)
	if nil != pongoError {
		return pongoError.Filename, pongoError.Line, pongoError.ErrorMsg
	}
	return "", 0, err.Error()
}

func (engine *PongoEngine) Lookup(templateName string) revel.Template {
	tpl := engine.templates[strings.ToLower(templateName)]
	if nil == tpl {
		return nil
	}
	return PongoTemplate{templateName, tpl, engine}
}

func init() {
	revel.TemplateEngines["pongo2"] = func(loader *revel.TemplateLoader) (revel.TemplateEngine, error) {
		return &PongoEngine{
			loader:                loader,
			templateSetBybasePath: map[string]*p2.TemplateSet{},
			templates:             map[string]*p2.Template{},
		}, nil
	}
}
