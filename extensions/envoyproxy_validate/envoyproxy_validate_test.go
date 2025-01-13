package extensions_test

import (
	"testing"

	"github.com/envoyproxy/protoc-gen-validate/validate"
	"github.com/moia-oss/protoc-gen-doc/extensions"
	. "github.com/moia-oss/protoc-gen-doc/extensions/lyft_validate"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestTransform(t *testing.T) {
	fieldRules := &validate.FieldRules{
		Type: &validate.FieldRules_String_{
			String_: &validate.StringRules{
				MinLen: proto.Uint64(1),
				NotIn:  []string{"invalid"},
			},
		},
	}
	dynamicFieldRules, err := extensions.ConvertToDynamicMessage(fieldRules)
	require.NoError(t, err)

	transformed := extensions.Transform(map[string]interface{}{"validate.rules": dynamicFieldRules})
	require.NotEmpty(t, transformed)

	rules := transformed["validate.rules"].(ValidateExtension).Rules()
	require.Equal(t, rules, []ValidateRule{
		{Name: "string.min_len", Value: uint64(1)},
		{Name: "string.not_in", Value: []string{"invalid"}},
	})
}
