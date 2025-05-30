import React, { useEffect, useState } from "react";
import axios from "axios";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

interface MonthlyTotal {
  year: number;
  month: number;
  type: string;
  total: number;
}

interface ChartData {
  month: string;
  expense: number;
  income: number;
  investment: number;
}

const MonthlyTotalsChart: React.FC = () => {
  const [chartData, setChartData] = useState<ChartData[]>([]);

  useEffect(() => {
    axios.get<MonthlyTotal[]>("/api/total")
      .then(res => {
        const raw = res.data;

        const grouped: Record<string, ChartData> = {};
        raw.forEach(item => {
          const key = `${item.year}-${String(item.month).padStart(2, "0")}`;
          if (!grouped[key]) {
            grouped[key] = { month: key, expense: 0, income: 0, investment: 0 };
          }
          if (item.type === "expense") grouped[key].expense = Math.abs(item.total);
          if (item.type === "income") grouped[key].income = item.total;
          if (item.type === "investment") grouped[key].investment = Math.abs(item.total);
        });

        setChartData(Object.values(grouped));
      })
      .catch(err => console.error("Error loading totals:", err));
  }, []);

  return (
    <ResponsiveContainer width="100%" height={400}>
      <BarChart data={chartData}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="month" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Bar dataKey="expense" fill="#FF8042" name="Expenses" />
        <Bar dataKey="income" fill="#00C49F" name="Income" />
        <Bar dataKey="investment" fill="#0088FE" name="Investment" />
      </BarChart>
    </ResponsiveContainer>
  );
};

export default MonthlyTotalsChart;

