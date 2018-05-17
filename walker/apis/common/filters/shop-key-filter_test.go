package filter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"btdxcx.com/walker/apis/common/filters"
	"github.com/emicklei/go-restful"
)

func TestShopKEYCenterFilter(t *testing.T) {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/dummy").To(dummy))
	restful.Add(ws)
	restful.Filter(filter.NCSACommonLogFormatLogger())
	restful.Filter(filter.ShopKEYFilter("center"))

	httpRequest, _ := http.NewRequest("OPTIONS", "http://here.io/dummy", nil)
	httpRequest.Header.Add("X-SHOP-KEY", "a9468c6895cd")
	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)
	shopid := httpWriter.Header().Get("X-Shop-Id")
	if "center" != shopid {
		t.Fatal("expected: X-Shop-Id but center:" + shopid)
	}
}

func TestShopKEYBackFilter(t *testing.T) {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/dummy").To(dummy))
	restful.Add(ws)
	restful.Filter(filter.NCSACommonLogFormatLogger())
	restful.Filter(filter.ShopKEYFilter("back"))

	httpRequest, _ := http.NewRequest("OPTIONS", "http://here.io/dummy", nil)
	httpRequest.Header.Add("X-SHOP-KEY", "a8428177b08fceb199fde1bb92ea23a539f94dd563c9a1b1616bdc9c5cae32ca6dd91f29c1d3c1399b")
	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)
	shopid := httpWriter.Header().Get("X-Shop-Id")
	if "0446d745-6333-4126-9dc0-57e3033613a6" != shopid {
		t.Fatal("expected: X-Shop-Id but back shopid:" + shopid)
	}
}

func TestShopKEYMiniFilter(t *testing.T) {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/dummy").To(dummy))
	restful.Add(ws)
	restful.Filter(filter.NCSACommonLogFormatLogger())
	restful.Filter(filter.ShopKEYFilter("mini"))

	httpRequest, _ := http.NewRequest("OPTIONS", "http://here.io/dummy", nil)
	httpRequest.Header.Add("X-SHOP-KEY", "a74a8c75b08fceb199fde1bb92ea23a59dd1971d85605c5fd06fb85b4a2b37cee1e10f1f6f5d630cd1")
	httpWriter := httptest.NewRecorder()
	restful.DefaultContainer.Dispatch(httpWriter, httpRequest)
	shopid := httpWriter.Header().Get("X-Shop-Id")
	if "0446d745-6333-4126-9dc0-57e3033613a6" != shopid {
		t.Fatal("expected: X-Shop-Id but back shopid:" + shopid)
	}
}

func dummy(i *restful.Request, o *restful.Response) {
}
