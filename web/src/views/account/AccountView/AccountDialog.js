import React, {useState} from 'react'
import {Button, Dialog, DialogActions, DialogContent, DialogTitle, TextField} from '@material-ui/core'
import PropTypes from 'prop-types'
import passwdHash from 'src/utils/passwdhash'
import RoleSelect from './RoleSelect'

const AccountDialog = ({
                         className,
                         roleList,
                         passwdOpen,
                         setPasswdOpen,
                         rolesOpen,
                         setRolesOpen,
                         selectRole,
                         setSelectRole,
                         updateUser,
                         ...rest
                       }) => {
  const [passwd, setPasswd] = useState('')

  const handlePasswdClose = () => {
    setPasswdOpen('')
  }

  const handlePasswdChange = (event, setup) => {
    setup(event.target.value)
  }

  const handleUpdatePasswd = (data) => {
    updateUser(data, {password: passwdHash(passwd)}).then(res => {
      handlePasswdClose()
    })
  }

  const handleRolesClose = () => {
    setRolesOpen({})
    setSelectRole([])
  }

  const handleUpdateRoles = (id) => {
    updateUser(id, {roleList: JSON.stringify(selectRole)}).then(res => {
      handleRolesClose()
    })
  }

  return (
    <div>
      <Dialog open={passwdOpen !== ''} onClose={handlePasswdClose} aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">重置密码</DialogTitle>
        <form>
          <DialogContent>
            <TextField
              autoFocus
              required
              margin="dense"
              id="password-update"
              label="密码"
              type="password"
              autoComplete="new-password"
              onChange={event => handlePasswdChange(event, setPasswd)}
              fullWidth
            />
          </DialogContent>
        </form>
        <DialogActions>
          <Button onClick={handlePasswdClose} color="primary">
            取消
          </Button>
          <Button onClick={() => handleUpdatePasswd(passwdOpen)} color="primary">
            确定
          </Button>
        </DialogActions>
      </Dialog>
      <Dialog open={Object.keys(rolesOpen).length !== 0} onClose={handleRolesClose} aria-labelledby="form-dialog-title">
        <DialogTitle id="form-dialog-title">调整角色</DialogTitle>
        <DialogContent>
          <RoleSelect roleProvider={roleList} selectRole={selectRole} setSelectRole={setSelectRole}/>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleRolesClose} color="primary">
            取消
          </Button>
          <Button onClick={() => handleUpdateRoles(rolesOpen.id)} color="primary">
            确定
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  )
}

AccountDialog.propTypes = {
  className: PropTypes.string
}

export default AccountDialog
