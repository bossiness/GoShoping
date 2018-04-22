package mongodb

import (
	"time"

	"btdxcx.com/micro/shop-srv/db"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	proto "btdxcx.com/micro/shop-srv/proto/shop"
	dproto "btdxcx.com/micro/shop-srv/proto/shop/details"
)

// Mongo DB
type Mongo struct {
	session *mgo.Session
}

var (
	// DBUrl mongodb URL
	DBUrl = "localhost:27017"
)

const (
	databaseName = "center"
	tblKey       = "shop-key"
	tblDetails   = "shop-details"
	tblOwner     = "shop-owner"
	tblMini      = "shop-mini"
	tblPhysical  = "shop-physical"
)

func init() {
	db.Register(new(Mongo))
}

// Init 数据库初始化
func (m *Mongo) Init() error {
	session, err := mgo.Dial(DBUrl)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	m.session = session
	return nil
}

// ShopKey DB
type ShopKey struct {
	ID     bson.ObjectId     `bson:"_id,omitempty"`
	UUID   string            `bson:"uuid"`
	TagKey map[string]string `bson:"tag_key"`
}

// ReadKey form uuid
func (m *Mongo) ReadKey(uuid string) (*proto.ShopTagKeys, error) {
	c := m.session.DB(databaseName).C(tblKey)

	result := &ShopKey{}
	if err := c.Find(bson.M{"uuid": uuid}).One(result); err != nil {
		return nil, err
	}

	shopKey := &proto.ShopTagKeys{
		Tagkeys: result.TagKey,
	}

	return shopKey, nil
}

// CreateKey form uuid
func (m *Mongo) CreateKey(uuid string, proto *proto.ShopTagKeys) error {
	c := m.session.DB(databaseName).C(tblKey)

	index := mgo.Index{
		Key:        []string{"uuid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return err
	}

	shopKey := &ShopKey{
		ID:     bson.NewObjectId(),
		UUID:   uuid,
		TagKey: proto.Tagkeys,
	}
	return c.Insert(shopKey)
}

// DeleteKey form uuid
func (m *Mongo) DeleteKey(uuid string) error {
	c := m.session.DB(databaseName).C(tblKey)
	return c.Remove(bson.M{"uuid": uuid})
}

// ShopDetails DB
type ShopDetails struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	UUID      string        `bson:"uuid"`
	CreateAt  int64         `bson:"create_at"`
	UpdateAt  int64         `bson:"update_at"`
	PeriodAt  int64         `bson:"period_at"`
	SubmitAt  int64         `bson:"submit_at"`
	State     dproto.State  `bson:"state"`
	Name      string        `bson:"name"`
	Logo      string        `bson:"logo"`
	Introduce string        `bson:"introduce"`

	Owner    mgo.DBRef `bson:"owner"`
	Mini     mgo.DBRef `bson:"mini"`
	Physical mgo.DBRef `bson:"physical"`
}

// ShopOwner owner
type ShopOwner struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Nickname string        `bson:"nickname"`
	Phone    string        `bson:"phone"`
}

// MiniApp 小程序信息
type MiniApp struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	IDWecht string        `bson:"wechat_id"`
}

// PhysicalStore 实体店信息
type PhysicalStore struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Contact  string        `bson:"contact"`
	Email    string        `bson:"email"`
	ZipCode  string        `bson:"zipCode"`
	Address  string        `bson:"address"`
	Location Location      `bson:"location"`
}

// Location 经纬度
type Location struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
}

