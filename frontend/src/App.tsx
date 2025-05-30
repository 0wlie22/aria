import ExpensesBarChart from "./components/ExpensesBarChart";
import TotalIncomeExpenseChart from "./components/TotalIncomeExpenseChart";
import CSVImport from "./components/CSVImport";

function App() {
    return (
        <div style={{ textAlign: "center", marginTop: "50px" }}>
            <h1>Finance Dashboard</h1>
            <ExpensesBarChart />
            <TotalIncomeExpenseChart/>
            <CSVImport />
        </div>
    );
}

export default App;
