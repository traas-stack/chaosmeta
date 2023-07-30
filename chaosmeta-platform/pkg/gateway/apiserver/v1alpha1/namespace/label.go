package namespace

import (
	"chaosmeta-platform/pkg/service/namespace"
	"context"
	"encoding/json"
	"fmt"
)

func (c *NamespaceController) ListLabel() {
	id, _ := c.GetInt(":id")
	name := c.GetString("name")
	orderBy := c.GetString("sort")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	namespace := &namespace.NamespaceService{}
	total, usrList, err := namespace.ListLabel(context.Background(), id, name, orderBy, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	fmt.Println(total, usrList)
}

func (c *NamespaceController) LabelCreate() {
	nsId, _ := c.GetInt(":id")
	namespace := &namespace.NamespaceService{}
	username := c.Ctx.Input.GetData("userName").(string)
	var reqBody LabelCreateRequest
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqBody); err != nil {
		c.Error(&c.Controller, err)
		return
	}

	id, err := namespace.CreateLabel(context.Background(), nsId, username, reqBody.Name)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, LabelCreateResponse{Id: id})
}

func (c *NamespaceController) LabelDelete() {
	nsId, _ := c.GetInt(":ns_id")
	id, _ := c.GetInt(":id")
	username := c.Ctx.Input.GetData("userName").(string)

	namespace := &namespace.NamespaceService{}
	if err := namespace.DeleteLabel(context.Background(), nsId, username, id); err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, "ok")
}
