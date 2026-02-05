package request_utils

import (
	"errors"
	"net/url"

	"github.com/gin-gonic/gin"
)

// ParamsToMapOneValue parses a url.Values into a map[string]interface{} with single values only.
type paramsBuilder struct {
	values  url.Values
	must    []string
	allowed []string
}

func ParamsToMapOneValue(values url.Values) *paramsBuilder {
	return &paramsBuilder{values: values}
}

func (b *paramsBuilder) Must(keys []string) *paramsBuilder    { b.must = keys; return b }
func (b *paramsBuilder) Allowed(keys []string) *paramsBuilder { b.allowed = keys; return b }

func (b *paramsBuilder) Result() (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for _, k := range b.must {
		if b.values.Get(k) == "" {
			return nil, errors.New(k + " is required")
		}
	}
	if len(b.allowed) > 0 {
		allow := map[string]bool{}
		for _, k := range b.allowed {
			allow[k] = true
		}
		for k := range b.values {
			if !allow[k] {
				continue
			}
			res[k] = b.values.Get(k)
		}
		return res, nil
	}
	for k := range b.values {
		res[k] = b.values.Get(k)
	}
	return res, nil
}

// QueriesToMapOneValue is a convenience wrapper to parse query params from gin.Context.
func QueriesToMapOneValue(c *gin.Context) *paramsBuilder {
	if c == nil || c.Request == nil {
		return &paramsBuilder{values: url.Values{}}
	}
	return ParamsToMapOneValue(c.Request.URL.Query())
}
