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

package itf

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (p Port) String() string {
	return fmt.Sprintf("%d/%s", p.Number, p.Protocol)
}

func ParseRestartStrategy(v string) (RestartStrategy, error) {
	for i := 0; i < len(restartStrategyStr); i++ {
		if restartStrategyStr[i] == v {
			return RestartStrategy(i), nil
		}
	}
	return 0, errors.New(fmt.Sprintf("unknown restart strategy '%s'", v))
}

func (s RestartStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(restartStrategyStr[s])
}

func (s *RestartStrategy) UnmarshalJSON(data []byte) (err error) {
	var v string
	if err = json.Unmarshal(data, &v); err != nil {
		return
	}
	*s, err = ParseRestartStrategy(v)
	return
}
