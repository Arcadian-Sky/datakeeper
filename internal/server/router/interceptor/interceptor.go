package interceptor

import (
	"context"
	"fmt"

	"github.com/Arcadian-Sky/datakkeeper/internal/server/router/jwtrule"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// we don't need to check the token for these methods.
	SkipCheckMethods = map[string]struct{}{
		"/proto.api.user.v1.UserService/Register":     {},
		"/proto.api.user.v1.UserService/Authenticate": {},
	}
)

func AuthCheckGRPC(log *logrus.Logger, secretKey string) grpc.UnaryServerInterceptor {
	fmt.Printf("\"!!!\": %v\n", "!!!")
	fmt.Printf(log.GetLevel().String())
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		log.Trace("--> unary interceptor: ", info.FullMethod)
		// check for method, which doesn't need to be intercepted
		_, ok := SkipCheckMethods[info.FullMethod]
		if ok {
			log.Trace("--> unary interceptor: pass")
			return handler(ctx, req)
		}
		log.Trace("--> unary interceptor: check")
		// check part
		token, err := grpc_auth.AuthFromMD(ctx, "bearer")
		if err != nil {
			return nil, err
		}

		log.Trace("--> unary interceptor: check")
		log.Debug("token: ", token)
		log.Debug(jwtrule.Validate(token, secretKey))
		jwToken, err := jwtrule.Validate(token, secretKey)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
		}
		log.Trace("--> unary interceptor UID: ", jwToken.Claims.UserID)

		return handler(jwtrule.SetUserIDToCTX(ctx, int(jwToken.Claims.UserID)), req)
	}
}

// AuthCheckGRPC interceptor verifies the authentication bearer token.
// func AuthCheckGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

// }

// // GetUnaryServerInterceptor returns server unary Interceptor to authenticate and authorize unary RPC
// func GetUnaryServerInterceptor(jwtKey string) grpc.UnaryServerInterceptor {
// 	return func(
// 		ctx context.Context,
// 		req interface{},
// 		info *grpc.UnaryServerInfo,
// 		handler grpc.UnaryHandler,
// 	) (interface{}, error) {

// 		if SkipCheckMethods[info.FullMethod] {
// 			token, err := parseJTW(ctx, jwtKey)
// 			if err != nil {
// 				return nil, err
// 			}
// 			ctx = metadata.AppendToOutgoingContext(ctx, "userid", fmt.Sprintf("%d", token.Claims.UserID))
// 		}

// 		return handler(ctx, req)
// 	}
// }

// // parseJWT  checks JWT and parses it
// func parseJTW(ctx context.Context, jwtKey string) (model.Jtoken, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return model.Jtoken{}, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
// 	}

// 	authHeader, ok := md["authorization"]
// 	if !ok {
// 		return model.Jtoken{}, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
// 	}

// 	token, err := jwtrule.Validate(authHeader[0], jwtKey)

// 	if err != nil {
// 		return model.Jtoken{}, status.Errorf(codes.Unauthenticated, err.Error())
// 	}
// 	return token, nil
// }
