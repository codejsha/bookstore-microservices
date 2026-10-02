package com.codejsha.bookstore.admin.domain.service

import com.codejsha.bookstore.admin.application.port.restclient.IdentityClient
import com.codejsha.bookstore.admin.application.usecase.UserUseCase
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import com.codejsha.platform.shared.data.ActorContext
import org.springframework.data.domain.Pageable
import org.springframework.stereotype.Service

@Service
class UserService(
    private val identityClient: IdentityClient,
) : UserUseCase {

    override fun findAllUsers(option: UserQueryOption, pageable: Pageable, context: ActorContext) =
        identityClient.findAllUsers(option, pageable)

    override fun findUser(uid: String, context: ActorContext) = identityClient.findUser(uid)

    override fun updateRoles(
        uid: String,
        roles: List<String>,
        actorUid: String,
        context: ActorContext,
    ): User {
        assertNotSelf(uid, actorUid, "administrators cannot change their own roles")
        return identityClient.updateRoles(uid, roles)
    }

    override fun suspendUser(uid: String, actorUid: String, context: ActorContext): User {
        assertNotSelf(uid, actorUid, "administrators cannot suspend their own account")
        return identityClient.suspendUser(uid)
    }

    override fun reactivateUser(uid: String, context: ActorContext) = identityClient.reactivateUser(uid)

    override fun deactivateUser(uid: String, actorUid: String, context: ActorContext): User {
        assertNotSelf(uid, actorUid, "administrators cannot deactivate their own account")
        return identityClient.deactivateUser(uid)
    }
}
