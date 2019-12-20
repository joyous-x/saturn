package wgrpc

import (
	"context"
	"github.com/joyous-x/saturn/common/xlog"
	"github.com/joyous-x/saturn/govern/wconsul"
	"github.com/joyous-x/saturn/wgrpc/protoc"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	grefl "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

// GrpcMethodCaller with conn, make real call
type GrpcMethodCaller func(ctx context.Context, conn *grpc.ClientConn) (proto.Message, error)

func grpcDialOpts(svcName string) []grpc.DialOption {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(UnaryClientInterceptorJaeger()),
	}
	return opts
}

// CallRaw generate a *grpc.ClientConn and pass to caller which make real call
func CallRaw(ctx context.Context, svcName, tag string, caller GrpcMethodCaller) (proto.Message, error) {
	svcData, err := wconsul.OneHealthSvcRandom(svcName, tag)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", svcData.Address, svcData.Port), grpcDialOpts(svcName)...)
	if err != nil {
		xlog.Error("CallRaw grpc.Dial: %v:%v err %v", svcData.Address, svcData.Port, err)
		return nil, err
	}
	defer conn.Close()

	return caller(ctx, conn)
}

// CallDispatch call a grpc service: lb_policy = random
func CallDispatch(ctx context.Context, svcName, tag string, in *protoc.DispatchReq) (*protoc.DispatchResp, error) {
	svcData, err := wconsul.OneHealthSvcRandom(svcName, tag)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", svcData.Address, svcData.Port), grpcDialOpts(svcName)...)
	if err != nil {
		xlog.Error("CallDispatch grpc.Dial: %v:%v err %v", svcData.Address, svcData.Port, err)
		return nil, err
	}
	defer conn.Close()
	cc := protoc.NewDispatchSvcClient(conn)

	resp, err := cc.Dispatch(ctx, in)
	if err != nil {
		return nil, err
	}

	xlog.Debug("DoDispatch succ host=%v:%v cmd=%v req_id=%v", svcData.Address, svcData.Port, in.Header.Cmd, in.Header.ReqId)
	return resp, nil
}

// CallRawRe call a grpc service: lb_policy = random
func CallRawRe(ctx context.Context, svcName, tag, svcGrpcName, method string, req proto.Message, resp ...proto.Message) error {
	svcData, err := wconsul.OneHealthSvcRandom(svcName, tag)
	if err != nil {
		return err
	}
	host := fmt.Sprintf("%v:%v", svcData.Address, svcData.Port)
	conn, err := grpc.Dial(host, grpcDialOpts(svcName)...)
	if err != nil {
		return err
	}
	defer conn.Close()

	reflectClent := grpcreflect.NewClient(ctx, grefl.NewServerReflectionClient(conn))
	srvDescriptor, err := reflectClent.ResolveService(svcGrpcName)
	if err != nil {
		return err
	}
	methodDesc := srvDescriptor.FindMethodByName(method)
	if methodDesc == nil {
		return fmt.Errorf("method:%v not found", method)
	}
	toOutgoingCtx := func(ctx context.Context) context.Context {
		md := make(metadata.MD)
		md.Append("CTX_TEST_KEY_A", "CTX_TEST_VALUE_A")
		return metadata.NewOutgoingContext(ctx, md)
	}
	retriveFromMD := func(ctx context.Context, md metadata.MD) context.Context {
		getMetaStr := func(ctx context.Context, meta metadata.MD, key string, quiet ...bool) string {
			value := meta.Get(key)
			if len(value) > 0 {
				return value[0]
			}
			return ""
		}
		ctx = context.WithValue(ctx, "CTX_TEST_KEY_A", getMetaStr(ctx, md, "CTX_TEST_VALUE_A"))
		return ctx
	}
	var header, trailer metadata.MD
	r, err := grpcdynamic.NewStub(conn).InvokeRpc(toOutgoingCtx(ctx), methodDesc, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		return err
	}
	ctx = retriveFromMD(ctx, header)

	if len(resp) > 0 {
		response := resp[0]
		dm := r.(*dynamic.Message)
		err = dm.ConvertTo(response)
		if err != nil {
			xlog.Error("Convert response failed. addr: %v, svc: %v, method: %v, resp: %v", host, srvDescriptor.GetFullyQualifiedName(), method, response)
		}
	}
	return err
}
