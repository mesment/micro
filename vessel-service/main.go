//提供货船服务

package main

import (
	pb "github.com/mesment/mirco/vessel-service/proto/vessel"
	"errors"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
	"log"
)

type Repository interface {
	FindAvailable(in *pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}

func (v *VesselRepository)FindAvailable(spe *pb.Specification) ( *pb.Vessel,error) {
	//选择一条容量、载重可以容纳的货轮
	for _, v := range v.vessels {
		if v.Capacity >= spe.Capacity && v.MaxWeight >= spe.MaxWeight {
			return v, nil
		}
	}

	return nil, errors.New("No vessel can be use")
}

type service struct {
	vesselRepository VesselRepository
}

func(s *service)FindAvailable(ctx context.Context, req *pb.Specification, resp *pb.Response) error {
	vessel, err := s.vesselRepository.FindAvailable(req)
	if err != nil {
		log.Println("FindAvailable failed,",err)
		return err
	}
	resp.Vessel = vessel
	return nil
}

func main(){
	server := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
		)
	server.Init()
	vessels :=[]*pb.Vessel{{Id: "vessel001", Name: "Boaty McBoatface", MaxWeight: 200000, Capacity: 500},}
	repo :=VesselRepository{vessels}


	pb.RegisterVesselServiceHandler(server.Server(),&service{repo})
	if err := server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