// CreateDetails create
func (m *Mongo) CreateDetails(req *dproto.CreateRequest) (*dproto.CreateResponse, error) {
	c := m.session.DB(databaseName).C(tblDetails)
	ownerC := m.session.DB(databaseName).C(tblOwner)
	miniC := m.session.DB(databaseName).C(tblMini)
	physicalC := m.session.DB(databaseName).C(tblPhysical)

	ownerRef := mgo.DBRef{}
	if req.Details.Owner != nil {

		index := mgo.Index{
			Key:        []string{"name"},
			Unique:     true,
			DropDups:   true,
			Background: true,
		}

		if err := ownerC.EnsureIndex(index); err != nil {
			return nil, err
		}

		ownerID := bson.NewObjectId()
		owner := &ShopOwner{
			ID:       ownerID,
			Name:     req.Details.Owner.Name,
			Nickname: req.Details.Owner.Nickname,
			Phone:    req.Details.Owner.Phone,
		}
		if err := ownerC.Insert(owner); err != nil {
			return nil, err
		}
		ownerRef.Collection = tblOwner
		ownerRef.Id = ownerID
		ownerRef.Database = databaseName
	}

	miniRef := mgo.DBRef{}
	if req.Details.Mimi != nil {
		miniID := bson.NewObjectId()
		mini := &MiniApp{
			ID:      miniID,
			IDWecht: req.Details.Mimi.WechatId,
		}
		if err := miniC.Insert(mini); err != nil {
			return nil, err
		}
		miniRef.Collection = tblMini
		miniRef.Id = miniID
		miniRef.Database = databaseName
	}

	physicalRef := mgo.DBRef{}
	if req.Details.Physical != nil {
		physicalID := bson.NewObjectId()
		var (
			latitude  float64
			longitude float64
		)
		if req.Details.Physical.Location != nil {
			latitude = req.Details.Physical.Location.Latitude
			longitude = req.Details.Physical.Location.Longitude
		}
		physical := &PhysicalStore{
			ID:      physicalID,
			Name:    req.Details.Physical.Name,
			Contact: req.Details.Physical.Contact,
			Email:   req.Details.Physical.Email,
			ZipCode: req.Details.Physical.ZipCode,
			Address: req.Details.Physical.Address,
			Location: Location{
				Latitude:  latitude,
				Longitude: longitude,
			},
		}
		if err := physicalC.Insert(physical); err != nil {
			return nil, err
		}
		physicalRef.Collection = tblPhysical
		physicalRef.Id = physicalID
		physicalRef.Database = databaseName
	}

	shopID := bson.NewObjectId()
	id, err2 := uuid.NewV4()
	if err2 != nil {
		return nil, err2
	}

	details := &ShopDetails{
		ID:        shopID,
		UUID:      id.String(),
		CreateAt:  time.Now().Unix(),
		UpdateAt:  time.Now().Unix(),
		State:     dproto.State_untreated,
		Name:      req.Details.Name,
		Logo:      req.Details.Logo,
		Introduce: req.Details.Introduce,
		Owner:     ownerRef,
		Mini:      miniRef,
		Physical:  physicalRef,
	}
	if err := c.Insert(details); err != nil {
		return nil, err
	}

	return &dproto.CreateResponse{ShopId: id.String()}, nil
}

// ReadDetails read
func (m *Mongo) ReadDetails(uuid string) (*dproto.ReadResponse, error) {
	c := m.session.DB(databaseName).C(tblDetails)

	details := &ShopDetails{}
	if err := c.Find(bson.M{"uuid": uuid}).One(details); err != nil {
		return nil, err
	}
	return m.readDetails(details)
}

func (m *Mongo) readDetails(details *ShopDetails) (*dproto.ReadResponse, error) {
	owner := &ShopOwner{}
	if details.Owner.Id != nil {
		clOwner := m.session.DB(details.Owner.Database).C(details.Owner.Collection)
		clOwner.FindId(details.Owner.Id).One(owner)
	}

	mini := &MiniApp{}
	if details.Mini.Id != nil {
		clMini := m.session.DB(details.Mini.Database).C(details.Mini.Collection)
		clMini.FindId(details.Mini.Id).One(mini)
	}

	physical := &PhysicalStore{}
	if details.Physical.Id != nil {
		clPhysical := m.session.DB(details.Physical.Database).C(details.Physical.Collection)
		clPhysical.FindId(details.Physical.Id).One(physical)
	}

	response := &dproto.ReadResponse{
		ShopId:   details.UUID,
		CreateAt: details.CreateAt,
		UpdateAt: details.UpdateAt,
		PeriodAt: details.PeriodAt,
		SubmitAt: details.SubmitAt,
		State:    details.State,
		Details: &dproto.ShopDetails{
			Name:      details.Name,
			Logo:      details.Logo,
			Introduce: details.Introduce,
			Owner: &dproto.ShopDetails_ShopOwner{
				Id:       owner.ID.Hex(),
				Name:     owner.Name,
				Nickname: owner.Nickname,
				Phone:    owner.Phone,
			},
			Mimi: &dproto.ShopDetails_MiniApp{
				Id:       mini.ID.Hex(),
				WechatId: mini.IDWecht,
			},
			Physical: &dproto.ShopDetails_PhysicalStore{
				Id:      physical.ID.Hex(),
				Name:    physical.Name,
				Contact: physical.Contact,
				Email:   physical.Email,
				ZipCode: physical.ZipCode,
				Address: physical.Address,
				Location: &dproto.ShopDetails_PhysicalStore_Location{
					Latitude:  physical.Location.Latitude,
					Longitude: physical.Location.Longitude,
				},
			},
		},
	}

	return response, nil
}

