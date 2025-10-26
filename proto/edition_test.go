package proto

import (
	"strings"
	"testing"

	cardv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/card/v1"
	ddv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/dd/v1"
	securityv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/security/v1"
	tachographv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/v1"
	vuv1 "github.com/way-platform/tachograph-go/proto/gen/go/wayplatform/connect/tachograph/vu/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Test_editions(t *testing.T) {
	var (
		_ = &tachographv1.File{}
		_ = &vuv1.VehicleUnitFile{}
		_ = &cardv1.DriverCardFile{}
		_ = &securityv1.Authentication{}
		_ = &ddv1.StringValue{}
	)
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		if !strings.HasPrefix(string(fd.Path()), "wayplatform/connect/tachograph/") {
			return true
		}
		if fd.Syntax() != protoreflect.Editions {
			t.Errorf("File %s has syntax %s, expected Editions", fd.Path(), fd.Syntax())
		}
		return true
	})
}
