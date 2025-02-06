import React from "react";
import PingTable from "./components/PingTable";
import "antd/dist/reset.css";

const App = () => {
  return (
    <div>
      <h1 style={{ textAlign: "center", margin: "20px 0" }}>Статистика пингов</h1>
      <PingTable />
    </div>
  );
};

export default App;