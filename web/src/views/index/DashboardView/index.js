import React from 'react'
import {Box, Divider, makeStyles, Typography} from '@material-ui/core'
import Page from 'src/components/Page'

const useStyles = makeStyles((theme) => ({
  root: {
    backgroundColor: theme.palette.background.dark,
    minHeight: '100%',
    paddingBottom: theme.spacing(3),
    paddingTop: theme.spacing(3)
  }
}))

const Dashboard = () => {
  const classes = useStyles()

  return (
    <Page
      className={classes.root}
      title="概览"
    >
      <Box
        alignItems="center"
        display="flex"
        flexDirection="column"
      >
        <Typography variant="h1">
          欢迎
        </Typography>
        <Divider/>
        <br/>
        <Typography variant="body1">
          这里是发布平台
        </Typography>
      </Box>
    </Page>
  )
}

export default Dashboard
