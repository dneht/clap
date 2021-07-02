import 'react-perfect-scrollbar/dist/css/styles.css'
import React, {useState} from 'react'
import {useRoutes} from 'react-router-dom'
import {Backdrop, CircularProgress, Snackbar, ThemeProvider} from '@material-ui/core'
import {Alert} from '@material-ui/lab'
import GlobalStyles from 'src/components/GlobalStyles'
import theme from 'src/theme'
import {routes} from 'src/routes'
import {makeStyles} from '@material-ui/core/styles'
import {CloseBackdrop, CloseSnackbar, Init} from 'src/utils/globalshow'
import http from 'src/requests'
import {currentBaseProp} from 'src/sessions'

const useStyles = makeStyles((theme) => ({
  backdrop: {
    zIndex: theme.zIndex.drawer + 1,
    color: '#fff',
  },
}))

const App = () => {
  const classes = useStyles()
  const routing = useRoutes(routes)
  const [error, setError] = useState(false)
  const [snackbar, setSnackbar] = useState({data: '', type: 'info', time: 5000})
  const [backdrop, setBackdrop] = useState(false)
  Init(setSnackbar, setBackdrop)

  if (!error && !currentBaseProp()) {
    http.initBase().catch((err) => {
      setError(true)
    })
  }

  return (
    <ThemeProvider theme={theme}>
      <Snackbar open={snackbar.data !== ''} autoHideDuration={snackbar.time} onClose={CloseSnackbar}>
        <Alert onClose={CloseSnackbar} severity={snackbar.type}>
          {snackbar.data}
        </Alert>
      </Snackbar>

      <Backdrop className={classes.backdrop} open={backdrop} onClick={CloseBackdrop}>
        <CircularProgress color="inherit"/>
      </Backdrop>

      <GlobalStyles/>
      {routing}
    </ThemeProvider>
  )
}

export default App
