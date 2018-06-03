package handler

import (
	"btdxcx.com/walker/apis/common/errors"
	"github.com/emicklei/go-restful"

	"btdxcx.com/walker/model"

	"net/http"

	"btdxcx.com/walker/service/auth"
	"golang.org/x/net/context"
)

// AuthHandler Auth handler
type AuthHandler struct {
	Service auth.IService
}

// Signin 登陆
func (h *AuthHandler) Signin(req *restful.Request, rsp *restful.Response) {
	in := new(model.AuthRequest)
	if err := req.ReadEntity(&in); err != nil {
		errors.Response(rsp, errors.BadRequest("handler.signup", "request error [%v]", err))
		return
	}

	if len(in.Username) == 0 {
		errors.Response(rsp, errors.BadRequest("handler.signup", "username empty"))
		return
	}

	if len(in.Password) == 0 {
		errors.Response(rsp, errors.BadRequest("handler.signup", "password empty."))
		return
	}

	in.Type = req.HeaderParameter("X-Tag")
	in.ShopID = req.HeaderParameter("X-Shop-Id")

	out := &model.Token{}
	if err := h.Service.Signin(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.signup", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
}

// Signup 注册
func (h *AuthHandler) Signup(req *restful.Request, rsp *restful.Response) {

	in := new(model.AuthRequest)
	if err := req.ReadEntity(&in); err != nil {
		errors.Response(rsp, errors.BadRequest("handler.signup", "request error [%v]", err))
		return
	}

	if len(in.Username) == 0 {
		errors.Response(rsp, errors.BadRequest("handler.signup", "username empty"))
		return
	}

	if len(in.Password) == 0 {
		errors.Response(rsp, errors.BadRequest("handler.signup", "password empty."))
		return
	}

	in.Type = req.HeaderParameter("X-Tag")
	in.ShopID = req.HeaderParameter("X-Shop-Id")

	out := &model.Token{}
	if err := h.Service.Create(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.signup", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
}

// Signout 登出
func (h *AuthHandler) Signout(req *restful.Request, rsp *restful.Response) {
	clientID := req.HeaderParameter("X-Username")

	in := &model.Introspect{
		Username: clientID,
	}
	out := &model.NoContent{}
	if err := h.Service.Signout(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.signout", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}

// Password 修改密码
func (h *AuthHandler) Password(req *restful.Request, rsp *restful.Response) {
}

// Token 刷新token
func (h *AuthHandler) Token(req *restful.Request, rsp *restful.Response) {

}
