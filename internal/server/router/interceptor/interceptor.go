package interceptor

import (
	"context"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
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
	PostProcessMethods = map[string]struct{}{
		"/proto.api.service.v1.DataKeeperService/UploadFile": {},
	}
)

func UnaryInterceptor(log *logrus.Logger, secretKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		log.Trace("--> unary interceptor: ", info.FullMethod)

		preProcess(ctx, info.FullMethod, log, secretKey)

		jwToken, err := checkAuth(&ctx, log, secretKey, info.FullMethod, nil)
		if err != nil {
			return ctx, err
		} else if jwToken == nil {
			resp, err = handler(ctx, req)
		} else {
			resp, err = handler(jwtrule.SetUserIDToCTX(ctx, int(jwToken.Claims.UserID)), req)
		}

		postProcess(ctx, info.FullMethod, log, err)

		log.Trace("--> finish unary interceptor: ", info.FullMethod)
		return resp, err
	}
}

func StreamInterceptor(log *logrus.Logger, secretKey string) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx := ss.Context()
		log.Trace("--> stream interceptor: ", info.FullMethod)

		preProcess(ctx, info.FullMethod, log, secretKey)

		jwToken, err := checkAuth(&ctx, log, secretKey, info.FullMethod, nil)
		log.Trace("--> jwToken: ", jwToken.Claims.UserID)
		log.Trace("--> err: ", err)
		if err != nil {
			return err
		} else if jwToken != nil {
			ctx = jwtrule.SetUserIDToCTX(ctx, int(jwToken.Claims.UserID))
			ss = &serverStreamWithContext{ServerStream: ss, ctx: ctx}
		}
		// Call the handler
		err = handler(srv, ss)

		postProcess(ctx, info.FullMethod, log, err)

		// Code to execute after the RPC
		log.Trace("--> finish stream interceptor: ", info.FullMethod)

		return err
	}
}

// serverStreamWithContext wraps the ServerStream to modify the context
type serverStreamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

// Context returns the modified context
func (s *serverStreamWithContext) Context() context.Context {
	return s.ctx
}

func checkAuth(ctx *context.Context, log *logrus.Logger, secretKey string, method string, validateFunc func(tokenString string, key string) (model.Jtoken, error)) (*model.Jtoken, error) {
	if validateFunc == nil {
		validateFunc = jwtrule.Validate
	}
	log.Trace("--> interceptor: ", method)
	// check for method, which doesn't need to be intercepted
	_, ok := SkipCheckMethods[method]
	if ok {
		log.Trace("--> interceptor: pass")
		return nil, nil
	}

	// check part
	token, err := grpc_auth.AuthFromMD(*ctx, "bearer")
	if err != nil {
		log.Trace("AuthFromMD")
		return nil, err
	}

	log.Trace("--> interceptor: check")
	jwToken, err := validateFunc(token, secretKey)
	if err != nil {
		log.Trace("--> interceptor: invalid auth token: ", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}
	log.Trace("--> interceptor UID: ", jwToken.Claims.UserID)

	return &jwToken, nil
}

func preProcess(ctx context.Context, info string, log *logrus.Logger, _ string) {
	log.Trace("--> interceptor: before executing:", info)
	userID := ctx.Value("userID")
	if userID != nil {
		log.Trace("--> pre-processing for user ID: ", userID)
	}
}

func postProcess(ctx context.Context, info string, log *logrus.Logger, err error) {
	log.Trace("--> interceptor: after executing:", info)
	if err != nil {
		log.Trace("--> interceptor: error occurred: ", err)
	}
	userID := ctx.Value("userID")
	if userID != nil {
		log.Trace("--> interceptor: post-processing for user ID: ", userID)
	}

	// check for method, which doesn't need to be intercepted
	_, ok := PostProcessMethods[info]
	if ok {
		log.Trace("--> interceptor: PostProcessMethods: ok")
	}

}
