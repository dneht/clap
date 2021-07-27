let setSnackbar
let setBackdrop

const Init = (snackbar, backdrop) => {
  setSnackbar = snackbar
  setBackdrop = backdrop
}

const ShowSnackbar = (msg, level, time = 5000) => {
  setSnackbar({data: msg, type: level, time: time})
}

const ShowBackdrop = () => {
  setBackdrop(true)
}

const CloseSnackbar = () => {
  setSnackbar({data: '', type: 'info', time: 5000})
}

const CloseBackdrop = () => {
  setBackdrop(false)
}

export {Init, ShowSnackbar, ShowBackdrop, CloseSnackbar, CloseBackdrop}
