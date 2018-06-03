package adminuser

import (
	"github.com/micro/go-micro/errors"
	"btdxcx.com/walker/model"
	"btdxcx.com/walker/repository/adminuser"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
)

// IService adminuser service
type IService interface {
	Create(ctx context.Context, req *model.AdminUser, resq *model.NoContent) error
	Reads(ctx context.Context, req *model.PageRequest, resq *model.AdminUsersPage) error
	Read(ctx context.Context, req *model.IDRequest, resq *model.AdminUsersRecord) error
	Update(ctx context.Context, req *model.AdminUser, resq *model.NoContent) error
	Delete(ctx context.Context, req *model.IDRequest, resq *model.NoContent) error
	ReadProfile(ctx context.Context, req *model.UsernameRequest, resq *model.AdminUsersRecord) error
}

// Service adminuser service
type Service struct {
	Session *mgo.Session
}

// GetRepo get  repository
func (s *Service) GetRepo() adminuser.IRepository {
	return &adminuser.Repository{
		Session: s.Session.Clone(),
	}
}

// Create a new admin user
func (s *Service) Create(ctx context.Context, req *model.AdminUser, resq *model.NoContent) error {
	repo := s.GetRepo()
	defer repo.Close()

	if len(req.Username) == 0 {
		return errors.BadRequest("walker.service.admin.create", "username Can't be empty")
	} 

	if err := repo.Create(req); err != nil {
		return errors.InternalServerError("walker.service.admin.create", "%v", err)
	}

	return nil
}

// Reads admin user list
func (s *Service) Reads(ctx context.Context, req *model.PageRequest, resq *model.AdminUsersPage) error {
	repo := s.GetRepo()
	defer repo.Close()

	users, total, err := repo.GetList(int(req.Offset), int(req.Limit))
	if err != nil {
		return errors.NotFound("walker.service.admin.reads", "%v", err)
	}

	resq.Offset = req.Offset
	resq.Limit = req.Limit
	resq.Total = int32(total)

	records := []*model.AdminUser{}
	for _, user := range *users {
		records = append(records, &user)
	}
	resq.Records = records

	return nil
}

func (s *Service) Read(ctx context.Context, req *model.IDRequest, resq *model.AdminUsersRecord) error {
	repo := s.GetRepo()
	defer repo.Close()

	user, err := repo.Get(req.ID)
	if err != nil {
		return errors.NotFound("walker.service.admin.read", "%v", err)
	}

	resq.Record = user
	return nil
}

// Update adminuser
func (s *Service) Update(ctx context.Context, req *model.AdminUser, resq *model.NoContent) error {
	repo := s.GetRepo()
	defer repo.Close()

	if err := repo.Update(req.ID.Hex(), req); err != nil {
		return errors.NotFound("walker.service.admin.update", "%v", err)
	}

	return nil
}

// Delete adminuser
func (s *Service) Delete(ctx context.Context, req *model.IDRequest, resq *model.NoContent) error {
	repo := s.GetRepo()
	defer repo.Close()

	if err := repo.Delete(req.ID); err != nil {
		return errors.NotFound("walker.service.admin.delete", "%v", err)
	}

	return nil
}

// ReadProfile read profile
func (s *Service) ReadProfile(ctx context.Context, req *model.UsernameRequest, resq *model.AdminUsersRecord) error {
	repo := s.GetRepo()
	defer repo.Close()

	if len(req.Username) == 0 {
		return errors.BadRequest("walker.service.admin.create", "username Can't be empty")
	} 

	user, err := repo.GetProfile(req.Username)
	if err != nil {
		return errors.NotFound("walker.service.admin.read", "%v", err)
	}

	resq.Record = user
	return nil
}
