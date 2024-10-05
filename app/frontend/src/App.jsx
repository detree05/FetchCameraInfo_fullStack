import React, { useState } from 'react';
import axios from 'axios';
import { Flex, List, Button } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import { InputForm } from './components/InputForm';
import CreateQuery from './components/CreateQuery';

const boxStyle = {
  width: '100%',
  height: '100%',
};

function App() {
  const [data, setData] = useState([])
  const [reqStatus, setReqStatus] = useState()
  const InitRequest = () => {
    axios.post('http://172.16.0.57:8888/getChannel', CreateQuery()).then(res => {
    // axios.post('http://localhost:8888/getChannel', CreateQuery()).then(res => {
      var resData = res.data;
      const resDataArray = [];
      for (var i = 0; i < resData.length; i++) {
        resDataArray[i] = `
        @name: ${resData[i]["Name"]}@ |
        @ext_id: ${resData[i]["ExtId"]}@ |
        @channel_id: ${resData[i]["ChannelId"]}@ |
        @control_host: ${resData[i]["ControlHost"]}@ |
        @channel_config: ${resData[i]["ChannelConfigLink"] ? resData[i]["ChannelConfigLink"] : "None"}@ |
        @channel_status: ${resData[i]["ChannelStatusLink"] ? resData[i]["ChannelStatusLink"] : "None"}@
        `
      }
      setData(resDataArray)
    }).catch(error => {
      setReqStatus(error.status)
    });
  }

  return (
    <div>
      <Flex
      vertical="vertical"
      gap="middle"
      style={boxStyle}
      justify="center"
      align="center"
      >
        <Button
        icon={<SearchOutlined />}
        iconPosition="end"
        onClick={() => InitRequest()}
        >Поиск</Button>
        {reqStatus ? `Error ${reqStatus}` : ``}
        <InputForm/>
        {data.length > 0 ? `Найдено ${data.length} совпадений` : ``}
        <List
        size="small"
        bordered
        dataSource={data}
        renderItem={(item) => <List.Item>{item}</List.Item>}
        />
      </Flex>
    </div>
  )
}

export { App }
