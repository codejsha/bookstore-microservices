package com.codejsha.bookstore.service.application.port.pb.authorpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.73.0)",
    comments = "Source: modules/author/v1/author.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class AuthorServiceGrpc {

  private AuthorServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "author.v1.AuthorService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp>
      getFindAllAuthorsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllAuthors",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp>
      getFindAllAuthorsMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp>
        getFindAllAuthorsMethod;
    if ((getFindAllAuthorsMethod = AuthorServiceGrpc.getFindAllAuthorsMethod) == null) {
      synchronized (AuthorServiceGrpc.class) {
        if ((getFindAllAuthorsMethod = AuthorServiceGrpc.getFindAllAuthorsMethod) == null) {
          AuthorServiceGrpc.getFindAllAuthorsMethod =
              getFindAllAuthorsMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.authorpb
                              .AuthorFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.authorpb
                              .AuthorFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllAuthors"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.authorpb
                                  .AuthorFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.authorpb
                                  .AuthorFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new AuthorServiceMethodDescriptorSupplier("FindAllAuthors"))
                      .build();
        }
      }
    }
    return getFindAllAuthorsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp>
      getFindAuthorMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAuthor",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp>
      getFindAuthorMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq,
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp>
        getFindAuthorMethod;
    if ((getFindAuthorMethod = AuthorServiceGrpc.getFindAuthorMethod) == null) {
      synchronized (AuthorServiceGrpc.class) {
        if ((getFindAuthorMethod = AuthorServiceGrpc.getFindAuthorMethod) == null) {
          AuthorServiceGrpc.getFindAuthorMethod =
              getFindAuthorMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.authorpb
                              .AuthorFindProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.authorpb
                              .AuthorFindProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAuthor"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.authorpb
                                  .AuthorFindProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.authorpb
                                  .AuthorFindProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new AuthorServiceMethodDescriptorSupplier("FindAuthor"))
                      .build();
        }
      }
    }
    return getFindAuthorMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static AuthorServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthorServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<AuthorServiceStub>() {
          @java.lang.Override
          public AuthorServiceStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new AuthorServiceStub(channel, callOptions);
          }
        };
    return AuthorServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static AuthorServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthorServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<AuthorServiceBlockingV2Stub>() {
          @java.lang.Override
          public AuthorServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new AuthorServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return AuthorServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AuthorServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthorServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<AuthorServiceBlockingStub>() {
          @java.lang.Override
          public AuthorServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new AuthorServiceBlockingStub(channel, callOptions);
          }
        };
    return AuthorServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static AuthorServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AuthorServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<AuthorServiceFutureStub>() {
          @java.lang.Override
          public AuthorServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new AuthorServiceFutureStub(channel, callOptions);
          }
        };
    return AuthorServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllAuthors(
        com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllAuthorsMethod(), responseObserver);
    }

    /** */
    default void findAuthor(
        com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getFindAuthorMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service AuthorService. */
  public abstract static class AuthorServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return AuthorServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service AuthorService. */
  public static final class AuthorServiceStub
      extends io.grpc.stub.AbstractAsyncStub<AuthorServiceStub> {
    private AuthorServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthorServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthorServiceStub(channel, callOptions);
    }

    /** */
    public void findAllAuthors(
        com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllAuthorsMethod(), getCallOptions()),
          request,
          responseObserver);
    }

    /** */
    public void findAuthor(
        com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAuthorMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service AuthorService. */
  public static final class AuthorServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<AuthorServiceBlockingV2Stub> {
    private AuthorServiceBlockingV2Stub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthorServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthorServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp
        findAllAuthors(
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllAuthorsMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp
        findAuthor(
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAuthorMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service AuthorService. */
  public static final class AuthorServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<AuthorServiceBlockingStub> {
    private AuthorServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthorServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthorServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp
        findAllAuthors(
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllAuthorsMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp
        findAuthor(
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAuthorMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service AuthorService. */
  public static final class AuthorServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<AuthorServiceFutureStub> {
    private AuthorServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AuthorServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AuthorServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoResp>
        findAllAuthors(
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllAuthorsMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoResp>
        findAuthor(
            com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAuthorMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_AUTHORS = 0;
  private static final int METHODID_FIND_AUTHOR = 1;

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
        case METHODID_FIND_ALL_AUTHORS:
          serviceImpl.findAllAuthors(
              (com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.authorpb
                          .AuthorFindAllProtoResp>)
                  responseObserver);
          break;
        case METHODID_FIND_AUTHOR:
          serviceImpl.findAuthor(
              (com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.authorpb
                          .AuthorFindProtoResp>)
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
            getFindAllAuthorsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.authorpb
                        .AuthorFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.authorpb
                        .AuthorFindAllProtoResp>(service, METHODID_FIND_ALL_AUTHORS)))
        .addMethod(
            getFindAuthorMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorFindProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.authorpb
                        .AuthorFindProtoResp>(service, METHODID_FIND_AUTHOR)))
        .build();
  }

  private abstract static class AuthorServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AuthorServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.authorpb.AuthorProto
          .getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("AuthorService");
    }
  }

  private static final class AuthorServiceFileDescriptorSupplier
      extends AuthorServiceBaseDescriptorSupplier {
    AuthorServiceFileDescriptorSupplier() {}
  }

  private static final class AuthorServiceMethodDescriptorSupplier
      extends AuthorServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    AuthorServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (AuthorServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new AuthorServiceFileDescriptorSupplier())
                      .addMethod(getFindAllAuthorsMethod())
                      .addMethod(getFindAuthorMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
