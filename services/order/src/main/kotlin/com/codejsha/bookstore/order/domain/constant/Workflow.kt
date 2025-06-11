package com.codejsha.bookstore.order.domain.constant

enum class Workflow(val wfName: String) {
    BOOKSTORE_ORDER_BOOK_WORKFLOW("bookstore_order_book_wf"),
    BOOKSTORE_ORDER_BOOK_FAILURE_WORKFLOW("bookstore_order_book_failure_wf"),
}
