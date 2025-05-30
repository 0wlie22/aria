import React, { useEffect, useState } from "react";
import axios from "axios";
import { PieChart, Pie, Cell, Tooltip, Legend } from "recharts";

interface DataPoint {
  name: string;
  value: number;
}

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042"];

const ExpensesChart: React.FC = () => {
  const [data, setData] = useState<DataPoint[]>([]);

  useEffect(() => {
    axios
      .get("/api/expenses", {
        params: { month: 5, year: 2025 },
      })
      .then((res) => {
        console.log("API response:", res.data);
        if (!Array.isArray(res.data)) {
          throw new Error("Expected array but got something else");
        }
        const transformed = res.data.map((item) => ({
          name: item.category,
          value: item.total,
        }));
        setData(transformed);
      })
      .catch(console.error);
  }, []);

  return (
    <PieChart width={400} height={400}>
      <Pie
        data={data}
        cx="50%"
        cy="50%"
        labelLine={false}
        outerRadius={150}
        dataKey="value"
      >
        {data.map((_, index) => (
          <Cell key={index} fill={COLORS[index % COLORS.length]} />
        ))}
      </Pie>
      <Tooltip />
      <Legend />
    </PieChart>
  );
};

export default ExpensesChart;
