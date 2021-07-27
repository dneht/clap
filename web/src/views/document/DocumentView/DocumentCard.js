import React, {useState} from 'react'
import clsx from 'clsx'
import PropTypes from 'prop-types'
import PerfectScrollbar from 'react-perfect-scrollbar'
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Button,
  Card,
  Collapse,
  FormControl,
  makeStyles,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Toolbar,
  Typography
} from '@material-ui/core'
import DocumentTable from './DocumentTable'
import Paper from '@material-ui/core/Paper'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import IfReactJson from './IfReactJson'
import {ShowSnackbar} from 'src/utils/globalshow'

const useStyles = makeStyles((theme) => ({
  root: {},
  formControl: {
    margin: theme.spacing(1),
    minWidth: 40,
  },
  formToken: {
    margin: theme.spacing(1),
    minWidth: 600,
  },
  tableCell: {
    whiteSpace: 'pre-wrap',
    wordWrap: 'break-word',
    minWidth: 160,
  },
}))

const DocumentCard = ({
                        className,
                        dataProvider,
                        doDocumentApiRequest,
                        spaceName,
                        paramList,
                        paramMockData,
                        returnList,
                        returnMockData,
                        thisRequest,
                        setThisRequest,
                        thisResponse,
                        setThisResponse,
                        ...rest
                      }) => {
  const classes = useStyles()
  const [tenant, setTenant] = useState('board')
  const [group, setGroup] = useState(spaceName)
  const [token, setToken] = useState('')
  const [startRequest, setStartRequest] = useState(false)
  const [openRequest, setOpenRequest] = useState(false)
  let returnComment = ''
  if (dataProvider.returnInfo && dataProvider.returnInfo.simpleTypeName) {
    returnComment += '返回参数说明: ' + dataProvider.returnInfo.simpleTypeName
    if (dataProvider.returnInfo.simpleComment) {
      returnComment += ' | ' + dataProvider.returnInfo.simpleComment
    }
    if (dataProvider.returnInfo.extendComment) {
      returnComment += ' | ' + dataProvider.returnInfo.extendComment
    }
  }

  const doStartRequest = () => {
    setStartRequest(true)
    setOpenRequest(true)
    try {
      doDocumentApiRequest(dataProvider.pathName, JSON.parse(thisRequest), tenant, group, token,
        (data) => {
          setStartRequest(false)
          setThisResponse(JSON.stringify(data))
        }, (err) => {
          setStartRequest(false)
        })
    } catch (ex) {
      setStartRequest(false)
      ShowSnackbar(ex, 'warn')
    }
  }

  return (
    <Card
      className={clsx(classes.root, className)}
      {...rest}
    >
      <PerfectScrollbar>
        <Box minWidth={1050}>
          <TableContainer component={Paper} className={clsx(classes.root, className)}
                          {...rest}>
            <Toolbar>
              <Typography variant="h2">说明</Typography>
            </Toolbar>
            <Table size="small" key="main">
              <TableHead>
                <TableRow>
                  <TableCell>作者</TableCell>
                  <TableCell>请求路径</TableCell>
                  <TableCell>组合路径</TableCell>
                  <TableCell>白名单</TableCell>
                  <TableCell>更新时间</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                <TableRow>
                  <TableCell>
                    <Typography variant="h5" className={classes.tableCell}>
                      {dataProvider.commentInfo ? dataProvider.commentInfo.author : ''}
                    </Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="h5" className={classes.tableCell}>{dataProvider.pathName}</Typography>
                  </TableCell>
                  <TableCell>
                    <Typography variant="h5" className={classes.tableCell}>{dataProvider.invokeName}</Typography>
                  </TableCell>
                  <TableCell>{dataProvider.isWhitelist ? '是' : '否'}</TableCell>
                  <TableCell>{dataProvider.updatedAt ? new Date(dataProvider.updatedAt).toLocaleString() : ''}</TableCell>
                </TableRow>
              </TableBody>
            </Table>
            <Accordion>
              <AccordionSummary expandIcon={<KeyboardArrowDownIcon/>}>
                <Typography variant="h3" className={classes.heading}>测试接口</Typography>
              </AccordionSummary>
              <AccordionDetails>
                <div>
                  <Toolbar>
                    <FormControl className={classes.formControl}>
                      <TextField id="tenant-basic" label="Tenant" onChange={(event) => {
                        setTenant(event.target.value)
                      }} value={tenant}/>
                    </FormControl>
                    <FormControl className={classes.formControl}>
                      <TextField id="group-basic" label="Group" onChange={(event) => {
                        setGroup(event.target.value)
                      }} value={group}/>
                    </FormControl>
                    <FormControl className={classes.formToken}>
                      <TextField id="token-basic" label="Token" onChange={(event) => {
                        setToken(event.target.value)
                      }} value={token}/>
                    </FormControl>
                    <Button color="primary" size="large" variant="contained" disabled={startRequest}
                            onClick={doStartRequest}>开始请求</Button>
                  </Toolbar>
                  <TextField multiline fullWidth variant="outlined" onChange={(event) => {
                    setThisRequest(event.target.value)
                  }} value={thisRequest}/>
                  <Collapse in={openRequest} timeout="auto" unmountOnExit>
                    <Toolbar>
                      <Typography variant="h4" className={classes.heading}>返回结果</Typography>
                    </Toolbar>
                    <IfReactJson mockData={thisResponse}/>
                  </Collapse>
                </div>
              </AccordionDetails>
            </Accordion>
          </TableContainer>
        </Box>
        <Box minWidth={1050} pt={4}>
          <DocumentTable title="请求参数" comment=""
                         paramList={paramList} paramMockData={paramMockData}/>
        </Box>
        <Box minWidth={1050} pt={4}>
          <DocumentTable title="返回参数" comment={returnComment}
                         paramList={returnList} paramMockData={returnMockData}/>
        </Box>
      </PerfectScrollbar>
    </Card>
  )
}

DocumentCard.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default DocumentCard
