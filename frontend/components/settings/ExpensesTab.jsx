"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import Modal, { modalStyles } from "@/components/Modal";
import ViewModal from "@/components/ViewModal";
import DataTable, { tableStyles } from "@/components/DataTable";
import { IconEye, IconPencil } from "@/components/Icons";
import { FREQ_MULT } from "@/lib/constants";
import ss from "@/app/settings/page.module.css";

export default function ExpensesTab({ expenses, onReload }) {
  const [showAdd, setShowAdd] = useState(false);
  const [addForm, setAddForm] = useState({ expense: "", cost: "", due_date: "", date: "", frequency: "monthly" });
  const [viewing, setViewing] = useState(null);
  const [editing, setEditing] = useState(null);
  const [editForm, setEditForm] = useState({});
  const [saving, setSaving] = useState(false);

  const totalExpenses = expenses.reduce((sum, e) => {
    return sum + e.cost * (FREQ_MULT[e.frequency] ?? 1);
  }, 0);

  const handleAdd = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      const payload = { expense: addForm.expense, cost: parseFloat(addForm.cost), frequency: addForm.frequency };
      if (addForm.frequency === "one-time") {
        payload.date = addForm.date || null;
        payload.due_date = null;
      } else {
        payload.due_date = addForm.due_date ? parseInt(addForm.due_date) : null;
        payload.date = null;
      }
      await api.addExpense(payload);
      setAddForm({ expense: "", cost: "", due_date: "", date: "", frequency: "monthly" });
      setShowAdd(false);
      onReload();
    } finally {
      setSaving(false);
    }
  };

  const openView = (exp) => setViewing(exp);

  const openEdit = (exp) => {
    setEditing(exp);
    setEditForm({ expense: exp.expense, cost: exp.cost, due_date: exp.due_date || "", date: exp.date || "", frequency: exp.frequency });
  };

  const handleSave = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      const payload = { expense: editForm.expense, cost: parseFloat(editForm.cost), frequency: editForm.frequency };
      if (editForm.frequency === "one-time") {
        payload.date = editForm.date || null;
        payload.due_date = null;
      } else {
        payload.due_date = editForm.due_date ? parseInt(editForm.due_date) : null;
        payload.date = null;
      }
      await api.updateExpense(editing.id, payload);
      setEditing(null);
      onReload();
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    setSaving(true);
    try {
      await api.deleteExpense(viewing.id);
      setViewing(null);
      onReload();
    } finally {
      setSaving(false);
    }
  };

  const viewFields = viewing ? [
    { label: "Name", value: viewing.expense },
    { label: "Cost", value: `$${viewing.cost.toFixed(2)}` },
    { label: "Due Date", value: viewing.due_date ? `Day ${viewing.due_date}` : "—" },
    { label: "Frequency", value: viewing.frequency },
  ] : [];

  return (
    <>
      <div className={ss.sectionBar}>
        <div className={ss.sectionTitleGroup}>
          <span className={ss.sectionTitle}>Expenses</span>
          <span className={ss.totalSubtitle}>${totalExpenses.toFixed(2)}/mo</span>
        </div>
        <button className={ss.addBtn} onClick={() => setShowAdd(true)}>+ Add</button>
      </div>

      <DataTable>
        <thead>
          <tr><th>Expense</th><th>Cost</th><th>Freq</th><th></th></tr>
        </thead>
        <tbody>
          {expenses.map((e) => (
            <tr key={e.id}>
              <td className={ss.tdTruncate}>{e.expense}</td>
              <td>${e.cost.toFixed(2)}</td>
              <td className={ss.tdMeta}>{e.frequency}</td>
              <td>
                <div className={tableStyles.actions}>
                  <button className={tableStyles.iconBtn} onClick={() => openView(e)} title="View" aria-label="View">
                    <IconEye />
                  </button>
                  <button className={tableStyles.iconBtn} onClick={() => openEdit(e)} title="Edit" aria-label="Edit">
                    <IconPencil />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </DataTable>

      {/* Add Modal */}
      <Modal isOpen={showAdd} onClose={() => setShowAdd(false)} title="Add Expense">
        <form className={modalStyles.form} onSubmit={handleAdd}>
          <input placeholder="Expense name" value={addForm.expense} onChange={(e) => setAddForm({ ...addForm, expense: e.target.value })} required />
          <div className={modalStyles.formRow}>
            <input type="number" step="0.01" placeholder="Cost" value={addForm.cost} onChange={(e) => setAddForm({ ...addForm, cost: e.target.value })} required />
            <input type="number" placeholder="Due day (1-31)" value={addForm.due_date} onChange={(e) => setAddForm({ ...addForm, due_date: e.target.value })} />
          </div>
          <select value={addForm.frequency} onChange={(e) => setAddForm({ ...addForm, frequency: e.target.value })}>
            <option value="monthly">Monthly</option>
            <option value="biweekly">Biweekly</option>
            <option value="weekly">Weekly</option>
            <option value="yearly">Yearly</option>
          </select>
          <button type="submit" className={modalStyles.submit} disabled={saving}>{saving ? "Saving..." : "Add Expense"}</button>
        </form>
      </Modal>

      {/* View Modal */}
      <ViewModal isOpen={!!viewing} onClose={() => setViewing(null)} title={viewing?.expense || ""} fields={viewFields} onEdit={() => { openEdit(viewing); setViewing(null); }} onDelete={handleDelete} />

      {/* Edit Modal */}
      <Modal isOpen={!!editing} onClose={() => setEditing(null)} title={`Edit — ${editing?.expense || ""}`}>
        {editing && (
          <form className={modalStyles.form} onSubmit={handleSave}>
            <input placeholder="Expense name" value={editForm.expense} onChange={(e) => setEditForm({ ...editForm, expense: e.target.value })} required />
            <div className={modalStyles.formRow}>
              <input type="number" step="0.01" placeholder="Cost" value={editForm.cost} onChange={(e) => setEditForm({ ...editForm, cost: e.target.value })} required />
              <input type="number" placeholder="Due day (1-31)" value={editForm.due_date} onChange={(e) => setEditForm({ ...editForm, due_date: e.target.value })} />
            </div>
            <select value={editForm.frequency} onChange={(e) => setEditForm({ ...editForm, frequency: e.target.value })}>
              <option value="monthly">Monthly</option>
              <option value="biweekly">Biweekly</option>
              <option value="weekly">Weekly</option>
              <option value="yearly">Yearly</option>
            </select>
            <button type="submit" className={modalStyles.submit} disabled={saving}>{saving ? "Saving..." : "Save"}</button>
          </form>
        )}
      </Modal>
    </>
  );
}
