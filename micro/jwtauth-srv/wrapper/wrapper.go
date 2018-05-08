package wrapper

import (
	"strings"
	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
	"context"
	"github.com/micro/go-micro/client"

	proto "btdxcx.com/micro/jwtauth-srv/proto/auth"
)

type wrapper struct {
	scope string
	client.Client
	proto.JwtAuthClient
}

func (w *wrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// get headers
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errors.Unauthorized("auth.wrapper", "loss authorization")
	}

	customScope := md["Custom-Scope"]
	if customScope == w.scope {
		return w.Client.Call(ctx, req, rsp, opts...)
	}

	// get Authorization val
	val := md["Authorization"]

	// noop on nil value
	if len(val) == 0 {
		return errors.Unauthorized("auth.wrapper", "loss authorization")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return errors.Unauthorized("auth.wrapper", "loss bearer authorization")
	}
	if vals[0] != "Bearer" {
		return errors.Unauthorized("auth.wrapper", "loss bearer authorization")
	}
	
	introspect, err := w.JwtAuthClient.Introspect(
		NewContextToken(ctx, vals[1]), 
		&proto.IntrospectRequest{})
	if err != nil {
		return err
	}

	isScope := false
	for _, scope := range introspect.Token.Scopes {
		if scope == w.scope {
			isScope = true
			break
		}
	}
	if !isScope {
		return errors.Unauthorized("auth.wrapper", "loss permissions")
	}
	
	// req.Request.Header
	md["Scopes"] = strings.Join(introspect.Token.Scopes, ";")
	for k, v := range introspect.Token.Metadata {
		md[k] = v
	}
	md["ClientId"] = introspect.ClientId
	// delete(md, "Authorization")
	// delete(md, "Token")

	newCtx := metadata.NewContext(ctx, md)

	return w.Client.Call(newCtx, req, rsp, opts...)
}

// NewClientWrapper is a wrapper which shards based on a header key value
func NewClientWrapper(scope string) client.Wrapper {
	jwtAuthClient := proto.NewJwtAuthClient("com.btdxcx.micro.srv.jwtauth", client.DefaultClient)
	return func(c client.Client) client.Client {
		return &wrapper{
			scope:    scope,
			Client: c,
			JwtAuthClient: jwtAuthClient,
		}
	}
}

// FromCtx from header get bearer authorization
func FromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		// noop
		return "", errors.Unauthorized("auth wrapper", "loss token")
	}
	return md["Token"], nil
}

func getClientIDFrom(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		// noop
		return "", errors.Unauthorized("auth wrapper", "loss token")
	}
	return md["Clientid"], nil
}

// GetClientIDFrom from header get shop id
func GetClientIDFrom(ctx context.Context, defID string) (string, error) {
	sid, err := getClientIDFrom(ctx)
	if err == nil {
		return sid, nil
	}
	if len(defID) == 0 {
		return "", err
	}
	return defID, nil
}

// NewContext new bearer authorization context
func NewContext(context context.Context, bearer string) context.Context {
	return newContex(context, "Authorization", bearer)
}

// NewContextToken new token authorization context
func NewContextToken(context context.Context, token string) context.Context {
	return newContex(context, "Token", token)
}

// NewContextCustomScope new custom scope context
func NewContextCustomScope(context context.Context, scope string) context.Context {
	return newContex(context, "Custom-Scope", scope)
}

func newContex(ctx context.Context, key string, val string) context.Context {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.Metadata{ key: val }
	} else {
		md[key] = val
	}
	return metadata.NewContext(ctx, md)
}
