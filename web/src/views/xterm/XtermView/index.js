import React from 'react'
import {makeStyles} from '@material-ui/core'
import {Terminal} from 'xterm'
import {FitAddon} from 'xterm-addon-fit'
import {SearchAddon} from 'xterm-addon-search'
import {WebLinksAddon} from 'xterm-addon-web-links'
import 'xterm/css/xterm.css'
import Requests from 'src/requests'
import {ShowSnackbar} from 'src/utils/globalshow'
import {Helmet} from 'react-helmet'
import {currentToken} from 'src/sessions'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
  },
  fullscreen: {
    height: '100%',
    minHeight: '100%',
  }
}))

const MainView = () => {
  const classes = useStyles()

  const term = new Terminal({
    cursorBlink: true
  })
  const fitAddon = new FitAddon()
  term.loadAddon(fitAddon)
  term.loadAddon(new SearchAddon())
  term.loadAddon(new WebLinksAddon())
  // term.open(document.getElementById('terminal'))
  const params = JSON.parse(atob(window.location.search.substring(4)))
  const resource = params[6]
  const connect = new WebSocket(Requests.wsUrl + '/select/' + resource + '/' + params[4]
    + '?env=' + params[0] + '&space=' + params[1] + '&project=' + params[2] + '&deploy=' + params[3]
    + '&container=' + params[5] + '&token=' + encodeURIComponent(currentToken()))

  const resize = () => {
    fitAddon.fit()
  }

  connect.onopen = function (event) {
    term.open(document.getElementById('terminal'))
    term.onResize(size => {
      connect.send(
        JSON.stringify({Op: 'resize', Rows: size.rows, Cols: size.cols})
      )
    })
    term.onData(msg => {
      connect.send(JSON.stringify({Op: 'stdin', Data: msg}))
    })
    resize()
    window.addEventListener('resize', resize)
  }

  connect.onclose = function (event) {
    term.writeln('\r\n\x1b[1;33mconnection close\x1B[0m')
    term.setOption('cursorBlink', false)
    window.removeEventListener('resize', resize)
  }

  connect.onmessage = function (event) {
    const message = JSON.parse(event.data)
    if (message.Op === 'stdout') {
      if (resource === 'attach') {
        term.write(message.Data.replace(/\r?\n/g, '\r\n'))
      } else {
        term.write(message.Data)
      }
    } else if (message.Op === 'toast') {
      ShowSnackbar(message.Data, message.ToastType, 30 * 1000)
    }
  }

  return (
    <div id="terminal" className={classes.fullscreen}>
      <Helmet>
        <title>终端</title>
      </Helmet>
    </div>
  )
}

export default MainView
