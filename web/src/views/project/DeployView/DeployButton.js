import React from 'react'
import {Button} from '@material-ui/core'
import {convertAppType} from 'src/utils/convertvalue'
import {ShowSnackbar} from 'src/utils/globalshow'
import {currentBaseProp, currentEnvName} from 'src/sessions'
import hiddenEle from 'src/utils/hiddenele'

class DeployButton extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      buttonStatus: props.dataProvider.deployStatus
    }
    this.dataProvider = props.dataProvider
    this.powerMap = props.powerMap
    this.deployId = props.dataProvider.id
    this.appId = props.dataProvider.appId
    this.spaceName = 'stable'
    if (props.dataProvider.spaceBase) {
      this.spaceName = props.dataProvider.spaceBase.spaceName
    }
    this.disablePack = props.dataProvider.isPackage === 0
    this.enableDocument = false
    const baseProps = currentBaseProp()
    if (baseProps && baseProps.document) {
      this.appDocKey = (currentEnvName() + '_' + convertAppType(props.dataProvider.appBase.appType)).toLowerCase()
      if (this.appDocKey in baseProps.document) {
        this.enableDocument = true
      }
    }
    this.appType = props.dataProvider.appBase.appType
    this.openPodDialog = props.openPodDialog
    this.openPropOpen = props.openPropOpen
    this.getBuildPods = props.getBuildPods
    this.gotoPackageApp = props.gotoPackageApp
    this.gotoPublishApp = props.gotoPublishApp
    this.navigateToDoc = props.navigateToDoc
    this.navigateToInner = props.navigateToInner
    this.gotoAppDocument = this.gotoAppDocument.bind(this)
    this.logPackageApp = this.logPackageApp.bind(this)
    this.reloadPackageStatus = this.reloadPackageStatus.bind(this)
    this.doPackageApp = this.doPackageApp.bind(this)
    this.doPublishApp = this.doPublishApp.bind(this)
  }

  gotoAppDocument() {
    this.navigateToDoc(this.dataProvider, this.appDocKey)
  }

  logPackageApp() {
    this.getBuildPods(this.deployId, (data) => {
      if (!data.pods || data.pods.length === 0) {
        ShowSnackbar('创建中，请稍后...', 'info')
      } else {
        this.navigateToInner(this.deployId, data.pods[0])
      }
    })
  }

  reloadPackageStatus() {
    this.getBuildPods(this.deployId, (data) => {
      if (!data.pods || data.pods.length === 0) {
        ShowSnackbar('创建中，请稍后...', 'info')
      } else {
        if (data.status.succeeded && data.status.succeeded > 0) {
          this.setState({
            buttonStatus: 2
          })
        } else {
          ShowSnackbar('点击左边按钮即可查看日志', 'info')
        }
      }
    })
  }

  doPackageApp() {
    this.gotoPackageApp(this.deployId, (data) => {
      this.setState({
        buttonStatus: 1
      })
      ShowSnackbar('打包中', 'error')
    })
  }

  doPublishApp() {
    this.gotoPublishApp(this.deployId, (data) => {
      this.setState({
        buttonStatus: 6
      })
      ShowSnackbar('发布中', 'error')
    })
  }

  render() {
    if (this.appType === 0) {
      return (
        <div>
          <Button variant="outlined" color="primary"
                  style={{display: hiddenEle(this.deployId, 'deployment', 'podLog', this.powerMap)}}
                  onClick={() => this.openPodDialog(this.dataProvider)}>
            查看
          </Button>
          <Button variant="outlined" color="primary"
                  style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                  onClick={this.doPublishApp}>
            发布
          </Button>
        </div>
      )
    } else {
      if (this.state.buttonStatus === 1) {
        return (
          <div>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(this.deployId, 'deployment', 'propView', this.powerMap)}}
                    onClick={() => this.openPropOpen(this.dataProvider)}>
              配置
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.logPackageApp}>
              打包中
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.reloadPackageStatus}>
              刷新状态
            </Button>
          </div>
        )
      } else if (this.state.buttonStatus === 2) {
        return (
          <div>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(this.deployId, 'deployment', 'propView', this.powerMap)}}
                    onClick={() => this.openPropOpen(this.dataProvider)}>
              配置
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(this.deployId, 'deployment', 'podLog', this.powerMap)}}
                    onClick={() => this.openPodDialog(this.dataProvider)}>
              查看
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.doPublishApp}>
              立即发布
            </Button>
          </div>
        )
      } else if (this.state.buttonStatus === 3) {
        return (
          <div>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(this.deployId, 'deployment', 'propView', this.powerMap)}}
                    onClick={() => this.openPropOpen(this.dataProvider)}>
              配置
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(this.deployId, 'deployment', 'podLog', this.powerMap)}}
                    onClick={() => this.openPodDialog(this.dataProvider)}>
              查看
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.logPackageApp}>
              打包失败
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.doPackageApp}>
              重新打包
            </Button>
          </div>
        )
      } else if (this.state.buttonStatus >= 6) {
        if (this.enableDocument) {
          return (
            <div>
              <Button variant="outlined" color="primary"
                      style={{display: hiddenEle(this.deployId, 'deployment', 'propView', this.powerMap)}}
                      onClick={() => this.openPropOpen(this.dataProvider)}>
                配置
              </Button>
              <Button variant="outlined" color="primary"
                      style={{display: hiddenEle(this.deployId, 'deployment', 'podLog', this.powerMap)}}
                      onClick={() => this.openPodDialog(this.dataProvider)}>
                查看
              </Button>
              <Button variant="outlined" color="primary"
                      style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                      onClick={this.doPackageApp}
                      disabled={this.disablePack}>
                打包
              </Button>
              <Button variant="outlined" color="primary"
                      style={{display: hiddenEle(this.deployId, 'deployment', 'docView', this.powerMap)}}
                      onClick={this.gotoAppDocument}>
                文档
              </Button>
            </div>
          )
        } else {
          return (
            <div>
              <Button variant="outlined" color="primary"
                      style={{display: hiddenEle(this.deployId, 'deployment', 'propView', this.powerMap)}}
                      onClick={() => this.openPropOpen(this.dataProvider)}>
                配置
              </Button>
              <Button variant="outlined" color="primary"
                      style={{display: hiddenEle(this.deployId, 'deployment', 'podLog', this.powerMap)}}
                      onClick={() => this.openPodDialog(this.dataProvider)}>
                查看
              </Button>
              <Button variant="outlined" color="primary"
                      style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                      onClick={this.doPackageApp}
                      disabled={this.disablePack}>
                打包
              </Button>
            </div>
          )
        }
      } else {
        return (
          <div>
            <Button variant="outlined" color="primary"
                    style={{display: hiddenEle(this.deployId, 'deployment', 'propView', this.powerMap)}}
                    onClick={() => this.openPropOpen(this.dataProvider)}>
              配置
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.doPackageApp}
                    disabled={this.disablePack}>
              打包
            </Button>
            <Button variant="outlined" color="primary"
                    style={{display: this.disablePack ? 'none' : hiddenEle(this.deployId, 'deployment', 'thisPack', this.powerMap)}}
                    onClick={this.doPublishApp}>
              发布
            </Button>
          </div>
        )
      }
    }
  }
}

export default DeployButton
