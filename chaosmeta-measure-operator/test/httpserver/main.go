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

package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	r := gin.Default()

	r.GET("/testget", func(c *gin.Context) {
		data := gin.H{
			"message": "Hello, World!",
			"test": map[string]any{
				"status": 0,
				"at":     "aa",
			},
		}

		c.JSON(http.StatusOK, data)
	})

	r.POST("/testpost", func(c *gin.Context) {
		var reqBody RequestBody
		if err := c.BindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 构造一个 JSON 返回体
		data := gin.H{
			"a": map[string]any{
				"b": map[string]any{
					"name": reqBody.Name,
					"age":  reqBody.Age,
				},
			},
		}

		c.JSON(http.StatusOK, data)
	})

	r.Run(":8080")
}
