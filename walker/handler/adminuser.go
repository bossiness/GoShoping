package handler

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"golang.org/x/net/context"
	"btdxcx.com/walker/model"
	"btdxcx.com/walker/service/adminuser"
	"github.com/emicklei/go-restful"
	"btdxcx.com/walker/apis/common/errors"

	"net/http"
)

// AdminUserHandler adminuser handler
type AdminUserHandler struct {
	Service adminuser.IService
}

// Create adminuser
func (h *AdminUserHandler) Create(req *restful.Request, rsp *restful.Response) {
	in := new(model.AdminUser)
	if err := req.ReadEntity(&in); err != nil {
		errors.Response(rsp, errors.BadRequest("handler.adminuser.create", "[%v]", err))
		return
	}

	out := &model.NoContent{}
	if err := h.Service.Create(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.adminuser.create", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}

	rsp.WriteHeader(http.StatusCreated)
}

// Reads read list
func (h *AdminUserHandler) Reads(req *restful.Request, rsp *restful.Response) {
	in := new(model.PageRequest)
	in.Offset, in.Limit = offsetlimit(req)

	out := &model.AdminUsersPage{}
	if err := h.Service.Reads(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.adminuser.reads", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}
}


func (h *AdminUserHandler) Read(req *restful.Request, rsp *restful.Response) {
	in := new(model.IDRequest)
	in.ID = req.PathParameter("id")

	out := &model.AdminUsersRecord{}
	if err := h.Service.Read(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		err1 := errors.InternalServerError("handler.adminuser.read", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}
}

// Update adminuser
func (h *AdminUserHandler) Update(req *restful.Request, rsp *restful.Response) {
	in := new(model.AdminUser)
	if err := req.ReadEntity(&in); err != nil {
		errors.Response(rsp, errors.BadRequest("handler.adminuser.create", "[%v]", err))
		return
	}
	in.ID = bson.ObjectIdHex(req.PathParameter("id"))

	out := &model.NoContent{}
	if err := h.Service.Update(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.adminuser.update", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}

// Delete a adminuser
func (h *AdminUserHandler) Delete(req *restful.Request, rsp *restful.Response) {
	in := new(model.IDRequest)
	in.ID = req.PathParameter("id")

	out := &model.NoContent{}
	if err := h.Service.Delete(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out); err != nil {
		err1 := errors.InternalServerError("handler.adminuser.delete", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}

	rsp.WriteHeader(http.StatusNoContent)
}

// ReadProfile read profile
func (h *AdminUserHandler) ReadProfile(req *restful.Request, rsp *restful.Response) {
	in := new(model.UsernameRequest)
	in.Username = req.PathParameter("username")

	out := &model.AdminUsersRecord{}
	if err := h.Service.ReadProfile(context.TODO(), in, out); err != nil {
		errors.Response(rsp, err)
		return
	}

	if err := rsp.WriteEntity(out.Record); err != nil {
		err1 := errors.InternalServerError("handler.adminuser.delete", "write entiry err: ", err)
		errors.Response(rsp, err1)
		return
	}
}

func offsetlimit(req *restful.Request) (int32, int32) {
	offset, err1 := strconv.Atoi(req.QueryParameter("offset"))
	if err1 != nil {
		offset = 0
	}
	limit, err2 := strconv.Atoi(req.QueryParameter("limit"))
	if err2 != nil {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	return int32(offset), int32(limit)
}
