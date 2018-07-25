package api

import (
	"context"
	"google.golang.org/grpc"
	pb "github.com/linkerd/linkerd2-examples/emojivoto/emojivoto-emoji-svc/gen/proto"
	"github.com/linkerd/linkerd2-examples/emojivoto/emojivoto-emoji-svc/emoji"
)

type EmojiServiceServer struct {
	allEmoji emoji.AllEmoji
}

func (svc *EmojiServiceServer) ListAll(ctx context.Context, req *pb.ListAllEmojiRequest) (*pb.ListAllEmojiResponse, error) {

	emoji := svc.allEmoji.List()

	list := make([]*pb.Emoji, 0)
	for _, e := range emoji {
		pbE := pb.Emoji{
			Unicode:   e.Unicode,
			Shortcode: e.Shortcode,
		}
		list = append(list, &pbE)
	}

	return &pb.ListAllEmojiResponse{List: list}, nil
}

func (svc *EmojiServiceServer)  FindByShortcode(ctx context.Context, req *pb.FindByShortcodeRequest) (*pb.FindByShortcodeResponse, error){
	var pbE *pb.Emoji
	foundEmoji := svc.allEmoji.WithShortcode(req.Shortcode)
	if foundEmoji != nil {
		pbE = &pb.Emoji{
			Unicode:   foundEmoji.Unicode,
			Shortcode: foundEmoji.Shortcode,
		}
	}
	return &pb.FindByShortcodeResponse{
		Emoji: pbE,
	}, nil
}

func NewGrpServer(grpcServer *grpc.Server, allEmoji emoji.AllEmoji) {
	pb.RegisterEmojiServiceServer(grpcServer, &EmojiServiceServer{
		allEmoji,
	})
}
