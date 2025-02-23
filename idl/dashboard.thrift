namespace go dashboard

struct PanelDataResponse {
    1: i64 serviceNum
    2: i64 appNum
    3: i64 currentQps
    4: i64 todayRequestNum
}

service dashboard {
    PanelDataResponse GetPanelData() (api.get="/dashboard/panel")
}