package mongodb

import (
	"btdxcx.com/micro/taxons-srv/db"
	proto "btdxcx.com/micro/taxons-srv/proto/taxons"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Mongo DB
type Mongo struct {
	session *mgo.Session
}

var (
	// URL mongodb URL
	URL = "localhost:27017"
)

const (
	taxonsCollectionName = "taxons"
	imageCollectionName  = "images"
)

// Image DB
type Image struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	Code string        `bson:"code"`
	Path string        `bson:"path"`
}

// Taxons DB
type Taxons struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Position    int32         `bson:"position"`
	Images      []mgo.DBRef   `bson:"images"`
	Parent      string        `bson:"parent_id"`
}

func init() {
	db.Register(new(Mongo))
}

// Init 数据库初始化
func (m *Mongo) Init() error {
	session, err := mgo.Dial(URL)
	if err != nil {
		return err
	}
	session.SetMode(mgo.Monotonic, true)
	m.session = session
	return nil
}

// Deinit 资源释放
func (m *Mongo) Deinit() {
	if m.session != nil {
		m.session.Close()
	}
}

// Read 读取数据
func (m *Mongo) Read(dbname string) (*proto.TaxonsMessage, error) {

	c := m.session.DB(dbname).C(taxonsCollectionName)
	ic := m.session.DB(dbname).C(imageCollectionName)

	return m.readRoot(c, ic)
}

func (m *Mongo) readRoot(c *mgo.Collection, ic *mgo.Collection) (*proto.TaxonsMessage, error) {

	root := &proto.TaxonsMessage{
		Code:        "root",
		Name:        "root",
		Description: "root taxons",
	}

	taxons, err := m.readRootTaxons(c)
	if err != nil {
		return nil, err
	}

	for _, t := range *taxons {
		root.Children = append(root.Children, taxons2message(&t, ic))
	}

	m.traverse(&root.Children, c, ic)

	return root, nil
}

func (m *Mongo) traverse(children *[]*proto.TaxonsMessage, c *mgo.Collection, ic *mgo.Collection) {

	for index := 0; index < len(*children); index++ {
		curr := (*children)[index]
		taxons, e := m.readTaxons(c, curr.Code)
		if e == nil {
			for _, t := range *taxons {
				curr.Children = append(curr.Children, taxons2message(&t, ic))
			}

			m.traverse(&curr.Children, c, ic)
		}
	}
}

func taxons2message(t *Taxons, ic *mgo.Collection) *proto.TaxonsMessage {
	result := &proto.TaxonsMessage{
		Code:        t.ID.Hex(),
		Name:        t.Name,
		Description: t.Description,
		Position:    t.Position,
	}
	for _, item := range t.Images {
		image := new(Image)
		err := ic.FindId(item.Id).One(image)
		if err == nil {
			id := image.ID.Hex()
			impImage := proto.Image{Id: id, Code: image.Code, Path: image.Path}
			result.Images = append(result.Images, &impImage)
		}
	}

	return result
}

func (m *Mongo) readRootTaxons(c *mgo.Collection) (*[]Taxons, error) {
	return m.readTaxons(c, "")
}

func (m *Mongo) readTaxons(c *mgo.Collection, parentID string) (*[]Taxons, error) {
	result := []Taxons{}
	err := c.Find(bson.M{"parent_id": parentID}).Sort("position").All(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *Mongo) readOneTaxons(c *mgo.Collection, parentID string) (*Taxons, error) {
	result := &Taxons{}

	err := c.Find(bson.M{"parent_id": parentID}).One(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Delete 删除数据
func (m *Mongo) Delete(dbname string, id string) error {
	c := m.session.DB(dbname).C(taxonsCollectionName)
	return c.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}

// Create 创造数据
func (m *Mongo) Create(dbname string, data *proto.TaxonsMessage) (string, error) {

	c := m.session.DB(dbname).C(taxonsCollectionName)
	ic := m.session.DB(dbname).C(imageCollectionName)

	imagesRef := []mgo.DBRef{}
	for index := 0; index < len(data.Images); index++ {
		_image := data.Images[index]
		tid := bson.NewObjectId()
		if err := ic.Insert(&Image{tid, _image.Code, _image.Path}); err != nil {
			return "", err
		}
		imagesRef = append(imagesRef, mgo.DBRef{imageCollectionName, tid, dbname})
	}

	index := mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return "", err
	}

	_id := bson.NewObjectId()
	if err := c.Insert(&Taxons{
		ID:          _id,
		Name:        data.Name,
		Description: data.Description,
		Position:    data.Position,
		Images:      imagesRef,
		Parent:      data.Code,
	}); err != nil {
		return "", err
	}

	return _id.Hex(), nil
}

// Update 更新数据
func (m *Mongo) Update(dbname string, data *proto.TaxonsMessage) error {

	c := m.session.DB(dbname).C(taxonsCollectionName)
	ic := m.session.DB(dbname).C(imageCollectionName)
	if err := m.update(c, dbname, data); err != nil {
		return err
	}
	return m.updateImages(c, ic, dbname, data)
}

func (m *Mongo) update(c *mgo.Collection, dbname string, data *proto.TaxonsMessage) error {

	selector := bson.M{"_id": bson.ObjectIdHex(data.Code)}

	updataData := bson.M{"$set": bson.M{
		"name":        data.Name,
		"position":    data.Position,
		"description": data.Description,
	}}

	err := c.Update(selector, updataData)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mongo) updateImages(c *mgo.Collection, ic *mgo.Collection, dbname string, data *proto.TaxonsMessage) error {

	if len(data.Images) == 0 {
		return nil
	}

	selector := bson.M{"_id": bson.ObjectIdHex(data.Code)}

	taxons := &Taxons{}
	if err := c.Find(selector).One(taxons); err != nil {
		return err
	}
	for _, image := range taxons.Images {
		if err := ic.RemoveId(image.Id); err != nil {
			return err
		}
	}

	imagesRef := []mgo.DBRef{}
	for index := 0; index < len(data.Images); index++ {
		_image := data.Images[index]
		tid := bson.NewObjectId()
		if err := ic.Insert(&Image{tid, _image.Code, _image.Path}); err != nil {
			return err
		}
		imagesRef = append(imagesRef, mgo.DBRef{imageCollectionName, tid, dbname})
	}

	updataData := bson.M{"$set": bson.M{"images": imagesRef}}
	err := c.Update(selector, updataData)
	if err != nil {
		return err
	}
	return nil
}
