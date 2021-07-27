import React, {useEffect, useState} from 'react'
import {
  Box,
  Button,
  ButtonGroup,
  Container,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  makeStyles,
  TextField
} from '@material-ui/core'
import http from 'src/requests'
import Page from 'src/components/Page'
import {ShowSnackbar} from 'src/utils/globalshow'
import AccountList from './AccountList'
import RoleSelect from './RoleSelect'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingBottom: theme.spacing(3),
    paddingTop: theme.spacing(3)
  }
}))

const MainView = () => {
  const classes = useStyles()
  const [open, setOpen] = React.useState(false)
  const [userList, setUserList] = useState({total: 0, size: 50, results: []})
  const [name, setName] = React.useState('')
  const [nickname, setNickname] = React.useState('')
  const [newPasswd, setNewPasswd] = React.useState('')
  const [roleList, setRoleList] = React.useState([])
  const [selectRole, setSelectRole] = useState([])

  useEffect(() => {
    http.getList('/api/user').then(data => setUserList(data))
    http.getSimple('/api/role').then(data => setRoleList(data))
  }, [])

  const handleOpen = () => {
    setOpen(true)
  }

  const handleClose = () => {
    setOpen(false)
  }

  const handleCreate = () => {
    if (name && nickname) {
      http.post('/api/user', {
        userName: name,
        nickname: nickname,
        password: newPasswd,
        roleList: JSON.stringify(selectRole)
      }).then(res => {
        http.getList('/api/user').then(data => setUserList(data))
        ShowSnackbar('用户创建成功', 'info')
        handleClose()
      }).catch(err => {
        ShowSnackbar('用户创建失败: ' + err, 'error')
      })
    } else {
      ShowSnackbar('请输入必选项', 'error')
    }
  }

  const updateUser = (id, info) => {
    return http.put('/api/user/' + id, info).then(res => {
      ShowSnackbar('用户更新成功', 'info')
    }).catch(err => {
      ShowSnackbar('用户更新失败: ' + err, 'error')
    })
  }

  const handleChange = (event, setup) => {
    setup(event.target.value)
  }

  return (
    <Page
      className={classes.root}
      title="用户列表"
    >
      <Container maxWidth={false}>
        <ButtonGroup color="primary" aria-label="outlined primary button group">
          <Button variant="outlined" onClick={handleOpen}>
            新增用户
          </Button>
        </ButtonGroup>
        <Dialog open={open} onClose={handleClose} aria-labelledby="form-dialog-title">
          <DialogTitle id="form-dialog-title">新增用户</DialogTitle>
          <DialogContent>
            <TextField
              autoFocus
              required
              margin="dense"
              id="name"
              label="用户名（唯一）"
              onChange={event => handleChange(event, setName)}
              fullWidth
            />
            <TextField
              autoFocus
              required
              margin="dense"
              id="nickname"
              label="昵称"
              onChange={event => handleChange(event, setNickname)}
              fullWidth
            />
            <TextField
              autoFocus
              required
              margin="dense"
              id="password"
              label="密码"
              type="password"
              autoComplete="new-password"
              onChange={event => handleChange(event, setNewPasswd)}
              fullWidth
            />
            <RoleSelect roleProvider={roleList} selectRole={selectRole} setSelectRole={setSelectRole}/>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="primary">
              取消
            </Button>
            <Button onClick={handleCreate} color="primary">
              确定
            </Button>
          </DialogActions>
        </Dialog>

        <Box mt={3}>
          <AccountList dataProvider={userList} roleList={roleList} updateUser={updateUser}/>
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
