import React, {useEffect, useState} from 'react'
import {Outlet, useLocation} from 'react-router-dom'
import PropTypes from 'prop-types'
import {Box, List, makeStyles} from '@material-ui/core'
import NavItem from './NavItem'
import routeItems from 'src/utils/routeitems'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.default,
    display: 'flex',
    height: '100%',
    overflow: 'hidden',
    width: '100%'
  },
  wrapper: {
    display: 'flex',
    flex: '1 1 auto',
    overflow: 'hidden',
    paddingTop: 0
  },
  contentContainer: {
    display: 'flex',
    flex: '1 1 auto',
    overflow: 'hidden'
  },
  content: {
    flex: '1 1 auto',
    height: '100%',
    overflow: 'auto'
  }
}))

const SimpleBar = ({onMobileClose, openMobile}) => {
  const classes = useStyles()
  const location = useLocation()
  const [routeData, setRouteData] = useState([])

  useEffect(async () => {
    const routeData = await routeItems()
    setRouteData(routeData)

    if (openMobile && onMobileClose) {
      onMobileClose()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname])

  return (
    <div className={classes.root}>
      <Box p={1}>
        <List>
          {routeData.map(item => (
            <NavItem
              href={item.href}
              key={item.title}
              title={item.title}
              icon={item.icon}
            />
          ))}
        </List>
      </Box>
      <div className={classes.wrapper}>
        <div className={classes.contentContainer}>
          <div className={classes.content}>
            <Outlet/>
          </div>
        </div>
      </div>
    </div>
  )
}

SimpleBar.propTypes = {
  onMobileClose: PropTypes.func,
  openMobile: PropTypes.bool
}

SimpleBar.defaultProps = {
  onMobileClose: () => {
  },
  openMobile: false
}

export default SimpleBar
