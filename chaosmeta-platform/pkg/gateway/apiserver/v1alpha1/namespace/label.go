/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package namespace

import (
	"chaosmeta-platform/pkg/service/namespace"
	"context"
	"encoding/json"
)

func (c *NamespaceController) ListLabel() {
	id, _ := c.GetInt(":id")
	name := c.GetString("name")
	orderBy := c.GetString("sort")
	creator := c.GetString("creator")
	page, _ := c.GetInt("page", 1)
	pageSize, _ := c.GetInt("page_size", 10)
	namespace := &namespace.NamespaceService{}
	total, labelList, err := namespace.ListLabel(context.Background(), id, name, creator, orderBy, page, pageSize)
	if err != nil {
		c.Error(&c.Controller, err)
		return
	}
	c.Success(&c.Controller, LabelListResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Labels:   labelList,
	})
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

	id, err := namespace.CreateLabel(context.Background(), nsId, username, reqBody.Name, reqBody.Color)
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

func (c *NamespaceController) LabelGet() {
	nsId, _ := c.GetInt(":ns_id")
	name := c.Ctx.Input.Param(":name")

	namespace := &namespace.NamespaceService{}
	label, err := namespace.GetLabelByName(context.Background(), nsId, name)
	if err != nil {
		c.Success(&c.Controller, nil)
		return
	}
	c.Success(&c.Controller, label)
}
