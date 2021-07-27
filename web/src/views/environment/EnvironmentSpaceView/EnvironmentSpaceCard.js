import React from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import {Box, Button, Card, CardContent, Divider, Grid, makeStyles, Typography} from '@material-ui/core'
import {currentEnvName, setCurrentSpaceId} from 'src/sessions'
import {useNavigate} from 'react-router-dom'
import hiddenEle from 'src/utils/hiddenele'

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
    flexDirection: 'column'
  },
  statsItem: {
    alignItems: 'center',
    display: 'flex'
  },
  statsIcon: {
    marginRight: theme.spacing(1)
  }
}))

const EnvironmentSpaceCard = ({className, dataProvider, powerMap, ...rest}) => {
  const classes = useStyles()
  const navigate = useNavigate()

  const navigateToCurrent = () => {
    setCurrentSpaceId(dataProvider.id)
    navigate('/app/current')
  }

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <CardContent>
        <Box
          display="flex"
          justifyContent="center"
          mb={3}
        />
        <Typography
          align="center"
          color="textPrimary"
          gutterBottom
          variant="h5"
        >
          {dataProvider.spaceName}
        </Typography>
        <Typography
          align="center"
          color="textPrimary"
          variant="body1"
        >
          {dataProvider.spaceDesc}
        </Typography>
      </CardContent>
      <Box flexGrow={1}/>
      <Divider/>
      <Box p={2}>
        <Grid
          container
          justify="space-between"
          spacing={2}
        >
          <Grid
            className={classes.statsItem}
            item
          >
            <Typography
              color="textSecondary"
              display="inline"
              variant="body2"
            >
              {`${currentEnvName()}:${dataProvider.spaceKeep}`}
            </Typography>
          </Grid>
          <Grid
            className={classes.statsItem}
            item
          >
            <Button variant="outlined" color="primary">
              详情
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, 'environment_space', 'editThis', powerMap)}}>
              编辑
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, 'environment_space', 'propView', powerMap)}}>
              配置
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, 'environment_space', 'podSpace', powerMap)}}
                    onClick={navigateToCurrent}>
              查看
            </Button>
          </Grid>
        </Grid>
      </Box>
    </Card>
  )
}

EnvironmentSpaceCard.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default EnvironmentSpaceCard
