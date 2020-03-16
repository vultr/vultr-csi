/*
Copyright 2020 Vultr Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"fmt"

	"github.com/vultr/govultr"
)

func GetVultrByName(client *govultr.Client, name string) (*govultr.Server, error) {
	instances, err := client.Server.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error while getting instance list: %s", err)
	}

	for _, v := range instances {
		if v.Label == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("could not retrieve instance: %s", name)
}
