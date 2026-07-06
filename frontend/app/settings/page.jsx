"use client";

import { Suspense, useCallback, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { api } from "@/lib/api";
import { cache } from "@/lib/cache";
import TabBar from "@/components/TabBar";
import IncomeTab from "@/components/settings/IncomeTab";
import ExpensesTab from "@/components/settings/ExpensesTab";
import DebtsTab from "@/components/settings/DebtsTab";
import AccountsTab from "@/components/settings/AccountsTab";
import AccountTab from "@/components/settings/AccountTab";

const TABS = [
  { id: "income", label: "Income" },
  { id: "expenses", label: "Expenses" },
  { id: "debts", label: "Debts" },
  { id: "accounts", label: "Accounts" },
  { id: "account", label: "Account" },
];

export default function Settings() {
  return (
    <Suspense fallback={null}>
      <SettingsContent />
    </Suspense>
  );
}

function SettingsContent() {
  const searchParams = useSearchParams();
  const initialTab = searchParams.get("tab") || "income";

  const [tab, setTab] = useState(initialTab);
  const [jobs, setJobs] = useState(() => cache.get("income") || []);
  const [liquid, setLiquid] = useState(() => cache.get("liquid") || []);
  const [expenses, setExpenses] = useState(() => cache.get("expenses") || []);
  const [debts, setDebts] = useState(() => cache.get("debts") || []);

  const load = useCallback(() => {
    api.getIncome().then(setJobs);
    api.getLiquid().then(setLiquid);
    api.getExpenses().then(setExpenses);
    api.getDebts().then(setDebts);
  }, []);

  useEffect(load, [load]);

  return (
    <div className="page">
      <h1>Settings</h1>
      <TabBar tabs={TABS} activeTab={tab} onTabChange={setTab} />

      {tab === "income" && <IncomeTab jobs={jobs} onReload={load} />}
      {tab === "expenses" && <ExpensesTab expenses={expenses} onReload={load} />}
      {tab === "debts" && <DebtsTab debts={debts} onReload={load} />}
      {tab === "accounts" && <AccountsTab liquid={liquid} onReload={load} />}
      {tab === "account" && <AccountTab />}
    </div>
  );
}