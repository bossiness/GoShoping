package shopkey

import (
	"github.com/micro/go-log"
	
	"github.com/micro/go-micro/errors"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"

	"context"

	shopproto "btdxcx.com/micro/shop-srv/proto/shop"
)

type shopkey struct {
	key string
	ctype string
	client.Client
	shopproto.ShopKeyClient
}

func (s *shopkey) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	// get headers
	md, ok := metadata.FromContext(ctx)
	log.Log("Received ShopKey.Call request")
	if !ok {
		// noop
		return errors.BadRequest("shop-key wrapper", "loss shop key")
	}

	// get key val
	val := md[s.key]

	// noop on nil value
	if len(val) == 0 {
		return errors.BadRequest("shop-key wrapper", "loss shop key")
	}

	introspect, err := s.ShopKeyClient.Introspect(ctx, &shopproto.IntrospectRequest{Key: val})
	if err != nil {
		return err
	}

	if introspect.Type != s.ctype {
		return errors.BadRequest("shop-key wrapper", "shop key type invalid")
	}

	// req.Request.Header
	md["X-SHOP-ID"] = introspect.Uuid
	// delete(md, "X-SHOP-KEY")

	newCtx := metadata.NewContext(ctx, md)

	return s.Client.Call(newCtx, req, rsp, opts...)
}

// NewClientWrapper is a wrapper which shards based on a header key value
func NewClientWrapper(key string, ctype string) client.Wrapper {
	shopCl := shopproto.NewShopKeyClient("com.btdxcx.micro.srv.shop", client.DefaultClient)
	return func(c client.Client) client.Client {
		return &shopkey{
			key:    key,
			ctype: ctype,
			Client: c,
			ShopKeyClient: shopCl,
		}
	}
}

// FromCtx from header get shop id
func FromCtx(ctx context.Context) (string, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		// noop
		return "", errors.BadRequest("shop-key wrapper", "loss shop key")
	}
	return md["X-Shop-Id"], nil
}

// NewNewContext new shop key context
func NewNewContext(context context.Context, shopKEY string) context.Context {

	md, ok := metadata.FromContext(context)
	if !ok {
		md = metadata.Metadata{
			"X-SHOP-KEY": shopKEY,
		}
	} else {
		md["X-SHOP-KEY"] = shopKEY
	}
	return metadata.NewContext(context, md)
}