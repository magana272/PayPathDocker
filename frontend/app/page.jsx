"use client";

import { useEffect, useState, useCallback } from "react";
import Link from "next/link";
import { IconSettings } from "@/components/Icons";
import { api } from "@/lib/api";
import { cache } from "@/lib/cache";
import SummaryCards from "@/components/dashboard/SummaryCards";
import LiquidSection from "@/components/dashboard/LiquidSection";
import ExpensesSection from "@/components/dashboard/ExpensesSection";
import DebtsSection from "@/components/dashboard/DebtsSection";
import { SummaryCardsSkeleton, ListSectionSkeleton, DebtsSectionSkeleton } from "@/components/Skeleton";
import ds from "./page.module.css";

export default function Dashboard() {
  const [jobs, setJobs] = useState(() => cache.get("income") || []);
  const [liquid, setLiquid] = useState(() => cache.get("liquid") || []);
  const [expenses, setExpenses] = useState(() => cache.get("expenses") || []);
  const [debts, setDebts] = useState(() => cache.get("debts") || []);
  const [summary, setSummary] = useState(() => cache.get("summary"));

  const loadData = useCallback(() => {
    api.getIncome().then(setJobs);
    api.getLiquid().then(setLiquid);
    api.getExpenses().then(setExpenses);
    api.getDebts().then(setDebts);
    api.getSummary().then(setSummary);
  }, []);

  useEffect(loadData, [loadData]);

  useEffect(() => {
    window.addEventListener("paypath:refresh", loadData);
    return () => window.removeEventListener("paypath:refresh", loadData);
  }, [loadData]);

  return (
    <div className="page">
      <div className={ds.pageHeader}>
        <h1>Dashboard</h1>

        <Link
          href="/settings"
          className={ds.settingsIcon}
          aria-label="Settings"
        >
          <IconSettings />
        </Link>
      </div>

      {jobs.length > 0 && (
        <div className={ds.occupation}>
          {jobs.map((j) => j.job).join(" / ")}
        </div>
      )}

      {summary ? (
        <SummaryCards summary={summary} />
      ) : (
        <SummaryCardsSkeleton />
      )}

      {liquid.length > 0 ? (
        <LiquidSection liquid={liquid} />
      ) : !summary ? (
        <ListSectionSkeleton rows={2} label={110} />
      ) : null}

      {expenses.length > 0 ? (
        <ExpensesSection expenses={expenses} />
      ) : !summary ? (
        <ListSectionSkeleton rows={3} label={70} />
      ) : null}

      {summary ? (
        <DebtsSection debts={debts} summary={summary} />
      ) : (
        <DebtsSectionSkeleton />
      )}
    </div>
  );
}