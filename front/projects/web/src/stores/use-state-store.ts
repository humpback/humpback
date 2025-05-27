import { cloneDeep, sortBy } from "lodash-es"
import { GroupInfo, NodeInfo } from "@/types"
import { PageGroupDetail, SessionStorageCurrentGroupId } from "@/models"

type TotalKey = PageGroupDetail.Services | PageGroupDetail.Nodes

interface StateGroupInfo extends GroupInfo {
  total: {
    [key in TotalKey]: number
  }
  nodeList: NodeInfo[]
}

export default defineStore("state", () => {
  const route = useRoute()
  const groups = ref<{ [key: string]: StateGroupInfo }>({})
  const services = ref<{ [key: string]: ServiceInfo }>({})
  const nodes = ref<{ [key: string]: NodeInfo }>({})

  function setGroup(id: string, groupInfo: GroupInfo) {
    sessionStorage.setItem(SessionStorageCurrentGroupId, id)
    const info = groups.value[id]
    const total = { services: info?.total.services || 0, nodes: groupInfo.nodes.length }
    groups.value[id] = { ...cloneDeep(groupInfo), total: total, nodeList: info?.nodeList || [] } as StateGroupInfo
  }

  function setGroupNodeList(id: string, nodeList: NodeInfo[]) {
    sessionStorage.setItem(SessionStorageCurrentGroupId, id)
    let groupInfo =
      groups.value[id] ||
      Object.assign({}, NewGroupEmptyInfo(), {
        groupId: id,
        total: { services: 0, nodes: 0 },
        nodeList: sortBy(nodeList, ["ipAddress"])
      })
    groupInfo.nodeList = sortBy(nodeList, ["ipAddress"])
    groups.value[id] = groupInfo
  }

  function setGroupTotal(id?: string, serviceTotal?: number, nodeTotal?: number) {
    const key = id || (route.params["groupId"] as string)
    const info =
      cloneDeep(groups.value[key]) ||
      Object.assign({}, NewGroupEmptyInfo(), {
        groupId: key,
        total: { services: 0, nodes: 0 },
        nodeList: []
      })
    info.total = {
      services: typeof serviceTotal === "undefined" ? info.total.services : serviceTotal,
      nodes: typeof nodeTotal === "undefined" ? info.total.nodes : nodeTotal
    }
    groups.value[key] = info
  }

  function getGroup(id?: string): StateGroupInfo | undefined {
    return groups.value[id || (route.params["groupId"] as string)]
  }

  function setService(id: string, serviceInfo: ServiceInfo) {
    services.value[id] = cloneDeep(serviceInfo)
  }

  function getService(id?: string): ServiceInfo | undefined {
    return services.value[id || (route.params["serviceId"] as string)]
  }

  function setNode(id: string, nodeInfo: NodeInfo) {
    nodes.value[id] = cloneDeep(nodeInfo)
  }

  function getNode(id?: string): NodeInfo | undefined {
    return nodes.value[id || (route.params["nodeId"] as string)]
  }

  return {
    setGroup,
    setGroupTotal,
    setGroupNodeList,
    getGroup,
    setService,
    getService,
    setNode,
    getNode
  }
})
