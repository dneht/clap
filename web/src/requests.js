import axios from 'axios'
import {currentToken, delCurrentToken, isDevMode, setCurrentBaseProp, setCurrentUserRes} from 'src/sessions'
import getCookie from 'src/utils/getcookie'
import {ShowSnackbar} from 'src/utils/globalshow'

const apiUrl = isDevMode()
  ? window.location.protocol + '//' + window.location.hostname + ':' + process.env.REACT_APP_API_PORT
  : window.location.origin
const wsUrl = isDevMode()
  ? 'ws://' + window.location.hostname + ':' + process.env.REACT_APP_API_PORT
  : 'ws://' + window.location.hostname

axios.defaults.timeout = 5000
axios.defaults.baseURL = apiUrl
axios.defaults.headers.post['Content-type'] = ''
axios.defaults.withCredentials = false

const buildQuery = (api, query = {}) => {
  if (query === {}) {
    return api
  }
  let all = api
  let sep = api.indexOf('?') >= 0 ? '&' : '?'
  for (const idx in query) {
    all += (sep + idx + '=' + query[idx])
    sep = '&'
  }
  return all
}

const buildHeader = (header = {}) => {
  if (!header) {
    header = {}
  }
  if (!('Authorization' in header) || !header['Authorization']) {
    const token = currentToken()
    header['Authorization'] = token ? ('Bearer ' + token) : ''
  }
  if (!('X-Csrf-Token' in header) || !header['X-Csrf-Token']) {
    header['X-Csrf-Token'] = getCookie('csrf_clap')
  }
  return header
}

const handleError = (err, reject) => {
  const has = err.response.data && err.response.data.message
  let msg = has ? err.response.data.message : 'Unknown error'
  if (err.response) {
    switch (err.response.status) {
      case 400:
        if (!has) {
          msg = 'Bad request'
        }
        break
      case 401:
        if (!has) {
          msg = 'Unauthorized'
        }
        delCurrentToken()
        window.location.href = '/login'
        break
      case 403:
        if (!has) {
          msg = 'No permission'
        }
        break
      case 404:
        if (!has) {
          msg = 'Not found'
        }
        break
      case 500:
        break
      case 502:
        break
      default:
    }
    ShowSnackbar(msg, 'error')
    reject(msg)
  } else {
    const msg = 'Error response'
    ShowSnackbar(msg, 'error')
    reject(msg)
  }
}

const http = {
  apiUrl: apiUrl,
  wsUrl: wsUrl,
  initRes: (func, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.get('/api/static', {headers: buildHeader(header)}).then(result => {
        setCurrentUserRes(result.data)
        resolve(func(result.data))
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  initProp: (header = {}) => {
    return new Promise((resolve, reject) => {
      axios.get('/config', {headers: buildHeader(header)}).then(result => {
        setCurrentBaseProp(result.data)
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  get: (api, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.get(buildQuery(api, query), {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  getOne: (pre, id, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.get(buildQuery(pre + '/' + id, query), {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  getSimple: (pre, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.get(buildQuery(pre + '/simple', query), {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  getList: (pre, refer = {}, page = {page: 1, size: 50, sort: {id: false}}, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      page.refer = refer
      axios.post(buildQuery(pre + '/list', query), page, {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  getMany: (pre, ids = [], query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.post(buildQuery(pre + '/list', query), {
        ids: [...new Set(ids)],
        size: -1
      }, {headers: buildHeader(header)}).then(result => {
        if (result.data === '') {
          reject('empty response')
        }
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  post: (api, data = {}, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.post(buildQuery(api, query), data, {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  put: (api, data = {}, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.put(buildQuery(api, query), data, {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  delete: (api, query = {}, header = {}) => {
    return new Promise((resolve, reject) => {
      axios.delete(buildQuery(api, query), {headers: buildHeader(header)}).then(result => {
        resolve(result.data)
      }).catch(err => {
        handleError(err, reject)
      })
    })
  },
  moreInfo: (apis = [{key: '', addr: '', field: ''}], results = [], header = {}) => {
    return new Promise((resolve, reject) => {
      const keyMap = new Map()
      if (results) {
        results.forEach((result) => {
          apis.forEach((api) => {
            if (result[api.key]) {
              let getVal = keyMap.get(api.key)
              const getKey = result[api.key]
              if (getKey) {
                if (getVal) {
                  getVal.push(result[api.key])
                } else {
                  getVal = [result[api.key]]
                }
                keyMap.set(api.key, getVal)
              }
            }
          })
        })
        Promise.all(apis.map((api) => http.getMany(api.addr, keyMap.get(api.key)))).then((listList) => {
          const valMap = new Map()
          listList.forEach((oneList, idx) => {
            if (oneList) {
              const api = apis[idx]
              oneList.forEach((one) => {
                valMap.set(api.key + one.id, one)
              })
            }
          })
          results.forEach((value) => {
            apis.forEach((api) => {
              value[api.field] = valMap.get(api.key + value[api.key])
            })
          })
          resolve(results)
        })
      } else {
        resolve(results)
      }
    })
  },
}

export default http
