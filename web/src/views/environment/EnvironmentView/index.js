import React, {useEffect, useState} from 'react'
import {Box, Container, Grid, makeStyles} from '@material-ui/core'
import {Pagination} from '@material-ui/lab'
import Page from 'src/components/Page'
import EnvironmentCard from './EnvironmentCard'
import http from 'src/requests'

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
  const [dataList, setDataList] = useState({results: []})
  useEffect(() => {
    http.getList('/api/env').then(data => setDataList(data))
  }, [])

  return (
    <Page
      className={classes.root}
      title="环境列表"
    >
      <Container maxWidth={false}>
        <Box mt={3}>
          <Grid
            container
            spacing={3}
          >
            {dataList.results.map((data) => (
              <Grid
                item
                key={data.id}
                lg={4}
                md={6}
                xs={12}
              >
                <EnvironmentCard
                  className={classes.dataProviderCard}
                  dataProvider={data}
                />
              </Grid>
            ))}
          </Grid>
        </Box>
        <Box
          mt={3}
          display="flex"
          justifyContent="center"
        >
          <Pagination
            color="primary"
            count={1}
            size="small"
          />
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
