"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import Modal, { modalStyles } from "@/components/Modal";
import ViewModal from "@/components/ViewModal";
import DataTable, { tableStyles } from "@/components/DataTable";
import { IconEye, IconPencil } from "@/components/Icons";
import ss from "@/app/settings/page.module.css";

export default function AccountsTab({ liquid, onReload }) {
  const [viewing, setViewing] = useState(null);
  const [editing, setEditing] = useState(null);
  const [editForm, setEditForm] = useState({ balance: "" });
  const [saving, setSaving] = useState(false);

  const openEdit = (l) => {
    setEditing(l);
    setEditForm({ balance: l.balance });
  };

  const handleUpdate = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      await api.updateLiquid(editing.id, { bank: editing.bank, balance: parseFloat(editForm.balance) });
      setEditing(null);
      onReload();
    } finally {
      setSaving(false);
    }
  };

  const viewFields = viewing ? [
    { label: "Bank", value: viewing.bank },
    { label: "Balance", value: `$${viewing.balance.toLocaleString()}` },
  ] : [];

  return (
    <>
      <div className={ss.sectionBar}>
        <span className={ss.sectionTitle}>Liquid Accounts</span>
      </div>
      <DataTable>
        <thead><tr><th>Bank</th><th>Balance</th><th></th></tr></thead>
        <tbody>
          {liquid.map((l) => (
            <tr key={l.id}>
              <td>{l.bank}</td>
              <td>${l.balance.toLocaleString()}</td>
              <td>
                <div className={tableStyles.actions}>
                  <button className={tableStyles.iconBtn} onClick={() => setViewing(l)} title="View" aria-label="View">
                    <IconEye />
                  </button>
                  <button className={tableStyles.iconBtn} onClick={() => openEdit(l)} title="Edit" aria-label="Edit">
                    <IconPencil />
                  </button>
                </div>
              </td>
            </tr>
          ))}
        </tbody>
      </DataTable>

      {/* View Modal */}
      <ViewModal isOpen={!!viewing} onClose={() => setViewing(null)} title={viewing?.bank || ""} fields={viewFields} onEdit={() => { openEdit(viewing); setViewing(null); }} editLabel="Edit Balance" />

      {/* Edit Modal */}
      <Modal isOpen={!!editing} onClose={() => setEditing(null)} title={`Update ${editing?.bank || ""}`}>
        {editing && (
          <form className={modalStyles.form} onSubmit={handleUpdate}>
            <input type="number" step="0.01" placeholder="Balance" value={editForm.balance} onChange={(e) => setEditForm({ balance: e.target.value })} required />
            <button type="submit" className={modalStyles.submit} disabled={saving}>{saving ? "Saving..." : "Update Balance"}</button>
          </form>
        )}
      </Modal>

    </>
  );
}
