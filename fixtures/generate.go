package fixtures

//go:generate protoc --descriptor_set_out=fileset.pb --include_imports --include_source_info -I. -I../tmp/googleapis -I../thirdparty Booking.proto Vehicle.proto nested/Book.proto

//go:generate protoc --descriptor_set_out=cookie.pb --include_imports --include_source_info -I. -I../thirdparty Cookie.proto
