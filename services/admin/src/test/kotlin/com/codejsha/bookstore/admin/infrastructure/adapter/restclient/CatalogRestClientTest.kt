package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.command.WorkCreateCommand
import com.codejsha.bookstore.admin.domain.model.command.WorkUpdateCommand
import com.codejsha.bookstore.admin.domain.model.external.Author
import com.codejsha.bookstore.admin.domain.model.external.Subject
import com.codejsha.bookstore.admin.domain.model.external.Work
import com.codejsha.bookstore.admin.domain.model.option.WorkQueryOption
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestInstance
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.client.HttpServerErrorException
import org.springframework.web.context.request.RequestContextHolder
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNull

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class CatalogRestClientTest {

    private val downstream = DownstreamStub()

    private val client = CatalogRestClient(
        downstream.config.catalogWorkApi(),
        downstream.config.catalogAuthorApi(),
        downstream.config.catalogSubjectApi(),
    )

    @AfterEach
    fun clearRequest() {
        RequestContextHolder.resetRequestAttributes()
    }

    @AfterAll
    fun stopDownstream() {
        downstream.stop()
    }

    @Test
    fun `countWorks_populatedCatalog_requestsOneItemAndRelaysTotal`() {
        downstream.respond("""{"total":137,"items":[$FULL_WORK_JSON]}""")

        val total = client.countWorks()

        assertEquals(137, total)
        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/works", request.path)
        assertEquals(mapOf("size" to "1"), request.query)
    }

    @Test
    fun `findAllWorks_filtersWithSortedPage_forwardsEveryParameterAndMapsItems`() {
        downstream.respond("""{"total":151,"items":[$FULL_WORK_JSON,$BARE_WORK_JSON]}""")

        val page = client.findAllWorks(
            WorkQueryOption(title = "Dune & Co", authorUid = "author-1", subjectUid = "subject-1"),
            SORTED_PAGE,
        )

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/works", request.path)
        assertEquals(
            mapOf("title" to "Dune & Co", "author_uid" to "author-1", "subject_uid" to "subject-1") +
                SORTED_PAGE_QUERY,
            request.query,
        )
        assertEquals(151, page.totalElements)
        assertEquals(SORTED_PAGE, page.pageable)
        assertEquals(listOf(FULL_WORK, BARE_WORK), page.content)
    }

    @Test
    fun `findAllWorks_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllWorks(WorkQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(0, page.totalElements)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findWork_existingUid_readsWorkByUid`() {
        downstream.respond(FULL_WORK_JSON)

        val work = client.findWork("work-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/works/work-1", request.path)
        assertEquals(emptyMap(), request.query)
        assertEquals(FULL_WORK, work)
    }

    @Test
    fun `findWork_optionalFieldsAbsent_mapsThemToNull`() {
        downstream.respond(BARE_WORK_JSON)

        val work = client.findWork("work-2")

        assertEquals(BARE_WORK, work)
        assertNull(work.updatedAt)
    }

    @Test
    fun `findWork_inboundBearerToken_replaysItDownstream`() {
        bindInboundAuthorization(CALLER_AUTHORIZATION)
        downstream.respond(FULL_WORK_JSON)

        client.findWork("work-1")

        assertEquals(CALLER_AUTHORIZATION, downstream.onlyRequest().authorization)
    }

    @Test
    fun `findWork_noInboundRequest_sendsNoAuthorization`() {
        downstream.respond(FULL_WORK_JSON)

        client.findWork("work-1")

        assertNull(downstream.onlyRequest().authorization)
    }

    @Test
    fun `findWork_unknownUid_propagatesNotFoundWithProblemBody`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        val failure = assertFailsWith<HttpClientErrorException.NotFound> { client.findWork("missing") }

        assertEquals(PROBLEM_NOT_FOUND, failure.responseBodyAsString)
    }

    @Test
    fun `findWork_downstreamUnavailable_propagatesServerError`() {
        downstream.respond(status = 503)

        val failure = assertFailsWith<HttpServerErrorException> { client.findWork("work-1") }

        assertEquals(HttpStatus.SERVICE_UNAVAILABLE, failure.statusCode)
    }

    @Test
    fun `createWork_fullCommand_postsSnakeCaseBodyWithoutCovers`() {
        downstream.respond(status = 201)

        client.createWork(
            WorkCreateCommand(
                title = "Dune",
                description = "Desert planet",
                firstPublishDate = "1965",
                olKey = "OL893415W",
                authorUids = listOf("author-1", "author-2"),
                subjectNames = listOf("Science fiction"),
            ),
        )

        val request = downstream.onlyRequest()
        assertEquals("POST", request.method)
        assertEquals("/api/v1/works", request.path)
        assertEquals(
            mapOf(
                "title" to "Dune",
                "description" to "Desert planet",
                "first_publish_date" to "1965",
                "ol_key" to "OL893415W",
                "author_uids" to listOf("author-1", "author-2"),
                "subject_names" to listOf("Science fiction"),
            ),
            request.json(),
        )
    }

    @Test
    fun `createWork_optionalFieldsNull_postsOnlyRequiredFields`() {
        downstream.respond(status = 201)

        client.createWork(
            WorkCreateCommand(
                title = "Untitled",
                description = null,
                firstPublishDate = null,
                olKey = null,
                authorUids = emptyList(),
                subjectNames = null,
            ),
        )

        assertEquals(
            mapOf("title" to "Untitled", "author_uids" to emptyList<String>()),
            downstream.onlyRequest().json(),
        )
    }

    @Test
    fun `createWork_duplicateOlKey_propagatesConflict`() {
        downstream.respond("""{"title":"Conflict","status":409}""", status = 409)

        val failure = assertFailsWith<HttpClientErrorException.Conflict> {
            client.createWork(
                WorkCreateCommand(
                    title = "Dune",
                    description = null,
                    firstPublishDate = null,
                    olKey = "OL893415W",
                    authorUids = listOf("author-1"),
                    subjectNames = null,
                ),
            )
        }

        assertEquals(HttpStatus.CONFLICT, failure.statusCode)
    }

    @Test
    fun `updateWork_fullCommand_putsSnakeCaseBodyAndMapsUpdatedWork`() {
        downstream.respond(FULL_WORK_JSON)

        val work = client.updateWork(
            "work-1",
            WorkUpdateCommand(
                title = "Dune",
                description = "Desert planet",
                firstPublishDate = "1965",
                olKey = "OL893415W",
                authorUids = listOf("author-1"),
                subjectNames = listOf("Science fiction"),
            ),
        )

        val request = downstream.onlyRequest()
        assertEquals("PUT", request.method)
        assertEquals("/api/v1/works/work-1", request.path)
        assertEquals(
            mapOf(
                "title" to "Dune",
                "description" to "Desert planet",
                "first_publish_date" to "1965",
                "ol_key" to "OL893415W",
                "author_uids" to listOf("author-1"),
                "subject_names" to listOf("Science fiction"),
            ),
            request.json(),
        )
        assertEquals(FULL_WORK, work)
    }

    @Test
    fun `updateWork_titleOnly_putsOnlyTheTitle`() {
        downstream.respond(BARE_WORK_JSON)

        val work = client.updateWork(
            "work-2",
            WorkUpdateCommand(
                title = "Renamed",
                description = null,
                firstPublishDate = null,
                olKey = null,
                authorUids = null,
                subjectNames = null,
            ),
        )

        assertEquals(mapOf("title" to "Renamed"), downstream.onlyRequest().json())
        assertEquals(BARE_WORK, work)
    }

    @Test
    fun `updateWork_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        assertFailsWith<HttpClientErrorException.NotFound> {
            client.updateWork(
                "missing",
                WorkUpdateCommand(
                    title = "Renamed",
                    description = null,
                    firstPublishDate = null,
                    olKey = null,
                    authorUids = null,
                    subjectNames = null,
                ),
            )
        }
    }

    @Test
    fun `findAllAuthors_nameFilterWithSortedPage_forwardsParametersAndMapsAuthors`() {
        downstream.respond(
            """{"total":107,"items":[
                {"uid":"author-1","name":"Frank Herbert","bio":"American author","ol_key":"OL79034A"},
                {"uid":"author-2","name":"Brian Herbert"}
            ]}""",
        )

        val page = client.findAllAuthors("Herbert", SORTED_PAGE)

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/authors", request.path)
        assertEquals(mapOf("name" to "Herbert") + SORTED_PAGE_QUERY, request.query)
        assertEquals(107, page.totalElements)
        assertEquals(
            listOf(Author("author-1", "Frank Herbert"), Author("author-2", "Brian Herbert")),
            page.content,
        )
    }

    @Test
    fun `findAllAuthors_noNameUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        client.findAllAuthors(null, Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
    }

    @Test
    fun `findAllSubjects_nameFilterWithSortedPage_forwardsParametersAndMapsSubjects`() {
        downstream.respond(
            """{"total":103,"items":[{"uid":"subject-1","name":"Science fiction"},{"uid":"subject-2","name":"Ecology"}]}""",
        )

        val page = client.findAllSubjects("sci", SORTED_PAGE)

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/subjects", request.path)
        assertEquals(mapOf("name" to "sci") + SORTED_PAGE_QUERY, request.query)
        assertEquals(103, page.totalElements)
        assertEquals(
            listOf(Subject("subject-1", "Science fiction"), Subject("subject-2", "Ecology")),
            page.content,
        )
    }

    @Test
    fun `findAllSubjects_noNameUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        client.findAllSubjects(null, Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
    }

    private companion object {
        private const val FULL_WORK_JSON =
            """{
                "uid":"work-1",
                "title":"Dune",
                "description":"Desert planet",
                "first_publish_date":"1965",
                "ol_key":"OL893415W",
                "cover_uids":["cover-1"],
                "authors":[{"uid":"author-1","name":"Frank Herbert","bio":"American author"}],
                "subjects":[{"uid":"subject-1","name":"Science fiction"},{"uid":"subject-2","name":"Ecology"}],
                "created_at":"2026-01-02T03:04:00Z",
                "updated_at":"2026-02-03T04:05:00Z"
            }"""

        private const val BARE_WORK_JSON =
            """{"uid":"work-2","title":"Untitled","authors":[],"subjects":[],"created_at":"2026-03-04T05:06:00Z"}"""

        private val FULL_WORK = Work(
            uid = "work-1",
            title = "Dune",
            description = "Desert planet",
            firstPublishDate = "1965",
            olKey = "OL893415W",
            authors = listOf(Author("author-1", "Frank Herbert")),
            subjects = listOf(Subject("subject-1", "Science fiction"), Subject("subject-2", "Ecology")),
            createdAt = utc(2026, 1, 2, 3, 4),
            updatedAt = utc(2026, 2, 3, 4, 5),
        )

        private val BARE_WORK = Work(
            uid = "work-2",
            title = "Untitled",
            description = null,
            firstPublishDate = null,
            olKey = null,
            authors = emptyList(),
            subjects = emptyList(),
            createdAt = utc(2026, 3, 4, 5, 6),
            updatedAt = null,
        )
    }
}
