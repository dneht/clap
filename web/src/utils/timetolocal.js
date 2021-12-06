const timeToLocal = (input) => {
  // YYYY-MM-DD HH:MM:SS
  if (input) {
    const utc = new Date(input)
    const year = utc.getFullYear()
    const month = utc.getMonth() + 1 < 10 ? `0${utc.getMonth() + 1}` : (utc.getMonth() + 1)
    const day = utc.getDate() < 10 ? `0${utc.getDate()}` : utc.getDate()
    const hour = utc.getHours() < 10 ? `0${utc.getHours()}` : utc.getHours()
    const minute = utc.getMinutes() < 10 ? `0${utc.getMinutes()}` : utc.getMinutes()
    const second = utc.getSeconds() < 10 ? `0${utc.getSeconds()}` : utc.getSeconds()
    return `${year}-${month}-${day} ${hour}:${minute}:${second}`
  }
  return ''
}
export default timeToLocal
