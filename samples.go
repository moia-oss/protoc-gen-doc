package gendoc

import (
	"encoding/json"
	"strconv"
	"strings"
)

const (
	BOOL   = "bool"
	BYTES  = "bytes"
	DOUBLE = "double"
	INT64  = "int64"
	INT32  = "int32"
	FLOAT  = "float"
	STRING = "string"
)

var samplesCache = map[string]interface{}{}
var scalarTypes map[string]struct{}
var wellKnownTypeDefaultSamples = map[string]interface{}{
	"google.protobuf.Any": "{{any}}",
	// google.protobuf.Api
	"google.protobuf.BoolValue": map[string]interface{}{
		"value": generateSampleScalarValue(BOOL),
	},
	"google.protobuf.BytesValue": map[string]interface{}{
		"value": generateSampleScalarValue(BYTES),
	},
	"google.protobuf.DoubleValue": map[string]interface{}{
		"value": generateSampleScalarValue(DOUBLE),
	},
	"google.protobuf.Duration": map[string]interface{}{
		"seconds": generateSampleScalarValue(INT64),
		"nanos":   generateSampleScalarValue(INT32),
	},
	"google.protobuf.Empty": struct{}{},
	"google.protobuf.FieldMask": map[string]string{
		"mask": "a.b.c,foo",
	},
	"google.protobuf.FloatValue": map[string]interface{}{
		"value": generateSampleScalarValue(FLOAT),
	},
	"google.protobuf.Int32Value": map[string]interface{}{
		"value": generateSampleScalarValue(INT32),
	},
	"google.protobuf.Int64Value": map[string]interface{}{
		"value": generateSampleScalarValue(INT64),
	},
	"google.protobuf.ListValue": map[string]interface{}{
		"value": []string{
			"{{any-value}}", // TODO: replace with Value
		},
	},
	// google.protobuf.Method
	// google.protobuf.Mixin
	// google.protobuf.NullValue
	// google.protobuf.Option
	// google.protobuf.SourceContext
	"google.protobuf.StringValue": map[string]interface{}{
		"value": generateSampleScalarValue(STRING),
	},
	// google.protobuf.Struct
	"google.protobuf.Timestamp": map[string]interface{}{
		"seconds": generateSampleScalarValue(INT64),
		"nanos":   generateSampleScalarValue(INT32),
	},
	//google.protobof.Type
	"google.protobuf.UInt32Value": map[string]interface{}{
		"value": generateSampleScalarValue("uint32"),
	},
	"google.protobuf.UInt64Value": map[string]interface{}{
		"value": generateSampleScalarValue("uint64"),
	},
	// google.protobuf.Value
}

func init() {
	scalars := makeScalars()
	scalarTypes = map[string]struct{}{}
	for _, scalar := range scalars {
		scalarTypes[scalar.ProtoType] = struct{}{}
	}
}

func findMessageInFiles(messageName string, files []*File) *Message {
	for _, file := range files {
		msg := file.GetMessage(messageName)
		if msg != nil {
			return msg
		}
	}
	return nil
}

func findEnumInFiles(enumName string, files []*File) *Enum {
	for _, file := range files {
		enum := file.GetEnum(enumName)
		if enum != nil {
			return enum
		}
	}
	return nil
}

// SampleGenerator generate sample json data for specified message
func SampleGenerator(messageName string, files []*File, indentWidth int) string {
	message := findMessageInFiles(messageName, files)
	sampleData := generateSample(message, files)
	indent := strings.Repeat(" ", indentWidth)
	sampleJSON, err := json.MarshalIndent(sampleData, "", indent)
	if err != nil {
		// TODO:
		panic(err)
	}
	return string(sampleJSON)
}

func generateSample(message *Message, files []*File) interface{} {
	if cached, found := samplesCache[message.FullName]; found {
		return cached
	}

	if !message.HasFields {
		return struct{}{}
	}

	// TODO: add support for one of fields

	result := map[string]interface{}{}
	for _, field := range message.Fields {
		// TODO: changed field.FullType to message.FullName
		fieldValue := generateFieldValue(field.FullType, field.IsMap, files)
		samplesCache[field.FullType] = result
		setMapField(result, field, fieldValue)
	}

	return result
}

func generateFieldValue(typeName string, isMap bool, files []*File) interface{} {
	if isScalarType(typeName) {
		return generateSampleScalarValue(typeName)
	}

	if sampleValue, ok := wellKnownTypeDefaultSamples[typeName]; ok {
		return sampleValue
	}
	if isMap {
		return generateSampleMap(typeName, files)
	}
	fieldMessage := findMessageInFiles(typeName, files)
	if fieldMessage != nil {
		return generateSample(fieldMessage, files)
	}
	fieldEnum := findEnumInFiles(typeName, files)
	if fieldEnum != nil {
		return generateEnum(fieldEnum)
	}
	// TODO: Lookup in imports
	return nil
}

func setMapField(m map[string]interface{}, field *MessageField, value interface{}) {
	if field.Label == "repeated" {
		m[field.Name] = []interface{}{value}
		return
	}
	m[field.Name] = value
}

func generateEnum(enum *Enum) interface{} {
	if len(enum.Values) == 0 {
		return 0
	}
	number := enum.Values[len(enum.Values)-1].Number
	integer, _ := strconv.Atoi(number)
	return integer
}

func generateSampleMap(typeName string, files []*File) interface{} {
	message := findMessageInFiles(typeName, files)
	var keyType, valueType string
	for _, field := range message.Fields {
		if field.Name == "key" {
			// According to proto3 definition, keytype can only be a string or integral
			keyType = field.FullType
		}
		if field.Name == "value" {
			valueType = field.FullType
		}
	}

	sampleValue := generateFieldValue(valueType, false, files)

	// we need to initialize the result in every case because
	// map[interface{}]interface{} cannot be marshalized.
	switch keyType {
	case INT32, "uint32", "sint32", "fixed32", "sfixed32":
		result := map[int]interface{}{
			123: sampleValue,
		}
		return result
	case INT64, "uint64", "sint64", "fixed64", "sfixed64":
		return "0"
	case BOOL:
		result := map[bool]interface{}{
			true: sampleValue,
		}
		return result
	case STRING:
		result := map[string]interface{}{
			"{{key}}": sampleValue,
		}
		return result
	}
	return nil
}

func generateSampleScalarValue(typeName string) interface{} {
	switch typeName {
	case DOUBLE, FLOAT:
		return 1.23
	case INT32, "uint32", "sint32", "fixed32", "sfixed32":
		return 123
	case INT64, "uint64", "sint64", "fixed64", "sfixed64":
		return "0"
	case BOOL:
		return false
	case STRING:
		return "{{string}}"
	case BYTES:
		return "{{binary bytes}}"
	default:
		return "{{unknown scalar value}}"
	}
}

func isScalarType(typeName string) bool {
	_, is := scalarTypes[typeName]
	return is
}
