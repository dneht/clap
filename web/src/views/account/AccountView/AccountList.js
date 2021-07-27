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
import AccountDialog from './AccountDialog'

const useStyles = makeStyles((theme) => ({
  root: {},
  statsItem: {
    alignItems: 'center',
    display: 'flex'
  }
}))

const AccountList = ({className, dataProvider, roleList, updateUser, ...rest}) => {
  const classes = useStyles()
  const [selectedListIds, setSelectedListIds] = useState([])
  const [limit, setLimit] = useState(50)
  const [page, setPage] = useState(0)
  const [passwdOpen, setPasswdOpen] = React.useState('')
  const [rolesOpen, setRolesOpen] = React.useState({})
  const [selectRole, setSelectRole] = useState([])

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

  const handlePasswdOpen = (id) => {
    setPasswdOpen(id)
  }

  const handleRolesOpen = (data) => {
    setRolesOpen(data)
    if (data.roleList) {
      setSelectRole(JSON.parse(data.roleList))
    }
  }

  const handleDisableUser = (id) => {
    updateUser(id, {isDisable: 1})
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
                  用户id
                </TableCell>
                <TableCell>
                  用户名
                </TableCell>
                <TableCell>
                  用户昵称
                </TableCell>
                <TableCell>
                  用户头像
                </TableCell>
                <TableCell>
                  角色管理
                </TableCell>
                <TableCell>
                  是否被禁用
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
                        {data.id}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    {data.userName}
                  </TableCell>
                  <TableCell>
                    {data.nickname}
                  </TableCell>
                  <TableCell>
                    {data.avatar}
                  </TableCell>
                  <TableCell>
                    {data.roleList}
                  </TableCell>
                  <TableCell>
                    {data.isDisable === 0 ? '否' : '是'}
                  </TableCell>
                  <TableCell>
                    <Grid
                      className={classes.statsItem}
                      item
                    >
                      <ButtonGroup color="primary" aria-label="outlined primary button group" data-id={data.id}>
                        <Button onClick={() => handlePasswdOpen(data.id)}>
                          密码
                        </Button>
                        <Button onClick={() => handleRolesOpen(data)}>
                          角色
                        </Button>
                        <Button onClick={() => handleDisableUser(data.id)}
                                disabled={data.isDisable !== 0}>
                          禁用
                        </Button>
                      </ButtonGroup>
                    </Grid>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <AccountDialog roleList={roleList}
                         passwdOpen={passwdOpen} setPasswdOpen={setPasswdOpen}
                         rolesOpen={rolesOpen} setRolesOpen={setRolesOpen}
                         selectRole={selectRole} setSelectRole={setSelectRole}
                         updateUser={updateUser}/>
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

AccountList.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default AccountList
