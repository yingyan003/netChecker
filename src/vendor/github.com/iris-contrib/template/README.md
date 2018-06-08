## Repository information

This repository contains the built'n and, future community's template engines support for the [Iris web framework](https://github.com/kataras/iris).


## Install

Install one template engine and all will be installed.

```sh
$ go get -u github.com/iris-contrib/template/$FOLDER
```


## Quick look

```go
// small example for html/template, same for all other

package main

import (
	"github.com/iris-contrib/template/html"
	"github.com/kataras/iris"
)

type mypage struct {
	Title   string
	Message string
}

func main() {
	iris.UseTemplate(html.New()).Directory("./templates", ".html")
	// for binary assets: .Directory(dir string, ext string).Binary(assetFn func(name string) ([]byte, error), namesFn func() []string)

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("mypage.html", mypage{"My Page title", "Hello world!"}, iris.Map{"gzip": true})
	})

	iris.Listen(":8080")
}

```

> Note: All template engines have optional configuration which can be passed within $engine.New($engine.Config{})

## How to use

- Docs [here](https://kataras.gitbooks.io/iris/content/template-engines.html)
- Examples [here](https://github.com/iris-contrib/examples/tree/master/template_engines)


## How can I make my own iris template engine?

Simply, you have to implement only **3  functions**, for load and execute the templates. One optionally (**Funcs() map[string]interface{}**) which is used to register the iris' helpers funcs like `{{ url }}` and `{{ urlpath }}`.

```go

type (
	// TemplateEngine the interface that all template engines must implement
	TemplateEngine interface {
		// LoadDirectory builds the templates, usually by directory and extension but these are engine's decisions
		LoadDirectory(directory string, extension string) error
		// LoadAssets loads the templates by binary
		// assetFn is a func which returns bytes, use it to load the templates by binary
		// namesFn returns the template filenames
		LoadAssets(virtualDirectory string, virtualExtension string, assetFn func(name string) ([]byte, error), namesFn func() []string) error

		// ExecuteWriter finds, execute a template and write its result to the out writer
		// options are the optional runtime options can be passed by user
		// an example of this is the "layout" or "gzip" option
		ExecuteWriter(out io.Writer, name string, binding interface{}, options ...map[string]interface{}) error
	}

	// TemplateEngineFuncs is optional interface for the TemplateEngine
	// used to insert the Iris' standard funcs, see var 'usedFuncs'
	TemplateEngineFuncs interface {
		// Funcs should returns the context or the funcs,
		// this property is used in order to register the iris' helper funcs
		Funcs() map[string]interface{}
	}
)

```

The simplest implementation, which you can look as example, is the Markdown Engine, which is located [here](https://github.com/iris-contrib/template/tree/master/markdown/markdown.go).

**Contributions are welcome, make your template engine and do a pr here!**

## License

This project is licensed under the MIT License.

License can be found [here](LICENSE).
