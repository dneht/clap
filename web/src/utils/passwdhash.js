import md5 from 'crypto-js/md5'
import Hex from 'crypto-js/enc-hex'

const passwdHash = (passwd) => {
  return Hex.stringify(md5(passwd))
}

export default passwdHash
