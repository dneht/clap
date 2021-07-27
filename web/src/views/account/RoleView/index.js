import React, {useEffect, useState} from 'react'
import {Box, Button, ButtonGroup, Container, makeStyles} from '@material-ui/core'
import Page from 'src/components/Page'
import RoleList from './RoleList'
import http from 'src/requests'

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
  const [roleList, setRoleList] = useState({total: 0, size: 50, results: []})
  useEffect(() => {
    http.getList('/api/role').then(data => setRoleList(data))
  }, [])

  return (
    <Page
      className={classes.root}
      title="角色列表"
    >
      <Container maxWidth={false}>
        <Box mt={3}>
          <RoleList dataProvider={roleList}/>
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
