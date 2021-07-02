import React, {useEffect, useState} from 'react'
import {Box, Container, makeStyles} from '@material-ui/core'
import Page from 'src/components/Page'
import CurrentCard from './CurrentCard'
import http from 'src/requests'
import {currentSpaceId} from 'src/sessions'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingBottom: theme.spacing(3),
    paddingTop: theme.spacing(3)
  }
}))

const ProjectListView = () => {
  const classes = useStyles()
  const [appList, setAppList] = useState([])
  useEffect(() => {
    http.get('/pod/space/' + currentSpaceId()).then(data => setAppList(data))
  }, [])

  return (
    <Page
      className={classes.root}
      title="项目列表"
    >
      <Container maxWidth={false}>
        <Box mt={3}>
          <CurrentCard dataProvider={appList}/>
        </Box>
      </Container>
    </Page>
  )
}

export default ProjectListView
