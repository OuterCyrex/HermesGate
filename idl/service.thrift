namespace go services

struct ServiceListResponse {
     1: i64 total
     2: list<ServiceInfoResponse> data
}

struct ServiceInfoResponse {
   1: i32 Id
   2: string serviceName
   3: string serviceDesc
   4: string loadType
   5: string serviceAddr
   6: i32 totalNode
   7: i64 qps
   8: i64 qpd
}

struct ServiceListRequest {
    1: string info (api.query="info")
    2: required i32 pageNum (api.query="pageNum" api.vd="$ >= 0")
    3: required i32 pageSize (api.query="pageSize"  api.vd="$ > 0")
}

struct ServiceDeleteRequest {
    1: required i32 ID (api.path="id" api.vd="$ > 0")
}

struct MessageResponse {
    1: string message
}

struct ServiceAddHTTPRequest {
    1: required string serviceName (api.body="service_name" api.vd="6 <= len($) && len($) <= 128")
    2: required string serviceDesc (api.body="service_desc" api.vd="len($) < 255 && len($) > 0")
    3: required i8 ruleType (api.body="rule_type" api.vd="$ <= 1 && $ >= 0")
    4: required string rule (api.body="rule")
    5: required i8 needHTTP(api.body="need_http" api.vd="$ <= 1 && $ >= 0")
    6: required i8 needStripUri (api.body="need_strip_uri" api.vd="$ <= 1 && $ >= 0")
    7: required i8 needWebsocket (api.body="need_websocket" api.vd="$ <= 1 && $ >= 0")
    8: string urlRewrite (api.body="url_rewrite")
    9: string headerTransfer (api.body="header_transfer")
    10: required i8 openAuth (api.body="open_auth" api.vd="$ <= 1 && $ >= 0")
    11: string blackList (api.body="black_list")
    12: string whiteList (api.body="white_list")
    13: i32 clientIPFlowLimit (api.body="client_ip_flow_limit")
    14: i32 serviceFlowLimit (api.body="service_flow_limit")
    15: required i8 roundType (api.body="round_type" api.vd="$ <= 2 && $ >= 0")
    16: required string ipList (api.body="ip_list")
    17: required string weightList (api.body="weightList")
    18: i32 upstreamConnectTimeout (api.body="upstream_connect_timeout")
    19: i32 upstreamHeaderTimeout (api.body="upstream_header_timeout")
    20: i32 upstreamIdleTimeout (api.body="upstream_idle_timeout")
    21: i32 upstreamMaxIdle (api.body="upstream_max_idle")
}

service services {
    ServiceListResponse ServiceList(1: ServiceListRequest req) (api.get="/service/list")
    MessageResponse ServiceDelete(1: ServiceDeleteRequest req) (api.delete="/service/delete")
    MessageResponse ServiceAddHTTP(1: ServiceAddHTTPRequest req) (api.post="/service/add/http")
}