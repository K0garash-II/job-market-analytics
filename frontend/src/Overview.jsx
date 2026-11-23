import React, { useState } from "react";
import { useMarket, dayLabel, number, sourceName } from "./api";
import { Icon, Ranking, Resource } from "./components";
import { JobList } from "./Vacancies";

function Publications({ days }) {
  const [range, setRange] = useState(30);
  const [active, setActive] = useState(null);
  const rows = days.slice(-range);
  const max = Math.max(1, ...rows.map((row) => row.count));
  const total = rows.reduce((sum, day) => sum + day.count, 0);
  const dateRange = new Intl.DateTimeFormat("ru-RU", {
    day: "numeric", month: "short", year: "numeric",
  }).formatRange(new Date(rows[0].name + "T12:00:00"), new Date(rows.at(-1).name + "T12:00:00"));
  return (
    <section className="paper">
      <div className="section-head">
        <h2>Динамика публикаций</h2>
        <div className="segment" aria-label="Период графика">
          {[7, 30].map((n) => (
            <button
              key={n}
              aria-pressed={n === range}
              onClick={() => {
                setRange(n);
                setActive(null);
              }}
            >
              {n} дней
            </button>
          ))}
        </div>
      </div>
      <div className="plot-stat">
        <strong>{number(active?.count ?? total)}</strong>
        <span>{active ? dayLabel(active.name) : dateRange}</span>
      </div>
      <div
        className="plot"
        onMouseLeave={() => setActive(null)}
        aria-label="Публикации по дням"
      >
        <div className="y-axis">
          <span>{number(max)}</span>
          <span>0</span>
        </div>
        <div className="bars">
          {rows.map((row) => (
            <div
              key={row.name}
              className={
                "bar-slot" + (active?.name === row.name ? " selected" : "")
              }
              onPointerEnter={(e) => {
                if (e.pointerType === "mouse") setActive(row);
              }}
              onClick={() => setActive(active?.name === row.name ? null : row)}
            >
              <i style={{ height: (row.count / max) * 100 + "%" }} />
            </div>
          ))}
        </div>
      </div>
      <div className="x-axis">
        {[rows[0], rows[Math.floor((rows.length - 1) / 2)], rows.at(-1)].map(
          (row, i) => (
            <span key={i}>{row && dayLabel(row.name)}</span>
          ),
        )}
      </div>
      <details>
        <summary>Данные по дням</summary>
        <dl>
          {rows.map((row) => (
            <div key={row.name}>
              <dt>{dayLabel(row.name)}</dt>
              <dd>{row.count}</dd>
            </div>
          ))}
        </dl>
      </details>
    </section>
  );
}
function WorkFormat({ remote, total }) {
  const percent = total ? Math.round((remote / total) * 100) : 0;
  return (
    <section className="panel format">
      <h2>Формат работы</h2>
      <div className="ring">
        <svg viewBox="0 0 120 120" aria-hidden="true">
          <circle className="ring-track" cx="60" cy="60" r="48" />
          {remote > 0 && (
            <circle
              className="ring-value"
              cx="60"
              cy="60"
              r="48"
              pathLength="100"
              strokeDasharray={(remote / total) * 100 + " 100"}
            />
          )}
        </svg>
        <div>
          <strong>{total ? percent + "%" : "—"}</strong>
          <small>удалённо</small>
        </div>
      </div>
      <dl>
        {[
          ["Удалённо", remote],
          ["Остальные", total - remote],
        ].map(([label, n], i) => (
          <div key={label}>
            <dt>
              <i className={"dot d" + i} />
              {label}
            </dt>
            <dd>{number(n)}</dd>
          </div>
        ))}
      </dl>
    </section>
  );
}
export default function Overview() {
  const overview = useMarket("analytics/overview?period=latest");
  const latest = useMarket("vacancies?limit=5");
  return (
    <Resource request={overview}>
      {(data) => (
        <>
          <dl className="metrics">
            {[
              ["Вакансии", data.total],
              ["Работодатели", data.companies],
              ["С зарплатой", data.with_salary],
              ["Источники", data.sources.length],
            ].map(([label, value]) => (
              <div key={label}>
                <dt>{label}</dt>
                <dd>{number(value)}</dd>
              </div>
            ))}
          </dl>
          <div className="charts">
            <Publications days={data.days} />
            <WorkFormat remote={data.remote} total={data.total} />
          </div>
          <div className="lower">
            <section className="panel">
              <div className="section-head">
                <h2>Последние вакансии</h2>
                <a
                  className="icon-link"
                  href="#vacancies"
                  aria-label="Смотреть все вакансии"
                >
                  <Icon name="arrow" />
                </a>
              </div>
              <Resource request={latest}>
                {(res) => <JobList jobs={res.data} />}
              </Resource>
            </section>
            <section className="panel">
              <h2>Источники</h2>
              <Ranking
                rows={data.sources}
                format={sourceName}
                total={data.total}
              />
            </section>
          </div>
        </>
      )}
    </Resource>
  );
}
