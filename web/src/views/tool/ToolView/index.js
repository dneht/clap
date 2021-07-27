import React from 'react'
import {Box, Button, ButtonGroup, Container, makeStyles} from '@material-ui/core'
import Page from 'src/components/Page'
import http from 'src/requests'
import {ShowSnackbar} from 'src/utils/globalshow'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingBottom: theme.spacing(3),
    paddingTop: theme.spacing(3)
  },
  dataProviderCard: {
    height: '100%'
  }
}))

const MainView = () => {
  const classes = useStyles()

  const cleanAll = () => {
    http.get('/api/clean').then(data => {
      if (data) {
        ShowSnackbar('清理成功', 'info')
      } else {
        ShowSnackbar('清理失败', 'warn')
      }
    })
  }

  return (
    <Page
      className={classes.root}
      title="工具箱"
    >
      <Container maxWidth={false}>
        <Box mt={3}>
          <ButtonGroup color="primary" aria-label="outlined primary button group">
            <Button color="primary" size="large" variant="contained" onClick={cleanAll}>清理缓存</Button>
          </ButtonGroup>
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
