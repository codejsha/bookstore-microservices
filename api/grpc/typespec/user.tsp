import "@typespec/protobuf";

using TypeSpec.Protobuf;

@package({
  name: "user.v1",
})
namespace ProtoUser;

// ─── Enum ───────────────────────────────────────────────────────────────────

enum UserStatus {
  USER_STATUS_UNSPECIFIED: 0,
  USER_STATUS_ACTIVE: 1,
  USER_STATUS_INACTIVE: 2,
  USER_STATUS_SUSPENDED: 3,
}

// ─── Entity ─────────────────────────────────────────────────────────────────

model User {
  @field(1) uid: string;
  @field(2) email: string;
  @field(3) first_name: string;
  @field(4) last_name: string;
  @field(5) phone?: string;
  @field(6) roles: string[];
  @field(7) status: UserStatus;
}

// ─── Messages ───────────────────────────────────────────────────────────────

model ListUsersRequest {
  @field(1) email?: string;
  @field(2) name?: string;
  @field(3) phone?: string;
  @field(4) status?: UserStatus;
  @field(5) page_size?: int32;
  @field(6) page_token?: string;
  @field(7) order_by?: string;
}

model ListUsersResponse {
  @field(1) users: User[];
  @field(2) next_page_token?: string;
  @field(3) total_size: int32;
}

model FindUserRequest {
  @field(1) uid: string;
}

model FindUserResponse {
  @field(1) user: User;
}

// ─── Service ────────────────────────────────────────────────────────────────

@Protobuf.service
interface UserService {
  ListUsers(...ListUsersRequest): ListUsersResponse;
  FindUser(...FindUserRequest): FindUserResponse;
}
