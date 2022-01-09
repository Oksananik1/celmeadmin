package blank

import (
	"celme/config"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"

	bin "celme/bindata/blank"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

var version = "dev"

// RenderData описывает данные для отрисовки базового шаблона
type RenderData struct {
	Title     string
	Container string
	Message   string
	Version   string
	Root      url.URL
}

// NewRenderData создаёт данные для шаблона с параметрами по умолчанию
func NewRenderData() RenderData {
	root, err := url.Parse("/admin/")
	if err != nil {
		panic(err)
	}
	return RenderData{
		Title:     "",
		Container: "container",
		Message:   "",
		Version:   version,
		Root:      *root,
	}
}

// Register регистрирует обработчик HTTP запросов для базового шаблона
func Register(conf config.Config, mux *http.ServeMux) {
	mux.Handle(
		path.Join("/blank", "static")+"/",
		http.FileServer(blankAssets))
}

// TemplateLoader создаёт загрузчик шаблонов для заданного префикса
func TemplateLoader(base string, assets *assetfs.AssetFS) func(names ...string) *template.Template {
	const blankRoot = "/blank/views/"
	var tmpl *template.Template
	load := func(t *template.Template, name string) *template.Template {
		fp := path.Join(base, name+".html")
		assets := assets // make local copy
		if strings.HasPrefix(name, blankRoot) {
			fp = name + ".html"
			name = strings.TrimPrefix(name, blankRoot)
			assets = blankAssets
		}
		if t == nil {
			t = template.New(name)
		}
		var tmpl *template.Template
		if t.Name() == name {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		contf, err := assets.Open(fp)
		if err != nil {
			log.Panic(err)
		}
		cont, err := ioutil.ReadAll(contf)
		contf.Close()
		if err != nil {
			log.Panic(err)
		}
		tmpl, err = tmpl.Parse(string(cont))
		if err != nil {
			log.Panic(err)
		}
		return tmpl
	}
	return func(names ...string) *template.Template {
		for _, name := range names {
			tmpl = load(tmpl, name)
		}
		return tmpl
	}
}

var blankAssets = &assetfs.AssetFS{
	Asset:     bin.Asset,
	AssetDir:  bin.AssetDir,
	AssetInfo: bin.AssetInfo,
	Prefix:    "",
}
