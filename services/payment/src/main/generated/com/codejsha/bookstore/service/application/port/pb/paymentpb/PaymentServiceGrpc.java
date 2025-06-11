package com.codejsha.bookstore.service.application.port.pb.paymentpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/** */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.73.0)",
    comments = "Source: modules/payment/v1/payment.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class PaymentServiceGrpc {

  private PaymentServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "payment.v1.PaymentService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp>
      getFindAllPaymentsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindAllPayments",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp
              .class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq,
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp>
      getFindAllPaymentsMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq,
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp>
        getFindAllPaymentsMethod;
    if ((getFindAllPaymentsMethod = PaymentServiceGrpc.getFindAllPaymentsMethod) == null) {
      synchronized (PaymentServiceGrpc.class) {
        if ((getFindAllPaymentsMethod = PaymentServiceGrpc.getFindAllPaymentsMethod) == null) {
          PaymentServiceGrpc.getFindAllPaymentsMethod =
              getFindAllPaymentsMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.paymentpb
                              .PaymentFindAllProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.paymentpb
                              .PaymentFindAllProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindAllPayments"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.paymentpb
                                  .PaymentFindAllProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.paymentpb
                                  .PaymentFindAllProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new PaymentServiceMethodDescriptorSupplier("FindAllPayments"))
                      .build();
        }
      }
    }
    return getFindAllPaymentsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp>
      getFindPaymentMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindPayment",
      requestType =
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq.class,
      responseType =
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq,
          com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp>
      getFindPaymentMethod() {
    io.grpc.MethodDescriptor<
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq,
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp>
        getFindPaymentMethod;
    if ((getFindPaymentMethod = PaymentServiceGrpc.getFindPaymentMethod) == null) {
      synchronized (PaymentServiceGrpc.class) {
        if ((getFindPaymentMethod = PaymentServiceGrpc.getFindPaymentMethod) == null) {
          PaymentServiceGrpc.getFindPaymentMethod =
              getFindPaymentMethod =
                  io.grpc.MethodDescriptor
                      .<com.codejsha.bookstore.service.application.port.pb.paymentpb
                              .PaymentFindProtoReq,
                          com.codejsha.bookstore.service.application.port.pb.paymentpb
                              .PaymentFindProtoResp>
                          newBuilder()
                      .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
                      .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindPayment"))
                      .setSampledToLocalTracing(true)
                      .setRequestMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.paymentpb
                                  .PaymentFindProtoReq.getDefaultInstance()))
                      .setResponseMarshaller(
                          io.grpc.protobuf.ProtoUtils.marshaller(
                              com.codejsha.bookstore.service.application.port.pb.paymentpb
                                  .PaymentFindProtoResp.getDefaultInstance()))
                      .setSchemaDescriptor(
                          new PaymentServiceMethodDescriptorSupplier("FindPayment"))
                      .build();
        }
      }
    }
    return getFindPaymentMethod;
  }

  /** Creates a new async stub that supports all call types for the service */
  public static PaymentServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceStub>() {
          @java.lang.Override
          public PaymentServiceStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PaymentServiceStub(channel, callOptions);
          }
        };
    return PaymentServiceStub.newStub(factory, channel);
  }

  /** Creates a new blocking-style stub that supports all types of calls on the service */
  public static PaymentServiceBlockingV2Stub newBlockingV2Stub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingV2Stub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingV2Stub>() {
          @java.lang.Override
          public PaymentServiceBlockingV2Stub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PaymentServiceBlockingV2Stub(channel, callOptions);
          }
        };
    return PaymentServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PaymentServiceBlockingStub newBlockingStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingStub>() {
          @java.lang.Override
          public PaymentServiceBlockingStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PaymentServiceBlockingStub(channel, callOptions);
          }
        };
    return PaymentServiceBlockingStub.newStub(factory, channel);
  }

  /** Creates a new ListenableFuture-style stub that supports unary calls on the service */
  public static PaymentServiceFutureStub newFutureStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceFutureStub> factory =
        new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceFutureStub>() {
          @java.lang.Override
          public PaymentServiceFutureStub newStub(
              io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
            return new PaymentServiceFutureStub(channel, callOptions);
          }
        };
    return PaymentServiceFutureStub.newStub(factory, channel);
  }

  /** */
  public interface AsyncService {

    /** */
    default void findAllPayments(
        com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.paymentpb
                    .PaymentFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindAllPaymentsMethod(), responseObserver);
    }

    /** */
    default void findPayment(
        com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp>
            responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(
          getFindPaymentMethod(), responseObserver);
    }
  }

  /** Base class for the server implementation of the service PaymentService. */
  public abstract static class PaymentServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override
    public final io.grpc.ServerServiceDefinition bindService() {
      return PaymentServiceGrpc.bindService(this);
    }
  }

  /** A stub to allow clients to do asynchronous rpc calls to service PaymentService. */
  public static final class PaymentServiceStub
      extends io.grpc.stub.AbstractAsyncStub<PaymentServiceStub> {
    private PaymentServiceStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceStub build(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceStub(channel, callOptions);
    }

    /** */
    public void findAllPayments(
        com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.paymentpb
                    .PaymentFindAllProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindAllPaymentsMethod(), getCallOptions()),
          request,
          responseObserver);
    }

    /** */
    public void findPayment(
        com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq request,
        io.grpc.stub.StreamObserver<
                com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp>
            responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindPaymentMethod(), getCallOptions()),
          request,
          responseObserver);
    }
  }

  /** A stub to allow clients to do synchronous rpc calls to service PaymentService. */
  public static final class PaymentServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<PaymentServiceBlockingV2Stub> {
    private PaymentServiceBlockingV2Stub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceBlockingV2Stub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp
        findAllPayments(
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllPaymentsMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp
        findPayment(
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindPaymentMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do limited synchronous rpc calls to service PaymentService. */
  public static final class PaymentServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<PaymentServiceBlockingStub> {
    private PaymentServiceBlockingStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceBlockingStub(channel, callOptions);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp
        findAllPayments(
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindAllPaymentsMethod(), getCallOptions(), request);
    }

    /** */
    public com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp
        findPayment(
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindPaymentMethod(), getCallOptions(), request);
    }
  }

  /** A stub to allow clients to do ListenableFuture-style rpc calls to service PaymentService. */
  public static final class PaymentServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<PaymentServiceFutureStub> {
    private PaymentServiceFutureStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceFutureStub(channel, callOptions);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoResp>
        findAllPayments(
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindAllPaymentsMethod(), getCallOptions()), request);
    }

    /** */
    public com.google.common.util.concurrent.ListenableFuture<
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoResp>
        findPayment(
            com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq
                request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindPaymentMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_FIND_ALL_PAYMENTS = 0;
  private static final int METHODID_FIND_PAYMENT = 1;

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
        case METHODID_FIND_ALL_PAYMENTS:
          serviceImpl.findAllPayments(
              (com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindAllProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.paymentpb
                          .PaymentFindAllProtoResp>)
                  responseObserver);
          break;
        case METHODID_FIND_PAYMENT:
          serviceImpl.findPayment(
              (com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentFindProtoReq)
                  request,
              (io.grpc.stub.StreamObserver<
                      com.codejsha.bookstore.service.application.port.pb.paymentpb
                          .PaymentFindProtoResp>)
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
            getFindAllPaymentsMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.paymentpb
                        .PaymentFindAllProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.paymentpb
                        .PaymentFindAllProtoResp>(service, METHODID_FIND_ALL_PAYMENTS)))
        .addMethod(
            getFindPaymentMethod(),
            io.grpc.stub.ServerCalls.asyncUnaryCall(
                new MethodHandlers<
                    com.codejsha.bookstore.service.application.port.pb.paymentpb
                        .PaymentFindProtoReq,
                    com.codejsha.bookstore.service.application.port.pb.paymentpb
                        .PaymentFindProtoResp>(service, METHODID_FIND_PAYMENT)))
        .build();
  }

  private abstract static class PaymentServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier,
          io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PaymentServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.service.application.port.pb.paymentpb.PaymentProto
          .getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("PaymentService");
    }
  }

  private static final class PaymentServiceFileDescriptorSupplier
      extends PaymentServiceBaseDescriptorSupplier {
    PaymentServiceFileDescriptorSupplier() {}
  }

  private static final class PaymentServiceMethodDescriptorSupplier
      extends PaymentServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    PaymentServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (PaymentServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor =
              result =
                  io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
                      .setSchemaDescriptor(new PaymentServiceFileDescriptorSupplier())
                      .addMethod(getFindAllPaymentsMethod())
                      .addMethod(getFindPaymentMethod())
                      .build();
        }
      }
    }
    return result;
  }
}
