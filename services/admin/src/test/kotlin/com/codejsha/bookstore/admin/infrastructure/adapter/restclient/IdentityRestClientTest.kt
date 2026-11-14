package com.codejsha.bookstore.admin.infrastructure.adapter.restclient

import com.codejsha.bookstore.admin.domain.model.UnsupportedValueException
import com.codejsha.bookstore.admin.domain.model.command.FlagRiskCommand
import com.codejsha.bookstore.admin.domain.model.external.RiskEntry
import com.codejsha.bookstore.admin.domain.model.external.User
import com.codejsha.bookstore.admin.domain.model.option.UserQueryOption
import org.junit.jupiter.api.AfterAll
import org.junit.jupiter.api.AfterEach
import org.junit.jupiter.api.Test
import org.junit.jupiter.api.TestInstance
import org.springframework.data.domain.Pageable
import org.springframework.http.HttpStatus
import org.springframework.web.client.HttpClientErrorException
import org.springframework.web.context.request.RequestContextHolder
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class IdentityRestClientTest {

    private val downstream = DownstreamStub()

    private val client = IdentityRestClient(
        downstream.config.identityUserApi(),
        downstream.config.identityRiskApi(),
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
    fun `findAllUsers_everyFilterWithSortedPage_forwardsParametersAndMapsUsers`() {
        downstream.respond("""{"total":140,"items":[$FULL_USER_JSON,$BARE_USER_JSON]}""")

        val page = client.findAllUsers(
            UserQueryOption(email = "paul@arrakis.example", name = "Paul", phone = "+82-10", status = "SUSPENDED"),
            SORTED_PAGE,
        )

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/users", request.path)
        assertEquals(
            mapOf(
                "email" to "paul@arrakis.example",
                "name" to "Paul",
                "phone" to "+82-10",
                "status" to "SUSPENDED",
            ) + SORTED_PAGE_QUERY,
            request.query,
        )
        assertEquals(140, page.totalElements)
        assertEquals(SORTED_PAGE, page.pageable)
        assertEquals(listOf(FULL_USER, BARE_USER), page.content)
    }

    @Test
    fun `findAllUsers_noFilterUnpaged_sendsNoQueryParameters`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val page = client.findAllUsers(UserQueryOption(), Pageable.unpaged())

        assertEquals(emptyMap(), downstream.onlyRequest().query)
        assertEquals(emptyList(), page.content)
    }

    @Test
    fun `findAllUsers_unknownStatus_failsWithoutCallingDownstream`() {
        downstream.respond("""{"total":0,"items":[]}""")

        val failure = assertFailsWith<UnsupportedValueException> {
            client.findAllUsers(UserQueryOption(status = "BANNED"), SORTED_PAGE)
        }

        assertEquals("unknown user status: BANNED", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `findUser_existingUid_readsUserByUid`() {
        downstream.respond(FULL_USER_JSON)

        val user = client.findUser("user-1")

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/users/user-1", request.path)
        assertEquals(FULL_USER, user)
    }

    @Test
    fun `findUser_optionalFieldsAbsent_mapsThemToNull`() {
        downstream.respond(BARE_USER_JSON)

        assertEquals(BARE_USER, client.findUser("user-2"))
    }

    @Test
    fun `findUser_unknownUid_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        val failure = assertFailsWith<HttpClientErrorException.NotFound> { client.findUser("missing") }

        assertEquals(PROBLEM_NOT_FOUND, failure.responseBodyAsString)
    }

    @Test
    fun `findUser_inboundBearerToken_replaysItDownstream`() {
        bindInboundAuthorization(CALLER_AUTHORIZATION)
        downstream.respond(FULL_USER_JSON)

        client.findUser("user-1")

        assertEquals(CALLER_AUTHORIZATION, downstream.onlyRequest().authorization)
    }

    @Test
    fun `updateRoles_knownRoles_putsRoleWireValuesAndMapsUser`() {
        downstream.respond(FULL_USER_JSON)

        val user = client.updateRoles("user-1", listOf("STAFF", "MANAGER"))

        val request = downstream.onlyRequest()
        assertEquals("PUT", request.method)
        assertEquals("/api/v1/users/user-1/roles", request.path)
        assertEquals(mapOf("roles" to listOf("STAFF", "MANAGER")), request.json())
        assertEquals(FULL_USER, user)
    }

    @Test
    fun `updateRoles_unrecognisedRole_failsWithoutCallingDownstream`() {
        downstream.respond(FULL_USER_JSON)

        val failure = assertFailsWith<UnsupportedValueException> {
            client.updateRoles("user-1", listOf("STAFF", "ROOT"))
        }

        assertEquals("unknown role: ROOT", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `updateRoles_unknownSentinelRole_failsWithoutCallingDownstream`() {
        downstream.respond(FULL_USER_JSON)

        val failure = assertFailsWith<UnsupportedValueException> {
            client.updateRoles("user-1", listOf("UNKNOWN"))
        }

        assertEquals("unknown role: UNKNOWN", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `updateRoles_downstreamForbidden_propagatesForbidden`() {
        downstream.respond("""{"title":"Forbidden","status":403}""", status = 403)

        val failure = assertFailsWith<HttpClientErrorException.Forbidden> {
            client.updateRoles("user-1", listOf("MANAGER"))
        }

        assertEquals(HttpStatus.FORBIDDEN, failure.statusCode)
    }

    @Test
    fun `suspendUser_existingUid_postsSuspendAndMapsUser`() {
        downstream.respond(FULL_USER_JSON)

        val user = client.suspendUser("user-1")

        val request = downstream.onlyRequest()
        assertEquals("POST", request.method)
        assertEquals("/api/v1/users/user-1/suspend", request.path)
        assertEquals("", request.body)
        assertEquals(FULL_USER, user)
    }

    @Test
    fun `reactivateUser_existingUid_postsReactivateAndMapsUser`() {
        downstream.respond(BARE_USER_JSON)

        val user = client.reactivateUser("user-2")

        val request = downstream.onlyRequest()
        assertEquals("POST", request.method)
        assertEquals("/api/v1/users/user-2/reactivate", request.path)
        assertEquals("", request.body)
        assertEquals(BARE_USER, user)
    }

    @Test
    fun `deactivateUser_existingUid_postsDeactivateAndMapsUser`() {
        downstream.respond(DEACTIVATED_USER_JSON)

        val user = client.deactivateUser("user-3")

        val request = downstream.onlyRequest()
        assertEquals("POST", request.method)
        assertEquals("/api/v1/users/user-3/deactivate", request.path)
        assertEquals("", request.body)
        assertEquals("DEACTIVATED", user.status)
        assertEquals("user-3", user.uid)
    }

    @Test
    fun `suspendUser_alreadySuspended_propagatesBadRequest`() {
        downstream.respond("""{"title":"Bad Request","status":400}""", status = 400)

        assertFailsWith<HttpClientErrorException.BadRequest> { client.suspendUser("user-1") }
    }

    @Test
    fun `listRisk_flaggedPrincipals_mapsEveryEntry`() {
        downstream.respond("""{"entries":[$BLOCK_ENTRY_JSON,$RESTRICT_ENTRY_JSON]}""")

        val entries = client.listRisk()

        val request = downstream.onlyRequest()
        assertEquals("GET", request.method)
        assertEquals("/api/v1/risk", request.path)
        assertEquals(emptyMap(), request.query)
        assertEquals(listOf(BLOCK_ENTRY, RESTRICT_ENTRY), entries)
    }

    @Test
    fun `listRisk_noFlags_yieldsEmptyList`() {
        downstream.respond("""{"entries":[]}""")

        assertEquals(emptyList(), client.listRisk())
    }

    @Test
    fun `flagRisk_blockWithTtl_putsWireLevelAndMapsEntry`() {
        downstream.respond(BLOCK_ENTRY_JSON)

        val entry = client.flagRisk("user-1", FlagRiskCommand(level = "block", reason = "card testing", ttlSeconds = 3600))

        val request = downstream.onlyRequest()
        assertEquals("PUT", request.method)
        assertEquals("/api/v1/risk/user-1", request.path)
        assertEquals(mapOf("level" to "block", "reason" to "card testing", "ttl_seconds" to 3600), request.json())
        assertEquals(BLOCK_ENTRY, entry)
    }

    @Test
    fun `flagRisk_restrictWithoutTtl_omitsTtlAndMapsEntry`() {
        downstream.respond(RESTRICT_ENTRY_JSON)

        val entry = client.flagRisk("user-2", FlagRiskCommand(level = "restrict", reason = "probing", ttlSeconds = null))

        assertEquals(mapOf("level" to "restrict", "reason" to "probing"), downstream.onlyRequest().json())
        assertEquals(RESTRICT_ENTRY, entry)
    }

    @Test
    fun `flagRisk_unknownLevel_failsWithoutCallingDownstream`() {
        downstream.respond(BLOCK_ENTRY_JSON)

        val failure = assertFailsWith<UnsupportedValueException> {
            client.flagRisk("user-1", FlagRiskCommand(level = "BLOCK", reason = "abuse", ttlSeconds = null))
        }

        assertEquals("unknown risk level: BLOCK", failure.message)
        assertEquals(0, downstream.requestCount())
    }

    @Test
    fun `unflagRisk_flaggedPrincipal_deletesTheFlag`() {
        downstream.respond(status = 204)

        client.unflagRisk("user-1")

        val request = downstream.onlyRequest()
        assertEquals("DELETE", request.method)
        assertEquals("/api/v1/risk/user-1", request.path)
    }

    @Test
    fun `unflagRisk_principalNotFlagged_propagatesNotFound`() {
        downstream.respond(PROBLEM_NOT_FOUND, status = 404)

        assertFailsWith<HttpClientErrorException.NotFound> { client.unflagRisk("user-9") }
    }

    private companion object {
        private const val FULL_USER_JSON =
            """{
                "uid":"user-1",
                "email":"paul@arrakis.example",
                "first_name":"Paul",
                "last_name":"Atreides",
                "phone":"+82-10-1234-5678",
                "status":"SUSPENDED",
                "roles":["USER","STAFF","MANAGER"],
                "last_login_at":"2026-05-06T07:08:00Z",
                "created_at":"2026-01-02T03:04:00Z",
                "updated_at":"2026-02-03T04:05:00Z"
            }"""

        private const val BARE_USER_JSON =
            """{
                "uid":"user-2",
                "email":"chani@arrakis.example",
                "first_name":"Chani",
                "last_name":"Kynes",
                "status":"ACTIVE",
                "roles":[],
                "created_at":"2026-03-04T05:06:00Z"
            }"""

        private const val DEACTIVATED_USER_JSON =
            """{
                "uid":"user-3",
                "email":"gone@arrakis.example",
                "first_name":"Gone",
                "last_name":"User",
                "status":"DEACTIVATED",
                "roles":["USER"],
                "created_at":"2026-03-04T05:06:00Z"
            }"""

        private const val BLOCK_ENTRY_JSON =
            """{
                "user_uid":"user-1",
                "level":"block",
                "reason":"card testing",
                "flagged_by":"admin-1",
                "flagged_at":"2026-06-07T08:09:00Z",
                "expires_at":"2026-06-07T09:09:00Z"
            }"""

        private const val RESTRICT_ENTRY_JSON =
            """{
                "user_uid":"user-2",
                "level":"restrict",
                "reason":"probing",
                "flagged_at":"2026-06-08T08:09:00Z",
                "expires_at":"2026-06-09T08:09:00Z"
            }"""

        private val FULL_USER = User(
            uid = "user-1",
            email = "paul@arrakis.example",
            firstName = "Paul",
            lastName = "Atreides",
            phone = "+82-10-1234-5678",
            status = "SUSPENDED",
            roles = listOf("USER", "STAFF", "MANAGER"),
            lastLoginAt = utc(2026, 5, 6, 7, 8),
            createdAt = utc(2026, 1, 2, 3, 4),
            updatedAt = utc(2026, 2, 3, 4, 5),
        )

        private val BARE_USER = User(
            uid = "user-2",
            email = "chani@arrakis.example",
            firstName = "Chani",
            lastName = "Kynes",
            phone = null,
            status = "ACTIVE",
            roles = emptyList(),
            lastLoginAt = null,
            createdAt = utc(2026, 3, 4, 5, 6),
            updatedAt = null,
        )

        private val BLOCK_ENTRY = RiskEntry(
            userUid = "user-1",
            level = "block",
            reason = "card testing",
            flaggedBy = "admin-1",
            flaggedAt = utc(2026, 6, 7, 8, 9),
            expiresAt = utc(2026, 6, 7, 9, 9),
        )

        private val RESTRICT_ENTRY = RiskEntry(
            userUid = "user-2",
            level = "restrict",
            reason = "probing",
            flaggedBy = null,
            flaggedAt = utc(2026, 6, 8, 8, 9),
            expiresAt = utc(2026, 6, 9, 8, 9),
        )
    }
}
