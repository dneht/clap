import React, {useEffect, useState} from 'react'
import {Outlet} from 'react-router-dom'
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

const SimpleLayout = () => {
  const classes = useStyles()
  const [routeData, setRouteData] = useState([])

  useEffect(async () => {
    const routeData = await routeItems()
    setRouteData(routeData)
  })

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

export default SimpleLayout
