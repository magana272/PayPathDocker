"use client";

import { useMemo, useDeferredValue, useRef, useState, useEffect, useLayoutEffect } from "react";
import { simulateAvalanche, downsample, calcMinPayment } from "@/lib/simulate";
import { GradientArea, StackedArea } from "@/components/charts";
import styles from "@/components/CardGrid.module.css";
import es from "@/app/explore/page.module.css";

export default function DebtPayoffTab({ payoff: initialPayoff, scenarios, debts, extraPayment = 0, onExtraPaymentChange, stickyOffset = 0, active = true }) {
  const baseBudget = initialPayoff?.budget || 0;
  const sliderMax = Math.max(2000, Math.round(baseBudget * 3 / 50) * 50);

  const totalMinPayments = useMemo(
    () => debts?.reduce((sum, d) => sum + calcMinPayment(d.balance, d.apy, d.type), 0) || 0,
    [debts]
  );

  const deferredExtra = useDeferredValue(extraPayment);

  const simulated = useMemo(() => {
    if (!debts?.length || baseBudget <= 0) return null;
    return simulateAvalanche(debts, baseBudget, deferredExtra);
  }, [deferredExtra, debts, baseBudget]);

  const payoff = simulated || initialPayoff;

  const hasHistory = payoff?.history?.length > 0;
  const negativeBudget = payoff && payoff.budget < 0;

  const debtNames = useMemo(
    () => initialPayoff?.history?.length
      ? Object.keys(initialPayoff.history[0]).filter((k) => k !== "month" && k !== "total" && k !== "interest")
      : [],
    [initialPayoff]
  );

  const chartHistory = useMemo(() => downsample(payoff?.history, 48), [payoff]);

  const hasInterest = chartHistory?.[0]?.interest != null;

  const interestLine = hasInterest
    ? [{ dataKey: "interest", name: "Cumulative Interest", color: "#c41e1e", width: 2.5, dashed: true }]
    : undefined;

  const totalDebtData = useMemo(
    () => chartHistory?.map((h) => ({ label: h.month, total: h.total })) ?? [],
    [chartHistory]
  );

  const didConverge = payoff?.months < 480;

  const savedMonths = initialPayoff?.months && payoff?.months && didConverge
    ? initialPayoff.months - payoff.months
    : 0;

  const sliderRef = useRef(null);
  const [sliderHeight, setSliderHeight] = useState(0);
  const showSlider = !!(payoff && !negativeBudget && debts?.length > 0);

  useLayoutEffect(() => {
    const el = sliderRef.current;
    if (!el) { setSliderHeight(0); return; }

    const measure = () => {
      const cs = getComputedStyle(el);
      const h  = parseFloat(cs.height) || 0;
      const pt = parseFloat(cs.paddingTop) || 0;
      const pb = parseFloat(cs.paddingBottom) || 0;
      const bt = parseFloat(cs.borderTopWidth) || 0;
      const bb = parseFloat(cs.borderBottomWidth) || 0;
      setSliderHeight(Math.ceil(bt + pt + h + pb + bb));
    };
    measure();

    const ro = new ResizeObserver(measure);
    ro.observe(el, { box: "border-box" });
    window.addEventListener("resize", measure);
    return () => { ro.disconnect(); window.removeEventListener("resize", measure); };
  }, [showSlider]);

  const stickyTop = { top: stickyOffset };
  const chartTop = stickyOffset + sliderHeight;

  const [balanceExpanded, setBalanceExpanded] = useState(false);
  const [totalDebtExpanded, setTotalDebtExpanded] = useState(true);
  const totalDebtRef = useRef(null);

  useEffect(() => {
    const el = totalDebtRef.current;
    if (!el || !active) return;
    const update = () => {
      const distance = el.getBoundingClientRect().top - chartTop;
      const range = window.innerHeight - chartTop;
      const t = 1 - Math.max(0, Math.min(distance / range, 1));
      el.style.opacity = t;
      el.style.transform = `translateY(${(1 - t) * 8}px)`;
    };
    update();
    window.addEventListener("scroll", update, { passive: true });
    return () => window.removeEventListener("scroll", update);
  }, [hasHistory, chartTop, active]);

  const chartStickyStyle = {
    top: chartTop,
    "--chart-top": `${chartTop}px`,
  };

  return (
    <>
      {payoff && (
        <div className={styles.grid}>
          <div className={`${styles.card}${negativeBudget ? ` ${styles.danger}` : ""}`}>
            <h3 className={styles.cardTitle}>Monthly Budget</h3>
            <p className="big-number">
              {negativeBudget ? "−" : ""}${Math.abs(payoff.budget).toLocaleString()}
              {extraPayment > 0 && (
                <span style={{ fontSize: 12, color: "var(--green)", marginLeft: 6 }}>+${extraPayment}</span>
              )}
            </p>
          </div>
          {hasHistory && didConverge ? (
            <div className={`${styles.card} ${styles.accent}`}>
              <h3 className={styles.cardTitle}>Debt-Free In</h3>
              <p className="big-number">
                {Math.floor(payoff.months / 12)}y {payoff.months % 12}m
                {savedMonths > 0 && (
                  <span style={{ fontSize: 11, color: "var(--green)", marginLeft: 6 }}>
                    ({savedMonths}mo faster)
                  </span>
                )}
              </p>
            </div>
          ) : (
            <div className={`${styles.card} ${styles.danger}`}>
              <h3 className={styles.cardTitle}>Debt-Free In</h3>
              <p className="big-number">40y+</p>
            </div>
          )}
          {totalMinPayments > 0 && (
            <div className={styles.card}>
              <h3 className={styles.cardTitle}>Min Payments</h3>
              <p className="big-number">${Math.round(totalMinPayments).toLocaleString()}/mo</p>
            </div>
          )}
          {didConverge && payoff.total_interest > 0 && (
            <div className={`${styles.card} ${styles.danger}`}>
              <h3 className={styles.cardTitle}>Total Interest</h3>
              <p className="big-number">${payoff.total_interest.toLocaleString()}</p>
            </div>
          )}
        </div>
      )}

      {scenarios.length > 0 && (
        <div className={es.chartContainer}>
          <h2 className={es.chartTitle}>Hourly Rate Needed by Payoff Target</h2>
          <GradientArea
            data={scenarios.map((s) => ({ label: `${s.months}mo`, rate: s.hourly_rate }))}
            dataKey="rate"
            height={240}
            color="#111"
            xLabel="Payoff Target"
            yLabel="Hourly Rate"
            noFill
          />
        </div>
      )}

      {payoff && !negativeBudget && debts?.length > 0 && (
        <div className={es.stickySection} style={{ ...stickyTop, zIndex: 41 }} ref={sliderRef}>
          <h2 className={es.chartTitle}>Extra Monthly Payment</h2>
          <div className={es.sliderRow}>
            <input
              type="range"
              min="0"
              max={sliderMax}
              step="25"
              value={extraPayment}
              onChange={(e) => onExtraPaymentChange(Number(e.target.value))}
              className={es.sliderInput}
            />
            <span className={es.sliderLabel} style={{ color: extraPayment > 0 ? "var(--green)" : "var(--text-muted)" }}>
              +${extraPayment.toLocaleString()}
            </span>
          </div>
          <div className={es.sliderRange}>
            <span>$0</span>
            <span>${sliderMax.toLocaleString()}</span>
          </div>
        </div>
      )}

      {negativeBudget && (
        <div className={`${styles.card} ${styles.danger}`} style={{ marginBottom: 6 }}>
          <p>Expenses exceed income by <strong>${Math.abs(payoff.budget).toLocaleString()}/mo</strong>. No payoff plan available.</p>
        </div>
      )}

      {hasHistory && (
        <>
          <div
            className={`${es.chartContainer} ${es.chartDropdown}`}
            onClick={() => setBalanceExpanded((v) => !v)}
          >
            <h2 className={`${es.chartTitle} ${es.chartToggle}`}>
              Balance Over Time (Avalanche) <span className={es.expandIcon}>{balanceExpanded ? "−" : "+"}</span>
            </h2>
          </div>

          {balanceExpanded && (
            <div className={es.chartContainer}>
              <StackedArea data={chartHistory} keys={debtNames} xKey="month" height={380} xLabel="Month" yLabel="Balance" lines={interestLine} dualAxis={hasInterest} rightYLabel="Cumulative Interest" />
              {hasInterest && (
                <div className={es.interestLegend}>
                  <span className={es.interestLegendLine} />
                  <span>Cumulative Interest (right axis)</span>
                </div>
              )}
            </div>
          )}

          <div ref={totalDebtRef} className={`${es.stickySection} ${es.chartDropdown}`} style={{ top: chartTop, "--chart-top": `${chartTop}px`, zIndex: 41 }} onClick={() => setTotalDebtExpanded((v) => !v)}>
            <h2 className={`${es.chartTitle} ${es.chartToggle}`}>
              Total Debt Over Time <span className={es.expandIcon}>{totalDebtExpanded ? "−" : "+"}</span>
            </h2>
          </div>
          {totalDebtExpanded && (
            <div className={es.chartContainer}>
              <GradientArea data={totalDebtData} dataKey="total" height={280} color="#c41e1e" xLabel="Month" yLabel="Total Owed" noFill />
            </div>
          )}
        </>
      )}

    </>
  );
}
