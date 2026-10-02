package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.domain.model.SelfManagementException

internal fun assertNotSelf(uid: String, actorUid: String, message: String) {
    if (uid.equals(actorUid, ignoreCase = true)) throw SelfManagementException(message)
}
