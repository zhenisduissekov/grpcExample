# grpcExample
dummy service


Instructions (url: https://www.youtube.com/watch?v=YudT0nHvkkE)
1. create three directories
2. go mod init example.com/go-usermgmt-grpc
3. create usermgmt.proto
4. protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto
5. if there is an error protoc-gen-go-grpc: program not found or is not executable
then:  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
6. file usermgmt.pb.go and usermgmt_grpc.pb.go created
7. write usermgmt_server/usermgmt_server.go and usermgt_client/usermgmt_client.go
8. run them go run usermgmt_client/usermgmt_client.go and go run usermgmt_server/usermgmt_server.go


modifying

9. modified usermgmt.proto: added messages and service rpc GetUsers(GetUserParams) returns (UserList) {};
10. modified _server.go
