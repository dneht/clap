import React, {useState} from 'react'
import {makeStyles} from '@material-ui/core'
import SimpleBar from './SimpleBar'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.default,
    display: 'flex',
    height: '100%',
    overflow: 'hidden',
    width: '100%'
  }
}))

const SimpleLayout = () => {
  const classes = useStyles()
  const [isMobileNavOpen, setMobileNavOpen] = useState(false)

  return (
    <div className={classes.root}>
      <SimpleBar
        onMobileClose={() => setMobileNavOpen(false)}
        openMobile={isMobileNavOpen}
      />
    </div>
  )
}

export default SimpleLayout
