import React, { useEffect, useState } from "react";
import axios from "axios";

interface MonthlyTotal {
  type: string;
  total: number;
  month: number;
  year: number;
}


const TotalsTable: React.FC = () => {
  const [totals, setTotals] = useState<MonthlyTotal[]>([]);

  useEffect(() => {
    axios.get("/api/total")
      .then(res => setTotals(res.data))
      .catch(err => console.error("Error loading totals:", err));
  }, []);

  return (
    <table border={1} cellPadding={8}>
      <thead>
        <tr>
          <th>Year</th>
          <th>Month</th>
          <th>Type</th>
          <th>Total</th>
        </tr>
      </thead>
      <tbody>
        {totals.map((t, idx) => (
          <tr key={idx}>
            <td>{t.year}</td>
            <td>{t.month}</td>
            <td>{t.type}</td>
            <td>{t.total.toFixed(2)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

export default TotalsTable;

