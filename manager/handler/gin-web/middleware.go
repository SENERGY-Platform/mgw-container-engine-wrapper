/*
 * Copyright 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package gin_web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func checkStaticOrigin(origins []string) gin.HandlerFunc {
	var isOk = func(so string) bool {
		for _, o := range origins {
			if so == o {
				return true
			}
		}
		return false
	}
	return func(c *gin.Context) {
		if so := c.GetHeader("X-Static-Origin"); so != "" {
			if !isOk(so) {
				c.AbortWithStatus(http.StatusForbidden)
			}
		} else {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}