// DeleteDetails delete
func (m *Mongo) DeleteDetails(uuid string) error {
	c := m.session.DB(databaseName).C(tblDetails)

	details := &ShopDetails{}
	if err := c.Find(bson.M{"uuid": uuid}).One(details); err != nil {
		return err
	}

	if cl := m.session.DB(details.Owner.Database).C(details.Owner.Collection); cl != nil {
		cl.RemoveId(details.Owner.Id)
	}

	if cl := m.session.DB(details.Mini.Database).C(details.Mini.Collection); cl != nil {
		cl.RemoveId(details.Mini.Id)
	}

	if cl := m.session.DB(details.Physical.Database).C(details.Physical.Collection); cl != nil {
		cl.RemoveId(details.Physical.Id)
	}

	return c.RemoveId(details.ID)
}

// UpdateDetails update
func (m *Mongo) UpdateDetails(req *dproto.UpdateRequest) error {
	c := m.session.DB(databaseName).C(tblDetails)

	selector := bson.M{"uuid": req.ShopId}

	perioAt, submitAt := int64(0), int64(0)
	if req.State == dproto.State_reviewing {
		submitAt = time.Now().Unix()
	} else if req.State == dproto.State_completed {
		d, _ := time.ParseDuration("24h")
		perioAt = time.Now().Add(d * 365).Unix()
	}

	shop := &ShopDetails{}
	c.Find(selector).One(shop)
	if req.Details != nil {
		if req.Details.Owner != nil {
			owner := req.Details.Owner
			if shop.Owner.Id != nil {
				if c := m.session.DB(shop.Owner.Database).C(shop.Owner.Collection); c != nil {
					updataData := bson.M{"$set": bson.M{
						"nickname": owner.Nickname,
						"phone":    owner.Phone,
					}}
					c.UpdateId(shop.Owner.Id, updataData)
				}
			}
		}
		if req.Details.Mimi != nil {
			mimi := req.Details.Mimi
			if shop.Mini.Id != nil {
				if c := m.session.DB(shop.Mini.Database).C(shop.Mini.Collection); c != nil {
					updataData := bson.M{"$set": bson.M{
						"wechat_id": mimi.WechatId,
					}}
					c.UpdateId(shop.Mini.Id, updataData)
				}
			}
		}
		if req.Details.Physical != nil {
			physical := req.Details.Physical
			if shop.Physical.Id != nil {
				if physical.Location == nil {
					physical.Location = &dproto.ShopDetails_PhysicalStore_Location{}
				}
				if c := m.session.DB(shop.Physical.Database).C(shop.Physical.Collection); c != nil {
					updataData := bson.M{"$set": bson.M{
						"name":               physical.Name,
						"contact":            physical.Contact,
						"email":              physical.Email,
						"zipCode":            physical.ZipCode,
						"address":            physical.Address,
						"location.latitude":  physical.Location.Latitude,
						"location.longitude": physical.Location.Longitude,
					}}
					c.UpdateId(shop.Physical.Id, updataData)
				}
			}
		}

		updataData := bson.M{"$set": bson.M{
			"update_at": time.Now().Unix(),
			"submit_at": submitAt,
			"period_at": perioAt,
			"state":     req.State,
			"name":      req.Details.Name,
			"logo":      req.Details.Logo,
			"introduce": req.Details.Introduce,
		}}
		return c.Update(selector, updataData)

	}

	updataData := bson.M{"$set": bson.M{
		"update_at": time.Now().Unix(),
		"submit_at": submitAt,
		"period_at": perioAt,
		"state":     req.State,
	}}
	return c.Update(selector, updataData)
	
}

// ListDetails list
func (m *Mongo) ListDetails(req *dproto.ListRequest) (*dproto.ListResponse, error) {
	c := m.session.DB(databaseName).C(tblDetails)

	start := int(req.Start)
	limit := int(req.Limit)
	result := &[]ShopDetails{}
	if err := c.Find(nil).Skip(start).Limit(limit).All(result); err != nil {
		return nil, err
	}

	items := []*dproto.ReadResponse{}
	for _, item := range *result {
		details, err := m.readDetails(&item)
		if err != nil {
			return nil, err
		}
		items = append(items, details)
	}

	response := &dproto.ListResponse{
		Start: req.Start,
		Limit: req.Limit,
		Total: int32(len(*result)),
		Items: items,
	}

	return response, nil
}
