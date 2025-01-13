package thirdparty

//go:generate mkdir -p github.com/envoyproxy/protoc-gen-validate/validate
//go:generate curl -fsSL https://github.com/envoyproxy/protoc-gen-validate/raw/master/validate/validate.proto -o github.com/envoyproxy/protoc-gen-validate/validate/validate.proto

//go:generate mkdir -p github.com/moia-oss/protokit/fixtures
//go:generate curl -fsSL https://github.com/moia-oss/protokit/raw/master/fixtures/extend.proto -o github.com/moia-oss/protokit/fixtures/extend.proto
