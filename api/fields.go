package api

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kataras/iris"
)

type dType int

const (
	dTypeInt dType = iota
	dTypeString
	dTypeUnsafeString
	dTypeRawString
	dTypeDate
	dTypeTags
)

type field struct {
	name           string
	required       bool
	needsPrivilege string
	dType          dType
	validator      func(*context, interface{}) bool
}

type fields struct {
	fields map[string]interface{}
}

func (f fields) mustGetString(key string) string {
	s, ok := f.getString(key)
	if !ok {
		panic(fmt.Sprintf("mustGetString: key %s not found", key))
	}

	return s
}

func (f fields) getString(key string) (string, bool) {
	val, ok := f.fields[key]
	if !ok {
		return "", false
	}

	s, ok := val.(string)
	if !ok {
		panic(fmt.Sprintf("field %s is %T but was requested as string", key, val))
	}

	return s, true
}

func (f fields) getInt(key string) (int, bool) {
	val, ok := f.fields[key]
	if !ok {
		return 0, false
	}

	i, ok := val.(int)
	if !ok {
		panic(fmt.Sprintf("field %s is %T but was requested as int", key, val))
	}

	return i, true
}

func (f fields) getDate(key string) (time.Time, bool) {
	val, ok := f.fields[key]
	if !ok {
		return time.Time{}, false
	}

	d, ok := val.(time.Time)
	if !ok {
		panic(fmt.Sprintf("field %s is %T but was requested as date", key, val))
	}

	return d, true
}

func (f fields) mustGetTags(key string) []string {
	t, ok := f.getTags(key)
	if !ok {
		panic(fmt.Sprintf("mustGetTags: key %s not found", key))
	}

	return t
}

func (f fields) getTags(key string) ([]string, bool) {
	val, ok := f.fields[key]
	if !ok {
		return nil, false
	}

	t, ok := val.([]string)
	if !ok {
		panic(fmt.Sprintf("field %s is %T but was requested as tags", key, val))
	}

	return t, true
}

func (a *API) withFields(fields []field) func(*context) {
	a.c.privileges.RLock()
	defer a.c.privileges.RUnlock()
	for _, f := range fields {
		if f.needsPrivilege != "" {
			has := a.c.privileges.l.Has(f.needsPrivilege)
			if !has {
				panic(fmt.Sprintf("withFields: cache doesn't know privilege %s", f.needsPrivilege))
			}
		}
	}

	return func(ctx *context) {
		var err error
		if ctx.Request().Header.Get("Content-Type") == "multipart/form-data" {
			err = ctx.Request().ParseMultipartForm(10 * 1024 * 1024)
		} else {
			err = ctx.Request().ParseForm()
		}
		if err != nil {
			ctx.Fail(userError(err, "invalid form"), iris.StatusBadRequest)
			return
		}

		for _, f := range fields {
			_, ok := ctx.Request().PostForm[f.name]
			if !ok {
				if f.required {
					ctx.Fail(fmt.Errorf("missing required field %s", f.name), iris.StatusBadRequest)
					return
				}
				continue
			}

			if f.needsPrivilege != "" {
				has, err := a.containsPrivilege(ctx.user.Privileges, f.needsPrivilege)
				if err != nil {
					ctx.Error(err, iris.StatusInternalServerError)
					return
				}
				if !has {
					ctx.Fail(fmt.Errorf("field %s provided, but missing necessary privilege %s", f.name, f.needsPrivilege), iris.StatusUnauthorized)
					return
				}
			}

			var parsed interface{}
			if f.dType == dTypeTags {
				tags, err := prepareTags(ctx.PostValues(f.name))
				if err != nil {
					ctx.Fail(userError(err, fmt.Sprintf("invalid tags for field %s", f.name)), iris.StatusBadRequest)
					return
				}
				if len(tags) == 0 {
					if f.required {
						ctx.Fail(fmt.Errorf("missing required field %s", f.name), iris.StatusBadRequest)
						return
					}
					continue
				}
				parsed = tags
			} else {
				raw := ctx.PostValue(f.name)
				trimmed := strings.TrimSpace(raw)
				if len(trimmed) == 0 && f.required {
					// we have to check this again - earlier we just checked if the parameter was found in the form at all.
					ctx.Fail(fmt.Errorf("missing required field %s", f.name), iris.StatusBadRequest)
					return
				}
				switch f.dType {
				case dTypeString:
					parsed = sanitizeString(trimmed)
				case dTypeUnsafeString:
					parsed = trimmed
				case dTypeRawString:
					parsed = raw
				case dTypeInt:
					i, err := strconv.Atoi(trimmed)
					if err != nil {
						ctx.Fail(userError(err, fmt.Sprintf("invalid value for field %s", f.name)), iris.StatusBadRequest)
						return
					}
					parsed = i
				case dTypeDate:
					d, err := time.Parse(time.RFC3339, trimmed)
					if err != nil {
						ctx.Fail(userError(err, fmt.Sprintf("invalid value for field %s", f.name)), iris.StatusBadRequest)
						return
					}
					parsed = d
				default:
					panic(fmt.Sprintf("invalid dType %d", f.dType))
				}
			}

			if f.validator != nil {
				valid := f.validator(ctx, parsed)
				if !valid {
					ctx.Fail(fmt.Errorf("invalid value for field %s", f.name), iris.StatusBadRequest)
					return
				}
			}

			if _, ok := ctx.fields.fields[f.name]; ok {
				panic(fmt.Sprintf("duplicate field value %s", f.name))
			}

			ctx.fields.fields[f.name] = parsed
		}

		ctx.Next()
	}
}
