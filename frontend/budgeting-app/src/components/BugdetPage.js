import React, { useState, useEffect } from 'react';
import './BudgetPage.css';
import { AuthProvider, useAuth } from './AuthContext';

// Helper functions for API calls
const API_BASE_URL = 'https://backend.shivikasingh.com/api' //'http://localhost:8080/api'; 

async function fetchIncome(userId, month, token) {
  const response = await fetch(`${API_BASE_URL}/income?userId=${userId}&month=${month}`, {
    method : 'GET', 
    headers: {
      'Authorization' : token
    }
  });
  console.log(response)
  if (!response.ok) throw new Error('Failed to fetch income');
  return response.json();
}

async function fetchBudget(userId, month, token) {
  const response = await fetch(`${API_BASE_URL}/budget?userId=${userId}&month=${month}`, {
    method: 'GET',
    headers: {
      'Authorization' : token
    }
  });
  if (!response.ok) throw new Error('Failed to fetch budget');
  return response.json();
}

async function addOrUpdateIncome(userId, month, item, editingIncomeIndex, token) {
  const method = editingIncomeIndex == true ? 'PUT' : 'POST';
  const url = editingIncomeIndex == true
    ? `${API_BASE_URL}/income/${userId}/${month}/${item.id}`
    : `${API_BASE_URL}/income`;
  const response = await fetch(url, {
    method,
    headers: { 
      'Content-Type': 'application/json',
      'Authorization' : token
    },
    body: editingIncomeIndex == true ? JSON.stringify({ newValue: item.incomeItemValue }) : JSON.stringify({ ...item, userId, month }),
  });
  if (!response.ok) throw new Error('Failed to add/update income');
  return response.json();
}

async function addOrUpdateBudget(userId, month, item, editingBudgetIndex, token) {
  const method = editingBudgetIndex == true ? 'PUT' : 'POST';
  const url = editingBudgetIndex == true 
    ? `${API_BASE_URL}/budget/${userId}/${month}/${item.id}`
    : `${API_BASE_URL}/budget`;
  const response = await fetch(url, {
    method,
    headers: { 
      'Content-Type': 'application/json',
      'Authorization' : token
    },
    body: editingBudgetIndex == true ? JSON.stringify({ newValue: item.budgetItemValue }) : JSON.stringify({ ...item, userId, month }),
  });
  if (!response.ok) throw new Error('Failed to add/update budget');
  return response.json();
}

async function deleteIncome(userId, month, itemId, token) {
  const response = await fetch(`${API_BASE_URL}/income/${userId}/${month}/${itemId}`, {
    method: 'DELETE',
    headers: {
      'Authorization' : token
    }
  });
  if (!response.ok) throw new Error('Failed to delete income');
}

async function deleteBudget(userId, month, itemId, token) {
  const response = await fetch(`${API_BASE_URL}/budget/${userId}/${month}/${itemId}`, {
    method: 'DELETE',
    headers: {
      'Authorization' : token
    }
  });
  if (!response.ok) throw new Error('Failed to delete budget');
}

