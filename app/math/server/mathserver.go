package server

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	math "github.com/sknv/microproto/app/math/rpc"
)

type MathServer struct{}

func (*MathServer) Circle(_ context.Context, args *math.CircleArgs) (*math.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, status.Error(codes.InvalidArgument, "radius must be a positive number")
	}

	pi := 3.1416 // there is math.Pi constant in the standard lib btw
	return &math.CircleReply{
		Length: 2 * pi * args.Radius,
		Square: pi * args.Radius * args.Radius,
	}, nil
}

func (*MathServer) Rect(_ context.Context, args *math.RectArgs) (*math.RectReply, error) {
	if args.Width <= 0 || args.Height <= 0 {
		return nil, status.Error(codes.InvalidArgument, "width and height must be positive numbers")
	}

	return &math.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}
