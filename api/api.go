package api

import (
	ctx "context"
	"fmt"
	"sync"

	"github.com/kataras/iris"
	"github.com/mutaborius/boiling-api/db"
)

type API struct {
	db  db.BoilingDB
	app *iris.Application
}

func New(db db.BoilingDB) (*API, error) {
	a := &API{db: db}

	app := iris.Default()
	a.app = app

	app.Post("/login", handler(a.postLogin))
	app.Post("/signup", handler(a.postSignup))

	withAuth := app.Party("/", handler(a.withLogin))
	withAuth.Get("/blogs", handler(a.getBlogs))
	withAuth.Post("/blogs", handler(a.postBlog))
	withAuth.Post("/blogs/{id}", handler(a.updateBlog))
	withAuth.Delete("/blogs/{id}", handler(a.deleteBlog))

	return a, nil
}

func (a *API) Run(runner iris.Runner) error {
	return a.app.Run(runner, iris.WithoutServerError(iris.ErrServerClosed))
}

func (a *API) Stop() error {
	return a.app.Shutdown(ctx.Background())
}

func handler(h func(*context)) iris.Handler {
	return func(original iris.Context) {
		if c, ok := original.(*context); ok {
			h(c)
			// we don't have to acquire or release anything here, the outermost
			// wrapper is gonna do that.
			return
		}

		c := acquire(original)
		h(c)
		release(c)
	}
}

type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type context struct {
	iris.Context
	user     db.User
	loggedIn bool
}

func (ctx *context) Next() {
	if ctx.IsStopped() {
		return
	}
	if n, handlers := ctx.HandlerIndex(-1)+1, ctx.Handlers(); n < len(handlers) {
		ctx.HandlerIndex(n)
		handlers[n](ctx)
	}
}

func (ctx *context) Error(e error, httpStatusCode int) {
	var ip, method, path string
	ip = ctx.RemoteAddr()
	method = ctx.Method()
	path = ctx.Path()
	ctx.Application().Logger().Error(fmt.Sprintf("--- --- %s %s %s %s", ip, method, path, e.Error()))

	ctx.StatusCode(httpStatusCode)
	ctx.JSON(Response{
		Status:  "error",
		Message: e.Error(),
	})
}

func (ctx *context) Fail(e error, httpStatusCode int) {
	var ip, method, path string
	ip = ctx.RemoteAddr()
	method = ctx.Method()
	path = ctx.Path()
	ctx.Application().Logger().Warn(fmt.Sprintf("--- --- %s %s %s %s", ip, method, path, e.Error()))

	ctx.StatusCode(httpStatusCode)
	ctx.JSON(Response{
		Status:  "fail",
		Message: e.Error(),
	})
}

func (ctx *context) Success(data interface{}) {
	ctx.SuccessWithCode(data, 200)
}

func (ctx *context) SuccessWithCode(data interface{}, code int) {
	ctx.StatusCode(code)
	ctx.JSON(Response{
		Status: "success",
		Data:   data,
	})
}

var contextPool = sync.Pool{New: func() interface{} { return &context{} }}

func acquire(original iris.Context) *context {
	if c, ok := original.(*context); ok {
		return c
	}

	c := contextPool.Get().(*context)
	c.Context = original
	return c
}

func release(ctx *context) {
	contextPool.Put(ctx)
}
