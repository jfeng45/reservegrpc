package main

import (
	"fmt"
	pb "github.com/jfeng45/reservegrpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"time"
)

type CacheService struct {
	storage map[string][]byte
}

func WithServerinterceptor() grpc.ServerOption{
	return grpc.UnaryInterceptor(serverInterceptor)
}

func serverInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,

) (interface{}, error ) {
	start:=time.Now()
	resp, err:=handler(ctx, req)
	//time.Sleep(time.Second)
	log.Printf("invoke server method=%s duratin=%s error=%v", info.FullMethod,
		time.Since(start), err)
	return resp, err
}

//var cs CacheService
func (s *CacheService) Get(ctx context.Context, req *pb.GetReq) (*pb.GetResp, error) {
	fmt.Println("get called")
	value, ok :=s.storage[req.Key]
	if  !ok {
		fmt.Println("not ok")
		return nil, status.Errorf(codes.NotFound, "key not found=%s", req.Key)
	}
	fmt.Println("ok, return value")
	return &pb.GetResp{Value:value}, nil
}


func (s *CacheService) Store(ctx context.Context, req *pb.StoreReq) (*pb.StoreResp,
	error) {
	fmt.Println("store called")
	//var c1 chan int
	c1:= make(chan int)
	defer close(c1)
	go func() {
		//logrus.Info("start testjob()")
		s.storage[req.Key] = req.Value
		//time.Sleep(10*time.Second)
		c1 <- 5
	}()
	select {
	case  <-c1:
		//s.storage[req.Key] = req.Value
		return &pb.StoreResp{}, nil
	case <-time.After(3*time.Second):
		return nil, status.Error(codes.DeadlineExceeded, "timeout")
	}
}

func (s *CacheService) Dump(req *pb.DumpReq, stream pb.CacheService_DumpServer) error {
	for k, v:=range s.storage {
		stream.Send(&pb.DumpItem{
			Key: k,
			Val:v,
		})
	}
	return nil
}

func runServer() error {
	fmt.Println("start runserver")
	//tlsCreds, err:=credentials.NewClientTLSFromFile("certificate.crt", "privatekey.key")
	//tlsCreds :=credentials.NewTLS(&tls.Config{InsecureSkipVerify:true})
	//if ( err != nil) {
	//	return err
	//}
	//srv:=grpc.NewServer(grpc.Creds(tlsCreds))
	srv:=grpc.NewServer(WithServerinterceptor())

	cs:= &CacheService{storage:make(map[string][]byte)}
	pb.RegisterCacheServiceServer(srv, cs)
	l, err:=net.Listen("tcp", "localhost:5051")

	if err!=nil {
		return err
	} else {
		fmt.Println("server listenig")
	}
	return srv.Serve(l)
}

func main () {
	//cs:= &CacheService{storage:make(map[string][]byte)}
	//cs.storage["java"]=[]byte("fini")
	//fmt.Println(cs)
	serverMain()
}
func serverMain() {
	if err := runServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run cache server: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Println("server started")
	}
}
