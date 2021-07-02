import {TextField} from '@material-ui/core';
import ReactJson from 'react-json-view';
import React from 'react';

const IfReactJson = ({mockData, ...rest}) => {
  if (mockData.indexOf('{') < 0 && mockData.indexOf('}') < 0) {
    return (
      <TextField fullWidth value={mockData}/>
    )
  } else {
    return (
      <ReactJson name={false} theme="summerfruit:inverted" displayDataTypes={false} src={JSON.parse(mockData)}/>
    )
  }
}

export default IfReactJson
