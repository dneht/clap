import React, {useState} from 'react'
import clsx from 'clsx'
import PropTypes from 'prop-types'
import PerfectScrollbar from 'react-perfect-scrollbar'
import {
  Box,
  Button,
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
import {convertAppType} from 'src/utils/convertvalue'
import hiddenEle from 'src/utils/hiddenele'
import PropertyView from 'src/views/property/PropertyView'

const useStyles = makeStyles((theme) => ({
  root: {},
  statsItem: {
    alignItems: 'center',
    display: 'flex'
  }
}))

const ProjectList = ({className, dataProvider, powerMap, ...rest}) => {
  const classes = useStyles()
  const [selectedListIds, setSelectedListIds] = useState([])
  const [limit, setLimit] = useState(50)
  const [page, setPage] = useState(0)
  const [propOpen, setPropOpen] = useState(false)
  const [projectData, setProjectData] = useState({})

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
                  项目id
                </TableCell>
                <TableCell>
                  项目名
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
                        {data.id}
                      </Typography>
                    </Box>
                  </TableCell>
                  <TableCell>
                    {data.appName}
                  </TableCell>
                  <TableCell>
                    {convertAppType(data.appType)}
                  </TableCell>
                  <TableCell>
                    <Grid
                      className={classes.statsItem}
                      item
                    >
                      <Button variant="outlined" color="primary">
                        详情
                      </Button>
                      <Button variant="outlined" color="primary"
                              style={{display: hiddenEle(data.id, 'project', 'thisEdit', powerMap)}}>
                        编辑
                      </Button>
                      <Button variant="outlined" color="primary"
                              style={{display: hiddenEle(data.id, 'project', 'propView', powerMap)}}
                              onClick={() => {
                                setPropOpen(true)
                                setProjectData(data)
                              }}>
                        配置
                      </Button>
                    </Grid>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
          <PropertyView dataProvider={projectData} powerMap={powerMap} inputType="project"
                        propOpen={propOpen} setPropOpen={setPropOpen}/>
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

ProjectList.propTypes = {
  className: PropTypes.string,
  dataProvider: PropTypes.object.isRequired
}

export default ProjectList
