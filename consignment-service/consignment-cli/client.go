package main 

import (
	"context"
	"encoding/json"
	"log"
	pb "github.com/mesment/mirco/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"io/ioutil"
)

const (
	address = "localhost:50051"
	jsonfile = "consignment.json"
)

func parseJSONFile(filename string) (*pb.Consignment, error){
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var consignment pb.Consignment
	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, err
	}
	return &consignment, nil
}

func main() {

	//连接gRPC服务器
	conn, err := grpc.Dial(address,grpc.WithInsecure())
	if err != nil {
		log.Fatal("connect to server failed,",err)
	}

	c := pb.NewShippingServiceClient(conn)
	consignment,err :=parseJSONFile(jsonfile)
	if err != nil {
		log.Fatal("parse file failed,",err)
	}
	resp,err := c.CreateConsignment(context.Background(),consignment)
	if err != nil {
		log.Printf("call CreateConsignment failed")
	}
	log.Printf("created:%v,",resp.Created)


	resp, err = c.GetConsignments(context.Background(),&pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignment %v",err)
	}

	for _, c := range resp.Consignments {
		log.Printf("%+v",c)
	}


}
