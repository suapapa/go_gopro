gen_pb:
	cd protobuf && protoc --go_out=../open_gopro --go_opt=paths=source_relative ./*.proto