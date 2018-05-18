package filter

import (
	"github.com/emicklei/go-restful"
	"btdxcx.com/walker/tools"
	"btdxcx.com/walker/apis/common/errors"
)

const(
	tagCenter = "center"
	xshopKEY = "X-SHOP-KEY"
	xshopID = "X-SHOP-ID"
)

// ShopKEYFilter filter shop key
func ShopKEYFilter(tag string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		shopkey := req.Request.Header.Get(xshopKEY)
		if len(shopkey) == 0 {
			errors.Response(resp, errors.Unauthorized("shopkey.filter", "Loss %s", xshopKEY))
			return
		}
		sk := tools.ShopKEY{}
		if tag == tagCenter {
			if err := filterCenter(sk, shopkey, req, resp); err != nil {
				errors.Response(resp, err)
				return
			}
		} else {
			if err := filterWithTag(sk, shopkey, tag, req, resp); err != nil {
				errors.Response(resp, err)
				return
			}
		}
		chain.ProcessFilter(req, resp)
	}
}

func filterCenter(sk tools.ShopKEY, shopkey string, req *restful.Request, resp *restful.Response) error {
	id, err := sk.Decoding(shopkey)
	if err != nil {
		return err
	}
	return filterShopTag(tagCenter, id, id, req, resp)
}

func filterWithTag(sk tools.ShopKEY, shopkey string, tag string, req *restful.Request, resp *restful.Response) error {
	shoptag, id, err := sk.DecodingWithTag(shopkey)
	if err != nil {
		return err
	}
	return filterShopTag(tag, shoptag, id, req, resp)
}

func filterShopTag(tag string, shoptag string, id string, req *restful.Request, resp *restful.Response) error {
	if tag != shoptag {
		return errors.Unauthorized("shopkey.filter", "tag %s id: %s", tag, id)
	}

	req.Request.Header.Del(xshopKEY)
	req.Request.Header.Add(xshopID, id)
	resp.AddHeader(xshopID, id)
	return nil
}