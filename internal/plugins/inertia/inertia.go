package inertia

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lemmego/api/app"
	"github.com/romsar/gonertia"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
)

const ViteHotPath = "./public/hot"
const InertiaRootTemplatePath = "resources/views/root.html"
const InertiaManifestPath = "./public/build/manifest.json"
const InertiaBuildPath = "/public/build/"

type Provider struct{}

type Flash struct {
	errors map[string]gonertia.ValidationErrors
}

type InertiaResponse struct {
	*Inertia
	filePath string
	props    map[string]any
	ctx      app.Context
}

type Inertia struct {
	inertia *gonertia.Inertia
}

func (i *Provider) Provide(a app.App) error {
	inertia := NewInertia(
		a,
		InertiaRootTemplatePath,
		gonertia.WithVersionFromFile(InertiaManifestPath),
		gonertia.WithSSR(),
		//inertia.WithVersion("1.0"),
		gonertia.WithFlashProvider(NewFlash()),
	)

	a.AddService(inertia)
	return nil
}

func Respond(c app.Context, filePath string, props map[string]any) error {
	return Get(c.App()).Respond(c, filePath, props)
}

func (i *Inertia) Respond(c app.Context, filePath string, props map[string]any) error {
	if errs := c.PopSession("errors"); errs != nil {
		if props == nil {
			props = map[string]any{}
		}

		props["errors"] = errs
	}
	ir := &InertiaResponse{Inertia: i, ctx: c, filePath: filePath, props: props}
	return c.Render(ir)
}

func (i *Inertia) Redirect(c app.Context, url string) {
	i.inertia.Redirect(c.ResponseWriter(), c.Request(), url)
	return
}

func (i *Inertia) Flash(c app.Context, message string, props map[string]any) {
	//
}

func (i *Inertia) Back(c app.Context) {
	i.inertia.Back(c.ResponseWriter(), c.Request(), c.Status())
}

func (ir *InertiaResponse) Render(w io.Writer) error {
	if ir.ctx.Status() == 0 {
		ir.ctx.SetStatus(http.StatusOK)
	}
	ir.ctx.WriteStatus(ir.ctx.Status())
	return ir.inertia.Render(w.(http.ResponseWriter), ir.ctx.Request(), ir.filePath, ir.props)
}

func NewInertia(a app.App, rootTemplatePath string, opts ...gonertia.Option) *Inertia {
	i, err := gonertia.NewFromFile(
		rootTemplatePath,
		opts...,
	)

	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(ViteHotPath)

	if err == nil {
		i.ShareTemplateFunc("vite", func(entry string) (string, error) {
			content, err := os.ReadFile(ViteHotPath)
			if err != nil {
				return "", err
			}
			url := strings.TrimSpace(string(content))
			if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
				url = url[strings.Index(url, ":")+1:]
			} else {
				url = "//localhost:5173"
			}
			if entry != "" && !strings.HasPrefix(entry, "/") {
				entry = "/" + entry
			}
			return url + entry, nil
		})
	} else {
		i.ShareTemplateFunc("vite", Vite(InertiaManifestPath, InertiaBuildPath))
	}

	i.ShareTemplateData("env", a.Config().Get("app.env"))

	return &Inertia{inertia: i}
}

func Vite(manifestPath, buildDir string) func(path string) (string, error) {
	f, err := os.Open(manifestPath)
	if err != nil {
		log.Fatalf("cannot open provided vite manifest file: %s", err)
	}
	defer f.Close()

	viteAssets := make(map[string]*struct {
		File   string `json:"file"`
		Source string `json:"src"`
	})
	err = json.NewDecoder(f).Decode(&viteAssets)
	// print content of viteAssets
	for k, v := range viteAssets {
		log.Printf("%s: %s\n", k, v.File)
	}

	if err != nil {
		log.Fatalf("cannot unmarshal vite manifest file to json: %s", err)
	}

	return func(p string) (string, error) {
		if val, ok := viteAssets[p]; ok {
			return path.Join("/", buildDir, val.File), nil
		}
		return "", fmt.Errorf("asset %q not found", p)
	}
}

func NewFlash() *Flash {
	return &Flash{errors: make(map[string]gonertia.ValidationErrors)}
}

func (p *Flash) FlashErrors(ctx context.Context, errors gonertia.ValidationErrors) error {
	if sessionID, ok := ctx.Value("sessionID").(string); ok {
		p.errors[sessionID] = errors
	}
	return nil
}

func (p *Flash) GetErrors(ctx context.Context) (gonertia.ValidationErrors, error) {
	var inertiaErrors gonertia.ValidationErrors
	if sessionID, ok := ctx.Value("sessionID").(string); ok {
		inertiaErrors = p.errors[sessionID]
		p.errors[sessionID] = nil
	}
	return inertiaErrors, nil
}

func Get(a app.App) *Inertia {
	return a.Service(reflect.TypeOf(&Inertia{})).(*Inertia)
}
