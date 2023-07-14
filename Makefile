gen_pb:
	cd open_gopro/protobuf && protoc --go_out=../ --go_opt=paths=source_relative ./*.proto