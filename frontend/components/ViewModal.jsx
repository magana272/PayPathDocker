"use client";

import Modal, { modalStyles } from "./Modal";
import ss from "@/app/settings/page.module.css";

export default function ViewModal({ isOpen, onClose, title, fields, onEdit, onDelete, editLabel = "Edit", children }) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={title}>
      <div>
        {fields?.map((f, i) => (
          <div key={i} className={ss.profileRow}>
            <span className={ss.profileLabel}>{f.label}</span>
            <span className={f.capitalize ? ss.profileValueCaps : ss.profileValue}>{f.value}</span>
          </div>
        ))}
        {children}
        {(onEdit || onDelete) && (
          <div className={modalStyles.actions}>
            {onEdit && (
              <button className={modalStyles.submitFlex} onClick={onEdit}>{editLabel}</button>
            )}
            {onDelete && (
              <button className={ss.dangerBtn} onClick={onDelete}>Delete</button>
            )}
          </div>
        )}
      </div>
    </Modal>
  );
}
