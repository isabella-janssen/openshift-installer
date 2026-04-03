//go:build !(okd || scos)

package rhcos

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/types"
)

const (
	// DefaultOSImageStream is the OS image stream used when the install-config
	// does not specify one.
	DefaultOSImageStream = types.OSImageStreamRHCOS9

	payloadImageStreamTagRHCOS9  = "rhel-coreos"
	payloadImageStreamTagRHCOS10 = "rhel-coreos-10"
)

func getStreamFileName(stream types.OSImageStream) string {
	if stream == "" {
		stream = DefaultOSImageStream
	}
	// All rhel-10 variants (rhel-10, rhel-10-nvidia, etc.) use the rhel-10 boot images
	if strings.HasPrefix(string(stream), "rhel-10") && stream != types.OSImageStreamRHCOS10 {
		stream = types.OSImageStreamRHCOS10
	}
	return fmt.Sprintf("coreos/coreos-%v.json", stream)
}

func getMarketplaceStreamFileName(stream types.OSImageStream) string {
	if stream == "" {
		stream = DefaultOSImageStream
	}
	// All rhel-10 variants (rhel-10, rhel-10-nvidia, etc.) use the rhel-10 boot images
	if strings.HasPrefix(string(stream), "rhel-10") && stream != types.OSImageStreamRHCOS10 {
		stream = types.OSImageStreamRHCOS10
	}
	return fmt.Sprintf("coreos/marketplace/coreos-%v.json", stream)
}

// GetPayloadImageStreamTag returns the payload image stream tag corresponding
// to the given OS image stream.
// Note: All rhel-10 variants (rhel-10-nvidia, etc.) use rhel-coreos-10 tag since they
// share the same boot images and only differ in OS stream configuration.
func GetPayloadImageStreamTag(stream types.OSImageStream) string {
	if stream == "" || stream == types.OSImageStreamRHCOS9 {
		return payloadImageStreamTagRHCOS9
	}
	// rhel-10 and all rhel-10-* variants use rhel-coreos-10 tag
	return payloadImageStreamTagRHCOS10
}
