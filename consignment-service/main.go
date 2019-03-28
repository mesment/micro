package main

import (
	pb "github.com/mesment/mirco/consignment-service/proto/consignment"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	port = ":50051"
)

//仓库接口
type IRepository interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

//实际的仓库
type Repository struct {
	consignments []*pb.Consignment
}

func (rep *Repository)Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	rep.consignments = append(rep.consignments, consignment)
	return consignment,nil
}

func (rep *Repository)GetAll() ([]*pb.Consignment) {
	return rep.consignments
}

//实现consignment.pb.go中的CreateConsignment接口
type service struct {
	repo Repository
}


func (s *service)CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {
	consignment,err := s.repo.Create(req)
	if err != nil {
		log.Printf("CreateConsignment failed ", err)
		return nil, err
	}
	resp := &pb.Response{
		Created:true,
		Consignment:consignment,
	}
	return resp,nil
}


func(s *service)GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	//目前所有托运的货物
	allConsignments := s.repo.GetAll()
	resp := &pb.Response{Consignments:allConsignments}
	log.Printf("%v",resp.Consignments)
	return resp, nil
}



func main(){

	listener,err := net.Listen("tcp",port)
	if err != nil {
		log.Fatalf("failed to listen:%v", err)
	}
	log.Printf("listen on :%s\n",port)

	//新建server
	server := grpc.NewServer()
	repo := Repository{}

	// 向 rRPC 服务器注册微服务
	// 把实现的微服务 service 与协议中的 ShippingServiceServer 绑定
	pb.RegisterShippingServiceServer(server,&service{repo})

	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to serve %v",err)
	}

}