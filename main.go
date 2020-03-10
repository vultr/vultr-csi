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
	_ "context"
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/vultr/govultr"
	driver "github.com/vultr/vultr-csi/vultr"
)

const (
	driverName = "vultrbs.csi.vultr.com"
)

var (
	vendorVersion string

	bsPrefix = flag.String("bs-prefix", "", "Vultr BlockStorage Volume label prefix")
	endpoint = flag.String("endpoint", "unix:///var/lib/kubelet/plugins/"+driverName+"/csi.sock", "Vultr CSI endpoint")
	node     = flag.String("node", "", "Vultr Hostname")
	token    = flag.String("token", "", "Vultr API Token")
	url      = flag.String("url", "https://vultr.com/api", "Vultr API URL")
)

func main() {
	rand.Seed(time.Now().UnixNano())
	flag.Parse()
	handle()
	os.Exit(0)
}

func handle() {
	// Check API token
	apiToken := os.Getenv("VULTR_API_KEY")
	if apiToken == "" {
		apiToken := *token
		if apiToken == "" {
			log.Fatalln("You must provide your Vultr API token")
		}
	}

	// Get Vultr client
	vultrClient := govultr.NewClient(nil, apiToken)

	// Initialize driver
	drv, err := driver.NewDriver(vultrClient, driverName, vendorVersion, *bsPrefix, *url, *endpoint)
	if err != nil {
		log.Fatalln(err)
	}

	// Run the service
	endpoint := "unix:///tmp/echo.sock" // to test locally
	drv.Run(endpoint)
}
