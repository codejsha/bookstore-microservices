package com.codejsha.bookstore.service.application.port.pb.publisherpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.73.0)",
    comments = "Source: modules/publisher/v1/publisher.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class PublisherServiceGrpc {

  private PublisherServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "publisher.v1.PublisherService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoResp>
      getFindAllPublishersMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllPublishers",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq
              .class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoResp
              .class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoResp>
      getFindAllPublishersMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.publisherpb
                .PublisherFindAllProtoResp>
        getFindAllPublishersMethod;
    if ((getFindAllPublishersMethod = PublisherServiceGrpc.getFindAllPublishersMethod) == null) {
      synchronized (PublisherServiceGrpc.class) {
        if ((getFindAllPublishersMethod = PublisherServiceGrpc.getFindAllPublishersMethod)
            == null) {
          PublisherServiceGrpc.getFindAllPublishersMethod =
              getFindAllPublishersMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.publisherpb
                              .PublisherFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.publisherpb
                              .PublisherFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllPublishers"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.publisherpb
                                  .PublisherFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.publisherpb
                                  .PublisherFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new PublisherServiceMethodDescriptorSupplier("FindAllPublishers"))
                      .build();
        }
      }
    }
    return getFindAllPublishersMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp>
      getFindPublisherMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindPublisher",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq
              .class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp
              .class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp>
      getFindPublisherMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq,
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp>
        getFindPublisherMethod;
    if ((getFindPublisherMethod = PublisherServiceGrpc.getFindPublisherMethod) == null) {
      synchronized (PublisherServiceGrpc.class) {
        if ((getFindPublisherMethod = PublisherServiceGrpc.getFindPublisherMethod) == null) {
          PublisherServiceGrpc.getFindPublisherMethod =
              getFindPublisherMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.publisherpb
                              .PublisherFindProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.publisherpb
                              .PublisherFindProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindPublisher"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.publisherpb
                                  .PublisherFindProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.publisherpb
                                  .PublisherFindProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new PublisherServiceMethodDescriptorSupplier("FindPublisher"))
                      .build();
        }
      }
    }
    return getFindPublisherMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static PublisherServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublisherServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PublisherServiceStub>() {
          @java.lang.Override
          public PublisherServiceStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PublisherServiceStub(channel, callOptions);
          }
        };
    return PublisherServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static PublisherServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublisherServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PublisherServiceBlockingV2Stub>() {
          @java.lang.Override
          public PublisherServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PublisherServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return PublisherServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PublisherServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublisherServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PublisherServiceBlockingStub>() {
          @java.lang.Override
          public PublisherServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PublisherServiceBlockingStub(channel, callOptions);
          }
        };
    return PublisherServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static PublisherServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PublisherServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PublisherServiceFutureStub>() {
          @java.lang.Override
          public PublisherServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PublisherServiceFutureStub(channel, callOptions);
          }
        };
    return PublisherServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllPublishers(
        com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq
            request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.publisherpb
                    .PublisherFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllPublishersMethod(), responseObserver);
    }

    /** */
    default void findPublisher(
        com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq
            request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.publisherpb
                    .PublisherFindProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindPublisherMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service PublisherService. */
  public abstract static class PublisherServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return PublisherServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service PublisherService. */
  public static final class PublisherServiceStub
      extends io.grpc.stub.AbstractAsyncStub<PublisherServiceStub> {
    private PublisherServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublisherServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublisherServiceStub(channel, callOptions);
    }

    /** */
    public void findAllPublishers(
        com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq
            request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.publisherpb
                    .PublisherFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllPublishersMethod(), getCallOptions()),
          request,
          responseObserver);
    }

    /** */
    public void findPublisher(
        com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq
            request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.publisherpb
                    .PublisherFindProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindPublisherMethod(), getCallOptions()),
          request,
          responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service PublisherService. */
  public static final class PublisherServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<PublisherServiceBlockingV2Stub> {
    private PublisherServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublisherServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublisherServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoResp
        findAllPublishers(
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllPublishersMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp
        findPublisher(
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindPublisherMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service PublisherService. */
  public static final class PublisherServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<PublisherServiceBlockingStub> {
    private PublisherServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublisherServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublisherServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoResp
        findAllPublishers(
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllPublishersMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp
        findPublisher(
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindPublisherMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service PublisherService. */
  public static final class PublisherServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<PublisherServiceFutureStub> {
    private PublisherServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PublisherServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PublisherServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.publisherpb
                .PublisherFindAllProtoResp>
        findAllPublishers(
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllPublishersMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoResp>
        findPublisher(
            com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindPublisherMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_PUBLISHERS = 0;
  private static final int METHODID_FIND_PUBLISHER = 1;

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
        case METHODID_FIND_ALL_PUBLISHERS:
          serviceImpl.findAllPublishers(
              (com.codejsha.bookstore.service.application.port.pb.publisherpb
                      .PublisherFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.publisherpb
                          .PublisherFindAllProtoResp>)
                  responseObserver);
          break;
        case METHODID_FIND_PUBLISHER:
          serviceImpl.findPublisher(
              (com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherFindProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.publisherpb
                          .PublisherFindProtoResp>)
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
            getFindAllPublishersMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.publisherpb
                        .PublisherFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.publisherpb
                        .PublisherFindAllProtoResp>(service, METHODID_FIND_ALL_PUBLISHERS)))
        .addMethod(
            getFindPublisherMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.publisherpb
                        .PublisherFindProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.publisherpb
                        .PublisherFindProtoResp>(service, METHODID_FIND_PUBLISHER)))
        .build();
  }

  private abstract static class PublisherServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PublisherServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.publisherpb.PublisherProto
          .getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("PublisherService");
    }
  }

  private static final class PublisherServiceFileDescriptorSupplier
      extends PublisherServiceBaseDescriptorSupplier {
    PublisherServiceFileDescriptorSupplier() {}
  }

  private static final class PublisherServiceMethodDescriptorSupplier
      extends PublisherServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    PublisherServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (PublisherServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new PublisherServiceFileDescriptorSupplier())
                      .addMethod(getFindAllPublishersMethod())
                      .addMethod(getFindPublisherMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
