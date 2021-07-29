import React, {useEffect, useState} from 'react'
import Button from '@material-ui/core/Button'
import {
  AppBar,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  Tab,
  Tabs,
  TextField,
  Typography
} from '@material-ui/core'
import http from 'src/requests'
import hiddenEle from 'src/utils/hiddenele'
import {ShowSnackbar} from 'src/utils/globalshow'

const PropertyView = ({
                        className,
                        dataProvider,
                        powerMap,
                        inputType,
                        propOpen,
                        setPropOpen,
                        ...rest
                      }) => {
  const [select, setSelect] = useState(0)
  const [readme, setReadme] = useState('')
  const [content, setContent] = useState('')
  const [propList, setPropList] = useState([])
  const [require, setRequire] = useState({open: false, keyword: ''})

  const handlePropClose = () => {
    setPropOpen(false)
    setReadme('')
    setContent('')
  }

  const handleRequireClose = () => {
    setRequire({open: false, keyword: ''})
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

  const handlePropUpdate = () => {
    const current = propList[select]
    http.put('/prop/' + inputType + '/' + current.id, {
      fileReadme: readme,
      fileContent: content,
    }).then(res => {
      ShowSnackbar('更新配置文件成功', 'info')
    }).catch(err => {
      ShowSnackbar('更新配置文件失败: ' + err, 'error')
    }).finally(() => {
      handleRequireClose()
      handlePropClose()
    })
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
                    onClick={() => setRequire({open: true, keyword: '更新', type: 'update'})}>
              更新
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(dataProvider.id, inputType, 'thisRollback', powerMap)}}
                    onClick={handlePropClose}>
              回滚
            </Button>
            <Button variant="outlined" color="primary"
                    onClick={handlePropClose}>
              取消
            </Button>
          </DialogActions>
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
