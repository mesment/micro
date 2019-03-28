//提供货运服务


package main

import (
	"fmt"
	pb "github.com/mesment/mirco/consignment-service/proto/consignment"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
	"log"
	veslpb "github.com/mesment/mirco/vessel-service/proto/vessel"
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

	//作为VesselService的客户端调用货船的服务
	vesselClient  veslpb.VesselService
}




func (s *service)CreateConsignment(ctx context.Context, req *pb.Consignment,resp *pb.Response) error {
	//检查是否有足够的货轮可以运输
	vReq :=  veslpb.Specification{
			Capacity:int32(len(req.Containers) ),
			MaxWeight:req.Weight,
	}
	veslresp, err := s.vesselClient.FindAvailable(context.Background(),&vReq)
	if err != nil {
		return err
	}
	// 货物有足够货轮可以运输
	log.Printf("found vessel:%s\n",veslresp.Vessel.Name)
	//设置货物的运输船ID
	req.VesselId = veslresp.Vessel.Id
	
	consignment,err := s.repo.Create(req)
	if err != nil {
		log.Printf("CreateConsignment failed ", err)
		return  err
	}
	resp.Created = true
	resp.Consignment = consignment
	log.Println("CreateConsignment success\n", consignment)
	return nil
}


func(s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error{
	//目前所有托运的货物
	allConsignments := s.repo.GetAll()
	resp.Consignments = allConsignments
	log.Println("reponse:")
	log.Printf("%v",resp.Consignments)
	return nil
}


func main(){

	// Create a new service. Optionally include some options here.
	serv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		)
	//解析命令行参数Init will parse the command line flags.
	serv.Init()
	repo := Repository{}

	//作为 vessel-service 的客户端
	vessClient := veslpb.NewVesselService("go.micro.srv.vessel",client.DefaultClient)
	fmt.Println("vessel client %v",vessClient)

	// Register handler
	pb.RegisterShippingServiceHandler(serv.Server(),&service{repo,vessClient})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve:%v",err)
	}


}