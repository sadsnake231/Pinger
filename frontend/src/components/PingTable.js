

import React, { useEffect, useState } from "react";
import { Table, Button, message } from "antd";
import axios from "axios";

const PingTable = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchData = async () => {
    setLoading(true);
    try {
      const response = await axios.get("http://localhost:5000/");
      setData(response.data);
    } catch (error) {
      message.error("Ошибка при загрузке данных!");
      console.error("Ошибка запроса:", error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchData();
  }, []);

  const columns = [
    {
      title: "Хост",
      dataIndex: "ip",
      key: "ip",
    },
    {
      title: "Мин. время (мс)",
      dataIndex: "min",
      key: "min ",
      render: (value) => value?.toFixed(2),
    },
    {
      title: "Макс. время (мс)",
      dataIndex: "max",
      key: "max",
      render: (value) => value?.toFixed(2),
    },
    {
      title: "Последний ответ",
      dataIndex: "last_up",
      key: "last_up",
    },
    {
      title: "Время пинга",
      dataIndex: "time",
      key: "time",
    },
  ];

  return (
    <div style={{ padding: 20 }}>
      <Button type="primary" onClick={fetchData} loading={loading} style={{ marginBottom: 16 }}>
        Обновить данные
      </Button>
      <Table columns={columns} dataSource={data} rowKey="host" />
    </div>
  );
};

export default PingTable;
