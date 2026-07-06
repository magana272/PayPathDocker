"use client";

import { useRef } from "react";
import styles from "./TabBar.module.css";

export default function TabBar({ tabs, activeTab, onTabChange, className }) {
  const tabsRef = useRef([]);

  const handleKeyDown = (e, index) => {
    let next;
    if (e.key === "ArrowRight") next = (index + 1) % tabs.length;
    else if (e.key === "ArrowLeft") next = (index - 1 + tabs.length) % tabs.length;
    else return;
    e.preventDefault();
    onTabChange(tabs[next].id);
    tabsRef.current[next]?.focus();
  };

  return (
    <div className={`${styles.tabs}${className ? ` ${className}` : ""}`} role="tablist">
      {tabs.map(({ id, label }, i) => (
        <button
          key={id}
          ref={(el) => { tabsRef.current[i] = el; }}
          className={`${styles.tab}${activeTab === id ? ` ${styles.active}` : ""}`}
          onClick={() => onTabChange(id)}
          onKeyDown={(e) => handleKeyDown(e, i)}
          role="tab"
          aria-selected={activeTab === id}
          tabIndex={activeTab === id ? 0 : -1}
        >
          {label}
        </button>
      ))}
    </div>
  );
}
