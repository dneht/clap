import React, {useState} from 'react'
import clsx from 'clsx'
import PropTypes from 'prop-types'
import PerfectScrollbar from 'react-perfect-scrollbar'
import {
  Box,
  Button,
  ButtonGroup,
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
  Typography
} from '@material-ui/core'
import DeployButton from './DeployButton'
import {useNavigate} from 'react-router-dom'
import navigateToTerm from 'src/utils/gototerm'
import {ShowSnackbar} from 'src/utils/globalshow'
import {convertAppType} from 'src/utils/convertvalue'

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
                      getDeployPods,
                      getBuildPods,
                      gotoPackageApp,
                      gotoPublishApp,
                      restartPod,
                      ...rest
                    }) => {
  const classes = useStyles()
  const navigate = useNavigate()
  const [selectedListIds, setSelectedListIds] = useState([])
  const [limit, setLimit] = useState(50)
  const [page, setPage] = useState(0)
  const [podDialog, setPodDialog] = useState(false)
  const [podList, setPodList] = useState([])
  const [selectDeploy, setSelectDeploy] = useState(false)
  const [require, setRequire] = useState({open: false, keyword: ''})

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

  const navigateToInner = (podData) => {
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

  const requirePackageApp = (deployId, func) => {
    setRequire({open: true, keyword: '打包', type: 'package', id: deployId, func: func})
  }

  const requirePublishApp = (deployId, func) => {
    setRequire({open: true, keyword: '发布', type: 'publish', id: deployId, func: func})
  }

  const handleRequireClose = () => {
    setRequire({open: false, keyword: ''})
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
                      <DeployButton dataProvider={data} openPodDialog={handleDialogOpen} navigateToDoc={navigateToDoc}
                                    navigateToInner={navigateToInner}
                                    getBuildPods={getBuildPods} gotoPackageApp={requirePackageApp}
                                    gotoPublishApp={requirePublishApp}/>
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
                            {podList.slice(0, limit).map((pod) => (
                              <TableRow
                                hover
                                key={`pod-select-${data.id}-${pod.podName}`}
                              >
                                <TableCell>
                                  {pod.podName}
                                </TableCell>
                                <TableCell>
                                  {pod.status.phase}
                                </TableCell>
                                <TableCell>
                                  {pod.status.startTime}
                                </TableCell>
                                <TableCell>
                                  <Grid
                                    className={classes.statsItem}
                                    item
                                  >
                                    <ButtonGroup color="primary" aria-label="outlined primary button group">
                                      <Button variant="outlined" onClick={() => navigateToAttach(selectDeploy, pod)}>
                                        日志
                                      </Button>
                                      <Button variant="outlined" onClick={() => navigateToExec(selectDeploy, pod)}
                                              disabled={selectDeploy.appBase && selectDeploy.appBase.isIngress === 0}>
                                        命令
                                      </Button>
                                      <Button onClick={() => restartPod(selectDeploy, pod)}>
                                        重启
                                      </Button>
                                    </ButtonGroup>
                                  </Grid>
                                </TableCell>
                              </TableRow>
                            ))}
                          </TableBody>
                        </Table>
                      </Dialog>
                    </Grid>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <Dialog
            open={require.open}
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
              <Button onClick={handleRequireClose} color="primary">
                取消
              </Button>
              <Button onClick={() => {
                if (require.type && require.id && require.func) {
                  if (require.type === 'publish') {
                    gotoPublishApp(require.id, require.func)
                  } else {
                    gotoPackageApp(require.id, require.func)
                  }
                  handleRequireClose()
                } else {
                  ShowSnackbar('Get info error', 'error')
                }
              }} color="primary" autoFocus>
                确定
              </Button>
            </DialogActions>
          </Dialog>
        </Box>
      </PerfectScrollbar>
      <TablePagination
        component="div"
        count={dataProvider.total}
        onChangePage={handlePageChange}
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
