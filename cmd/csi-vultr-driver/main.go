/*
Copyright 2020 Vultr.
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

package main

import (
	"flag"
	"log"

	"github.com/vultr/vultr-csi/driver"
)

var version string

func main() {

	var (
		endpoint   = flag.String("endpoint", "unix:///var/lib/kubelet/plugins/"+driver.DefaultDriverName+"/csi.sock", "CSI endpoint")
		token      = flag.String("token", "", "Vultr API Token")
		apiURL     = flag.String("api-url", "", "Vultr API URL")
		driverName = flag.String("driver-name", driver.DefaultDriverName, "Name of driver")
		userAgent  = flag.String("user-agent", "", "Custom user agent")
	)
	flag.Parse()

	if version == "" {
		log.Fatal("version must be defined at compilation")
	}

	d, err := driver.NewDriver(*endpoint, *token, *driverName, version, *userAgent, *apiURL)
	if err != nil {
		log.Fatalln(err)
	}

	d.Run()
}
