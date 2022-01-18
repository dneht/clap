import React, {useEffect, useState} from 'react'
import {Box, Container, FormControl, makeStyles, TextField} from '@material-ui/core'
import Page from 'src/components/Page'
import http from 'src/requests'
import DocumentCard from 'src/views/document/DocumentView/DocumentCard'
import {currentBaseProp} from 'src/sessions'
import {ShowSnackbar} from 'src/utils/globalshow'
import {Autocomplete} from '@material-ui/lab'
import {useParams} from 'react-router'
import sha256 from 'crypto-js/sha256'

const useStyles = makeStyles((theme) => ({
  formControl: {
    margin: theme.spacing(1),
    minWidth: 480,
  },
}))

const MainView = () => {
  const classes = useStyles()
  const [docList, setDocList] = useState([])
  const [docMethod, setDocMethod] = useState([])
  const [docDetail, setDocDetail] = useState({commentInfo: {}, parameterInfo: [], returnInfo: {}})
  const [paramList, setParamList] = useState([])
  const [paramMockData, setParamMockData] = useState('')
  const [returnList, setReturnList] = useState([])
  const [returnMockData, setReturnMockData] = useState('')
  const [thisRequest, setThisRequest] = useState('')
  const [thisResponse, setThisResponse] = useState('')
  const params = useParams()
  const appId = params.app
  const docKey = params.key
  const spaceName = params.space
  const baseProps = currentBaseProp()
  const docProps = (baseProps && baseProps.document && docKey in baseProps.document) ? baseProps.document[docKey] : {}
  const docApiBase = docProps.api_base

  const doDocumentApiRequest = (reqUrl, param, tenant, group, token, func, err) => {
    if (appId && spaceName && docApiBase) {
      const time = new Date().getTime()
      http.post(docApiBase + reqUrl, {
        param: param,
        tenant: tenant,
        timestamp: time,
        token: token,
        sign: String(sha256(JSON.stringify(param) + tenant + time + token))
      }, {
        group: group
      }).then(data => {
        func(data)
      }).then(data => {
        err(data)
      })
    } else {
      ShowSnackbar('document api not found')
    }
  }

  const buildDocField = (inputList, resultList = [], addSize = 0) => {
    if (!inputList) {
      return resultList
    }
    for (const idx in inputList) {
      const inputOne = inputList[idx]
      if (inputOne && 'fieldName' in inputOne) {
        if (addSize > 0) {
          inputOne.fieldName = '\t'.repeat(addSize) + inputOne.fieldName
          inputOne.simpleTypeName = '\t'.repeat(addSize) + inputOne.simpleTypeName
        }
        resultList.push(inputOne)
        if ('fieldList' in inputOne) {
          resultList = buildDocField(inputOne.fieldList, resultList, addSize + 1)
        }
      }
    }
    return resultList
  }

  const initDocDetail = (data) => {
    setParamList(buildDocField(data.parameterInfo))
    if (data.paramMock) {
      setParamMockData(data.paramMock)
      setThisRequest(JSON.stringify(JSON.parse(data.paramMock), null, 4))
    } else {
      setThisRequest('')
    }

    setReturnList(buildDocField(data.returnInfo.fieldList))
    if (data.returnMock) {
      setReturnMockData(data.returnMock)
    } else {
      setReturnMockData('')
    }
  }

  const getDocDetail = (invokeName, invokeLength) => {
    if (appId && spaceName && docApiBase && docProps) {
      const time = String(new Date().getTime())
      http.get(docApiBase + docProps.api_method,
        {system: appId, group: spaceName, invokeName: invokeName, invokeLength: invokeLength},
        {
          'X-Clap-Time': String(new Date().getTime()),
          'X-Clap-Sign': String(sha256(`system=${appId}&group=${spaceName}&invokeName=${invokeName}&invokeLength=${invokeLength}${time}${docProps.token}`))
        }).then(data => {
        if (data) {
          setDocDetail(data)
          initDocDetail(data)
        } else {
          ShowSnackbar('document get error')
        }
      })
    } else {
      ShowSnackbar('document api not found')
    }
  }

  const docListOptions = {
    options: docList,
    getOptionLabel: (option) => option.simpleComment ? `${option.simpleComment}(${option.simpleName})` : option.simpleName,
  }
  const docListChange = (event, value, reason) => {
    if (value && value.methodInfo) {
      setDocMethod(value.methodInfo)
    }
  }
  const docMethodOptions = {
    options: docMethod,
    getOptionLabel: (option) => option.simpleComment ? `${option.simpleComment}(${option.simpleName})` : option.simpleName,
  }
  const docMethodChange = (event, value, reason) => {
    if (value && value.invokeName) {
      getDocDetail(value.invokeName, value.invokeLength)
    }
  }

  useEffect(() => {
    if (appId && spaceName && docApiBase) {
      const time = String(new Date().getTime())
      http.get(docApiBase + docProps.api_clazz,
        {system: appId, group: spaceName},
        {
          'X-Clap-Time': time,
          'X-Clap-Sign': String(sha256(`system=${appId}&group=${spaceName}${time}${docProps.token}`))
        }).then(data => {
        if (data) {
          setDocList(data)
        } else {
          ShowSnackbar('document get error')
        }
      })
    } else {
      ShowSnackbar('document api not found')
    }
  }, [])

  return (
    <Page
      className={classes.root}
      title="文档"
    >
      <Container maxWidth={false}>
        <FormControl className={classes.formControl}>
          <Autocomplete
            {...docListOptions}
            id="doc-list"
            clearOnEscape
            onChange={docListChange}
            renderInput={(params) => <TextField {...params} label="类目" margin="normal"/>}
          />
        </FormControl>
        <FormControl className={classes.formControl}>
          <Autocomplete
            {...docMethodOptions}
            id="doc-method"
            clearOnEscape
            onChange={docMethodChange}
            renderInput={(params) => <TextField {...params} label="子类目" margin="normal"/>}
          />
        </FormControl>
        <Box mt={3}>
          <DocumentCard dataProvider={docDetail} docProps={docProps}
                        doDocumentApiRequest={doDocumentApiRequest} spaceName={spaceName}
                        paramList={paramList} paramMockData={paramMockData}
                        returnList={returnList} returnMockData={returnMockData}
                        thisRequest={thisRequest} setThisRequest={setThisRequest}
                        thisResponse={thisResponse} setThisResponse={setThisResponse}/>
        </Box>
      </Container>
    </Page>
  )
}

export default MainView
