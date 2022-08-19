import React, {useState} from 'react'
import clsx from 'clsx'
import PropTypes from 'prop-types'
import PerfectScrollbar from 'react-perfect-scrollbar'
import {
  Box,
  Button,
  Card,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  Grid,
  makeStyles,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
  Typography
} from '@material-ui/core'
import DeployButton from './DeployButton'
import PropertyView from 'src/views/property/PropertyView'
import {useNavigate} from 'react-router-dom'
import navigateToTerm from 'src/utils/gototerm'
import {ShowSnackbar} from 'src/utils/globalshow'
import {convertAppType} from 'src/utils/convertvalue'
import hiddenEle from 'src/utils/hiddenele'
import timeToLocal from 'src/utils/timetolocal'

const useStyles = makeStyles((theme) => ({
  root: {},
  statsItem: {
    alignItems: 'center',
    display: 'flex'
  }
}))

const DeployList = ({
                      className,
                      dataProvider,
                      powerMap,
                      getDeployProps,
                      getDeployPods,
                      getBuildPods,
                      getDeploySnaps,
                      downloadPod,
                      restartPod,
                      gotoPackageApp,
                      gotoAutoDeployApp,
                      gotoPublishApp,
                      gotoCancelApp,
                      gotoRollbackApp,
                      ...rest
                    }) => {
  const classes = useStyles()
  const navigate = useNavigate()
  const [selectedListIds, setSelectedListIds] = useState([])
  const [limit, setLimit] = useState(50)
  const [page, setPage] = useState(0)
  const [podDialog, setPodDialog] = useState(false)
  const [podList, setPodList] = useState([])
  const [selectDeploy, setSelectDeploy] = useState({})
  const [branchName, setBranchName] = useState('')
  const [snapList, setSnapList] = useState([])
  const [propOpen, setPropOpen] = useState(false)
  const [require, setRequire] = useState({open: false, more: false, keyword: '', type: ''})
  const [rollback, setRollback] = useState({open: false, id: 0})

  const navigateToAttach = (data, podData) => {
    podData.envId = data.spaceBase.envId
    podData.spaceId = data.spaceBase.id
    podData.deployId = data.id
    navigateToTerm(navigate, 'attach', podData)
  }

  const navigateToExec = (data, podData) => {
    podData.envId = data.spaceBase.envId
    podData.spaceId = data.spaceBase.id
    podData.deployId = data.id
    navigateToTerm(navigate, 'exec', podData)
  }

  const navigateToInner = (id, podData) => {
    podData.deployId = id
    navigateToTerm(navigate, 'inner', podData)
  }

  const navigateToDoc = (data, docKey) => {
    navigate(`/app/document/${data.appBase.id}/${docKey}/${data.spaceBase.spaceName}`)
  }

  const dataResults = dataProvider.results
  const handleSelectAll = (event) => {
    let newSelectedListIds

    if (event.target.checked) {
      newSelectedListIds = dataResults.map((data) => data.id)
    } else {
      newSelectedListIds = []
    }

    setSelectedListIds(newSelectedListIds)
  }

  const handleSelectOne = (event, id) => {
    const selectedIndex = selectedListIds.indexOf(id)
    let newSelectedListIds = []

    if (selectedIndex === -1) {
      newSelectedListIds = newSelectedListIds.concat(selectedListIds, id)
    } else if (selectedIndex === 0) {
      newSelectedListIds = newSelectedListIds.concat(selectedListIds.slice(1))
    } else if (selectedIndex === selectedListIds.length - 1) {
      newSelectedListIds = newSelectedListIds.concat(selectedListIds.slice(0, -1))
    } else if (selectedIndex > 0) {
      newSelectedListIds = newSelectedListIds.concat(
        selectedListIds.slice(0, selectedIndex),
        selectedListIds.slice(selectedIndex + 1)
      )
    }

    setSelectedListIds(newSelectedListIds)
  }

  const handleLimitChange = (event) => {
    setLimit(event.target.value)
  }

  const handlePageChange = (event, newPage) => {
    setPage(newPage)
  }

  const handleDialogOpen = (deploy) => {
    getDeployPods(deploy.id, setPodList)
    setSelectDeploy(deploy)
    setPodDialog(true)
  }

  const handleDialogClose = () => {
    setPodDialog(false)
  }

  const handlePropOpen = (deploy) => {
    setSelectDeploy(deploy)
    setPropOpen(true)
  }

  const handleRollbackOpen = (id, deploy, list) => {
    setSelectDeploy(deploy)
    setSnapList(list)
    setRollback({open: true, id: id})
  }

  const handleRollbackClose = () => {
    setRollback({open: false, id: 0})
  }

  const requirePackageApp = (deployId, canBranch, func) => {
    setRequire({open: true, more: canBranch, keyword: '打包', type: 'Package', id: deployId, func: func})
  }

  const requirePublishApp = (deployId, func) => {
    setRequire({open: true, more: false, keyword: '发布', type: 'Publish', id: deployId, func: func})
  }

  const cancelPackageApp = (deployId, func) => {
    setRequire({open: true, more: false, keyword: '取消', type: 'Cancel', id: deployId, func: func})
  }

  const handleRequireClose = () => {
    setRequire({open: false, more: false, keyword: '', type: ''})
  }

  const updateBranchName = (event) => {
    setBranchName(event.target.value)
    dataResults.forEach(dataOne => {
      if (dataOne.id === require.id) {
        dataOne.branchName = event.target.value
      }
    })
  }

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <PerfectScrollbar>
        <Box minWidth={1050}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell padding="checkbox">
                  <Checkbox
                    checked={selectedListIds.length === dataResults.length}
                    color="primary"
                    indeterminate={
                      selectedListIds.length > 0
                      && selectedListIds.length < dataResults.length
                    }
                    onChange={handleSelectAll}
                  />
                </TableCell>
                <TableCell>
                  发布名
                </TableCell>
                <TableCell>
                  项目id
                </TableCell>
                <TableCell>
                  项目名
                </TableCell>
                <TableCell>
                  所在空间
                </TableCell>
                <TableCell>
                  发布分支
                </TableCell>
                <TableCell>
                  类型
                </TableCell>
                <TableCell>
                  操作
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {dataResults.slice(0, limit).map((data) => (
                <TableRow
                  hover
                  key={data.id}
                  selected={selectedListIds.indexOf(data.id) !== -1}
                >
                  <TableCell padding="checkbox">
                    <Checkbox
                      checked={selectedListIds.indexOf(data.id) !== -1}
                      onChange={(event) => handleSelectOne(event, data.id)}
                      value="true"
                    />
                  </TableCell>
                  <TableCell>
                    <Box
                      alignItems="center"
                      display="flex"
                    >
                      <Typography
                        color="textPrimary"
                        variant="body1"
                      >
                        {data.deployName}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    {data.appBase.id}
                  </TableCell>
                  <TableCell>
                    {data.appBase.appName}
                  </TableCell>
                  <TableCell>
                    {data.spaceBase.spaceName}
                  </TableCell>
                  <TableCell>
                    {data.branchName}
                  </TableCell>
                  <TableCell>
                    {convertAppType(data.appBase.appType)}
                  </TableCell>
                  <TableCell>
                    <Grid
                      className={classes.statsItem}
                      item
                    >
                      <DeployButton dataProvider={data} powerMap={powerMap}
                                    openDeployRollback={handleRollbackOpen}
                                    openPodDialog={handleDialogOpen}
                                    openPropOpen={handlePropOpen}
                                    navigateToDoc={navigateToDoc}
                                    navigateToInner={navigateToInner}
                                    getBuildPods={getBuildPods}
                                    getDeploySnaps={getDeploySnaps}
                                    gotoPackageApp={requirePackageApp}
                                    gotoPublishApp={requirePublishApp}
                                    gotoCancelApp={cancelPackageApp}/>
                    </Grid>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <Dialog maxWidth="lg" open={podDialog} onClose={handleDialogClose}
                  aria-labelledby="pod-dialog-title">
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                    名称
                  </TableCell>
                  <TableCell>
                    状态
                  </TableCell>
                  <TableCell>
                    创建时间
                  </TableCell>
                  <TableCell>
                    操作
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {podList.map((pod) => (
                  <TableRow
                    hover
                    key={`pod-select-${pod.podName}`}
                  >
                    <TableCell>
                      {pod.podName}
                    </TableCell>
                    <TableCell>
                      {pod.status.phase}
                    </TableCell>
                    <TableCell>
                      {timeToLocal(pod.status.startTime)}
                    </TableCell>
                    <TableCell>
                      <Grid
                        className={classes.statsItem}
                        item
                      >
                        <Button variant="outlined" color="primary"
                                style={{display: hiddenEle(selectDeploy.id, 'deployment', 'podLog', powerMap)}}
                                onClick={() => navigateToAttach(selectDeploy, pod)}>
                          日志
                        </Button>
                        <Button variant="outlined" color="primary"
                                style={{
                                  display: selectDeploy.appBase && selectDeploy.appBase.isIngress === 0 ? 'none'
                                    : hiddenEle(selectDeploy.id, 'deployment', 'podExec', powerMap)
                                }}
                                onClick={() => navigateToExec(selectDeploy, pod)}
                                disabled={selectDeploy.appBase && selectDeploy.appBase.isIngress === 0}>
                          命令
                        </Button>
                        <Button variant="outlined" color="primary"
                                style={{display: hiddenEle(selectDeploy.id, 'deployment', 'podDownload', powerMap)}}
                                onClick={() => downloadPod(selectDeploy, pod)}>
                          下载
                        </Button>
                        <Button variant="outlined" color="primary"
                                style={{display: hiddenEle(selectDeploy.id, 'deployment', 'podRestart', powerMap)}}
                                onClick={() => restartPod(selectDeploy, pod)}>
                          重启
                        </Button>
                      </Grid>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Dialog>
          <Dialog
            open={require.open && require.more}
            onClose={handleRequireClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
          >
            <DialogContent>
              <DialogContentText id="alert-dialog-description">
                {`确定继续${require.keyword}?`}
              </DialogContentText>
              <Typography id="alert-dialog-readme" variant="overline" mt={2}>
                仓库分支留空则使用默认分支，否则将会使用并把输入分支作为默认分支
              </Typography>
              <TextField id="alert-dialog-branch-name" variant="outlined"
                         label="仓库分支名" defaultValue=""
                         onChange={updateBranchName}/>
            </DialogContent>
            <DialogActions>
              <Button variant="outlined" color="primary"
                      onClick={handleRequireClose}>
                取消
              </Button>
              <Button variant="outlined" color="primary" autoFocus
                      onClick={() => {
                        if (require.type && require.id && require.func) {
                          gotoPackageApp(require.id, branchName, require.func)
                        } else {
                          ShowSnackbar('Select type not found', 'warn')
                        }
                        handleRequireClose()
                      }}>
                确定
              </Button>
              <Button variant="outlined" color="primary" autoFocus
                      onClick={() => {
                        if (require.type && require.id && require.func) {
                          gotoAutoDeployApp(require.id, branchName, require.func)
                        } else {
                          ShowSnackbar('Select type not found', 'warn')
                        }
                        handleRequireClose()
                      }}>
                自动发布
              </Button>
            </DialogActions>
          </Dialog>
          <Dialog maxWidth="lg" open={rollback.open} onClose={handleRollbackClose}
                  aria-labelledby="rollback-dialog-title">
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>
                    创建者
                  </TableCell>
                  <TableCell>
                    发布时间
                  </TableCell>
                  <TableCell>
                    分支名
                  </TableCell>
                  <TableCell>
                    发布版本
                  </TableCell>
                  <TableCell>
                    操作
                  </TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {snapList.map((snap) => (
                  <TableRow
                    hover
                    key={`rollback-select-${snap.id}`}
                  >
                    <TableCell>
                      {snap.userInfo && snap.userInfo.nickname ? snap.userInfo.nickname : '已禁用'}
                    </TableCell>
                    <TableCell>
                      {timeToLocal(snap.createdAt)}
                    </TableCell>
                    <TableCell>
                      {snap.branchName}
                    </TableCell>
                    <TableCell>
                      {snap.deployTag}
                    </TableCell>
                    <TableCell>
                      <Grid
                        className={classes.statsItem}
                        item
                      >
                        <Button variant="outlined" color="primary"
                                style={{display: hiddenEle(snap.deployId, 'deployment', 'thisRollback', powerMap)}}
                                disabled={snap.deployTag === selectDeploy.deployTag}
                                onClick={() => setRequire({
                                  open: true,
                                  more: false,
                                  keyword: '回滚',
                                  type: 'Rollback',
                                  id: snap.id,
                                  func: (data) => {
                                    selectDeploy.deployTag = snap.deployTag
                                    handleRollbackClose()
                                    ShowSnackbar('回滚中', 'error')
                                  }
                                })}>
                          回滚
                        </Button>
                      </Grid>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </Dialog>
          <Dialog
            open={require.open && !require.more}
            onClose={handleRequireClose}
            aria-labelledby="alert-dialog-title"
            aria-describedby="alert-dialog-description"
          >
            <DialogContent>
              <DialogContentText id="alert-dialog-description">
                {`确定继续${require.keyword}?`}
              </DialogContentText>
            </DialogContent>
            <DialogActions>
              <Button variant="outlined" color="primary"
                      onClick={handleRequireClose}>
                取消
              </Button>
              <Button variant="outlined" color="primary" autoFocus
                      onClick={() => {
                        if (require.type && require.id && require.func) {
                          switch (require.type) {
                            case 'Publish': {
                              gotoPublishApp(require.id, require.func)
                              break
                            }
                            case 'Package': {
                              gotoPackageApp(require.id, branchName, require.func)
                              break
                            }
                            case 'Cancel': {
                              gotoCancelApp(require.id, require.func)
                              break
                            }
                            case 'Rollback': {
                              gotoRollbackApp(require.id, require.func)
                              break
                            }
                            default: {
                              ShowSnackbar('Select type not found', 'warn')
                            }
                          }
                          handleRequireClose()
                        } else {
                          ShowSnackbar('Get info error', 'error')
                        }
                      }}>
                确定
              </Button>
            </DialogActions>
          </Dialog>
          <PropertyView dataProvider={selectDeploy} powerMap={powerMap} inputType="deployment"
                        propOpen={propOpen} setPropOpen={setPropOpen}/>
        </Box>
      </PerfectScrollbar>
      <TablePagination
        component="div"
        count={dataProvider.total}
        onChangeRowsPerPage={handleLimitChange}
        page={page}
        rowsPerPage={limit}
        rowsPerPageOptions={[10, 25, 50]}
      />
    </Card>
  )
}

DeployList.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default DeployList
