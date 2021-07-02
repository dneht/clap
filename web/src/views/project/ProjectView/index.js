import React, {useEffect, useState} from 'react'
import {Box, Container, makeStyles} from '@material-ui/core'
import Page from 'src/components/Page'
import ProjectList from './ProjectList'
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
  const [appList, setAppList] = useState({total: 0, size: 10, results: []})
  useEffect(() => {
    http.getList('/api/app').then(data => setAppList(data))
  }, [])

  return (
    <Page
      className={classes.root}
      title="项目列表"
    >
      <Container maxWidth={false}>
        <Box mt={3}>
          <ProjectList dataProvider={appList}/>
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
