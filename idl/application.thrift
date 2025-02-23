namespace go application

struct AppAddHttpRequest {
    1: required string appID (api.body="app_id" api.vd="len($) < 128 && len($) > 0")
    2: required string name (api.body="name"  api.vd="len($) < 128 && len($) > 0")
    3: string secret (api.body="secret")
    4: string whiteIPS (api.body="white_ips")
    5: i64 qpd (api.body="qpd" api.vd="$ > 0")
    6: i64 qps (api.body="qps" api.vd="$ > 0")
}

struct AppDetailRequest {
    1: required i32 ID (api.path="id" api.vd="$ > 0")
}

struct AppDetailResponse {
    1: i32 ID
    2: string appID
    3: string name
    4: string secret
    5: string whiteIPs
    6: i64 qpd
    7: i64 qps
}

struct MessageResponse {
    1: string message
}

service application {
    MessageResponse ApplicationAddHTTP (1: AppAddHttpRequest req) (api.post="/application/add/http")
    AppDetailResponse ApplicationDetail (1: AppDetailRequest req) (api.get="/application/detail/:id")
}