package controllers

import (
	"encoding/json"

	"github.com/beego/beego/utils/pagination"
	"github.com/casibase/casibase/object"
	"github.com/casibase/casibase/util"
)

// GetTagLinks
// @Title GetTagLinks
// @Tag TagLink API
// @Description get tag links
// @Success 200 {array} object.TagLink The Response object
// @router /get-tag-links [get]
func (c *ApiController) GetTagLinks() {
	owner := c.Input().Get("owner")
	tag := c.Input().Get("tag")
	limit := c.Input().Get("pageSize")
	page := c.Input().Get("p")

	if limit == "" || page == "" {
		links, err := object.GetTagLinks(owner, tag)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}

		c.ResponseOk(links)
	} else {
		limitInt := util.ParseInt(limit)
		count, err := object.GetTagLinkCount(owner, tag)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		paginator := pagination.SetPaginator(c.Ctx, limitInt, count)
		links, err := object.GetPaginationTagLinks(owner, paginator.Offset(), limitInt, tag)
		if err != nil {
			c.ResponseError(err.Error())
			return
		}
		c.ResponseOk(links, paginator.Nums())
	}
}

// AddTagLink
// @Title AddTagLink
// @Tag TagLink API
// @Description add tag link
// @Param body body object.TagLink true "The details of the tag link"
// @Success 200 {object} controllers.Response The Response object
// @router /add-tag-link [post]
func (c *ApiController) AddTagLink() {
	var link object.TagLink
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &link)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	success, err := object.AddTagLink(&link)
	if err != nil {
		c.ResponseError(err.Error())
		return
	}

	c.ResponseOk(success)
}
