import React, {useEffect, useState} from 'react'
import Button from '@material-ui/core/Button'
import {
  AppBar,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  Grid,
  makeStyles,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
  Tabs,
  TextField,
  Typography
} from '@material-ui/core'
import http from 'src/requests'
import hiddenEle from 'src/utils/hiddenele'
import {ShowSnackbar} from 'src/utils/globalshow'
import timeToLocal from 'src/utils/timetolocal'

const useStyles = makeStyles((theme) => ({
  statsItem: {
    alignItems: 'center',
    display: 'flex'
  }
}))

const PropertyView = ({
                        className,
                        dataProvider,
                        powerMap,
                        inputType,
                        propOpen,
                        setPropOpen,
                        ...rest
                      }) => {
  const classes = useStyles()
  const [select, setSelect] = useState(0)
  const [readme, setReadme] = useState('')
  const [content, setContent] = useState('')
  const [propList, setPropList] = useState([])
  const [snapList, setSnapList] = useState([])
  const [require, setRequire] = useState({open: false, keyword: '', type: ''})
  const [rollback, setRollback] = useState({open: false, id: 0})

  const handlePropClose = () => {
    setPropOpen(false)
    setReadme('')
    setContent('')
  }

  const handleRequireClose = () => {
    setRequire({open: false, keyword: '', type: ''})
  }

  const handleRollbackClose = () => {
    setRollback({open: false, id: 0, deploy: 0})
  }

  const handleTableChange = (event, select) => {
    setReadme(propList[select].readme)
    setContent(propList[select].content)
    setSelect(select)
  }

  const getDeployProps = () => {
    http.get('/prop/' + inputType + '/' + dataProvider.id).then(data => {
      setPropList(data)
    }).catch(err => {
      ShowSnackbar('获取配置文件失败: ' + err, 'error')
    })
  }

  const getDeployPropSnaps = (id) => {
    http.get('/snaps/config/' + id).then(data => {
      http.moreInfo([
        {key: 'userId', addr: '/api/user', field: 'userInfo'},
      ], data).then(results => {
        setSnapList(results)
      })
    }).catch(err => {
      ShowSnackbar('获取配置历史失败: ' + err, 'error')
    })
  }

  const handleRollbackOpen = (id, deploy) => {
    setRollback({open: true, id: id, deploy: deploy})
    getDeployPropSnaps(id)
  }

  const handlePropUpdate = () => {
    switch (require.type) {
      case 'Update': {
        const current = propList[select]
        http.post('/prop/' + inputType + '/' + current.id, {
          fileContent: content,
        }).then(res => {
          ShowSnackbar('更新配置文件成功', 'info')
        }).catch(err => {
          ShowSnackbar('更新配置文件失败: ' + err, 'error')
        }).finally(() => {
          handleRequireClose()
          handlePropClose()
        })
        break
      }
      case 'Rollback': {
        http.post('/snaps/config/' + require.id).then(res => {
          ShowSnackbar('回滚配置文件成功', 'info')
        }).catch(err => {
          ShowSnackbar('回滚配置文件失败: ' + err, 'error')
        }).finally(() => {
          handleRequireClose()
          handleRollbackClose()
          handlePropClose()
        })
        break
      }
      default: {
        ShowSnackbar('操作类型未找到', 'info')
      }
    }
  }

  useEffect(() => {
    if (propOpen) {
      getDeployProps()
    }
  }, [propOpen])

  useEffect(() => {
    if (propList.length > 0) {
      setReadme(propList[0].readme)
      setContent(propList[0].content)
    }
  }, [propList])

  if (propList.length === 0) {
    return (
      <div>
        <Dialog fullWidth maxWidth="xs"
                open={propOpen} onClose={handlePropClose}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
        >
          <DialogContent>
            <DialogContentText id="alert-dialog-description">
              暂无配置
            </DialogContentText>
          </DialogContent>
        </Dialog>
      </div>
    )
  } else {
    return (
      <div>
        <Dialog fullWidth maxWidth="md"
                open={propOpen} onClose={handlePropClose}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
        >
          <AppBar position="static">
            <Tabs value={select} onChange={handleTableChange}>
              {propList.map((prop) => (
                <Tab key={'prop-select-' + prop.id} label={prop.name}/>
              ))}
            </Tabs>
          </AppBar>
          <DialogContent>
            <Typography variant="body1">{readme}</Typography>
            <TextField multiline fullWidth
                       autoFocus required
                       margin="dense"
                       id="prop-update"
                       inputProps={{step: 300}}
                       onChange={event => setContent(event.target.value)}
                       variant="outlined" size="medium" rowsMax={25} rows={25}
                       value={content}/>
          </DialogContent>
          <DialogActions>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, inputType, 'propEdit', powerMap)}}
                    onClick={() => setRequire({open: true, keyword: '更新', type: 'Update'})}>
              更新
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, inputType, 'propEdit', powerMap)}}
                    onClick={() => handleRollbackOpen(propList[select].id, dataProvider.id)}>
              回滚
            </Button>
            <Button variant="outlined" color="primary"
                    onClick={handlePropClose}>
              取消
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
                  更新时间
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
                    <Grid
                      className={classes.statsItem}
                      item
                    >
                      <Button variant="outlined" color="primary"
                              style={{display: hiddenEle(rollback.deploy, inputType, 'propEdit', powerMap)}}
                              onClick={() => setRequire({open: true, keyword: '回滚', type: 'Rollback', id: snap.id})}>
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
            <Button variant="outlined" color="primary"
                    onClick={handleRequireClose}>
              取消
            </Button>
            <Button variant="outlined" color="primary" autoFocus
                    onClick={handlePropUpdate}>
              确定
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    )
  }
}

export default PropertyView
