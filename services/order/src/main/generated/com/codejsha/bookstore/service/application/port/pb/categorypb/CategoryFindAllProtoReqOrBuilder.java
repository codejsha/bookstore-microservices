// Generated by the protocol buffer compiler.  DO NOT EDIT!
// NO CHECKED-IN PROTOBUF GENCODE
// source: modules/category/v1/category.proto
// Protobuf Java Version: 4.29.3

package com.codejsha.bookstore.service.application.port.pb.categorypb;

public interface CategoryFindAllProtoReqOrBuilder
    extends
    // @@protoc_insertion_point(interface_extends:category.v1.CategoryFindAllProtoReq)
    com.google.protobuf.MessageOrBuilder {

  /**
   * <code>optional .category.v1.QueryFilter filter = 1 [json_name = "filter"];</code>
   *
   * @return Whether the filter field is set.
   */
  boolean hasFilter();

  /**
   * <code>optional .category.v1.QueryFilter filter = 1 [json_name = "filter"];</code>
   *
   * @return The filter.
   */
  com.codejsha.bookstore.service.application.port.pb.categorypb.QueryFilter getFilter();

  /** <code>optional .category.v1.QueryFilter filter = 1 [json_name = "filter"];</code> */
  com.codejsha.bookstore.service.application.port.pb.categorypb.QueryFilterOrBuilder
      getFilterOrBuilder();
}
