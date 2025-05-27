package types

type ResourceTotalInfo struct {
    Services          int                             `json:"services"`
    Nodes             int                             `json:"nodes"`
    Users             int                             `json:"users"`
    Groups            int                             `json:"groups"`
    OwnGroups         int                             `json:"ownGroups"`
    OwnServices       int                             `json:"ownServices"`
    EnableOwnServices int                             `json:"enableOwnServices"`
    EnableOwnNodes    int                             `json:"enableOwnNodes"`
    ExceptionServices []*ResourceExceptionServiceInfo `json:"exceptionServices"`
    AbnormalNodes     []*Node                         `json:"abnormalNodes"`
}

type ResourceExceptionServiceInfo struct {
    GroupName string `json:"groupName"`
    *Service
}
