import React, { useEffect, useState } from "react";
import { Table, Button, message } from "antd";
import axios from "axios";

const PingTable = ({ onUnauthorized }) => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(false);

  // Добавляем определение колонок
  const columns = [
    {
      title: "Хост",
      dataIndex: "ip",
      key: "ip",
    },
    {
      title: "Мин. время (мс)",
      dataIndex: "min",
      key: "min",
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

  const fetchData = async () => {
    setLoading(true);
    try {
      const response = await axios.get("http://localhost:5000/", {
        withCredentials: true
      });
      setData(response.data);
    } catch (error) {
      if (error.response?.status === 401) {
        message.error("Требуется авторизация!");
        onUnauthorized();
      } else {
        message.error("Ошибка при загрузке данных!");
      }
      console.error("Ошибка запроса:", error);
    }
    setLoading(false);
  };

  useEffect(() => {
    fetchData();
  }, []);

  return (
    <div style={{ padding: 20 }}>
      <Button 
        type="primary" 
        onClick={fetchData} 
        loading={loading} 
        style={{ marginBottom: 16 }}
      >
        Обновить данные
      </Button>
      <Table 
        columns={columns} 
        dataSource={data} 
        rowKey="ip" // Лучше использовать ip вместо host, если в данных есть ip
      />
    </div>
  );
};

export default PingTable;