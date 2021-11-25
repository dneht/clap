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
import navigateToTerm from 'src/utils/gototerm'
import {useNavigate} from 'react-router-dom'

const useStyles = makeStyles((theme) => ({
  root: {},
  statsItem: {
    alignItems: 'center',
    display: 'flex'
  }
}))

const CurrentCard = ({className, dataProvider, ...rest}) => {
  const classes = useStyles()
  const navigate = useNavigate()
  const [selectedListIds, setSelectedListIds] = useState([])
  const [limit, setLimit] = useState(50)
  const [page, setPage] = useState(0)

  const handleSelectAll = (event) => {
    let newSelectedListIds

    if (event.target.checked) {
      newSelectedListIds = dataProvider.map((data) => data.id)
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
                    checked={selectedListIds.length === dataProvider.length}
                    color="primary"
                    indeterminate={
                      selectedListIds.length > 0
                      && selectedListIds.length < dataProvider.length
                    }
                    onChange={handleSelectAll}
                  />
                </TableCell>
                <TableCell>
                  容器名
                </TableCell>
                <TableCell>
                  命名空间
                </TableCell>
                <TableCell>
                  所属项目
                </TableCell>
                <TableCell>
                  环境
                </TableCell>
                <TableCell>
                  空间
                </TableCell>
                <TableCell>
                  操作
                </TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {dataProvider.slice(0, limit).map((data) => (
                <TableRow
                  hover
                  key={data.podName}
                  selected={selectedListIds.indexOf(data.podName) !== -1}
                >
                  <TableCell padding="checkbox">
                    <Checkbox
                      checked={selectedListIds.indexOf(data.podName) !== -1}
                      onChange={(event) => handleSelectOne(event, data.podName)}
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
                        {data.podName}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    {data.namespace}
                  </TableCell>
                  <TableCell>
                    {data.appEnv}
                  </TableCell>
                  <TableCell>
                    {data.appGroup}
                  </TableCell>
                  <TableCell>
                    <Grid
                      className={classes.statsItem}
                      item
                    >
                      <ButtonGroup color="primary" aria-label="outlined primary button group">
                        <Button onClick={() => {
                          navigateToTerm(navigate, 'attach', data)
                        }}>
                          日志
                        </Button>
                        <Button onClick={() => {
                          navigateToTerm(navigate, 'exec', data)
                        }}>
                          命令
                        </Button>
                      </ButtonGroup>
                    </Grid>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </Box>
      </PerfectScrollbar>
      <TablePagination
        component="div"
        count={dataProvider.length}
        onChangePage={handlePageChange}
        onChangeRowsPerPage={handleLimitChange}
        page={page}
        rowsPerPage={limit}
        rowsPerPageOptions={[10, 25, 50]}
      />
    </Card>
  )
}

CurrentCard.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.array.isRequired
}

export default CurrentCard
