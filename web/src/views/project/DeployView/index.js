import React, {useEffect, useState} from 'react'
import {Box, Container, makeStyles} from '@material-ui/core'
import {saveAs} from 'file-saver'
import Page from 'src/components/Page'
import DeployList from './DeployList'
import http from 'src/requests'
import Toolbar from './Toolbar'
import {currentEnvId, currentSpaceId, setCurrentEnvId, setCurrentEnvName, setCurrentSpaceId} from 'src/sessions'
import {ShowSnackbar} from 'src/utils/globalshow'

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
  const [envSelect, setEnvSelect] = useState(0)
  const [envList, setEnvList] = useState([])
  const [spaceSelect, setSpaceSelect] = useState(0)
  const [spaceList, setSpaceList] = useState([])
  const [deployList, setDeployList] = useState({total: 0, size: 10, results: []})
  const [powerMap, setPowerMap] = useState({})

  const getDeployList = (spaceId) => {
    setSpaceSelect(spaceId)
    setCurrentSpaceId(spaceId)
    http.getList('/api/deploy', {spaceId: spaceId}).then(data => {
      if (data.results) {
        http.postSimple('/api/pow', {type: 'deployment'},
          data.results.map(e => e.id)).then(pow => {
          setPowerMap(pow)
          http.moreInfo([
            {key: 'appId', addr: '/api/app', field: 'appBase'},
            {key: 'spaceId', addr: '/api/space', field: 'spaceBase'}
          ], data.results).then(results => {
            data.results = results
            setDeployList(data)
          })
        }).catch(err => {
          ShowSnackbar('get pow err:' + err, 'error')
        })
      } else {
        setDeployList({total: 0, size: 10, results: []})
      }
    })
  }

  const getSpaceList = (envId) => {
    setEnvSelect(envId)
    setCurrentEnvId(envId)
    envList.forEach(one => {
      if (one.id === envId) {
        setCurrentEnvName(one.env)
      }
    })
    const getSpaceId = currentSpaceId()
    http.getSimple('/api/space', {eid: envId}).then(data => {
      if (data) {
        let spaceId = data[0].id
        if (getSpaceId) {
          data.forEach(one => {
            if (one.id === getSpaceId) {
              spaceId = one.id
            }
          })
        }
        setCurrentSpaceId(spaceId)
        setSpaceSelect(spaceId)
        setSpaceList(data)
        getDeployList(spaceId)
      }
    })
  }

  useEffect(() => {
    let getEnvId = currentEnvId()
    if (!getEnvId) {
      getEnvId = 1
    }
    http.getSimple('/api/env').then(data => {
      if (data) {
        let envId = data[0].id
        let envName = data[0].env
        data.forEach(one => {
          if (one.id === getEnvId) {
            envId = one.id
            envName = one.env
          }
        })
        setCurrentEnvId(envId)
        setCurrentEnvName(envName)
        setEnvSelect(envId)
        setEnvList(data)
        getSpaceList(envId)
      }
    })
  }, [])

  const getDeployPods = (deployId, func) => {
    http.get('/pod/deploy/' + deployId).then(data => func(data))
  }
  const getBuildPods = (deployId, func) => {
    http.get('/deploy/check/' + deployId).then(data => func(data))
  }
  const gotoPackageApp = (deployId, func) => {
    http.get('/deploy/build/' + deployId).then(data => func(data))
  }
  const gotoPublishApp = (deployId, func) => {
    http.get('/deploy/deploy/' + deployId).then(data => func(data))
  }
  const downloadPod = (data, podData) => {
    http.get('/pod/download/' + data.id,
      {
        'sid': data.spaceBase.id,
        'aid': data.appBase.id,
        'pod': podData.podName,
        'container': podData.containers[0].name,
        'file': '/logs/debug.log'
      }).then(data => {
      saveAs(new File([data], 'debug.log.tar', {type: 'application/tar'}))
    })
  }
  const restartPod = (data, podData) => {
    http.get('/pod/restart/' + data.id,
      {
        'sid': data.spaceBase.id,
        'aid': data.appBase.id,
        'pod': podData.podName
      }).then(data => {
      ShowSnackbar('重启成功', 'info')
      window.location.reload()
    })
  }

  return (
    <Page
      className={classes.root}
      title="部署列表"
    >
      <Container maxWidth={false}>
        <Toolbar envProvider={envList} envSelect={envSelect} getDeployList={getDeployList}
                 spaceProvider={spaceList} spaceSelect={spaceSelect} getSpaceList={getSpaceList}/>
        <Box mt={3}>
          <DeployList dataProvider={deployList} powerMap={powerMap}
                      getDeployPods={getDeployPods} getBuildPods={getBuildPods} downloadPod={downloadPod}
                      restartPod={restartPod}
                      gotoPackageApp={gotoPackageApp} gotoPublishApp={gotoPublishApp}/>
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
