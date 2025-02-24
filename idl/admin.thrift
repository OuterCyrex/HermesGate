namespace go admin

struct AdminLoginRequest {
    1: string username (api.body="username", api.vd="len($) >= 2 && len($) <= 20")
    2: string password (api.body="password", api.vd="len($) >= 6 && len($) <= 16")
}

struct AdminLoginResponse {
    1: string token
    2: string message
}

struct AdminInfoResponse {
    1: i32 id
    2: string username
    3: i64 loginTime
    4: string avatar
    5: string introduction
    6: list<string> roles
}

struct MessageResponse {
    1: string message
}

struct ChangePasswordRequest {
    1: string password (api.body="password", api.vd="len($) >= 6 && len($) <= 16")
}

service Admin {
    AdminLoginResponse AdminLogin(1: AdminLoginRequest req) (api.post="/admin/login")
    AdminInfoResponse AdminInfo() (api.get="/admin/info")
    MessageResponse AdminLogout() (api.get="/admin/logout")
    MessageResponse ChangePassword(1: ChangePasswordRequest req) (api.post="/admin/pwd")
}