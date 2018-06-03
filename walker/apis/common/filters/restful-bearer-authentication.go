package filter

import (
	"btdxcx.com/walker/apis/common/errors"
	"btdxcx.com/walker/model"
	"btdxcx.com/walker/service/auth"
	"github.com/emicklei/go-restful"
	"golang.org/x/net/context"
	"strings"
)

// BearerAuthenticate filter bearer authenticate
func BearerAuthenticate(service auth.IService) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		beare := req.HeaderParameter("Authorization")
		shopID := req.HeaderParameter("X-Shop-Id")

		// noop on nil value
		if len(beare) == 0 {
			errors.Response(resp, errors.Unauthorized("auth.wrapper", "loss authorization"))
			return
		}

		vals := strings.Split(beare, " ")
		if len(vals) != 2 {
			errors.Unauthorized("auth.wrapper", "loss bearer authorization")
			return
		}
		if vals[0] != "Bearer" {
			errors.Unauthorized("auth.wrapper", "loss bearer authorization")
			return
		}

		in := &model.IntrospectRequest{
			AccessToken: vals[1],
			ShopID:      shopID,
		}
		out := &model.Introspect{}
		if err := service.Introspect(context.TODO(), in, out); err != nil {
			errors.Response(resp, errors.Unauthorized("bearer.filter", "Introspect %v", err))
			return
		}

		if len(out.Username) == 0 {
			errors.Response(resp, errors.Unauthorized("bearer.filter", "username error"))
			return
		}

		req.Request.Header.Del("Authorization")
		req.Request.Header.Add("X-Username", out.Username)

		chain.ProcessFilter(req, resp)
	}
}
