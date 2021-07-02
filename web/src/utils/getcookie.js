const getcookie = (name) => {
  const cookies = document.cookie
  const idx = cookies.indexOf(name + '=')
  if (idx < 0) {
    return ''
  }
  const last = cookies.indexOf(';', idx)
  return cookies.substring(idx + name.length + 1, last)
}

export default getcookie
