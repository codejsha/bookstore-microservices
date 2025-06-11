package com.codejsha.bookstore.service.application.port.pb.categorypb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.73.0)",
    comments = "Source: modules/category/v1/category.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class CategoryServiceGrpc {

  private CategoryServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "category.v1.CategoryService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp>
      getFindAllCategoriesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllCategories",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq
              .class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp
              .class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp>
      getFindAllCategoriesMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp>
        getFindAllCategoriesMethod;
    if ((getFindAllCategoriesMethod = CategoryServiceGrpc.getFindAllCategoriesMethod) == null) {
      synchronized (CategoryServiceGrpc.class) {
        if ((getFindAllCategoriesMethod = CategoryServiceGrpc.getFindAllCategoriesMethod) == null) {
          CategoryServiceGrpc.getFindAllCategoriesMethod =
              getFindAllCategoriesMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.categorypb
                              .CategoryFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.categorypb
                              .CategoryFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllCategories"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.categorypb
                                  .CategoryFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.categorypb
                                  .CategoryFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new CategoryServiceMethodDescriptorSupplier("FindAllCategories"))
                      .build();
        }
      }
    }
    return getFindAllCategoriesMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static CategoryServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CategoryServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<CategoryServiceStub>() {
          @java.lang.Override
          public CategoryServiceStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new CategoryServiceStub(channel, callOptions);
          }
        };
    return CategoryServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static CategoryServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CategoryServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<CategoryServiceBlockingV2Stub>() {
          @java.lang.Override
          public CategoryServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new CategoryServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return CategoryServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static CategoryServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CategoryServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<CategoryServiceBlockingStub>() {
          @java.lang.Override
          public CategoryServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new CategoryServiceBlockingStub(channel, callOptions);
          }
        };
    return CategoryServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static CategoryServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<CategoryServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<CategoryServiceFutureStub>() {
          @java.lang.Override
          public CategoryServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new CategoryServiceFutureStub(channel, callOptions);
          }
        };
    return CategoryServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllCategories(
        com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq
            request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.categorypb
                    .CategoryFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllCategoriesMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service CategoryService. */
  public abstract static class CategoryServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return CategoryServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service CategoryService. */
  public static final class CategoryServiceStub
      extends io.grpc.stub.AbstractAsyncStub<CategoryServiceStub> {
    private CategoryServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CategoryServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CategoryServiceStub(channel, callOptions);
    }

    /** */
    public void findAllCategories(
        com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq
            request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.categorypb
                    .CategoryFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllCategoriesMethod(), getCallOptions()),
          request,
          responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service CategoryService. */
  public static final class CategoryServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<CategoryServiceBlockingV2Stub> {
    private CategoryServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CategoryServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CategoryServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp
        findAllCategories(
            com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllCategoriesMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service CategoryService. */
  public static final class CategoryServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<CategoryServiceBlockingStub> {
    private CategoryServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CategoryServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CategoryServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp
        findAllCategories(
            com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllCategoriesMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service CategoryService. */
  public static final class CategoryServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<CategoryServiceFutureStub> {
    private CategoryServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected CategoryServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new CategoryServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoResp>
        findAllCategories(
            com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllCategoriesMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_CATEGORIES = 0;

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
        case METHODID_FIND_ALL_CATEGORIES:
          serviceImpl.findAllCategories(
              (com.codejsha.bookstore.service.application.port.pb.categorypb
                      .CategoryFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.categorypb
                          .CategoryFindAllProtoResp>)
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
            getFindAllCategoriesMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.categorypb
                        .CategoryFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.categorypb
                        .CategoryFindAllProtoResp>(service, METHODID_FIND_ALL_CATEGORIES)))
        .build();
  }

  private abstract static class CategoryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    CategoryServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.categorypb.CategoryProto
          .getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("CategoryService");
    }
  }

  private static final class CategoryServiceFileDescriptorSupplier
      extends CategoryServiceBaseDescriptorSupplier {
    CategoryServiceFileDescriptorSupplier() {}
  }

  private static final class CategoryServiceMethodDescriptorSupplier
      extends CategoryServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    CategoryServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (CategoryServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new CategoryServiceFileDescriptorSupplier())
                      .addMethod(getFindAllCategoriesMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
