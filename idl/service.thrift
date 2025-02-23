namespace go services

struct ServiceListResponse {
     1: i64 total
     2: list<ServiceListItemResponse> data
}

struct ServiceListItemResponse {
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
    5: required i8 needHTTPS(api.body="need_https" api.vd="$ <= 1 && $ >= 0")
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

struct ServiceUpdateHTTPRequest {
    1: required i32 ID (api.path="id" api.vd="$ > 0")
    2: required i8 needHTTPS(api.body="need_https" api.vd="$ <= 1 && $ >= 0")
    3: required i8 needStripUri (api.body="need_strip_uri" api.vd="$ <= 1 && $ >= 0")
    4: required i8 needWebsocket (api.body="need_websocket" api.vd="$ <= 1 && $ >= 0")
    5: string urlRewrite (api.body="url_rewrite")
    6: string headerTransfer (api.body="header_transfer")
    7: required i8 openAuth (api.body="open_auth" api.vd="$ <= 1 && $ >= 0")
    8: string blackList (api.body="black_list")
    9: string whiteList (api.body="white_list")
    10: i32 clientIPFlowLimit (api.body="client_ip_flow_limit")
    11: i32 serviceFlowLimit (api.body="service_flow_limit")
    12: required i8 roundType (api.body="round_type" api.vd="$ <= 2 && $ >= 0")
    13: required string ipList (api.body="ip_list")
    14: required string weightList (api.body="weightList")
    15: i32 upstreamConnectTimeout (api.body="upstream_connect_timeout")
    16: i32 upstreamHeaderTimeout (api.body="upstream_header_timeout")
    17: i32 upstreamIdleTimeout (api.body="upstream_idle_timeout")
    18: i32 upstreamMaxIdle (api.body="upstream_max_idle")
}

struct ServiceAddGrpcRequest {
    1: required string serviceName (api.body="service_name" api.vd="6 <= len($) && len($) <= 128")
    2: required string serviceDesc (api.body="service_desc" api.vd="len($) < 255 && len($) > 0")
    3: required i32 port (api.body="port" api.vd = "$ > 8000 && $ < 65535")
    4: string headerTransfer (api.body="header_transfer")
    5: required i8 openAuth (api.body="open_auth")
    6: string blackList (api.body="black_list")
    7: string whiteList (api.body="white_list")
    8: string whiteHostName (api.body="white_host_name")
    9: i32 clientIPFlowLimit (api.body="client_ip_flow_limit")
    10: i32 serviceFlowLimit (api.body="service_flow_limit")
    11: required i8 roundType (api.body="round_type" api.vd="$ <= 2 && $ >= 0")
    12: required string ipList (api.body="ip_list")
    13: required string weightList (api.body="weightList")
    14: string forbidList (api.body="forbid_list")
}

struct ServiceUpdateGrpcRequest {
    1: required i32 ID (api.path="id" api.vd="$ > 0")
    2: string headerTransfer (api.body="header_transfer")
    3: required i8 openAuth (api.body="open_auth")
    4: string blackList (api.body="black_list")
    5: string whiteList (api.body="white_list")
    6: string whiteHostName (api.body="white_host_name")
    7: i32 clientIPFlowLimit (api.body="client_ip_flow_limit")
    8: i32 serviceFlowLimit (api.body="service_flow_limit")
    9: required i8 roundType (api.body="round_type" api.vd="$ <= 2 && $ >= 0")
    10: required string ipList (api.body="ip_list")
    11: required string weightList (api.body="weightList")
    12: string forbidList (api.body="forbid_list")
}

struct ServiceAddTcpRequest {
    1: required string serviceName (api.body="service_name" api.vd="6 <= len($) && len($) <= 128")
    2: required string serviceDesc (api.body="service_desc" api.vd="len($) < 255 && len($) > 0")
    3: required i32 port (api.body="port" api.vd = "$ > 8000 && $ < 65535")
    4: required i8 openAuth (api.body="open_auth")
    5: string blackList (api.body="black_list")
    6: string whiteList (api.body="white_list")
    7: string whiteHostName (api.body="white_host_name")
    8: i32 clientIPFlowLimit (api.body="clientip_flow_limit")
    9: i32 serviceFlowLimit (api.body="service_flow_limit")
    10: required i8 roundType (api.body="round_type" api.vd="$ <= 2 && $ >= 0")
    11: required string ipList (api.body="ip_list")
    12: required string weightList (api.body="weight_list")
    13: string forbidList (api.body="forbid_list")
}

struct ServiceUpdateTcpRequest {
    1: required i32 ID (api.path="id" api.vd="$ > 0")
    2: required i8 openAuth (api.body="open_auth")
    3: string blackList (api.body="black_list")
    4: string whiteList (api.body="white_list")
    5: string whiteHostName (api.body="white_host_name")
    6: i32 clientIPFlowLimit (api.body="clientip_flow_limit")
    7: i32 serviceFlowLimit (api.body="service_flow_limit")
    8: required i8 roundType (api.body="round_type" api.vd="$ <= 2 && $ >= 0")
    9: required string ipList (api.body="ip_list")
    10: required string weightList (api.body="weight_list")
    11: string forbidList (api.body="forbid_list")
}

// the interfaces below all belong to serviceDetail
// DO NOT get mistaken

struct ServiceInfoPart {
    1: i32 ID
    2: i8 LoadType
    3: string ServiceName
    4: string ServiceDesc
}

struct ServiceHttpRulePart {
    1: i32 ID
    2: i32 ServiceID
    3: i8 RuleType
    4: string Rule
    5: i8 NeedHttps
    6: i8 NeedWebsocket
    7: i8 NeedStripUri
    8: string UrlRewrite
    9: string HeaderTransfer
}

struct ServiceGRPCRulePart {
    1: i32 ID
    2: i32 ServiceID
    3: i32 Port
    4: string HeaderTrans
}

struct ServiceTcpRulePart {
    1: i32 ID
    2: i32 ServiceID
    3: i32 Port
}

struct ServiceLoadBalancePart {
    1: i32 ID
    2: i32 ServiceID
    3: i32 CheckMethod
    4: i32 CheckTimeout
    5: i32 CheckInterval
    6: i8 RoundType
    7: string IpList
    8: string WeightList
    9: string ForbidList
    10: i32 UpstreamConnectTimeout
    11: i32 UpstreamHeaderTimeout
    12: i32 UpstreamIdleTimeout
    13: i32 UpstreamMaxIdle
}

struct ServiceAccessControlPart {
    1: i32 ID
    2: i32 ServiceID
    3: i8 OpenAuth
    4: string BlackList
    5: string WhiteList
    6: string WhiteHostName
    7: i32 ClientIPFlowLimit
    8: i32 ServiceFlowLimit
}

struct ServiceDetailResponse {
    1: optional ServiceInfoPart Info
    2: optional ServiceHttpRulePart Http
    3: optional ServiceTcpRulePart Tcp
    4: optional ServiceGRPCRulePart Grpc
    5: optional ServiceLoadBalancePart LoadBalance
    6: optional ServiceAccessControlPart AccessControl
}

struct ServiceDetailRequest {
    1: i32 ID (api.path="id" api.vd="$ > 0")
}

struct ServiceStaticResponse {
    1: list<i64> today
    2: list<i64> yesterday
}

struct ServiceStaticRequest {
    1: i32 ID (api.path="id" api.vd="$ > 0")
}

service services {
    ServiceListResponse ServiceList(1: ServiceListRequest req) (api.get="/service/list")
    MessageResponse ServiceDelete(1: ServiceDeleteRequest req) (api.delete="/service/delete/:id")
    MessageResponse ServiceAddHTTP(1: ServiceAddHTTPRequest req) (api.post="/service/add/http")
    MessageResponse ServiceUpdateHTTP(1: ServiceUpdateHTTPRequest req) (api.put="/service/update/http/:id")
    ServiceDetailResponse ServiceDetail (1: ServiceDetailRequest req) (api.get="/service/detail/:id")
    ServiceStaticResponse ServiceStatic (1: ServiceStaticRequest req) (api.get="/service/static/:id")
    MessageResponse ServiceAddGRPC (1: ServiceAddGrpcRequest req) (api.post="/service/add/grpc")
    MessageResponse ServiceAddTCP (1: ServiceAddTcpRequest req) (api.post="/service/add/tcp")
    MessageResponse ServiceUpdateGRPC (1: ServiceUpdateGrpcRequest req) (api.put="/service/update/grpc/:id")
    MessageResponse ServiceUpdateTCP (1: ServiceUpdateTcpRequest req) (api.put="/service/update/tcp/:id")
}

