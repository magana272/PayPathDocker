"use client";

import styles from "./DataTable.module.css";

export default function DataTable({ children }) {
  return (
    <div className={styles.scroll}>
      <table className={styles.table}>
        {children}
      </table>
    </div>
  );
}

export { styles as tableStyles };
