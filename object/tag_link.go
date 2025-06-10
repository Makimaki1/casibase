package object

import (
	"fmt"

	"github.com/casibase/casibase/util"
	"xorm.io/core"
)

// TagLink connects a Casibase tag with an external resource
// such as a GitHub issue, bookmark or forum thread.
type TagLink struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Tag         string `xorm:"varchar(100) index" json:"tag"`
	Type        string `xorm:"varchar(100)" json:"type"`
	Url         string `xorm:"varchar(500)" json:"url"`
	Description string `xorm:"varchar(500)" json:"description"`
}

func GetTagLinks(owner, tag string) ([]*TagLink, error) {
	links := []*TagLink{}
	session := adapter.engine.Desc("created_time")
	if owner != "" {
		session = session.Where("owner = ?", owner)
	}
	if tag != "" {
		session = session.And("tag = ?", tag)
	}
	err := session.Find(&links)
	if err != nil {
		return links, err
	}
	return links, nil
}

func GetPaginationTagLinks(owner string, offset, limit int, tag string) ([]*TagLink, error) {
	links := []*TagLink{}
	session := adapter.engine.Desc("created_time")
	if owner != "" {
		session = session.Where("owner = ?", owner)
	}
	if tag != "" {
		session = session.And("tag = ?", tag)
	}
	err := session.Limit(limit, offset).Find(&links)
	if err != nil {
		return links, err
	}
	return links, nil
}

func getTagLink(owner, name string) (*TagLink, error) {
	if owner == "" || name == "" {
		return nil, nil
	}

	link := TagLink{Owner: owner, Name: name}
	existed, err := adapter.engine.Get(&link)
	if err != nil {
		return &link, err
	}

	if existed {
		return &link, nil
	} else {
		return nil, nil
	}
}

func GetTagLink(id string) (*TagLink, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getTagLink(owner, name)
}

func GetTagLinkCount(owner, tag string) (int64, error) {
	session := adapter.engine.Where("1=1")
	if owner != "" {
		session = session.And("owner = ?", owner)
	}
	if tag != "" {
		session = session.And("tag = ?", tag)
	}
	return session.Count(&TagLink{})
}

func UpdateTagLink(id string, link *TagLink) (bool, error) {
	owner, name := util.GetOwnerAndNameFromId(id)
	_, err := adapter.engine.ID(core.PK{owner, name}).AllCols().Update(link)
	if err != nil {
		return false, err
	}
	return true, nil
}

func AddTagLink(link *TagLink) (bool, error) {
	affected, err := adapter.engine.Insert(link)
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func DeleteTagLink(link *TagLink) (bool, error) {
	affected, err := adapter.engine.ID(core.PK{link.Owner, link.Name}).Delete(&TagLink{})
	if err != nil {
		return false, err
	}

	return affected != 0, nil
}

func (link *TagLink) GetId() string {
	return fmt.Sprintf("%s/%s", link.Owner, link.Name)
}
