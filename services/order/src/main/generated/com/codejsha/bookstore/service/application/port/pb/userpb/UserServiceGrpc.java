package com.codejsha.bookstore.service.application.port.pb.userpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.72.0)",
    comments = "Source: modules/user/v1/user.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class UserServiceGrpc {

  private UserServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "user.v1.UserService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>
      getFindAllUsersMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllUsers",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>
      getFindAllUsersMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>
        getFindAllUsersMethod;
    if ((getFindAllUsersMethod = UserServiceGrpc.getFindAllUsersMethod) == null) {
      synchronized (UserServiceGrpc.class) {
        if ((getFindAllUsersMethod = UserServiceGrpc.getFindAllUsersMethod) == null) {
          UserServiceGrpc.getFindAllUsersMethod =
              getFindAllUsersMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllUsers"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new UserServiceMethodDescriptorSupplier("FindAllUsers"))
                      .build();
        }
      }
    }
    return getFindAllUsersMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>
      getFindUserMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindUser",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>
      getFindUserMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq,
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>
        getFindUserMethod;
    if ((getFindUserMethod = UserServiceGrpc.getFindUserMethod) == null) {
      synchronized (UserServiceGrpc.class) {
        if ((getFindUserMethod = UserServiceGrpc.getFindUserMethod) == null) {
          UserServiceGrpc.getFindUserMethod =
              getFindUserMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserFindProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindUser"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserFindProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserFindProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new UserServiceMethodDescriptorSupplier("FindUser"))
                      .build();
        }
      }
    }
    return getFindUserMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>
      getUpdateUserMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UpdateUser",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>
      getUpdateUserMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq,
            com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>
        getUpdateUserMethod;
    if ((getUpdateUserMethod = UserServiceGrpc.getUpdateUserMethod) == null) {
      synchronized (UserServiceGrpc.class) {
        if ((getUpdateUserMethod = UserServiceGrpc.getUpdateUserMethod) == null) {
          UserServiceGrpc.getUpdateUserMethod =
              getUpdateUserMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserUpdateProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserUpdateProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "UpdateUser"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserUpdateProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserUpdateProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new UserServiceMethodDescriptorSupplier("UpdateUser"))
                      .build();
        }
      }
    }
    return getUpdateUserMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>
      getDeleteUserMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DeleteUser",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq,
          com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>
      getDeleteUserMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq,
            com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>
        getDeleteUserMethod;
    if ((getDeleteUserMethod = UserServiceGrpc.getDeleteUserMethod) == null) {
      synchronized (UserServiceGrpc.class) {
        if ((getDeleteUserMethod = UserServiceGrpc.getDeleteUserMethod) == null) {
          UserServiceGrpc.getDeleteUserMethod =
              getDeleteUserMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserDeleteProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.userpb
                              .UserDeleteProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "DeleteUser"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserDeleteProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.userpb
                                  .UserDeleteProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(new UserServiceMethodDescriptorSupplier("DeleteUser"))
                      .build();
        }
      }
    }
    return getDeleteUserMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static UserServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UserServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<UserServiceStub>() {
          @java.lang.Override
          public UserServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new UserServiceStub(channel, callOptions);
          }
        };
    return UserServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static UserServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UserServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<UserServiceBlockingV2Stub>() {
          @java.lang.Override
          public UserServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new UserServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return UserServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static UserServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UserServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<UserServiceBlockingStub>() {
          @java.lang.Override
          public UserServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new UserServiceBlockingStub(channel, callOptions);
          }
        };
    return UserServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static UserServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<UserServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<UserServiceFutureStub>() {
          @java.lang.Override
          public UserServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new UserServiceFutureStub(channel, callOptions);
          }
        };
    return UserServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllUsers(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllUsersMethod(), responseObserver);
    }

    /** */
    default void findUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getFindUserMethod(), responseObserver);
    }

    /** */
    default void updateUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getUpdateUserMethod(), responseObserver);
    }

    /** */
    default void deleteUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getDeleteUserMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service UserService. */
  public abstract static class UserServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return UserServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service UserService. */
  public static final class UserServiceStub
      extends io.grpc.stub.AbstractAsyncStub<UserServiceStub> {
    private UserServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UserServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UserServiceStub(channel, callOptions);
    }

    /** */
    public void findAllUsers(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllUsersMethod(), getCallOptions()),
          request,
          responseObserver);
    }

    /** */
    public void findUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindUserMethod(), getCallOptions()), request, responseObserver);
    }

    /** */
    public void updateUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getUpdateUserMethod(), getCallOptions()), request, responseObserver);
    }

    /** */
    public void deleteUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getDeleteUserMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service UserService. */
  public static final class UserServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<UserServiceBlockingV2Stub> {
    private UserServiceBlockingV2Stub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UserServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UserServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp
        findAllUsers(
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllUsersMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp findUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindUserMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp updateUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUpdateUserMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp deleteUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDeleteUserMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service UserService. */
  public static final class UserServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<UserServiceBlockingStub> {
    private UserServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UserServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UserServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp
        findAllUsers(
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllUsersMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp findUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindUserMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp updateUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getUpdateUserMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp deleteUser(
        com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getDeleteUserMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service UserService. */
  public static final class UserServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<UserServiceFutureStub> {
    private UserServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected UserServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new UserServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>
        findAllUsers(
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllUsersMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>
        findUser(
            com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindUserMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>
        updateUser(
            com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getUpdateUserMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>
        deleteUser(
            com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getDeleteUserMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_USERS = 0;
  private static final int METHODID_FIND_USER = 1;
  private static final int METHODID_UPDATE_USER = 2;
  private static final int METHODID_DELETE_USER = 3;

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
        case METHODID_FIND_ALL_USERS:
          serviceImpl.findAllUsers(
              (com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.userpb
                          .UserFindAllProtoResp>)
                  responseObserver);
          break;
        case METHODID_FIND_USER:
          serviceImpl.findUser(
              (com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq) request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>)
                  responseObserver);
          break;
        case METHODID_UPDATE_USER:
          serviceImpl.updateUser(
              (com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.userpb
                          .UserUpdateProtoResp>)
                  responseObserver);
          break;
        case METHODID_DELETE_USER:
          serviceImpl.deleteUser(
              (com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.userpb
                          .UserDeleteProtoResp>)
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
            getFindAllUsersMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserFindAllProtoResp>(
                    service, METHODID_FIND_ALL_USERS)))
        .addMethod(
            getFindUserMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserFindProtoResp>(
                    service, METHODID_FIND_USER)))
        .addMethod(
            getUpdateUserMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserUpdateProtoResp>(
                    service, METHODID_UPDATE_USER)))
        .addMethod(
            getDeleteUserMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.userpb.UserDeleteProtoResp>(
                    service, METHODID_DELETE_USER)))
        .build();
  }

  private abstract static class UserServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    UserServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.userpb.UserProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("UserService");
    }
  }

  private static final class UserServiceFileDescriptorSupplier
      extends UserServiceBaseDescriptorSupplier {
    UserServiceFileDescriptorSupplier() {}
  }

  private static final class UserServiceMethodDescriptorSupplier
      extends UserServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    UserServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (UserServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new UserServiceFileDescriptorSupplier())
                      .addMethod(getFindAllUsersMethod())
                      .addMethod(getFindUserMethod())
                      .addMethod(getUpdateUserMethod())
                      .addMethod(getDeleteUserMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
