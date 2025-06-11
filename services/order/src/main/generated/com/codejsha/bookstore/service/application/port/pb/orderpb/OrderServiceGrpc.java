package com.codejsha.bookstore.service.application.port.pb.orderpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.73.0)",
    comments = "Source: modules/order/v1/order.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class OrderServiceGrpc {

  private OrderServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "order.v1.OrderService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp>
      getFindAllOrdersMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllOrders",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp>
      getFindAllOrdersMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp>
        getFindAllOrdersMethod;
    if ((getFindAllOrdersMethod = OrderServiceGrpc.getFindAllOrdersMethod) == null) {
      synchronized (OrderServiceGrpc.class) {
        if ((getFindAllOrdersMethod = OrderServiceGrpc.getFindAllOrdersMethod) == null) {
          OrderServiceGrpc.getFindAllOrdersMethod =
              getFindAllOrdersMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.orderpb
                              .OrderFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.orderpb
                              .OrderFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllOrders"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.orderpb
                                  .OrderFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.orderpb
                                  .OrderFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new OrderServiceMethodDescriptorSupplier("FindAllOrders"))
                      .build();
        }
      }
    }
    return getFindAllOrdersMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>
      getFindOrderMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindOrder",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>
      getFindOrderMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq,
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>
        getFindOrderMethod;
    if ((getFindOrderMethod = OrderServiceGrpc.getFindOrderMethod) == null) {
      synchronized (OrderServiceGrpc.class) {
        if ((getFindOrderMethod = OrderServiceGrpc.getFindOrderMethod) == null) {
          OrderServiceGrpc.getFindOrderMethod =
              getFindOrderMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.orderpb
                              .OrderFindProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.orderpb
                              .OrderFindProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindOrder"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.orderpb
                                  .OrderFindProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.orderpb
                                  .OrderFindProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new OrderServiceMethodDescriptorSupplier("FindOrder"))
                      .build();
        }
      }
    }
    return getFindOrderMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static OrderServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<OrderServiceStub>() {
          @java.lang.Override
          public OrderServiceStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new OrderServiceStub(channel, callOptions);
          }
        };
    return OrderServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static OrderServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<OrderServiceBlockingV2Stub>() {
          @java.lang.Override
          public OrderServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new OrderServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return OrderServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static OrderServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<OrderServiceBlockingStub>() {
          @java.lang.Override
          public OrderServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new OrderServiceBlockingStub(channel, callOptions);
          }
        };
    return OrderServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static OrderServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<OrderServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<OrderServiceFutureStub>() {
          @java.lang.Override
          public OrderServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new OrderServiceFutureStub(channel, callOptions);
          }
        };
    return OrderServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllOrders(
        com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllOrdersMethod(), responseObserver);
    }

    /** */
    default void findOrder(
        com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getFindOrderMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service OrderService. */
  public abstract static class OrderServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return OrderServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service OrderService. */
  public static final class OrderServiceStub
      extends io.grpc.stub.AbstractAsyncStub<OrderServiceStub> {
    private OrderServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderServiceStub(channel, callOptions);
    }

    /** */
    public void findAllOrders(
        com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllOrdersMethod(), getCallOptions()),
          request,
          responseObserver);
    }

    /** */
    public void findOrder(
        com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindOrderMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service OrderService. */
  public static final class OrderServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<OrderServiceBlockingV2Stub> {
    private OrderServiceBlockingV2Stub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp
        findAllOrders(
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllOrdersMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp findOrder(
        com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindOrderMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service OrderService. */
  public static final class OrderServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<OrderServiceBlockingStub> {
    private OrderServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp
        findAllOrders(
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllOrdersMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp findOrder(
        com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindOrderMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service OrderService. */
  public static final class OrderServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<OrderServiceFutureStub> {
    private OrderServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected OrderServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new OrderServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoResp>
        findAllOrders(
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllOrdersMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>
        findOrder(
            com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindOrderMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_ORDERS = 0;
  private static final int METHODID_FIND_ORDER = 1;

  private static final class MethodHandlers<Req, Resp>
      implements io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
          io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
          io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
          io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_FIND_ALL_ORDERS:
          serviceImpl.findAllOrders(
              (com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.orderpb
                          .OrderFindAllProtoResp>)
                  responseObserver);
          break;
        case METHODID_FIND_ORDER:
          serviceImpl.findOrder(
              (com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.orderpb
                          .OrderFindProtoResp>)
                  responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
            getFindAllOrdersMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.orderpb
                        .OrderFindAllProtoResp>(service, METHODID_FIND_ALL_ORDERS)))
        .addMethod(
            getFindOrderMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.orderpb.OrderFindProtoResp>(
                    service, METHODID_FIND_ORDER)))
        .build();
  }

  private abstract static class OrderServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    OrderServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.orderpb.OrderProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("OrderService");
    }
  }

  private static final class OrderServiceFileDescriptorSupplier
      extends OrderServiceBaseDescriptorSupplier {
    OrderServiceFileDescriptorSupplier() {}
  }

  private static final class OrderServiceMethodDescriptorSupplier
      extends OrderServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    OrderServiceMethodDescriptorSupplier(java.lang.String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (OrderServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new OrderServiceFileDescriptorSupplier())
                      .addMethod(getFindAllOrdersMethod())
                      .addMethod(getFindOrderMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
