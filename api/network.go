/*
 * Copyright 2023 InfAI (CC SES)
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

package api

import (
	"container-engine-wrapper/model"
	"context"
)

func (a *Api) GetNetworks(ctx context.Context) ([]model.Network, error) {
	return a.ceHandler.ListNetworks(ctx)
}

func (a *Api) GetNetwork(ctx context.Context, id string) (model.Network, error) {
	return a.ceHandler.NetworkInfo(ctx, id)
}

func (a *Api) CreateNetwork(ctx context.Context, net model.Network) error {
	return a.ceHandler.NetworkCreate(ctx, net)
}

func (a *Api) RemoveNetwork(ctx context.Context, id string) error {
	return a.ceHandler.NetworkRemove(ctx, id)
}
