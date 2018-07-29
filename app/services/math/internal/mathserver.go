package internal

import (
	"context"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sknv/microproto/app/services/math/rpc"
)

type MathServer struct{}

func (s *MathServer) Rect(_ context.Context, args *rpc.RectArgs) (*rpc.RectReply, error) {
	if args.Width <= 0 {
		return nil, status.Error(codes.InvalidArgument, "width must be a positive number")
	}
	if args.Height <= 0 {
		return nil, status.Error(codes.InvalidArgument, "height must be a positive number")
	}

	return &rpc.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}

func (s *MathServer) Circle(_ context.Context, args *rpc.CircleArgs) (*rpc.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, status.Error(codes.InvalidArgument, "radius must be a positive number")
	}

	return &rpc.CircleReply{
		Length: 2 * math.Pi * args.Radius,
		Square: math.Pi * args.Radius * args.Radius,
	}, nil
}
