export function calcMinPayment(balance, apy, debtType) {
  if (balance <= 0) return 0;
  const r = apy / 100 / 12;
  const interest = balance * r;

  if (debtType === "credit_card") {
    return Math.max(balance * 0.01 + interest, 25);
  }

  if (r === 0) return balance / 60;

  let n = 120;
  if (debtType === "car") n = 60;

  return (balance * r * Math.pow(1 + r, n)) / (Math.pow(1 + r, n) - 1);
}

export function simulateAvalanche(debts, budget, extraPayment) {
  if (!debts?.length || budget <= 0) return null;

  let balances = debts
    .filter((d) => d.balance > 0)
    .map((d) => ({
      name: d.name,
      balance: d.balance,
      apy: d.apy,
      rate: d.apy / 100 / 12,
      debtType: d.type,
      fixedMin: d.type !== "credit_card" ? calcMinPayment(d.balance, d.apy, d.type) : 0,
    }));

  if (!balances.length) return null;

  balances.sort((a, b) => b.apy - a.apy);

  const history = [];
  const totalBudget = budget + extraPayment;
  let cumulativeInterest = 0;

  for (let month = 1; month <= 480; month++) {
    let monthInterest = 0;
    for (const d of balances) {
      if (d.balance <= 0) { d.min = 0; continue; }
      d.min = d.debtType === "credit_card"
        ? calcMinPayment(d.balance, d.apy, d.debtType)
        : d.fixedMin;
      const interest = d.balance * d.rate;
      monthInterest += interest;
      cumulativeInterest += interest;
      d.balance += interest;
      d.min = Math.min(d.min, d.balance);
    }

    const sumMins = balances.reduce((s, d) => s + d.min, 0);
    let remaining = Math.max(totalBudget, sumMins);
    let monthPaid = 0;

    for (const d of balances) {
      if (d.balance <= 0) continue;
      const payment = Math.min(d.min, remaining);
      d.balance -= payment;
      remaining -= payment;
      monthPaid += payment;
      if (d.balance < 0.01) d.balance = 0;
    }

    for (const d of balances) {
      if (remaining <= 0) break;
      if (d.balance <= 0) continue;
      const payment = Math.min(remaining, d.balance);
      d.balance -= payment;
      remaining -= payment;
      monthPaid += payment;
      if (d.balance < 0.01) d.balance = 0;
    }

    const row = { month };
    let total = 0;
    for (const d of balances) {
      const bal = Math.round(d.balance * 100) / 100;
      row[d.name] = bal;
      total += d.balance;
    }
    row.total = Math.round(total * 100) / 100;
    row.interest = Math.round(cumulativeInterest * 100) / 100;
    row.monthInterest = Math.round(monthInterest * 100) / 100;
    row.monthPrincipal = Math.round((monthPaid - monthInterest) * 100) / 100;
    history.push(row);

    if (total <= 0) {
      return { budget, months: month, total_interest: Math.round(cumulativeInterest * 100) / 100, history };
    }
  }

  return { budget, months: 480, total_interest: Math.round(cumulativeInterest * 100) / 100, history };
}

export function downsample(data, maxPoints = 80) {
  if (!data || data.length <= maxPoints) return data;
  const step = (data.length - 1) / (maxPoints - 1);
  const out = [];
  for (let i = 0; i < maxPoints - 1; i++) {
    out.push(data[Math.round(i * step)]);
  }
  out.push(data[data.length - 1]);
  return out;
}
