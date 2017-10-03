package api

import (
	ctx "context"
	"fmt"
	"strings"
	"sync"

	"github.com/kataras/iris"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"

	"github.com/boilingrip/boiling-api/db"
)

type API struct {
	db  db.BoilingDB
	app *iris.Application

	c Cache
}

func New(db db.BoilingDB) (*API, error) {
	a := &API{db: db}
	log.Infoln("Building cache...")
	c, err := NewCache(db)
	if err != nil {
		return nil, err
	}
	a.c = c

	a.app = iris.Default()
	a.makeRoutes()

	return a, nil
}

func (a *API) makeRoutes() {
	a.app.Post("/login", handler(a.withFields([]field{
		{
			name:     "username",
			required: true,
			dType:    dTypeString,
		},
		{
			name:     "password",
			required: true,
			dType:    dTypeUnsafeString,
		},
	})), handler(a.postLogin))
	a.app.Post("/signup",
		handler(a.withFields([]field{
			{
				name:     "username",
				required: true,
				dType:    dTypeUnsafeString,
			},
			{
				name:     "email",
				required: true,
				dType:    dTypeUnsafeString,
			},
			{
				name:     "password",
				required: true,
				dType:    dTypeRawString, // postSignup checks if this contains spaces before or after, hence it has to be a raw (non-trimmed) string
			},
		})),
		handler(a.postSignup))

	withAuth := a.app.Party("/", handler(a.withLogin))
	withAuth.Get("/blogs", handler(a.withPrivilege("get_blogs")), handler(a.getBlogs))
	withAuth.Post("/blogs", handler(a.withPrivilege("post_blog")),
		handler(a.withFields([]field{
			{
				name:     "title",
				required: true,
				dType:    dTypeString,
			},
			{
				name:     "content",
				required: true,
				dType:    dTypeString,
			},
			{
				name:     "tags",
				required: true,
				dType:    dTypeTags,
			},
			{
				name:           "author",
				dType:          dTypeInt,
				needsPrivilege: "post_blog_override_author",
				validator: func(_ *context, v interface{}) bool {
					author := v.(int)
					return author >= 0
				},
			},
			{
				name:           "posted_at",
				dType:          dTypeDate,
				needsPrivilege: "post_blog_override_posted_at",
			},
		})),
		handler(a.postBlog))
	withAuth.Post("/blogs/{id}", handler(a.withPrivilege("update_blog")),
		handler(a.withFields([]field{
			{
				name:     "title",
				required: true,
				dType:    dTypeString,
			},
			{
				name:     "content",
				required: true,
				dType:    dTypeString,
			},
			{
				name:     "tags",
				required: true,
				dType:    dTypeTags,
			},
			{
				name:           "author",
				dType:          dTypeInt,
				needsPrivilege: "update_blog_override_author",
				validator: func(_ *context, v interface{}) bool {
					author := v.(int)
					return author >= 0
				},
			},
			{
				name:           "posted_at",
				dType:          dTypeDate,
				needsPrivilege: "update_blog_override_posted_at",
			},
		})),
		handler(a.updateBlog))
	withAuth.Delete("/blogs/{id}", handler(a.withPrivilege("delete_blog")), handler(a.deleteBlog))

	withAuth.Get("/users", handler(a.getUserSelf))
	withAuth.Get("/users/{id}", handler(a.getUser))

	withAuth.Get("/artists/{id}", handler(a.withPrivilege("get_artist")), handler(a.getArtist))
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

func sanitizeString(s string) string {
	return strings.TrimSpace(bluemonday.StrictPolicy().Sanitize(s))
}

type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type errorWithMessage struct {
	original error
	message  string
}

func (e errorWithMessage) Error() string {
	return e.message
}

func userError(original error, message string) error {
	return errorWithMessage{
		original: original,
		message:  message,
	}
}

type context struct {
	iris.Context
	user db.User

	fields
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

const internalServerError = "internal server error"

func (ctx *context) Error(e error, httpStatusCode int) {
	var ip, method, path string
	ip = ctx.RemoteAddr()
	method = ctx.Method()
	path = ctx.Path()
	ctx.Application().Logger().Error(fmt.Sprintf("%d --- %s %s %s %s", httpStatusCode, ip, method, path, e.Error()))

	ctx.StatusCode(httpStatusCode)
	ctx.JSON(Response{
		Status:  "error",
		Message: internalServerError,
	})
}

func (ctx *context) Fail(e error, httpStatusCode int) {
	var ip, method, path string
	ip = ctx.RemoteAddr()
	method = ctx.Method()
	path = ctx.Path()
	ue, ok := e.(errorWithMessage)
	if !ok {
		ctx.Application().Logger().Warn(fmt.Sprintf("%d --- %s %s %s %s", httpStatusCode, ip, method, path, e.Error()))
	} else {
		ctx.Application().Logger().Warn(fmt.Sprintf("%d --- %s %s %s %s (%s)", httpStatusCode, ip, method, path, ue.original.Error(), ue.message))
	}

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
	c.fields.fields = make(map[string]interface{})
	return c
}

func release(ctx *context) {
	contextPool.Put(ctx)
}
