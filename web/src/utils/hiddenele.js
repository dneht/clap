import md5 from 'crypto-js/md5'
import Hex from 'crypto-js/enc-hex'
import {currentToken} from 'src/sessions'

const hiddenEle = (id, get, key, powerMap) => {
  const token = currentToken()
  const adKey = Hex.stringify(md5(`0${get}1${token}`))
  if (powerMap[adKey] !== undefined && Object.keys(powerMap).length === 1) {
    return 'inline-block'
  }

  const cdKey = Hex.stringify(md5(`0${get}1${token}2${id}`))
  const powerKey = powerMap[cdKey]
  if (powerKey === undefined || Object.keys(powerKey).length === 0) {
    return 'none'
  }
  return powerKey[key] ? 'inline-block' : 'none'
}
export default hiddenEle
