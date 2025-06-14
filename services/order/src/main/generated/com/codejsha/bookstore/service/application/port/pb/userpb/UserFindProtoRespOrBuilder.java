// Generated by the protocol buffer compiler.  DO NOT EDIT!
// NO CHECKED-IN PROTOBUF GENCODE
// source: modules/user/v1/user.proto
// Protobuf Java Version: 4.29.3

package com.codejsha.bookstore.service.application.port.pb.userpb;

public interface UserFindProtoRespOrBuilder
    extends
    // @@protoc_insertion_point(interface_extends:user.v1.UserFindProtoResp)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>string id = 1 [json_name = "id"];</code>
   *
   * @return The id.
   */
  java.lang.String getId();

  /**
   * <code>string id = 1 [json_name = "id"];</code>
   *
   * @return The bytes for id.
   */
  com.google.protobuf.ByteString getIdBytes();

  /**
   * <code>string email = 2 [json_name = "email"];</code>
   *
   * @return The email.
   */
  java.lang.String getEmail();

  /**
   * <code>string email = 2 [json_name = "email"];</code>
   *
   * @return The bytes for email.
   */
  com.google.protobuf.ByteString getEmailBytes();

  /**
   * <code>string first_name = 3 [json_name = "firstName"];</code>
   *
   * @return The firstName.
   */
  java.lang.String getFirstName();

  /**
   * <code>string first_name = 3 [json_name = "firstName"];</code>
   *
   * @return The bytes for firstName.
   */
  com.google.protobuf.ByteString getFirstNameBytes();

  /**
   * <code>string last_name = 4 [json_name = "lastName"];</code>
   *
   * @return The lastName.
   */
  java.lang.String getLastName();

  /**
   * <code>string last_name = 4 [json_name = "lastName"];</code>
   *
   * @return The bytes for lastName.
   */
  com.google.protobuf.ByteString getLastNameBytes();

  /**
   * <code>optional string phone = 5 [json_name = "phone"];</code>
   *
   * @return Whether the phone field is set.
   */
  boolean hasPhone();

  /**
   * <code>optional string phone = 5 [json_name = "phone"];</code>
   *
   * @return The phone.
   */
  java.lang.String getPhone();

  /**
   * <code>optional string phone = 5 [json_name = "phone"];</code>
   *
   * @return The bytes for phone.
   */
  com.google.protobuf.ByteString getPhoneBytes();

  /**
   * <code>repeated .user.v1.AuthRole roles = 6 [json_name = "roles"];</code>
   *
   * @return A list containing the roles.
   */
  java.util.List<com.codejsha.bookstore.service.application.port.pb.userpb.AuthRole> getRolesList();

  /**
   * <code>repeated .user.v1.AuthRole roles = 6 [json_name = "roles"];</code>
   *
   * @return The count of roles.
   */
  int getRolesCount();

  /**
   * <code>repeated .user.v1.AuthRole roles = 6 [json_name = "roles"];</code>
   *
   * @param index The index of the element to return.
   * @return The roles at the given index.
   */
  com.codejsha.bookstore.service.application.port.pb.userpb.AuthRole getRoles(int index);

  /**
   * <code>repeated .user.v1.AuthRole roles = 6 [json_name = "roles"];</code>
   *
   * @return A list containing the enum numeric values on the wire for roles.
   */
  java.util.List<java.lang.Integer> getRolesValueList();

  /**
   * <code>repeated .user.v1.AuthRole roles = 6 [json_name = "roles"];</code>
   *
   * @param index The index of the value to return.
   * @return The enum numeric value on the wire of roles at the given index.
   */
  int getRolesValue(int index);
}
