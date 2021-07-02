import {currentEnvId, currentSpaceId} from 'src/sessions'
import http from 'src/requests'


const buildTermPath = (eid, sid, pid, did, pod, cname, res) => {
  return `/term/?eid=${eid}&sid=${sid}&pid=${pid}&did=${did}&pod=${pod}&cname=${cname}&select=${res}`
}

const navigateToTerm = (nav, res, data, inner = true) => {
  let eid = currentEnvId()
  let sid = currentSpaceId()
  let pid = ''
  let did = ''
  if (res !== 'inner') {
    eid = data.envId ? data.envId : eid
    sid = data.spaceId ? data.spaceId : sid
    pid = data.appId ? data.appId : ''
    did = data.deployId ? data.deployId : ''
  }
  if (inner) {
    nav('/app/terminal?in=' + btoa(JSON.stringify([eid, sid, pid, did, data.podName, data.containers[0].name, res])))
  } else {
    window.location.href = http.apiUrl + buildTermPath(eid, sid, pid, did, data.podName, data.containers[0].name, res)
  }
}

export default navigateToTerm
