package tools_test

import (
	"testing"

	"btdxcx.com/walker/tools"
	"github.com/satori/go.uuid"
)

func TestCodingCenter(t *testing.T) {
	sk := tools.ShopKEY{}
	ret, err := sk.Coding("center")
	if err != nil {
		t.Error("shop key coding center, err ", err)
	}
	if ret != "a9468c6895cd" {
		t.Error("shop key coding center, got", ret)
	}
}

func TestDecodingCenter(t *testing.T) {
	sk := tools.ShopKEY{}
	ret, err := sk.Decoding("a9468c6895cd")
	if err != nil {
		t.Error("shop key decoding center, err ", err)
	}
	if ret != "center" {
		t.Error("shop key decoding center, got", ret)
	}
}

func TestCodingTags(t *testing.T) {
	sk := tools.ShopKEY{}
	uid, err := uuid.NewV4()
	if err != nil {
		t.Error("shop key coding tags, new uuid ", err)
		return
	}
	ret, err := sk.CodingWithTags([]string{"back", "mini"}, uid.String())
	if err != nil {
		t.Error("shop key coding tags, err ", err)
	}
	t.Log(ret)
}

func TestDecodingBack(t *testing.T) {
	sk := tools.ShopKEY{}
	const key = "a8428177b08fceb199fde1bb92ea23a539f94dd563c9a1b1616bdc9c5cae32ca6dd91f29c1d3c1399b"
	back, _, err := sk.DecodingWithTag(key)
	if err != nil {
		t.Error("shop key decoding tags, err ", err)
	}
	if back != "back" {
		t.Error("shop key decoding tags, back ", back)
	}
}

func TestDecodingMini(t *testing.T) {
	sk := tools.ShopKEY{}
	const key = "a74a8c75b08fceb199fde1bb92ea23a59dd1971d85605c5fd06fb85b4a2b37cee1e10f1f6f5d630cd1"
	mini, _, err := sk.DecodingWithTag(key)
	if err != nil {
		t.Error("shop key decoding tags, err ", err)
	}
	if mini != "mini" {
		t.Error("shop key decoding tags, mini ", mini)
	}
}
