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
)

// CreateVolume provisions a new volume on behalf of the user
func CreateVolume() {
	fmt.Print("IMPLEMENT ME")
}

// DeleteVolume deletes a volume created by CreateVolume
func DeleteVolume() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerPublishVolume makes a volume available on a specified node. This will attach a volume to the node
func ControllerPublishVolume() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerUnpublishVolume makes the volume unavailable on a given node. Makes a call to detach a volume from a node
func ControllerUnpublishVolume() {
	fmt.Print("IMPLEMENT ME")
}

// ValidateVolumeCapabilities returns capabilities of a volume
func ValidateVolumeCapabilities() {
	fmt.Print("IMPLEMENT ME")
}

// ListVolumes returns all available volumes
func ListVolumes() {
	fmt.Print("IMPLEMENT ME")
}

// GetCapacity returns capacity of total storage pool available
func GetCapacity() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerGetCapabilities returns the capabilities of the Controller plugin
func ControllerGetCapabilities() {
	fmt.Print("IMPLEMENT ME")
}

// CreateSnapshot creates a snapshot
func CreateSnapshot() {
	fmt.Print("IMPLEMENT ME")
}

// DeleteSnapshot deletes a snapshot
func DeleteSnapshot() {
	fmt.Print("IMPLEMENT ME")
}

// ListSnapshots lists snapshots
func ListSnapshots() {
	fmt.Print("IMPLEMENT ME")
}

// ControllerExpandVolume ...
func ControllerExpandVolume() {
	fmt.Print("IMPLEMENT ME")
}
