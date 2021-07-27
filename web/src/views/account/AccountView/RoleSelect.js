import React from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import {Checkbox, FormControl, Input, InputLabel, ListItemText, makeStyles, MenuItem, Select} from '@material-ui/core'

const useStyles = makeStyles((theme) => ({
  root: {},
  formControl: {
    paddingBottom: theme.spacing(3),
    paddingTop: theme.spacing(3),
    minHeight: '100%',
    minWidth: 500,
    maxWidth: '100%',
  },
}))

const MenuProps = {
  PaperProps: {
    style: {
      maxHeight: 500,
      width: 250,
    }
  }
}

const RoleSelect = ({
                      className,
                      roleProvider,
                      selectRole,
                      setSelectRole,
                      ...rest
                    }) => {
  const classes = useStyles()

  const getSelectRole = (event) => {
    setSelectRole(event.target.value)
  }

  return (
    <div
      className={clsx(classes.root, className)}
      {...rest}
    >
      <FormControl className={classes.formControl}>
        <InputLabel id="env-select-label">选择角色</InputLabel>
        <Select
          labelId="role-select-label"
          id="role-select"
          multiple
          input={<Input/>}
          onChange={getSelectRole}
          renderValue={(selected) => selected.join(', ')}
          label="role"
          MenuProps={MenuProps}
          value={selectRole}
        >
          {roleProvider.map((data, idx) => (
            <MenuItem key={data.id} value={data.id}>
              <Checkbox checked={selectRole.indexOf(data.id) > -1}/>
              <ListItemText value={data.id} primary={data.roleName}/>
            </MenuItem>
          ))}
        </Select>
      </FormControl>
    </div>
  )
}

RoleSelect.propTypes = {
  className: PropTypes.string
}

export default RoleSelect
