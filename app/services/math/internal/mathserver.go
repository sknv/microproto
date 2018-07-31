package internal

import (
	"context"
	"math"

	"github.com/twitchtv/twirp"

	"github.com/sknv/microproto/app/services/math/rpc"
)

type MathServer struct{}

func (*MathServer) Rect(_ context.Context, args *rpc.RectArgs) (*rpc.RectReply, error) {
	if args.Width <= 0 || args.Height <= 0 {
		return nil, twirp.InvalidArgumentError("width and height", "must be positive numbers")
	}

	return &rpc.RectReply{
		Perimeter: 2*args.Width + 2*args.Height,
		Square:    args.Width * args.Height,
	}, nil
}

func (*MathServer) Circle(_ context.Context, args *rpc.CircleArgs) (*rpc.CircleReply, error) {
	if args.Radius <= 0 {
		return nil, twirp.InvalidArgumentError("radius", "must be a positive number")
	}

	return &rpc.CircleReply{
		Length: 2 * math.Pi * args.Radius,
		Square: math.Pi * args.Radius * args.Radius,
	}, nil
}
