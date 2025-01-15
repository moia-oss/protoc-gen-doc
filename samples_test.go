package gendoc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_generateSampleWithRepeatedField(t *testing.T) {
	mapEntryMessage := &Message{
		FullName:  "a.b.c.MapEntry",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "key", FullType: "string"},
			{Name: "value", FullType: "string"},
		},
	}

	baseMsg := &Message{
		FullName:  "a.b.c.Base",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "double_field", FullType: "double"},
			{Name: "float_field", FullType: "float"},
			{Name: "map_field", IsMap: true, FullType: "a.b.c.MapEntry"},
		},
	}

	simpleMsg := &Message{
		FullName:  "a.b.c.Simple",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "double_field", FullType: "double"},
		},
	}

	repeatedMsg := &Message{
		FullName:  "a.b.c.RepeatedField",
		HasFields: true,
		Fields: []*MessageField{
			{Name: "repeated_strings", FullType: "string", Label: "repeated"},
		},
	}

	baseEnum := &Enum{
		FullName: "a.b.c.Enum",
		Values: []*EnumValue{
			{Name: "value1", Number: "1"},
			{Name: "value2", Number: "2"},
			{Name: "value3", Number: "3"},
		},
	}

	file := &File{
		Package:     "a.b.c",
		HasEnums:    true,
		HasMessages: true,
		Enums:       orderedEnums{baseEnum},
		Messages:    orderedMessages{mapEntryMessage, baseMsg, simpleMsg, repeatedMsg},
	}

	message := &Message{
		FullName:  "a.b.c.Test",
		HasFields: true,
		HasOneofs: true,
		Fields: []*MessageField{
			{Name: "double_field", FullType: "double"},
			{Name: "float_field", FullType: "float"},
			{Name: "int64_field", FullType: "int64"},
			{Name: "uint32_field", FullType: "uint32"},
			{Name: "uint64_field", FullType: "uint64"},
			{Name: "sint32_field", FullType: "sint32"},
			{Name: "sint64_field", FullType: "sint64"},
			{Name: "fixed32_field", FullType: "fixed32"},
			{Name: "fixed64_field", FullType: "fixed64"},
			{Name: "sfixed32_field", FullType: "sfixed32"},
			{Name: "sfixed64_field", FullType: "sfixed64"},
			{Name: "bool_field", FullType: "bool", IsOneof: true, OneofDecl: "test_oneof"},
			{Name: "string_field", FullType: "string", IsOneof: true, OneofDecl: "test_oneof"},
			{Name: "bytes_field", FullType: "bytes"},
			{Name: "message_field", FullType: "a.b.c.Base"},
			{Name: "enum_field", FullType: "a.b.c.Enum"},
			{Name: "repeated_field", FullType: "string", Label: "repeated"},
			{Name: "repeated_msg_field", FullType: "a.b.c.Simple", Label: "repeated"},
		},
	}

	got := generateSample(message, []*File{file})
	expect := map[string]interface{}{
		"double_field":   1.23,
		"float_field":    1.23,
		"int64_field":    "0",
		"uint32_field":   123,
		"uint64_field":   "0",
		"sint32_field":   123,
		"sint64_field":   "0",
		"fixed32_field":  123,
		"fixed64_field":  "0",
		"sfixed32_field": 123,
		"sfixed64_field": "0",
		"string_field":   "{{string}}",
		"bool_field":     false,
		//"string_field":   "string", // this field is omitted because of oneof
		"bytes_field": "{{binary bytes}}",
		"message_field": map[string]interface{}{
			"double_field": 1.23,
			"float_field":  1.23,
			"map_field": map[string]interface{}{
				"{{key}}": "{{string}}",
			},
		},
		"enum_field": 3,
		"repeated_field": []interface{}{
			"{{string}}",
		},
		"repeated_msg_field": []interface{}{
			map[string]interface{}{
				"double_field": 1.23,
			},
		},
	}

	// test := map[interface{}]interface{}{
	// 	"map_field": map[string]interface{}{
	// 		"{{string}}": "{{string}}",
	// 	},
	// }
	// indent := strings.Repeat(" ", 2)
	// sampleJSON, err := json.MarshalIndent(test, "", indent)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(sampleJSON))

	require.Equal(t, expect, got)
}
