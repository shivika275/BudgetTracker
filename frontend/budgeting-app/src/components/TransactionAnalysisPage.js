import React, { useState, useEffect } from 'react';
import * as XLSX from 'xlsx';
import './TransactionAnalysisPage.css';
import { AuthProvider, useAuth } from './AuthContext';

const API_BASE_URL = 'https://backend.shivikasingh.com/api' //'http://localhost:8080/api';

function TransactionAnalysisPage() {
  const [transactions, setTransactions] = useState([]);
  const [currentMonth, setCurrentMonth] = useState(new Date().toISOString().slice(0, 7));
  //const userId = 'user123'; // Replace with actual user ID from authentication

  const { userId, token } = useAuth();

  useEffect(() => {
    fetchTransactions();
  }, [currentMonth]);

  const fetchTransactions = async () => {
    try {
      const response = await fetch(`${API_BASE_URL}/expense?userId=${userId}&month=${currentMonth}`, {
        method: 'GET',
        headers: {
          'Authorization' : token
        }
      });
      if (!response.ok) throw new Error('Failed to fetch transactions');
      const data = await response.json();
      setTransactions(data);
    } catch (error) {
      console.error('Error fetching transactions:', error);
    }
  };

  const handleFileUpload = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();
    reader.onload = (evt) => {
      const bstr = evt.target.result;
      const wb = XLSX.read(bstr, { type: 'binary' });
      const wsname = wb.SheetNames[0];
      const ws = wb.Sheets[wsname];
      const data = XLSX.utils.sheet_to_json(ws);
      const formattedData = data.map(item => ({
        userId,
        expenseItemName: item.name,
        month: currentMonth,
        expenseItemValue: parseFloat(item.amount),
        expenseTags: [],
      }));
      
      // Deduplicate and merge with existing transactions
      const mergedTransactions = deduplicateTransactions([...transactions, ...formattedData]);
      setTransactions(mergedTransactions);
      saveTransactions(mergedTransactions);
    };
    reader.readAsBinaryString(file);
  };

  const deduplicateTransactions = (allTransactions) => {
    const uniqueTransactions = {};
    allTransactions.forEach(transaction => {
      uniqueTransactions[transaction.expenseItemName] = transaction;
    });
    return Object.values(uniqueTransactions);
  };

  const saveTransactions = async (transactionsToSave) => {
    try {
      const response = await fetch(`${API_BASE_URL}/expense`, {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json' ,
          'Authorization' : token
        },
        body: JSON.stringify({ expenses: transactionsToSave }),
      });
      if (!response.ok) throw new Error('Failed to save transactions');
      console.log('Transactions saved successfully');
    } catch (error) {
      console.error('Error saving transactions:', error);
    }
  };

  const handleCategoryChange = async (index, category) => {
    const updatedTransactions = [...transactions];
    updatedTransactions[index].expenseTags = [category];
    setTransactions(updatedTransactions);
    await updateTransaction(updatedTransactions[index]);
  };

  const handleValueChange = async (index, newValue) => {
    const updatedTransactions = [...transactions];
    updatedTransactions[index].expenseItemValue = parseFloat(newValue);
    setTransactions(updatedTransactions);
    await updateTransaction(updatedTransactions[index]);
  };

  const updateTransaction = async (transaction) => {
    try {
      const response = await fetch(`${API_BASE_URL}/expense/${userId}/${currentMonth}/${transaction.expenseItemName}`, {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json' ,
          'Authorization' : token
        },
        body: JSON.stringify({
          newValue: transaction.expenseItemValue,
          newTags: transaction.expenseTags
        }),
      });
      if (!response.ok) throw new Error('Failed to update transaction');
      console.log('Transaction updated successfully');
    } catch (error) {
      console.error('Error updating transaction:', error);
    }
  };

  const deleteTransaction = async (index) => {
    const transactionToDelete = transactions[index];
    try {
      const response = await fetch(`${API_BASE_URL}/expense/${userId}/${currentMonth}/${transactionToDelete.expenseItemName}`, {
        method: 'DELETE',
        headers: {
          'Authorization' : token
        }
      });
      if (!response.ok) throw new Error('Failed to delete transaction');
      console.log('Transaction deleted successfully');
      setTransactions(transactions.filter((_, i) => i !== index));
    } catch (error) {
      console.error('Error deleting transaction:', error);
    }
  };

  return (
    <div className="transaction-analysis-page">
      <div className="month-selector">
        <label htmlFor="month-select">Select Month: </label>
        <input 
          type="month" 
          id="month-select" 
          value={currentMonth} 
          onChange={(e) => setCurrentMonth(e.target.value)}
        />
      </div>

      <div className="file-upload">
        <input type="file" onChange={handleFileUpload} accept=".xlsx, .xls" />
      </div>
      
      {transactions.length > 0 && (
        <table className="transactions-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Amount</th>
              <th>Category</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {transactions.map((transaction, index) => (
              <tr key={index}>
                <td>{transaction.expenseItemName}</td>
                <td>
                  <input
                    type="number"
                    value={transaction.expenseItemValue}
                    onChange={(e) => handleValueChange(index, e.target.value)}
                  />
                </td>
                <td>
                  <select
                    className="category-select"
                    value={(transaction.expenseTags && transaction.expenseTags.length > 0) ? transaction.expenseTags[0] : ''}
                    onChange={(e) => handleCategoryChange(index, e.target.value)}
                  >
                    <option value="">Select Category</option>
                    <option value="Food">Food</option>
                    <option value="Transportation">Transportation</option>
                    <option value="Entertainment">Entertainment</option>
                    <option value="Utilities">Utilities</option>
                  </select>
                </td>
                <td>
                  <button onClick={() => deleteTransaction(index)}>Delete</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default TransactionAnalysisPage;
