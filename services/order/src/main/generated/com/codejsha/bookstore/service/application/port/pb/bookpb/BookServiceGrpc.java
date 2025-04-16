package com.codejsha.bookstore.service.application.port.pb.bookpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.72.0)",
    comments = "Source: modules/book/v1/book.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class BookServiceGrpc {

  private BookServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "book.v1.BookService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>
      getFindAllBooksMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllBooks",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>
      getFindAllBooksMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>
        getFindAllBooksMethod;
    if ((getFindAllBooksMethod = BookServiceGrpc.getFindAllBooksMethod) == null) {
      synchronized (BookServiceGrpc.class) {
        if ((getFindAllBooksMethod = BookServiceGrpc.getFindAllBooksMethod) == null) {
          BookServiceGrpc.getFindAllBooksMethod =
              getFindAllBooksMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.bookpb
                              .BookFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.bookpb
                              .BookFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllBooks"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.bookpb
                                  .BookFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.bookpb
                                  .BookFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new BookServiceMethodDescriptorSupplier("FindAllBooks"))
                      .build();
        }
      }
    }
    return getFindAllBooksMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>
      getFindBookMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindBook",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>
      getFindBookMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq,
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>
        getFindBookMethod;
    if ((getFindBookMethod = BookServiceGrpc.getFindBookMethod) == null) {
      synchronized (BookServiceGrpc.class) {
        if ((getFindBookMethod = BookServiceGrpc.getFindBookMethod) == null) {
          BookServiceGrpc.getFindBookMethod =
              getFindBookMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.bookpb
                              .BookFindProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindBook"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.bookpb
                                  .BookFindProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.bookpb
                                  .BookFindProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new BookServiceMethodDescriptorSupplier("FindBook"))
                      .build();
        }
      }
    }
    return getFindBookMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static BookServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BookServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<BookServiceStub>() {
          @java.lang.Override
          public BookServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new BookServiceStub(channel, callOptions);
          }
        };
    return BookServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static BookServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BookServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<BookServiceBlockingV2Stub>() {
          @java.lang.Override
          public BookServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new BookServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return BookServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BookServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BookServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<BookServiceBlockingStub>() {
          @java.lang.Override
          public BookServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new BookServiceBlockingStub(channel, callOptions);
          }
        };
    return BookServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static BookServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<BookServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<BookServiceFutureStub>() {
          @java.lang.Override
          public BookServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new BookServiceFutureStub(channel, callOptions);
          }
        };
    return BookServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllBooks(
        com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllBooksMethod(), responseObserver);
    }

    /** */
    default void findBook(
        com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getFindBookMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service BookService. */
  public abstract static class BookServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return BookServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service BookService. */
  public static final class BookServiceStub
      extends io.grpc.stub.AbstractAsyncStub<BookServiceStub> {
    private BookServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BookServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BookServiceStub(channel, callOptions);
    }

    /** */
    public void findAllBooks(
        com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllBooksMethod(), getCallOptions()),
          request,
          responseObserver);
    }

    /** */
    public void findBook(
        com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindBookMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service BookService. */
  public static final class BookServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<BookServiceBlockingV2Stub> {
    private BookServiceBlockingV2Stub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BookServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BookServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp
        findAllBooks(
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllBooksMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp findBook(
        com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindBookMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service BookService. */
  public static final class BookServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<BookServiceBlockingStub> {
    private BookServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BookServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BookServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp
        findAllBooks(
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllBooksMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp findBook(
        com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindBookMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service BookService. */
  public static final class BookServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<BookServiceFutureStub> {
    private BookServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BookServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new BookServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>
        findAllBooks(
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllBooksMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>
        findBook(
            com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindBookMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_BOOKS = 0;
  private static final int METHODID_FIND_BOOK = 1;

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
        case METHODID_FIND_ALL_BOOKS:
          serviceImpl.findAllBooks(
              (com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.bookpb
                          .BookFindAllProtoResp>)
                  responseObserver);
          break;
        case METHODID_FIND_BOOK:
          serviceImpl.findBook(
              (com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq) request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>)
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
            getFindAllBooksMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindAllProtoResp>(
                    service, METHODID_FIND_ALL_BOOKS)))
        .addMethod(
            getFindBookMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.bookpb.BookFindProtoResp>(
                    service, METHODID_FIND_BOOK)))
        .build();
  }

  private abstract static class BookServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BookServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.bookpb.BookProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("BookService");
    }
  }

  private static final class BookServiceFileDescriptorSupplier
      extends BookServiceBaseDescriptorSupplier {
    BookServiceFileDescriptorSupplier() {}
  }

  private static final class BookServiceMethodDescriptorSupplier
      extends BookServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    BookServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (BookServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new BookServiceFileDescriptorSupplier())
                      .addMethod(getFindAllBooksMethod())
                      .addMethod(getFindBookMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
