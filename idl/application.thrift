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

struct AppListItemResponse {
    1: i32 ID
    2: string appID
    3: string name
    4: string secret
    5: string whiteIPs
    6: i64 qpd
    7: i64 qps
    8: i64 realQps
    9: i64 realQpd
}

struct AppListResponse {
    1: i64 total
    2: list<AppListItemResponse> data
}

struct AppListRequest {
    1: string info (api.query="info")
    2: required i32 pageNum (api.query="page_num")
    3: required i32 pageSize (api.query="page_size")
}

struct AppUpdateRequest {
    1: required i32 ID (api.path="id" api.vd="$ > 0")
    2: required string name (api.body="name"  api.vd="len($) < 128 && len($) > 0")
    3: string secret (api.body="secret")
    4: string whiteIPS (api.body="white_ips")
    5: i64 qpd (api.body="qpd" api.vd="$ > 0")
    6: i64 qps (api.body="qps" api.vd="$ > 0")
}

struct MessageResponse {
    1: string message
}

struct AppDeleteRequest {
    1: i32 ID (api.path="id" api.vd="$ > 0")
}

struct AppStaticRequest {
    1: i32 ID (api.path="id")
}

struct AppStaticResponse {
    1: list<i64> yesterday
    2: list<i64> today
}

service application {
    MessageResponse ApplicationAddHTTP (1: AppAddHttpRequest req) (api.post="/application/add/http")
    AppDetailResponse ApplicationDetail (1: AppDetailRequest req) (api.get="/application/detail/:id")
    MessageResponse AppUpdate (1: AppUpdateRequest req) (api.put="/application/update/:id")
    MessageResponse AppDelete (1: AppDeleteRequest req) (api.delete="/application/delete/:id")
    AppListResponse AppList (1: AppListRequest req) (api.get="/application/list")
    AppStaticResponse AppStatic (1: AppStaticRequest req) (api.get="/application/static/:id")
}