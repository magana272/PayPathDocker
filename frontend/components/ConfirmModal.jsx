"use client";

import Modal, { modalStyles } from "./Modal";

export default function ConfirmModal({ isOpen, onClose, title, message, confirmLabel = "Confirm", onConfirm, loading = false }) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={title}>
      <p className={modalStyles.message}>{message}</p>
      <div className={modalStyles.actions}>
        <button className={modalStyles.btnDanger} onClick={onConfirm} disabled={loading}>
          {loading ? "..." : confirmLabel}
        </button>
        <button className={modalStyles.btnSecondary} onClick={onClose}>Cancel</button>
      </div>
    </Modal>
  );
}
