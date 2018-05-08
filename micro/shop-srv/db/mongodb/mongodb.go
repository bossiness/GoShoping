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
	tblWX        = "shop-wx"
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
	WeiXin   mgo.DBRef `bson:"weixin"`
	Physical mgo.DBRef `bson:"physical"`
}

// ShopOwner owner
type ShopOwner struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Name     string        `bson:"name"`
	Nickname string        `bson:"nickname"`
	Phone    string        `bson:"phone"`
}

// WeiXin 小程序信息
type WeiXin struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	WechatID   string        `bson:"wechat_id,omitempty"`
	Appid      string        `bson:"appid,omitempty"`
	AppSecret  string        `bson:"app_secret,omitempty"`
	PartnerID  string        `bson:"partner_id,omitempty"`
	PartnerKey string        `bson:"partner_key,omitempty"`
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
	wxC := m.session.DB(databaseName).C(tblWX)
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

	wxRef := mgo.DBRef{}
	if req.Details.Weixin != nil {
		id := bson.NewObjectId()
		weixin := &WeiXin{
			ID:         id,
			WechatID:   req.Details.Weixin.WechatId,
			Appid:      req.Details.Weixin.Appid,
			AppSecret:  req.Details.Weixin.AppSecret,
			PartnerID:  req.Details.Weixin.PartnerId,
			PartnerKey: req.Details.Weixin.PartnerKey,
		}
		if err := wxC.Insert(weixin); err != nil {
			return nil, err
		}
		wxRef.Collection = tblWX
		wxRef.Id = id
		wxRef.Database = databaseName
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
		WeiXin:    wxRef,
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

	mini := &WeiXin{}
	if details.WeiXin.Id != nil {
		clMini := m.session.DB(details.WeiXin.Database).C(details.WeiXin.Collection)
		clMini.FindId(details.WeiXin.Id).One(mini)
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
			Weixin: &dproto.ShopDetails_WeiXin{
				Id:         mini.ID.Hex(),
				WechatId:   mini.WechatID,
				Appid:      mini.Appid,
				AppSecret:  mini.AppSecret,
				PartnerId:  mini.PartnerID,
				PartnerKey: mini.PartnerKey,
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

	if cl := m.session.DB(details.WeiXin.Database).C(details.WeiXin.Collection); cl != nil {
		cl.RemoveId(details.WeiXin.Id)
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
	if err := c.Find(selector).One(shop); err != nil {
		return err
	}
	if req.Details != nil {
		m.updateDetails(shop, req.Details, req.State,
			perioAt, submitAt, selector)

	}

	updataData := bson.M{"$set": bson.M{
		"update_at": time.Now().Unix(),
		"submit_at": submitAt,
		"period_at": perioAt,
		"state":     req.State,
	}}
	return c.Update(selector, updataData)

}

func (m *Mongo) updateDetails(
	shop *ShopDetails, details *dproto.ShopDetails, state dproto.State,
	perioAt int64, submitAt int64,
	selector bson.M,
) error {
	c := m.session.DB(databaseName).C(tblDetails)

	m.updateDetailsOwner(shop, details.Owner, selector)
	m.updateDetailsWeixin(shop, details.Weixin, selector)
	m.updateDetailsPhysical(shop, details.Physical, selector)

	updataData := bson.M{"$set": bson.M{
		"update_at": time.Now().Unix(),
		"submit_at": submitAt,
		"period_at": perioAt,
		"state":     state,
		"name":      details.Name,
		"logo":      details.Logo,
		"introduce": details.Introduce,
	}}
	return c.Update(selector, updataData)
}

func (m *Mongo) updateDetailsOwner(shop *ShopDetails, sowner *dproto.ShopDetails_ShopOwner, selector bson.M) {
	if sowner == nil {
		return
	}

	if shop.Owner.Id != nil && len(shop.Owner.Database) > 0 && len(shop.Owner.Collection) > 0 {
		if oc := m.session.DB(shop.Owner.Database).C(shop.Owner.Collection); oc != nil {
			updataData := bson.M{"$set": bson.M{
				"nickname": sowner.Nickname,
				"phone":    sowner.Phone,
			}}
			oc.UpdateId(shop.Owner.Id, updataData)
		}
		return
	}

	ownerC := m.session.DB(databaseName).C(tblOwner)
	c := m.session.DB(databaseName).C(tblDetails)

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := ownerC.EnsureIndex(index); err != nil {
		return
	}

	ownerID := bson.NewObjectId()
	owner := &ShopOwner{
		ID:       ownerID,
		Name:     sowner.Name,
		Nickname: sowner.Nickname,
		Phone:    sowner.Phone,
	}
	if err := ownerC.Insert(owner); err != nil {
		return
	}
	ownerRef := mgo.DBRef{
		Collection: tblOwner,
		Id:         ownerID,
		Database:   databaseName,
	}

	updataData := bson.M{"$set": bson.M{
		"owner": ownerRef,
	}}

	c.Update(selector, updataData)

}

func (m *Mongo) updateDetailsWeixin(shop *ShopDetails, sweixin *dproto.ShopDetails_WeiXin, selector bson.M) {
	if sweixin == nil {
		return
	}

	if shop.WeiXin.Id != nil && len(shop.WeiXin.Database) > 0 && len(shop.WeiXin.Collection) > 0 {
		if wc := m.session.DB(shop.WeiXin.Database).C(shop.WeiXin.Collection); wc != nil {
			updataData := bson.M{"$set": bson.M{
				"wechat_id":   sweixin.WechatId,
				"appid":       sweixin.Appid,
				"app_secret":  sweixin.AppSecret,
				"partner_id":  sweixin.PartnerId,
				"partner_key": sweixin.PartnerKey,
			}}
			wc.UpdateId(shop.WeiXin.Id, updataData)
		}
		return
	}

	wxC := m.session.DB(databaseName).C(tblWX)
	c := m.session.DB(databaseName).C(tblDetails)

	id := bson.NewObjectId()
	weixin := &WeiXin{
		ID:         id,
		WechatID:   sweixin.WechatId,
		Appid:      sweixin.Appid,
		AppSecret:  sweixin.AppSecret,
		PartnerID:  sweixin.PartnerId,
		PartnerKey: sweixin.PartnerKey,
	}
	if err := wxC.Insert(weixin); err != nil {
		return
	}
	wxRef := mgo.DBRef{
		Collection: tblWX,
		Id:         id,
		Database:   databaseName,
	}

	updataData := bson.M{"$set": bson.M{
		"weixin": wxRef,
	}}

	c.Update(selector, updataData)
}

func (m *Mongo) updateDetailsPhysical(shop *ShopDetails, sphysical *dproto.ShopDetails_PhysicalStore, selector bson.M) {
	if sphysical == nil {
		return
	}

	if shop.Physical.Id != nil && len(shop.Physical.Database) > 0 && len(shop.Physical.Collection) > 0 {
		if sphysical.Location == nil {
			sphysical.Location = &dproto.ShopDetails_PhysicalStore_Location{}
		}
		if c := m.session.DB(shop.Physical.Database).C(shop.Physical.Collection); c != nil {
			updataData := bson.M{"$set": bson.M{
				"name":               sphysical.Name,
				"contact":            sphysical.Contact,
				"email":              sphysical.Email,
				"zipCode":            sphysical.ZipCode,
				"address":            sphysical.Address,
				"location.latitude":  sphysical.Location.Latitude,
				"location.longitude": sphysical.Location.Longitude,
			}}
			c.UpdateId(shop.Physical.Id, updataData)
		}
		return
	}

	physicalC := m.session.DB(databaseName).C(tblPhysical)
	c := m.session.DB(databaseName).C(tblDetails)

	physicalID := bson.NewObjectId()
	var (
		latitude  float64
		longitude float64
	)
	if sphysical.Location != nil {
		latitude = sphysical.Location.Latitude
		longitude = sphysical.Location.Longitude
	}
	physical := &PhysicalStore{
		ID:      physicalID,
		Name:    sphysical.Name,
		Contact: sphysical.Contact,
		Email:   sphysical.Email,
		ZipCode: sphysical.ZipCode,
		Address: sphysical.Address,
		Location: Location{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}
	if err := physicalC.Insert(physical); err != nil {
		return
	}

	physicalRef := mgo.DBRef{
		Collection: tblPhysical,
		Id:         physicalID,
		Database:   databaseName,
	}

	updataData := bson.M{"$set": bson.M{
		"physical": physicalRef,
	}}

	c.Update(selector, updataData)
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
