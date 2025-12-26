package v1

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/usememos/memos/internal/profile"
	"github.com/usememos/memos/internal/util"
	v1pb "github.com/usememos/memos/proto/gen/api/v1"
	"github.com/usememos/memos/store"
)

type APIV1Service struct {
	grpc_health_v1.UnimplementedHealthServer

	v1pb.UnimplementedWorkspaceServiceServer
	v1pb.UnimplementedWorkspaceSettingServiceServer
	v1pb.UnimplementedAuthServiceServer
	v1pb.UnimplementedUserServiceServer
	v1pb.UnimplementedMemoServiceServer
	v1pb.UnimplementedResourceServiceServer
	v1pb.UnimplementedShortcutServiceServer
	v1pb.UnimplementedInboxServiceServer
	v1pb.UnimplementedActivityServiceServer
	v1pb.UnimplementedWebhookServiceServer
	v1pb.UnimplementedMarkdownServiceServer
	v1pb.UnimplementedIdentityProviderServiceServer

	Secret  string
	Profile *profile.Profile
	Store   *store.Store

	grpcServer *grpc.Server
}

func NewAPIV1Service(secret string, profile *profile.Profile, store *store.Store, grpcServer *grpc.Server) *APIV1Service {
	grpc.EnableTracing = true
	apiv1Service := &APIV1Service{
		Secret:     secret,
		Profile:    profile,
		Store:      store,
		grpcServer: grpcServer,
	}
	grpc_health_v1.RegisterHealthServer(grpcServer, apiv1Service)
	v1pb.RegisterWorkspaceServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterWorkspaceSettingServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterAuthServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterUserServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterMemoServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterResourceServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterShortcutServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterInboxServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterActivityServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterWebhookServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterMarkdownServiceServer(grpcServer, apiv1Service)
	v1pb.RegisterIdentityProviderServiceServer(grpcServer, apiv1Service)
	reflection.Register(grpcServer)
	return apiv1Service
}

// RegisterGateway registers the gRPC-Gateway with the given Echo instance.
func (s *APIV1Service) RegisterGateway(ctx context.Context, echoServer *echo.Echo) error {
	var target string
	if len(s.Profile.UNIXSock) == 0 {
		target = fmt.Sprintf("%s:%d", s.Profile.Addr, s.Profile.Port)
	} else {
		target = fmt.Sprintf("unix:%s", s.Profile.UNIXSock)
	}
	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
	)
	if err != nil {
		return err
	}

	gwMux := runtime.NewServeMux()
	if err := v1pb.RegisterWorkspaceServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterWorkspaceSettingServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterAuthServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterUserServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterMemoServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterResourceServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterShortcutServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterInboxServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterActivityServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterWebhookServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterMarkdownServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	if err := v1pb.RegisterIdentityProviderServiceHandler(ctx, gwMux, conn); err != nil {
		return err
	}
	gwGroup := echoServer.Group("")
	gwGroup.Use(middleware.CORS())

	// Register ticket routes directly to Echo group with Auth middleware
	// Register these BEFORE the gRPC-gateway Any wildcard to ensure they take precedence
	ticketGroup := echoServer.Group("/api/v1")
	ticketGroup.Use(s.AuthMiddleware)
	s.RegisterTicketRoutes(ticketGroup)
	s.RegisterNotificationRoutes(ticketGroup)

	handler := echo.WrapHandler(gwMux)
	gwGroup.Any("/api/v1/*", handler)
	gwGroup.Any("/file/*", handler)

	// GRPC web proxy.
	options := []grpcweb.Option{
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(_ string) bool {
			return true
		}),
	}
	wrappedGrpc := grpcweb.WrapServer(s.grpcServer, options...)
	echoServer.Any("/memos.api.v1.*", echo.WrapHandler(wrappedGrpc))

	// Register SSE notification stream endpoint
	// This uses raw http.Handler to support SSE streaming properly
	echoServer.GET("/api/v1/notifications/stream", func(c echo.Context) error {
		userID, ok := c.Get(getUserIDContextKey()).(int32)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing user ID")
		}
		s.NotificationStreamHandler(c.Response().Writer, c.Request(), userID)
		return nil
	}, s.AuthMiddleware)

	return nil
}

func (s *APIV1Service) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		accessToken := ""

		// Check header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				accessToken = parts[1]
			}
		} else {
			// Check cookie
			cookie, err := c.Cookie(AccessTokenCookieName)
			if err == nil {
				accessToken = cookie.Value
			}
		}

		if accessToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing access token")
		}

		// Validate token
		claims := &ClaimsMessage{}
		_, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Name {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			if kid, ok := t.Header["kid"].(string); ok {
				if kid == KeyID {
					return []byte(s.Secret), nil
				}
			}
			return nil, fmt.Errorf("unexpected kid: %v", t.Header["kid"])
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
		}

		userID, err := util.ConvertStringToInt32(claims.Subject)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token subject")
		}

		// Get user to ensure exists and active
		user, err := s.Store.GetUser(ctx, &store.FindUser{ID: &userID})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user").SetInternal(err)
		}
		if user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "User not found")
		}
		if user.RowStatus == store.Archived {
			return echo.NewHTTPError(http.StatusUnauthorized, "User is archived")
		}

		// Validate token against DB tokens
		accessTokens, err := s.Store.GetUserAccessTokens(ctx, user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user access tokens").SetInternal(err)
		}
		isValid := false
		for _, t := range accessTokens {
			if t.AccessToken == accessToken {
				isValid = true
				break
			}
		}
		if !isValid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Token revoked or invalid")
		}

		c.Set(getUserIDContextKey(), userID)
		return next(c)
	}
}
