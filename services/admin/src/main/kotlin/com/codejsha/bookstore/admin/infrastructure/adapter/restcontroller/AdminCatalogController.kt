package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.CatalogUseCase
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.admin.infrastructure.support.auth.Principal
import com.codejsha.bookstore.admin.infrastructure.support.auth.assertAdmin
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminCatalogApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminAuthorFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSubjectFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkUpdateRequest
import com.codejsha.platform.shared.data.ActorContext
import com.codejsha.platform.shared.data.ActorType
import kotlinx.coroutines.runBlocking
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
class AdminCatalogController(
    private val catalogUseCase: CatalogUseCase,
    private val principalResolver: HttpPrincipalResolver,
) : AdminCatalogApi {

    override fun adminCatalogListWorks(
        title: String?,
        authorUid: String?,
        subjectUid: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminWorkFindAllResponse> = runBlocking {
        val principal = requireAdmin()
        val option = WorkQueryOption(title = title, authorUid = authorUid, subjectUid = subjectUid)
        val context = buildContext(principal)
        val result = catalogUseCase.findAllWorks(option, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminWorkFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminWorkResponse(it) },
            )
        )
    }

    override fun adminCatalogReadWork(uid: String): ResponseEntity<AdminWorkResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        ResponseEntity.ok(toAdminWorkResponse(catalogUseCase.findWork(uid, context)))
    }

    override fun adminCatalogCreateWork(requestBody: AdminWorkCreateRequest): ResponseEntity<Unit> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        val command = WorkCreateCommand(
            title = requestBody.title,
            description = requestBody.description,
            firstPublishDate = requestBody.firstPublishDate,
            olKey = requestBody.olKey,
            authorUids = requestBody.authorUids,
            subjectNames = requestBody.subjectNames,
        )
        catalogUseCase.createWork(command, context)
        ResponseEntity.status(HttpStatus.CREATED).build()
    }

    override fun adminCatalogUpdateWork(
        uid: String,
        requestBody: AdminWorkUpdateRequest,
    ): ResponseEntity<AdminWorkResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        val command = WorkUpdateCommand(
            title = requestBody.title,
            description = requestBody.description,
            firstPublishDate = requestBody.firstPublishDate,
            olKey = requestBody.olKey,
            authorUids = requestBody.authorUids,
            subjectNames = requestBody.subjectNames,
        )
        ResponseEntity.ok(toAdminWorkResponse(catalogUseCase.updateWork(uid, command, context)))
    }

    override fun adminCatalogListAuthors(
        name: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminAuthorFindAllResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        val result = catalogUseCase.findAllAuthors(name, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminAuthorFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminAuthorItem(it) },
            )
        )
    }

    override fun adminCatalogListSubjects(
        name: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminSubjectFindAllResponse> = runBlocking {
        val principal = requireAdmin()
        val context = buildContext(principal)
        val result = catalogUseCase.findAllSubjects(name, pageable ?: Pageable.unpaged(), context)
        ResponseEntity.ok(
            AdminSubjectFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminSubjectItem(it) },
            )
        )
    }

    private fun requireAdmin(): Principal = principalResolver.require().also { it.assertAdmin() }

    private fun buildContext(principal: Principal) =
        ActorContext(actorId = 0L, ActorType.USER)
}