function BudgetPage() {
  const [currentMonth, setCurrentMonth] = useState(new Date().toISOString().slice(0, 7));
  const [incomeSources, setIncomeSources] = useState([]);
  const [budgetCategories, setBudgetCategories] = useState([]);
  const [newIncomeSource, setNewIncomeSource] = useState({ incomeItemName: '', incomeItemValue: 0 });
  const [newBudgetCategory, setNewBudgetCategory] = useState({ budgetItemName: '', budgetItemValue: 0 });
  const [editingIncomeIndex, setEditingIncomeIndex] = useState(null);
  const [editingBudgetIndex, setEditingBudgetIndex] = useState(null);

  // Assume we have a userId from authentication
  //const userId = 'user123'; // Replace with actual user ID from auth
  const { userId, token } = useAuth();

  useEffect(() => {
    loadData();
  }, [currentMonth]);

  async function loadData() {
    try {
      const [incomeData, budgetData] = await Promise.all([
        fetchIncome(userId, currentMonth, token),
        fetchBudget(userId, currentMonth, token)
      ]);
      setIncomeSources(incomeData);
      setBudgetCategories(budgetData);
    } catch (error) {
      console.error('Failed to load data:', error);
      // Handle error (e.g., show error message to user)
    }
  }

  const addIncomeSource = async (e) => {
    e.preventDefault();
    try {
      const updatedItem = await addOrUpdateIncome(userId, currentMonth, newIncomeSource, editingIncomeIndex, token);
      if (editingIncomeIndex !== null) {
        setIncomeSources(prev => prev.map((item, index) => 
          index === editingIncomeIndex ? newIncomeSource : item
        ));
        setEditingIncomeIndex(null);
      } else {
        setIncomeSources(prev => [...prev, newIncomeSource]);
      }
      setNewIncomeSource({ incomeItemName: '', incomeItemValue: '' });
    } catch (error) {
      console.error('Failed to add/update income:', error);
      // Handle error
    }
  };

  const addBudgetCategory = async (e) => {
    e.preventDefault();
    try {
      const updatedItem = await addOrUpdateBudget(userId, currentMonth, newBudgetCategory, editingBudgetIndex, token);
      if (editingBudgetIndex !== null) {
        setBudgetCategories(prev => prev.map((item, index) => 
          index === editingBudgetIndex ? newBudgetCategory : item
        ));
        setEditingBudgetIndex(null);
      } else {
        setBudgetCategories(prev => [...prev, newBudgetCategory]);
      }
      setNewBudgetCategory({ budgetItemName: '', budgetItemValue: '' });
    } catch (error) {
      console.error('Failed to add/update budget:', error);
      // Handle error
    }
  };

  const editIncomeSource = (index) => {
    setNewIncomeSource(incomeSources[index]);
    setEditingIncomeIndex(index);
  };

  const editBudgetCategory = (index) => {
    setNewBudgetCategory(budgetCategories[index]);
    setEditingBudgetIndex(index);
  };

  const deleteIncomeSource = async (index) => {
    try {
      await deleteIncome(userId, currentMonth, incomeSources[index].id, token);
      setIncomeSources(prev => prev.filter((_, i) => i !== index));
    } catch (error) {
      console.error('Failed to delete income:', error);
      // Handle error
    }
  };

  const deleteBudgetCategory = async (index) => {
    try {
      await deleteBudget(userId, currentMonth, budgetCategories[index].id, token);
      setBudgetCategories(prev => prev.filter((_, i) => i !== index));
    } catch (error) {
      console.error('Failed to delete budget:', error);
      // Handle error
    }
  };

  const totalIncome = incomeSources.reduce((sum, source) => sum + parseFloat(source.incomeItemValue || 0), 0);
  const totalBudget = budgetCategories.reduce((sum, category) => sum + parseFloat(category.budgetItemValue || 0), 0);
  const remaining = totalIncome - totalBudget;

  return (
    <div className="income-budget-page">
      <div className="month-selector">
        <label htmlFor="month-select">Select Month: </label>
        <input 
          type="month" 
          id="month-select" 
          value={currentMonth} 
          onChange={(e) => setCurrentMonth(e.target.value)}
        />
      </div>

      <section className="income-section">
        <h2>Income Sources</h2>
        <form className="income-form" onSubmit={addIncomeSource}>
          <input
            type="text"
            placeholder="Source Name"
            value={newIncomeSource.incomeItemName}
            onChange={(e) => setNewIncomeSource({ ...newIncomeSource, incomeItemName: e.target.value })}
          />
          <input
            type="number"
            placeholder="Amount"
            value={newIncomeSource.incomeItemValue}
            onChange={(e) => setNewIncomeSource({ ...newIncomeSource, incomeItemValue: parseFloat(e.target.value) })}
          />
          <button type="submit">{editingIncomeIndex !== null ? 'Update' : 'Add'} Income</button>
        </form>
        <ul className="income-list">
          {incomeSources.map((source, index) => (
            <li>
              {source.incomeItemName}: ${source.incomeItemValue}
              <button onClick={() => editIncomeSource(index)}>Edit</button>
              <button onClick={() => deleteIncomeSource(index)}>Delete</button>
            </li>
          ))}
        </ul>
      </section>

      <section className="budget-section">
        <h2>Budget Categories</h2>
        <form className="budget-form" onSubmit={addBudgetCategory}>
          <input
            type="text"
            placeholder="Category Name"
            value={newBudgetCategory.budgetItemName}
            onChange={(e) => setNewBudgetCategory({ ...newBudgetCategory, budgetItemName: e.target.value })}
          />
          <input
            type="number"
            placeholder="Amount"
            value={newBudgetCategory.budgetItemValue}
            onChange={(e) => setNewBudgetCategory({ ...newBudgetCategory, budgetItemValue: parseFloat(e.target.value) })}
          />
          <button type="submit">{editingBudgetIndex !== null ? 'Update' : 'Add'} Category</button>
        </form>
        <ul className="budget-list">
          {budgetCategories.map((category, index) => (
            <li key={category.id}>
              {category.budgetItemName}: ${category.budgetItemValue}
              <button onClick={() => editBudgetCategory(index)}>Edit</button>
              <button onClick={() => deleteBudgetCategory(index)}>Delete</button>
            </li>
          ))}
        </ul>
      </section>

      <section className="summary-section">
        <h2>Summary for {new Date(currentMonth).toLocaleString('default', { month: 'long', year: 'numeric' })}</h2>
        <p>Total Income: ${totalIncome.toFixed(2)}</p>
        <p>Total Budgeted: ${totalBudget.toFixed(2)}</p>
        <p>Remaining: ${remaining.toFixed(2)}</p>
      </section>
    </div>
  );
}

export default BudgetPage;
