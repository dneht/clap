import {currentBaseProp} from 'src/sessions'

const convertAppType = (appType) => {
  const props = currentBaseProp()
  if (props && props.type) {
    const get = props.type[appType]
    if (get) {
      return get
    }
    return 'None'
  }
  return 'Unknown'
}

export {convertAppType}
