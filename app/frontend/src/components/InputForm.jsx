import { useState } from 'react';
import { Flex, Input, Select, Checkbox, Space } from 'antd';
import { setVerbose, setProjectName, setCameraName, setUsername, setPassword } from './DataHandler';

function InputForm() {
  const projects = [
    {
      value: 'bryansk',
      label: 'Брянск',
    },
    {
      value: 'codd',
      label: 'ЦОДД',
    },
    {
      value: 'gmc',
      label: 'ГМЦ',
    },
    {
      value: 'komi',
      label: 'Коми',
    },
    {
      value: 'mo',
      label: 'МО',
    },
    {
      value: 'murmansk',
      label: 'Мурманск',
    },
    {
      value: 'novosibirsk',
      label: 'Новосибирск',
    },
    {
      value: 'perm',
      label: 'Пермь',
    },
    {
      value: 'primorye',
      label: 'Цифровое Приморье',
    },
    {
      value: 'rostov',
      label: 'Ростов',
    },
    {
      value: 'ryazan',
      label: 'Рязань',
    },
    {
      value: 'sirius',
      label: 'Сириус',
    },
    {
      value: 'tyumen',
      label: 'Тюмень',
    },
    {
      value: 'volgograd',
      label: 'Волгоград',
    },
    {
      value: 'yakutiya',
      label: 'Якутия',
    },
    {
      value: 'yanao',
      label: 'ЯНАО',
    },
    {
      value: 'yaroslavl',
      label: 'Ярославль',
    },
  ];

  const [checked, setChecked] = useState(null);
  const onCheckChange = (e) => {
    setChecked(e.target.checked);
    setVerbose(!checked);
  };

  const onSelectChange = (value) => {
    setProjectName(value);
  };

  const onCameraChange = (e) => {
    setCameraName(e.target.value)
  };

  const onUsernameChange = (e) => {
    setUsername(e.target.value)
  };

  const onPasswordChange = (e) => {
    setPassword(e.target.value)
  };

  return (
    <div>
      <Flex vertical="vertical" gap="small">
        <Space>
          <Space.Compact>
            <Select
            id="projectNameInput"
            placeholder="Проект"
            options={projects}
            onChange={onSelectChange}
            showSearch
            optionFilterProp="label"
            ></Select>
            <Input
            id="cameraNameInput"
            placeholder="Название камеры"
            onChange={onCameraChange}
            allowClear
            ></Input>
          </Space.Compact>
          <Checkbox
          id="verboseParam"
          checked={checked}
          onChange={onCheckChange}
          >channel.config и /info</Checkbox>
        </Space>
        <Space.Compact>
          <Input
          id="usernameInput"
          placeholder="Username"
          type={checked ? "text" : "hidden"}
          onChange={onUsernameChange}
          ></Input>
          <Input
          id="passwordInput"
          placeholder="Password"
          type={checked ? "password" : "hidden"}
          onChange={onPasswordChange}
          ></Input>
        </Space.Compact>
      </Flex>
    </div>
  )
}

export { InputForm }
