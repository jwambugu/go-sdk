gen_proto:
	protoc -I elarian-proto/ elarian-proto/web.proto \
	--go_out=plugins=grpc:com_elarian_hera_proto
	protoc -I elarian-proto/ elarian-proto/common.proto \
	--go_out=plugins=grpc:com_elarian_hera_proto


run:
	go run .