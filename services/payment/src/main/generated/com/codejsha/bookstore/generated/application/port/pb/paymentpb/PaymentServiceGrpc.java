package com.codejsha.bookstore.generated.application.port.pb.paymentpb;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@io.grpc.stub.annotations.GrpcGenerated
public final class PaymentServiceGrpc {

  private PaymentServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "payment.v1.PaymentService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse> getListPaymentsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListPayments",
      requestType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest.class,
      responseType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse> getListPaymentsMethod() {
    io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse> getListPaymentsMethod;
    if ((getListPaymentsMethod = PaymentServiceGrpc.getListPaymentsMethod) == null) {
      synchronized (PaymentServiceGrpc.class) {
        if ((getListPaymentsMethod = PaymentServiceGrpc.getListPaymentsMethod) == null) {
          PaymentServiceGrpc.getListPaymentsMethod = getListPaymentsMethod =
              io.grpc.MethodDescriptor.<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListPayments"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PaymentServiceMethodDescriptorSupplier("ListPayments"))
              .build();
        }
      }
    }
    return getListPaymentsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse> getFindPaymentMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindPayment",
      requestType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest.class,
      responseType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse> getFindPaymentMethod() {
    io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse> getFindPaymentMethod;
    if ((getFindPaymentMethod = PaymentServiceGrpc.getFindPaymentMethod) == null) {
      synchronized (PaymentServiceGrpc.class) {
        if ((getFindPaymentMethod = PaymentServiceGrpc.getFindPaymentMethod) == null) {
          PaymentServiceGrpc.getFindPaymentMethod = getFindPaymentMethod =
              io.grpc.MethodDescriptor.<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindPayment"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PaymentServiceMethodDescriptorSupplier("FindPayment"))
              .build();
        }
      }
    }
    return getFindPaymentMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse> getListRefundsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ListRefunds",
      requestType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest.class,
      responseType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse> getListRefundsMethod() {
    io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse> getListRefundsMethod;
    if ((getListRefundsMethod = PaymentServiceGrpc.getListRefundsMethod) == null) {
      synchronized (PaymentServiceGrpc.class) {
        if ((getListRefundsMethod = PaymentServiceGrpc.getListRefundsMethod) == null) {
          PaymentServiceGrpc.getListRefundsMethod = getListRefundsMethod =
              io.grpc.MethodDescriptor.<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "ListRefunds"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PaymentServiceMethodDescriptorSupplier("ListRefunds"))
              .build();
        }
      }
    }
    return getListRefundsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse> getFindRefundMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "FindRefund",
      requestType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest.class,
      responseType = com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest,
      com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse> getFindRefundMethod() {
    io.grpc.MethodDescriptor<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse> getFindRefundMethod;
    if ((getFindRefundMethod = PaymentServiceGrpc.getFindRefundMethod) == null) {
      synchronized (PaymentServiceGrpc.class) {
        if ((getFindRefundMethod = PaymentServiceGrpc.getFindRefundMethod) == null) {
          PaymentServiceGrpc.getFindRefundMethod = getFindRefundMethod =
              io.grpc.MethodDescriptor.<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest, com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "FindRefund"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse.getDefaultInstance()))
              .setSchemaDescriptor(new PaymentServiceMethodDescriptorSupplier("FindRefund"))
              .build();
        }
      }
    }
    return getFindRefundMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static PaymentServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceStub>() {
        @java.lang.Override
        public PaymentServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PaymentServiceStub(channel, callOptions);
        }
      };
    return PaymentServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static PaymentServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingV2Stub>() {
        @java.lang.Override
        public PaymentServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PaymentServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return PaymentServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static PaymentServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceBlockingStub>() {
        @java.lang.Override
        public PaymentServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PaymentServiceBlockingStub(channel, callOptions);
        }
      };
    return PaymentServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static PaymentServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<PaymentServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<PaymentServiceFutureStub>() {
        @java.lang.Override
        public PaymentServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new PaymentServiceFutureStub(channel, callOptions);
        }
      };
    return PaymentServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     */
    default void listPayments(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListPaymentsMethod(), responseObserver);
    }

    /**
     */
    default void findPayment(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getFindPaymentMethod(), responseObserver);
    }

    /**
     */
    default void listRefunds(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getListRefundsMethod(), responseObserver);
    }

    /**
     */
    default void findRefund(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getFindRefundMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service PaymentService.
   */
  public static abstract class PaymentServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return PaymentServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service PaymentService.
   */
  public static final class PaymentServiceStub
      extends io.grpc.stub.AbstractAsyncStub<PaymentServiceStub> {
    private PaymentServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceStub(channel, callOptions);
    }

    /**
     */
    public void listPayments(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListPaymentsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void findPayment(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindPaymentMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void listRefunds(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getListRefundsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void findRefund(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest request,
        io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getFindRefundMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service PaymentService.
   */
  public static final class PaymentServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<PaymentServiceBlockingV2Stub> {
    private PaymentServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse listPayments(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getListPaymentsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse findPayment(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getFindPaymentMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse listRefunds(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getListRefundsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse findRefund(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest request) throws io.grpc.StatusException {
      return io.grpc.stub.ClientCalls.blockingV2UnaryCall(
          getChannel(), getFindRefundMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service PaymentService.
   */
  public static final class PaymentServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<PaymentServiceBlockingStub> {
    private PaymentServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse listPayments(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListPaymentsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse findPayment(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindPaymentMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse listRefunds(com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getListRefundsMethod(), getCallOptions(), request);
    }

    /**
     */
    public com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse findRefund(com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getFindRefundMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service PaymentService.
   */
  public static final class PaymentServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<PaymentServiceFutureStub> {
    private PaymentServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected PaymentServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new PaymentServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse> listPayments(
        com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListPaymentsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse> findPayment(
        com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindPaymentMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse> listRefunds(
        com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getListRefundsMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse> findRefund(
        com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getFindRefundMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LIST_PAYMENTS = 0;
  private static final int METHODID_FIND_PAYMENT = 1;
  private static final int METHODID_LIST_REFUNDS = 2;
  private static final int METHODID_FIND_REFUND = 3;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
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
        case METHODID_LIST_PAYMENTS:
          serviceImpl.listPayments((com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest) request,
              (io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse>) responseObserver);
          break;
        case METHODID_FIND_PAYMENT:
          serviceImpl.findPayment((com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest) request,
              (io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse>) responseObserver);
          break;
        case METHODID_LIST_REFUNDS:
          serviceImpl.listRefunds((com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest) request,
              (io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse>) responseObserver);
          break;
        case METHODID_FIND_REFUND:
          serviceImpl.findRefund((com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest) request,
              (io.grpc.stub.StreamObserver<com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse>) responseObserver);
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
          getListPaymentsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsRequest,
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListPaymentsResponse>(
                service, METHODID_LIST_PAYMENTS)))
        .addMethod(
          getFindPaymentMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentRequest,
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindPaymentResponse>(
                service, METHODID_FIND_PAYMENT)))
        .addMethod(
          getListRefundsMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsRequest,
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.ListRefundsResponse>(
                service, METHODID_LIST_REFUNDS)))
        .addMethod(
          getFindRefundMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundRequest,
              com.codejsha.bookstore.generated.application.port.pb.paymentpb.FindRefundResponse>(
                service, METHODID_FIND_REFUND)))
        .build();
  }

  private static abstract class PaymentServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    PaymentServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return com.codejsha.bookstore.generated.application.port.pb.paymentpb.V1Proto.getDescriptor();
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
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new PaymentServiceFileDescriptorSupplier())
              .addMethod(getListPaymentsMethod())
              .addMethod(getFindPaymentMethod())
              .addMethod(getListRefundsMethod())
              .addMethod(getFindRefundMethod())
              .build();
        }
      }
    }
    return result;
  }
}
