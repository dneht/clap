import React, {useEffect} from 'react'
import {useNavigate} from 'react-router-dom'
import Avatar from '@material-ui/core/Avatar'
import Button from '@material-ui/core/Button'
import CssBaseline from '@material-ui/core/CssBaseline'
import TextField from '@material-ui/core/TextField'
import Grid from '@material-ui/core/Grid'
import LockOutlinedIcon from '@material-ui/icons/LockOutlined'
import {makeStyles} from '@material-ui/core/styles'
import Container from '@material-ui/core/Container'
import http from 'src/requests'
import passwdHash from 'src/utils/passwdhash'
import {ShowSnackbar} from 'src/utils/globalshow'
import {currentToken, setCurrentToken} from 'src/sessions'

const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main,
  },
  form: {
    width: '100%',
    marginTop: theme.spacing(3),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}))

const LoginView = () => {
  const classes = useStyles()
  const [name, setName] = React.useState('')
  const [passwd, setPasswd] = React.useState('')
  const navigate = useNavigate()

  useEffect(() => {
    if (currentToken()) {
      navigate('/app/dashboard')
    }
  }, [])

  const handleLogin = () => {
    if (name && passwd) {
      http.post('/login', {
        userName: name,
        password: passwdHash(passwd),
      }).then(token => {
        setCurrentToken(token)
        navigate('/app/dashboard')
      }).catch(err => {
        ShowSnackbar(err, 'error')
      })
    } else {
      ShowSnackbar('请输入必选项', 'error')
    }
  }

  const handleChange = (event, setup) => {
    setup(event.target.value)
  }

  return (
    <Container component="main" maxWidth="xs">
      <CssBaseline/>
      <div className={classes.paper}>
        <Avatar className={classes.avatar}>
          <LockOutlinedIcon/>
        </Avatar>
        <form>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                id="username"
                name="用户名"
                label="用户名"
                autoComplete="current-username"
                onChange={event => handleChange(event, setName)}
              />
            </Grid>
            <Grid item xs={12}>
              <TextField
                variant="outlined"
                required
                fullWidth
                id="password"
                name="密码"
                label="密码"
                type="password"
                autoComplete="current-password"
                onChange={event => handleChange(event, setPasswd)}
              />
            </Grid>
          </Grid>
        </form>
        <Button
          type="submit"
          fullWidth
          variant="contained"
          color="primary"
          className={classes.submit}
          onClick={handleLogin}
        >
          登录
        </Button>
      </div>
    </Container>
  )
}

export default LoginView
