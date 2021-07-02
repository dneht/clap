import React from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import {FormControl, InputLabel, makeStyles, MenuItem, Select} from '@material-ui/core'

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}))

const Toolbar = ({
                   className,
                   envProvider,
                   envSelect,
                   getSpaceList,
                   spaceProvider,
                   spaceSelect,
                   getDeployList,
                   ...rest
                 }) => {
  const classes = useStyles()

  const getSpaceAndDeployList = (event) => {
    getSpaceList(event.target.value)
  }

  const getOnlyDeployList = (event) => {
    getDeployList(event.target.value)
  }

  return (
    <div
      className={clsx(classes.root, className)}
      {...rest}
    >
      <FormControl className={classes.formControl}>
        <InputLabel id="env-select-label">环境</InputLabel>
        <Select
          labelId="env-select-label"
          id="env-select"
          onChange={getSpaceAndDeployList}
          label="env"
          value={envSelect}
        >
          {envProvider.map((data, idx) => (
            <MenuItem key={`env-${data.id}`} value={data.id}>{data.env}</MenuItem>
          ))}
        </Select>
      </FormControl>
      <FormControl className={classes.formControl}>
        <InputLabel id="space-select-label">空间</InputLabel>
        <Select
          labelId="space-select-label"
          id="space-select"
          onChange={getOnlyDeployList}
          label="space"
          value={spaceSelect}
        >
          {spaceProvider.map((data, idx) => (
            <MenuItem key={`space-${data.id}`} value={data.id}>{data.spaceName}</MenuItem>
          ))}
        </Select>
      </FormControl>
    </div>
  )
}

Toolbar.propTypes = {
  className: PropTypes.string
}

export default Toolbar
