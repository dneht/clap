import React, {useState} from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import {Box, Button, Card, CardContent, Divider, Grid, makeStyles, Typography} from '@material-ui/core'
import {useNavigate} from 'react-router-dom'
import PropertyView from 'src/views/property/PropertyView'
import {setCurrentEnvId, setCurrentEnvName} from 'src/sessions'
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

const EnvironmentCard = ({className, dataProvider, powerMap, ...rest}) => {
  const classes = useStyles()
  const navigate = useNavigate()
  const [propOpen, setPropOpen] = useState(false)

  const navigateToSpace = () => {
    setCurrentEnvId(dataProvider.id)
    setCurrentEnvName(dataProvider.env)
    navigate('/app/spaces')
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
          {dataProvider.envName}
        </Typography>
        <Typography
          align="center"
          color="textPrimary"
          variant="body1"
        >
          {dataProvider.envDesc}
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
              {dataProvider.env}
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
                    style={{display: hiddenEle(dataProvider.id, 'environment', 'thisEdit', powerMap)}}>
              编辑
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, 'environment', 'propView', powerMap)}}
                    onClick={() => setPropOpen(true)}>
              配置
            </Button>
            <Button variant="outlined" color="primary"
                    onClick={navigateToSpace}>
              空间
            </Button>
          </Grid>
        </Grid>
        <PropertyView dataProvider={dataProvider} powerMap={powerMap} inputType="environment"
                      propOpen={propOpen} setPropOpen={setPropOpen}/>
      </Box>
    </Card>
  )
}

EnvironmentCard.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default EnvironmentCard
