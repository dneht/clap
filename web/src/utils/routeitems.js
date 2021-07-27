import {
  AlarmOn as PlanIcon,
  Apps as AppsIcon,
  Build as ToolIcon,
  Dashboard as DashboardIcon,
  GroupAdd as RoleIcon,
  InsertChart as DeployIcon,
  Kitchen as EnvIcon,
  People as UserIcon,
} from '@material-ui/icons'
import http from 'src/requests'
import {currentUserRes} from 'src/sessions'
import {ShowSnackbar} from 'src/utils/globalshow'

const iconMap = new Map()
iconMap.set('/app/dashboard', DashboardIcon)
iconMap.set('/app/environment', EnvIcon)
iconMap.set('/app/projects', AppsIcon)
iconMap.set('/app/deploys', DeployIcon)
iconMap.set('/app/plans', PlanIcon)
iconMap.set('/app/accounts', UserIcon)
iconMap.set('/app/roles', RoleIcon)
iconMap.set('/app/tools', ToolIcon)

const userResource = currentUserRes()
const handleResource = function (data) {
  const routeMap = new Map()
  const routeData = []
  if (data.memuList) {
    data.memuList.forEach(one => {
      if (one) {
        one.icon = iconMap.get(one.href)
        if (!routeMap.get(one.href)) {
          routeMap.set(one.href, true)
          routeData.push(one)
        }
      }
    })
  }
  if (routeData.length === 0) {
    if (!routeMap.get('/app/dashboard')) {
      routeMap.set('/app/dashboard', true)
      routeData.push({
        href: '/app/dashboard',
        icon: DashboardIcon,
        title: '面板'
      })
    }
  }
  return routeData
}

const routeItems = () => {
  return userResource ? new Promise((resolve) => {
    resolve(handleResource(userResource))
  }) : http.initRes((data) => {
    return handleResource(data)
  }).catch((err) => {
    ShowSnackbar('cannot get res: ' + err, 'error')
  })
}
export default routeItems
