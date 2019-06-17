package main

import (
	"fmt"
	pb "github.com/jfeng45/reservegrpc"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"time"
)

func WithClientInterceptor() grpc.DialOption {
	return grpc.WithUnaryInterceptor(clientInterceptor)
}

//func clientInterceptor (
//	ctx context.Context,
//		method string,
//			req interface{},
//			reply interface{},
//			cc *grpc.ClientConn,
//			invoker grpc.UnaryInvoker,
//			opts ...grpc.CallOption,
//) error {
//	start:=time.Now()
//	err:=invoker(ctx, method, req, reply,cc, opts...)
//	log.Printf("invoke remote metho=%s durtion=%s error=%v", method, time.Since(start), err)
//	return err
//}

func clientInterceptor (
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	var err error
	retry :=RetryCount{3,0}
	c:=make ( chan int)
	go retry.exeRetry(c)
	defer close(c)
	for {
		select {
			case <-ctx.Done():
				log.Println("clientInterceptor():ctx time out")
				err = status.Errorf(codes.DeadlineExceeded, "cliemt timeout")
			case <-c:
				log.Println("clientInterceptor(): exe retyr")
				startAttempt:=time.Now()
				err = invoker(ctx, method, req, reply, cc, opts...)
				if ( err!=nil ) {
					//log.Println(" %v time retey", retry.ExeCount)
					log.Printf("clientInterceptor():invoke remote method %v=%s durtion=%s error=%v", retry.ExeCount,method, time.Since(startAttempt), err)
					go retry.exeRetry(c)

					continue
				} else {
					log.Printf("clientInterceptor():invoke remote method %v=%s durtion=%s error=%v", retry.ExeCount,method, time.Since(startAttempt), err)
					break

			}
		}
		break
	}
	log.Println( "clientInterceptor():rety finished in interceptor")
	//log.Printf("invoke remote metho=%s durtion=%s error=%v", method, time.Since(start), err)
	return err
	//start:=time.Now()
	//err:=invoker(ctx, method, req, reply,cc, opts...)
	//log.Printf("invoke remote metho=%s durtion=%s error=%v", method, time.Since(start), err)
	//return err
}

type RetryCount struct {
	RetryCount int
	ExeCount int
}
type  ContRetryer interface {
	//contRetry()  (bool, error)
	exeRetry()
}

func (r *RetryCount) exeRetry( c chan int) (bool, error) {
	log.Println("exeRetry(): r.RetryCount=%v", r.RetryCount)
	if r.RetryCount >r.ExeCount {
		r.ExeCount++
		log.Println("exeRetry(): r.ExeCount=%v", r.ExeCount)
		c<-1
		return true, nil
	} else {
		log.Println("exeRetry(): no more retrys")
		return false, errors.New("no more retrys")
	}
}

//func (r RetryCont) contRetry() (bool, error) {
//	if r.RetryCount >r.ExeCount {
//		return true, nil
//	} else {
//		return false, nil
//	}
//}

func callStore (cache pb.CacheServiceClient) {
	ctx := context.Background()
	key:= "go"

	//value:="abc"
	//_, err= cache.Store(ctx, &pb.StoreReq{key,[]byte("con"), a, []byte("abc"),1})
	timeoutCtx, _ := context.WithTimeout(ctx, 10*time.Second)
	_, err:= cache.Store(timeoutCtx, &pb.StoreReq{Key:key,Value:[]byte("con")})

	if err != nil {
		//fmt.Errorf("failed to store: %v", err)
		fmt.Println(err)
	} else {
		fmt.Println("store called")
	}
}
func callGet(cache pb.CacheServiceClient) {
	//_, err:=cache.Get(context.Background(), &pb.GetReq{Key:"go", })
	resp, err:=cache.Get(context.Background(), &pb.GetReq{Key:"go"})

	if err != nil {
		//fmt.Errorf("failed to get: %v", err)
		fmt.Println(err)
		//fmt.Errorf()
	} else {
		fmt.Printf("Got cached value %s\n", string(resp.Value))
	}
}

func callDump(cache pb.CacheServiceClient) {
	//_, err:=cache.Get(context.Background(), &pb.GetReq{Key:"go", })
	stream, err:=cache.Dump(context.Background(),&pb.DumpReq{})
	if err != nil {
		fmt.Errorf("failed to dump: %v", err)
	}
	for {
		item, err:=stream.Recv()
		if err == io.EOF {
			break
		}
		if err!= nil {
			fmt.Errorf("failed to stream item: %v", err)
		}
		fmt.Fprintf(os.Stdout,"item=%v", item)
	}

	//if err != nil {
	//	//fmt.Errorf("failed to get: %v", err)
	//	fmt.Println(err)
	//	//fmt.Errorf()
	//} else {
	//	for _, item:=range resp.Items {
	//		fmt.Println("key=%s,value=%", item.Key,string(item.Val))
	//	}
	//	//fmt.Printf("Got cached value %s\n", string(resp.Items))
	//}
}
func main() {
	//retryBackoff
	//grpc.BackoffConfig{}
	//tlsCreds, err:=credentials.NewClientTLSFromFile("certificate.crt", "privatekey.key")
	//tlsCreds :=credentials.NewTLS(&tls.Config{InsecureSkipVerify:true})
	//if ( err != nil) {
	//	fmt.Println(err)
	//}
	//conn, err:=grpc.Dial("localhost:5051", grpc.WithTransportCredentials(tlsCreds))
	//conn, err:=grpc.Dial("localhost:5051", grpc.WithInsecure())
	conn, err:=grpc.Dial("localhost:5051", grpc.WithInsecure(), WithClientInterceptor())
	if err != nil {
		fmt.Errorf("failed to dial server: %v", err)
	}
	cache :=pb.NewCacheServiceClient(conn)
	fmt.Println("client strated")

	callStore(cache)
	callGet(cache)
	callDump(cache)


	//fmt.Printf("abcwids \n")

}
