package sd

import (
	"context"
	"fmt"
	"time"

	"github.com/PolarPanda611/trinitygo/util"
	"github.com/coreos/etcd/clientv3"
	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"google.golang.org/grpc"
	"google.golang.org/grpc/naming"
)

// ServiceMesh interface
type ServiceMesh interface {
	GetClient() interface{}
	RegService(projectName string, projectVersion string, serviceIP string, servicePort int, Tags []string, timeout int) error
	DeRegService(projectName string, projectVersion string, serviceIP string, servicePort int, timeout int) error
}

// ServiceMeshEtcdImpl consul register
type ServiceMeshEtcdImpl struct {
	// config
	Address string // consul address
	Port    int

	// runtime
	client *clientv3.Client
}

// NewEtcdRegister New consul register
func NewEtcdRegister(address string, port int) (ServiceMesh, error) {
	s := &ServiceMeshEtcdImpl{
		Address: address,
		Port:    port,
	}

	cli, err := clientv3.NewFromURL(fmt.Sprintf("http://%v:%v", s.Address, s.Port))

	if err != nil {
		return nil, err
	}
	s.client = cli
	return s, nil
}

// GetClient get etcd client
func (s *ServiceMeshEtcdImpl) GetClient() interface{} {
	return s.client
}

// RegService register etcd service
func (s *ServiceMeshEtcdImpl) RegService(projectName string, projectVersion string, serviceIP string, servicePort int, Tags []string, timeout int) error {
	r := &etcdnaming.GRPCResolver{Client: s.client}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	err := r.Update(ctx, util.GetServiceName(projectName), naming.Update{Op: naming.Add, Addr: fmt.Sprintf("%v:%v", serviceIP, servicePort), Metadata: fmt.Sprintf("%v", Tags)})
	if err != nil {
		return err
	}
	return nil
}

// DeRegService deregister service
func (s *ServiceMeshEtcdImpl) DeRegService(projectName string, projectVersion string, serviceIP string, servicePort int, timeout int) error {
	r := &etcdnaming.GRPCResolver{Client: s.client}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()
	err := r.Update(ctx, util.GetServiceName(projectName), naming.Update{Op: naming.Delete, Addr: fmt.Sprintf("%v:%v", serviceIP, servicePort)})
	if err != nil {
		return err
	}
	return nil
}

// NewEtcdClientConn new etcd client connection
func NewEtcdClientConn(address string, port int, serviceName string, timeout int) (*grpc.ClientConn, error) {
	cli, err := clientv3.NewFromURL(fmt.Sprintf("http://%v:%v", address, port))
	if err != nil {
		return nil, fmt.Errorf("failed to conn etcd client , %v", err)
	}
	r := &etcdnaming.GRPCResolver{Client: cli}
	b := grpc.RoundRobin(r)

	ctx1, cel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cel()
	conn, err := grpc.DialContext(ctx1, serviceName, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithBalancer(b))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
