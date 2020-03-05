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
	"fmt"
	// "github.com/container-storage-interface/spec/lib/go/csi"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

func CreateVolume() {
	fmt.Print("IMPLEMENT ME")
}

func DeleteVolume() {
	fmt.Print("IMPLEMENT ME")
}

func ControllerPublishVolume() {
	fmt.Print("IMPLEMENT ME")
}

func ControllerUnpublishVolume() {
	fmt.Print("IMPLEMENT ME")
}

func ValidateVolumeCapabilities() {
	fmt.Print("IMPLEMENT ME")
}

func ListVolumes() {
	fmt.Print("IMPLEMENT ME")
}

func GetCapacity() {
	fmt.Print("IMPLEMENT ME")
}

func ControllerGetCapabilities() {
	fmt.Print("IMPLEMENT ME")
}

func CreateSnapshot() {
	fmt.Print("IMPLEMENT ME")
}

func DeleteSnapshot() {
	fmt.Print("IMPLEMENT ME")
}

func ListSnapshots() {
	fmt.Print("IMPLEMENT ME")
}

func ControllerExpandVolume() {
	fmt.Print("IMPLEMENT ME")
}
