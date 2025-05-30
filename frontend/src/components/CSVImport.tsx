import React, { useState, useEffect } from "react";
import axios from "axios";

interface Transaction {
  id: number;
  date: string;
  amount: number;
  narrative: string;
  category_id: number | null;
}

interface Category {
  id: number;
  name: string;
  type: string;
}

const TransactionClassifier: React.FC = () => {
  const [file, setFile] = useState<File | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);

  // Load unclassified transactions
  useEffect(() => {
    axios.get("/api/unclassified").then(res => setTransactions(res.data));
    axios.get("/api/categories").then(res => setCategories(res.data));
  }, []);

  const handleUpload = async () => {
    if (!file) return;
    const formData = new FormData();
    formData.append("file", file);
    await axios.post("/api/upload", formData, { headers: { "Content-Type": "multipart/form-data" }});
    alert("File uploaded successfully");
  };

  const handleCategoryChange = (id: number, category_id: number) => {
    setTransactions(prev => prev.map(tx => tx.id === id ? { ...tx, category_id } : tx));
  };

  const saveClassification = async (id: number, category_id: number) => {
    await axios.put(`/api/classify/${id}`, { category_id });
  };

  return (
    <div>
      <h2>Upload Bank Statement CSV</h2>
      <input type="file" onChange={e => setFile(e.target.files?.[0] || null)} />
      <button onClick={handleUpload}>Upload</button>

      <h3>Unclassified Transactions</h3>
      <table border={1} cellPadding={8}>
        <thead>
          <tr>
            <th>Date</th>
            <th>Amount</th>
            <th>Narrative</th>
            <th>Category</th>
            <th>Save</th>
          </tr>
        </thead>
        <tbody>
          {transactions.map(tx => (
            <tr key={tx.id}>
              <td>{tx.date}</td>
              <td>{(tx.amount / 100).toFixed(2)}</td>
              <td>{tx.narrative}</td>
              <td>
                <select
                  value={tx.category_id || ""}
                  onChange={e => handleCategoryChange(tx.id, Number(e.target.value))}
                >
                  <option value="">-- Select --</option>
                  {categories.map(cat => (
                    <option key={cat.id} value={cat.id}>
                      {cat.type} - {cat.name}
                    </option>
                  ))}
                </select>
              </td>
              <td>
                <button onClick={() => saveClassification(tx.id, tx.category_id!)}>Save</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default TransactionClassifier;

