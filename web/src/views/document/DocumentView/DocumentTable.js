import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  makeStyles,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Toolbar,
  Typography
} from '@material-ui/core'
import clsx from 'clsx'
import Paper from '@material-ui/core/Paper'
import KeyboardArrowDownIcon from '@material-ui/icons/KeyboardArrowDown'
import PropTypes from 'prop-types'
import React from 'react'
import IfReactJson from './IfReactJson'

const useStyles = makeStyles((theme) => ({
  root: {},
  formControl: {
    margin: theme.spacing(1),
    minWidth: 30,
  },
  heading: {
    margin: theme.spacing(2),
    minWidth: 300,
  },
  tableCell: {
    whiteSpace: 'pre-wrap',
    wordWrap: 'break-word',
    minWidth: 160,
  }
}))

const DocumentTable = ({className, title, comment, paramList, paramMockData, ...rest}) => {
  const classes = useStyles()

  return (
    <TableContainer component={Paper} className={clsx(classes.root, className)}
                    {...rest}>
      <Toolbar>
        <Typography variant="h2">{title}</Typography>
      </Toolbar>
      <Typography variant="h5" className={classes.heading}>{comment}</Typography>
      <Table size="small" key="main">
        <TableHead>
          <TableRow>
            <TableCell>参数名</TableCell>
            <TableCell>类型</TableCell>
            <TableCell>说明</TableCell>
            <TableCell>详情</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {paramList.map((data, idx) => (
            <TableRow key={'param-' + data.fieldName.trim() + idx}>
              <TableCell align="left">
                <Typography className={classes.tableCell}>{data.fieldName}</Typography>
              </TableCell>
              <TableCell align="left">
                <Typography className={classes.tableCell}>{data.simpleTypeName}</Typography>
              </TableCell>
              <TableCell align="left">{data.simpleComment ? data.simpleComment : data.simpleName}</TableCell>
              <TableCell align="left">{data.extendComment ? data.extendComment : ''}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Accordion>
        <AccordionSummary expandIcon={<KeyboardArrowDownIcon/>}>
          <Typography variant="h4" className={classes.heading}>显示示例</Typography>
        </AccordionSummary>
        <AccordionDetails>
          <IfReactJson mockData={paramMockData}/>
        </AccordionDetails>
      </Accordion>
    </TableContainer>
  )
}

DocumentTable.propTypes = {
  className: PropTypes.string,
}

export default DocumentTable
