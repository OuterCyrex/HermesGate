namespace go dashboard

struct PanelDataResponse {
    1: i64 serviceNum
    2: i64 appNum
    3: i64 currentQps
    4: i64 todayRequestNum
}

struct FlowStatResponse {
    1: list<i64> today
    2: list<i64> yesterday
}

struct DashServiceStatItem {
    1: string name
    2: i64 value
}

struct DashServiceStatResponse {
    1: i64 total
    2: list<DashServiceStatItem> data
}

service dashboard {
    PanelDataResponse GetPanelData() (api.get="/dashboard/panel")
    FlowStatResponse GetFlowStatistics () (api.get="/dashboard/stat/flow")
    DashServiceStatResponse GetDashServiceStat () (api.get="/dashboard/stat/service")
}