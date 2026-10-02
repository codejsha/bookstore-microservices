package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType

internal fun buildContext() = ActorContext(actorId = 0L, ActorType.USER)
