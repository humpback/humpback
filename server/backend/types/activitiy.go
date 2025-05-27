package types

type ActivityAction string

const (
    ActivityActionLogin             ActivityAction = "Login"
    ActivityActionLogout            ActivityAction = "Logout"
    ActivityActionUpdate            ActivityAction = "Update"
    ActivityActionAdd               ActivityAction = "Add"
    ActivityActionDelete            ActivityAction = "Delete"
    ActivityActionChangePassword    ActivityAction = "ChangePassword"
    ActivityActionUpdateLabel       ActivityAction = "UpdateLabel"
    ActivityActionEnable            ActivityAction = "Enable"
    ActivityActionDisable           ActivityAction = "Disable"
    ActivityActionAddNode           ActivityAction = "AddNode"
    ActivityActionRemoveNode        ActivityAction = "RemoveNode"
    ActivityActionUpdateBasic       ActivityAction = "UpdateBasic"
    ActivityActionUpdateApplication ActivityAction = "UpdateApplication"
    ActivityActionUpdateDeployment  ActivityAction = "UpdateDeployment"
    ActivityActionStart             ActivityAction = "Start"
    ActivityActionRestart           ActivityAction = "Restart"
    ActivityActionStop              ActivityAction = "Stop"
)

type ActivityInfo struct {
    ActivityId   string         `json:"activityId"`
    Action       ActivityAction `json:"action"`
    Description  string         `json:"description"`
    OldContent   any            `json:"oldContent"`
    NewContent   any            `json:"newContent"`
    Operator     string         `json:"operator"`
    OperatorId   string         `json:"operatorId"`
    OperateAt    int64          `json:"operateAt"`
    ResourceId   string         `json:"resourceId"`
    ResourceName string         `json:"resourceName"`
    InstanceName string         `json:"instanceName,omitempty"`
    GroupName    string         `json:"groupName,omitempty"`
}
