package gendoc_test

import (
	"testing"

	. "github.com/moia-oss/protoc-gen-doc"
	"github.com/moia-oss/protokit/utils"
	"github.com/stretchr/testify/require"
)

func BenchmarkParseCodeRequest(b *testing.B) {
	set, _ := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	req := utils.CreateGenRequest(set, "Booking.proto", "Vehicle.proto")
	plugin := new(Plugin)

	for i := 0; i < b.N; i++ {
		_, err := plugin.Generate(req)
		require.NoError(b, err)
	}
}
