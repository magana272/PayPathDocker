"use client";

import Link from "next/link";
import ds from "@/app/page.module.css";
import styles from "@/components/CardGrid.module.css";

export default function DebtsSection({ debts, summary }) {
  const totalDebt = debts.reduce((s, d) => s + d.balance, 0);

  return (
    <div className={ds.section}>
      <div className={ds.sectionHeader}>
        <span>Debts</span>
        <Link href="/settings?tab=debts" className={ds.manageLink}>Manage →</Link>
      </div>
      <div className={styles.grid} style={{ marginBottom: 6 }}>
        <div className={styles.card}>
          <h3 className={styles.cardTitle}>Total Owed</h3>
          <p className="big-number red">${totalDebt.toLocaleString()}</p>
        </div>
        <div className={styles.card}>
          <h3 className={styles.cardTitle}>Accounts</h3>
          <p className="big-number">{debts.length}</p>
        </div>
        {summary && (
          <>
            <div className={styles.card}>
              <h3 className={styles.cardTitle}>Monthly Interest</h3>
              <p className="big-number red">${summary.monthly_interest.toLocaleString()}</p>
            </div>
            <div className={styles.card}>
              <h3 className={styles.cardTitle}>DTI</h3>
              <p className={`big-number ${summary.dti > 80 ? "red" : ""}`}>{summary.dti}%</p>
            </div>
          </>
        )}
      </div>
      <div className={ds.list}>
        {debts.map((d) => (
          <div key={d.id} className={ds.listItem}>
            <span>{d.name} <span className={ds.freq}>{d.bank} · {d.apy}% APR</span></span>
            <span className={ds.listValueRed}>${d.balance.toLocaleString()}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
