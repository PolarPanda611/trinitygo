package runtime

import (
	"context"
	"fmt"

	"github.com/PolarPanda611/trinitygo/application"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// New  for runtime
// if required is true , and runtime key is not set , return codes.InvalidArgument
// if required is fale , runtime key is null , set default ""
// initial runtime info , shoould be first interceptor
func New(app application.Application) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		for _, v := range app.RuntimeKeys() {
			if _, ok := md[v.GetKeyName()]; !ok {
				if v.GetRequired() {
					line := fmt.Sprintf("%v %v %v ", app.Conf().GetProjectName(), app.Conf().GetProjectVersion(), info.FullMethod)
					errMessage := fmt.Sprintf("%v required but not found", v.GetKeyName())
					line += fmt.Sprintf("%v %v ", "Error", errMessage)
					app.Logger().Error(line)
					return nil, status.Errorf(codes.InvalidArgument, errMessage)
				}
				md.Append(v.GetKeyName(), "")
			}
		}
		resp, err := handler(ctx, req)
		return resp, err
	}
}
