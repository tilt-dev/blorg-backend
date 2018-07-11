//go:generate protoc -I ./proto --go_out=plugins=grpc:./proto ./proto/backend.proto

package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/windmilleng/blorg-backend/golink"
	pb "github.com/windmilleng/blorg-backend/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

const USER = "blorger"

var dbAddr = flag.String("dbAddr", "localhost:26257", "address of the blorg database")
var db *sql.DB

type server struct {
	gl *golink.Golink
}

func (s *server) Pong(ctx context.Context, in *pb.PongRequest) (*pb.PongResponse, error) {
	return &pb.PongResponse{}, nil
}

func (s *server) GetGolink(ctx context.Context, in *pb.GetGolinkRequest) (*pb.Golink, error) {
	link, err := s.gl.LinkFromName(in.Name)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "error getting link from name: %s", err)
	}

	if link == "" {
		return nil, grpc.Errorf(codes.NotFound, "Link not found for %s", in.Name)
	}

	return &pb.Golink{
		Name:    in.Name,
		Address: link,
	}, nil
}

func (s *server) CreateGolink(ctx context.Context, in *pb.Golink) (*pb.Golink, error) {
	l := &golink.Link{
		Name:    in.Name,
		Address: in.Address,
	}
	err := s.gl.WriteLink(l)
	if err != nil {
		return &pb.Golink{}, grpc.Errorf(codes.Internal, "Error writing link: %s", err)
	}
	return &pb.Golink{}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	setupDatabase()
	gl := golink.NewGolink(db)

	s := grpc.NewServer()
	pb.RegisterBackendServer(s, &server{gl: gl})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupDatabase() {
	cs := fmt.Sprintf("postgresql://%s@%s/golink?sslmode=disable", USER, *dbAddr)
	db2, err := sql.Open("postgres", cs)
	db = db2
	if err != nil {
		log.Fatal("error opening database: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}
}
