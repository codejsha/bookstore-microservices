// Generated by the protocol buffer compiler.  DO NOT EDIT!
// NO CHECKED-IN PROTOBUF GENCODE
// source: modules/payment/v1/payment.proto
// Protobuf Java Version: 4.29.3

package com.codejsha.bookstore.service.application.port.pb.paymentpb;

public interface PaymentFindProtoRespOrBuilder
    extends
    // @@protoc_insertion_point(interface_extends:payment.v1.PaymentFindProtoResp)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>int64 id = 1 [json_name = "id"];</code>
   *
   * @return The id.
   */
  long getId();

  /**
   * <code>int64 order_id = 2 [json_name = "orderId"];</code>
   *
   * @return The orderId.
   */
  long getOrderId();

  /**
   * <code>string user_id = 3 [json_name = "userId"];</code>
   *
   * @return The userId.
   */
  java.lang.String getUserId();

  /**
   * <code>string user_id = 3 [json_name = "userId"];</code>
   *
   * @return The bytes for userId.
   */
  com.google.protobuf.ByteString getUserIdBytes();

  /**
   * <code>string payment_type = 4 [json_name = "paymentType"];</code>
   *
   * @return The paymentType.
   */
  java.lang.String getPaymentType();

  /**
   * <code>string payment_type = 4 [json_name = "paymentType"];</code>
   *
   * @return The bytes for paymentType.
   */
  com.google.protobuf.ByteString getPaymentTypeBytes();

  /**
   * <code>string card_number = 5 [json_name = "cardNumber"];</code>
   *
   * @return The cardNumber.
   */
  java.lang.String getCardNumber();

  /**
   * <code>string card_number = 5 [json_name = "cardNumber"];</code>
   *
   * @return The bytes for cardNumber.
   */
  com.google.protobuf.ByteString getCardNumberBytes();

  /**
   * <code>string amount = 6 [json_name = "amount"];</code>
   *
   * @return The amount.
   */
  java.lang.String getAmount();

  /**
   * <code>string amount = 6 [json_name = "amount"];</code>
   *
   * @return The bytes for amount.
   */
  com.google.protobuf.ByteString getAmountBytes();

  /**
   * <code>.google.protobuf.Timestamp payment_date = 7 [json_name = "paymentDate"];</code>
   *
   * @return Whether the paymentDate field is set.
   */
  boolean hasPaymentDate();

  /**
   * <code>.google.protobuf.Timestamp payment_date = 7 [json_name = "paymentDate"];</code>
   *
   * @return The paymentDate.
   */
  com.google.protobuf.Timestamp getPaymentDate();

  /** <code>.google.protobuf.Timestamp payment_date = 7 [json_name = "paymentDate"];</code> */
  com.google.protobuf.TimestampOrBuilder getPaymentDateOrBuilder();
}
