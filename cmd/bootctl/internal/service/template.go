package service

var protoText = `syntax = "proto3";

option go_package = "./pb";

package pb;

service {{name}}Service {
}`

var mainText = `package main

import (
	"flag"
	"fmt"

	"{{project}}/internal/service/{{name}}/server"
	"github.com/trainking/goboot/pkg/boot"
)

var (
	name       = flag.String("name", "{{Name}}Service", "service name")
	addr       = flag.String("addr", "{{addr}}", "{{name}} service listen address")
	configPath = flag.String("config", "configs/{{name}}.service.yml", "config file path")
	instanceId = flag.Int64("instance", {{id}}, "run instance id")
)

func main() {
	flag.Parse()

	instance := server.New(*name, *configPath, *addr, *instanceId)

	fmt.Println("server start listen: ", *addr)
	if err := boot.BootServe(instance); err != nil {
		fmt.Println("server start failed, Error: ", err)
		return
	}
}`

var serverText = `package server

import (
	"{{project}}/internal/pb"
	"github.com/trainking/goboot/pkg/boot"
	"github.com/trainking/goboot/pkg/log"
	"github.com/trainking/goboot/pkg/service"
	"github.com/trainking/goboot/pkg/utils"
)

type (
	Server struct {
		pb.Unimplemented{{Name}}ServiceServer
		service.BaseService
	}
)

// New 创建Server
func New(name string, configPath string, addr string, instanceId int64) boot.Instance {
	v, err := utils.LoadConfigFileViper(configPath)
	if err != nil {
		panic(err)
	}

	server := new(Server)
	server.Name = name
	server.Config = v
	server.Addr = addr
	server.IntanceID = instanceId
	return server
}

// Init 服务初始化
func (s *Server) Init() error {
	var err error
	if err = s.BaseService.Init(); err != nil {
		log.Errorf("BaseServcie error %v", err)
		return err
	}

	pb.Register{{Name}}ServiceServer(s.GrpcServer, s)

	return nil
}`

var clientText = `package client

import (
	"github.com/spf13/viper"
	"{{project}}/internal/pb"
	"github.com/trainking/goboot/pkg/service"
)

// New{{Name}}Service 创建{{Name}}Servcie客户端
func New{{Name}}Service(serviceName string, config *viper.Viper) (pb.{{Name}}ServiceClient, error) {
	conn, err := service.NewGrpcClientConn(serviceName, config)
	if err != nil {
		return nil, err
	}
	return pb.New{{Name}}ServiceClient(conn), nil
}`
