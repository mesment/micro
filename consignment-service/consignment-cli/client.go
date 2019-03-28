package main

import (
	"context"
	"encoding/json"
	pb "github.com/mesment/mirco/consignment-service/proto/consignment"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"io/ioutil"
	"log"
	"time"
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

	cmd.Init()
	c := pb.NewShippingService("go.micro.srv.consignment",client.DefaultClient)


	consignment,err :=parseJSONFile(jsonfile)
	if err != nil {
		log.Fatal("parse file failed,",err)
	}
	resp,err := c.CreateConsignment(context.Background(),consignment)
	if err != nil {
		log.Fatalf("call CreateConsignment failed,",err)
		return
	}
	log.Printf("created:%t\n",resp.Created)


	time.Sleep(10*time.Second)
	resp, err = c.GetConsignments(context.Background(),&pb.GetRequest{})
	if err != nil {
		log.Fatalf("failed to list consignment %v",err)
	}

	for _, c := range resp.Consignments {
		log.Printf("%+v",c)
	}


}
