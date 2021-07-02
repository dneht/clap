import React from 'react'
import PropTypes from 'prop-types'
import clsx from 'clsx'
import {Box, Button, ButtonGroup, Card, CardContent, Divider, Grid, makeStyles, Typography} from '@material-ui/core'
import {useNavigate} from 'react-router-dom'
import {setCurrentEnvId, setCurrentEnvName} from 'src/sessions'

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

const EnvironmentCard = ({className, dataProvider, ...rest}) => {
  const classes = useStyles()
  const navigate = useNavigate()

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
            <ButtonGroup color="primary" aria-label="outlined primary button group">
              <Button>
                详情
              </Button>
              <Button>
                编辑
              </Button>
              <Button onClick={navigateToSpace}>
                空间
              </Button>
            </ButtonGroup>
          </Grid>
        </Grid>
      </Box>
    </Card>
  )
}

EnvironmentCard.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default EnvironmentCard
