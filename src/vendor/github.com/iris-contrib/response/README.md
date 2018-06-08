## Repository information

This repository contains the built'n and, future community's response engines support for the [Iris web framework](https://github.com/kataras/iris).


## Install

Install one response engine and all will be installed.

```sh
$ go get -u github.com/iris-contrib/response/$FOLDER
```


## Quick look

```go
package main

import (
	"github.com/iris-contrib/response/json"
	"github.com/kataras/iris"
)

type mystruct {
  Name string `json:"name"`
}

func main() {
  cfg:= json.DefaultConfig()
	// you don't have to import json, jsonp, xml, data, text, markdown, iris uses these by default if no other response engine is registered for these content types.  
	iris.UseResponse(json.New(cfg),json.ContentType) // you can still register with your own preferred custom key!
  // custom type use
	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("application/json",  mystruct{Name:"iris"}, iris.Map{"gzip": true,"charset":"UTF-8"}) // gzip is false by default, charset is UTF-8 by default
    // or ctx.RenderWithStatus(iris.StatusOK,...)
	})

  // but for the standard/default 'rest' types like JSON, JSONP, XML, Data, Text, Markdown you can use these context's functions also:
  iris.Get("/defaults", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, mystruct{Name:"iris"}) // options are optional parameter
    	// ctx.JSONP(...)
    	// ctx.XML(...)
    	// ctx.Data(...)
    	// ctx.Text(...)
		// ctx.Markdown(...)
	})


	iris.Listen(":8080")
}

```

> Note: All response engines have optional configuration which can be passed within $engine.New($engine.Config{})

## How to use

- Docs [here](https://kataras.gitbooks.io/iris/content/response-engines.html)
- Examples [here](https://github.com/iris-contrib/examples/tree/master/response_engines)


## How can I make my own iris response engine?

Simply, you have to implement only **one function**.
```go

// ResponseEngine is the interface which all response engines should implement to send responses
// ResponseEngine(s) can be registered with,for example: iris.UseResponse( json.New(), "application/json")
ResponseEngine interface {
  Response(interface{}, ...map[string]interface{}) ([]byte, error)
}
// ResponseEngineFunc is the alternative way to implement a ResponseEngine using a simple function
ResponseEngineFunc func(interface{}, ...map[string]interface{}) ([]byte, error)

```

Register with

```go
// UseResponse accepts a ResponseEngine and the key or content type on which the developer wants to register this response engine
// the gzip and charset are automatically supported by Iris, by passing the iris.RenderOptions{} map on the context.Render
// context.Render renders this response or a template engine if no response engine with the 'key' found
// with these engines you can inject the context.JSON,Text,Data,JSONP,XML also
// to do that just register with UseResponse(myEngine,"application/json") and so on
// look at the https://github.com/iris-contrib/response for examples
//
// if more than one respone engine with the same key/content type exists then the results will be appended to the final request's body
// this allows the developer to be able to create 'middleware' responses engines
//
// Note: if you pass an engine which contains a dot('.') as key, then the engine will not be registered.
// you don't have to import and use github.com/iris-contrib/json, jsonp, xml, data, text, markdown
// because iris uses these by default if no other response engine is registered for these content types
UseResponse(e ResponseEngine, forContentTypesOrKeys ...string)
```

The simplest implementation, which you can look as example, is the Markdown Engine, which is located [here](https://github.com/iris-contrib/response/tree/master/text/text.go).

**Contributions are welcome, make your response engine and do a pr here!**

## License

This project is licensed under the MIT License.

License can be found [here](LICENSE).
