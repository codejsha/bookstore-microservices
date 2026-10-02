package com.codejsha.bookstore.admin.infrastructure.adapter.restcontroller

import com.codejsha.bookstore.admin.application.usecase.CatalogUseCase
import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import com.codejsha.bookstore.admin.infrastructure.support.auth.HttpPrincipalResolver
import com.codejsha.bookstore.generated.application.port.openapi.api.AdminCatalogApi
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminAuthorFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminSubjectFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkCreateRequest
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkFindAllResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkResponse
import com.codejsha.bookstore.generated.application.port.openapi.model.AdminWorkUpdateRequest
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
    ): ResponseEntity<AdminWorkFindAllResponse> {
        principalResolver.requireStaff()
        val option = WorkQueryOption(title = title, authorUid = authorUid, subjectUid = subjectUid)
        val context = buildContext()
        val result = catalogUseCase.findAllWorks(option, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminWorkFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminWorkResponse(it) },
            )
        )
    }

    override fun adminCatalogReadWork(uid: String): ResponseEntity<AdminWorkResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        return ResponseEntity.ok(toAdminWorkResponse(catalogUseCase.findWork(uid, context)))
    }

    override fun adminCatalogCreateWork(requestBody: AdminWorkCreateRequest): ResponseEntity<Unit> {
        principalResolver.requireManager()
        val context = buildContext()
        val command = WorkCreateCommand(
            title = requestBody.title,
            description = requestBody.description,
            firstPublishDate = requestBody.firstPublishDate,
            olKey = requestBody.olKey,
            authorUids = requestBody.authorUids,
            subjectNames = requestBody.subjectNames,
        )
        catalogUseCase.createWork(command, context)
        return ResponseEntity.status(HttpStatus.CREATED).build()
    }

    override fun adminCatalogUpdateWork(
        uid: String,
        requestBody: AdminWorkUpdateRequest,
    ): ResponseEntity<AdminWorkResponse> {
        principalResolver.requireManager()
        val context = buildContext()
        val command = WorkUpdateCommand(
            title = requestBody.title,
            description = requestBody.description,
            firstPublishDate = requestBody.firstPublishDate,
            olKey = requestBody.olKey,
            authorUids = requestBody.authorUids,
            subjectNames = requestBody.subjectNames,
        )
        return ResponseEntity.ok(toAdminWorkResponse(catalogUseCase.updateWork(uid, command, context)))
    }

    override fun adminCatalogListAuthors(
        name: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminAuthorFindAllResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        val result = catalogUseCase.findAllAuthors(name, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminAuthorFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminAuthorItem(it) },
            )
        )
    }

    override fun adminCatalogListSubjects(
        name: String?,
        pageable: Pageable?,
    ): ResponseEntity<AdminSubjectFindAllResponse> {
        principalResolver.requireStaff()
        val context = buildContext()
        val result = catalogUseCase.findAllSubjects(name, pageable ?: Pageable.unpaged(), context)
        return ResponseEntity.ok(
            AdminSubjectFindAllResponse(
                total = result.totalElements,
                items = result.content.map { toAdminSubjectItem(it) },
            )
        )
    }
}
