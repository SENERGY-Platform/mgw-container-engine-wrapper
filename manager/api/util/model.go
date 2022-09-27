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

package util

type ErrResponse struct {
	Status string   `json:"status"`
	Errors []string `json:"errors"`
}

type ContainersQuery struct {
	Name  string   `form:"name"`
	State string   `form:"state"`
	Label []string `form:"label"`
}

type ContainerLogQuery struct {
	MaxLines int    `form:"max_lines"`
	Since    string `form:"since"`
	Until    string `form:"until"`
}

type ImagesQuery struct {
	Label []string `form:"label"`
}

type PostContainerResponse struct {
	ID string `json:"id"`
}
