const isDevMode = () => {
  return process.env.NODE_ENV === 'development'
}

const getObjectCurrent = (key) => {
  const get = sessionStorage.getItem(key)
  if (get) {
    return JSON.parse(get)
  }
  return null
}
const setObjectCurrent = (key, value) => {
  if (value) {
    sessionStorage.setItem(key, JSON.stringify(value))
  }
}

const currentUserRes = () => {
  return getObjectCurrent('currentUserRes')
}
const setCurrentUserRes = (value) => {
  setObjectCurrent('currentUserRes', value)
}

const currentBaseProp = () => {
  return getObjectCurrent('currentBaseProp')
}
const setCurrentBaseProp = (value) => {
  setObjectCurrent('currentBaseProp', value)
}

const currentEnvId = () => {
  return parseInt(sessionStorage.getItem('currentEnvId'))
}
const setCurrentEnvId = (value) => {
  sessionStorage.setItem('currentEnvId', value)
}
const currentEnvName = () => {
  return sessionStorage.getItem('currentEnvName')
}
const setCurrentEnvName = (value) => {
  sessionStorage.setItem('currentEnvName', value)
}

const currentSpaceId = () => {
  return parseInt(sessionStorage.getItem('currentSpaceId'))
}
const setCurrentSpaceId = (value) => {
  sessionStorage.setItem('currentSpaceId', value)
}

const currentDocumentDetail = (id) => {
  return getObjectCurrent('currentDocument:' + id)
}
const setCurrentDocumentDetail = (id, value) => {
  setObjectCurrent('currentDocument:' + id, value)
}

const currentToken = () => {
  return localStorage.getItem('.token')
}
const setCurrentToken = (token) => {
  localStorage.setItem('.token', token)
}

const delCurrentToken = () => {
  sessionStorage.clear()
  localStorage.removeItem('.token')
}

export {
  isDevMode, currentBaseProp, setCurrentBaseProp, currentUserRes, setCurrentUserRes,
  currentEnvId, currentEnvName, setCurrentEnvId, setCurrentEnvName,
  currentSpaceId, setCurrentSpaceId,
  currentDocumentDetail, setCurrentDocumentDetail,
  currentToken, setCurrentToken, delCurrentToken
}
