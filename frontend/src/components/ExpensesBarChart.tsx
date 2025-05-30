import React, { useEffect, useState } from "react";
import axios from "axios";
import { BarChart, Bar, XAxis, YAxis, Tooltip, Legend, CartesianGrid } from "recharts";

interface DataPoint {
  category: string;
  total: number;
}

interface MonthYear {
  month: number;
  year: number;
}

const monthNames = ["January","February","March","April","May","June","July","August","September","October","November","December"];

const ExpensesBarChart: React.FC = () => {
  const [data, setData] = useState<DataPoint[]>([]);
  const [months, setMonths] = useState<MonthYear[]>([]);
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    axios.get("/api/available-months")
      .then((res) => {
        if (Array.isArray(res.data)) {
          setMonths(res.data);
          setCurrentIndex(res.data.length - 1); // default → latest month
        }
      })
      .catch(console.error);
  }, []);

  useEffect(() => {
    if (months.length === 0) return;
    const { month, year } = months[currentIndex];
    axios.get("/api/expenses", { params: { month, year } })
      .then((res) => setData(res.data))
      .catch(console.error);
  }, [months, currentIndex]);

  const handlePrev = () => {
    if (currentIndex > 0) setCurrentIndex(currentIndex - 1);
  };

  const handleNext = () => {
    if (currentIndex < months.length - 1) setCurrentIndex(currentIndex + 1);
  };

  if (months.length === 0) return <p>Loading...</p>;

  const { month, year } = months[currentIndex];

  return (
    <div style={{ marginTop: "50px" }}>
      <h2>Expenses by Category</h2>

      <div style={{ display: "flex", justifyContent: "center", alignItems: "center", marginBottom: "10px" }}>
        <button onClick={handlePrev} disabled={currentIndex === 0}>◀</button>
        <span style={{ margin: "0 10px", fontWeight: "bold" }}>{`${monthNames[month-1]} ${year}`}</span>
        <button onClick={handleNext} disabled={currentIndex === months.length - 1}>▶</button>
      </div>

      <BarChart width={600} height={400} data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="category" />
        <YAxis />
        <Tooltip />
        <Legend />
        <Bar dataKey="total" fill="#8884d8" />
      </BarChart>
    </div>
  );
};

export default ExpensesBarChart;
