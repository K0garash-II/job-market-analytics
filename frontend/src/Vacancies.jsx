import React, { useState } from "react";
import { number, sourceName, useMarket } from "./api";
import { Resource, Icon } from "./components";

function salary(job) {
  const min = job.salary_min,
    max = job.salary_max;
  if (min == null && max == null) return "";
  const range =
    min != null && max != null
      ? number(min) + "–" + number(max)
      : min != null
        ? "от " + number(min)
        : "до " + number(max);
  return range + " " + job.currency;
}
export function JobList({ jobs }) {
  if (!jobs.length)
    return <p className="muted">Вакансий по этому запросу пока нет.</p>;
  return (
    <ul className="jobs">
      {jobs.map((job) => (
        <li key={job.id}>
          <div>
            {/^https?:\/\//i.test(job.url) ? (
              <a href={job.url} target="_blank" rel="noreferrer">
                {job.title}
                <Icon name="external" />
              </a>
            ) : (
              <strong>{job.title}</strong>
            )}
            <small>
              {[
                job.company,
                job.location,
                job.remote && "Удалённо",
                sourceName(job.source),
              ]
                .filter(Boolean)
                .join(" · ")}
            </small>
          </div>
          {salary(job) && <span>{salary(job)}</span>}
        </li>
      ))}
    </ul>
  );
}
export default function Vacancies() {
  const [filters, setFilters] = useState({ q: "", remote: "", source: "" });
  const [page, setPage] = useState(1);
  const params = new URLSearchParams({ ...filters, page, limit: 20 });
  const result = useMarket("vacancies?" + params);
  function change(event) {
    setPage(1);
    setFilters({ ...filters, [event.target.name]: event.target.value });
  }
  return (
    <section className="panel">
      <div className="filters">
        <label>
          Название или компания
          <input
            type="search"
            name="q"
            value={filters.q}
            onChange={change}
            placeholder="Например, Go"
          />
        </label>
        <label>
          Формат
          <select name="remote" value={filters.remote} onChange={change}>
            <option value="">Любой формат</option>
            <option value="true">Удалённо</option>
            <option value="false">Остальные</option>
          </select>
        </label>
        <label>
          Источник
          <select name="source" value={filters.source} onChange={change}>
            <option value="">Все источники</option>
            {["habr", "remoteok", "hh"].map((source) => (
              <option key={source} value={source}>
                {sourceName(source)}
              </option>
            ))}
          </select>
        </label>
      </div>
      <Resource request={result}>
        {(res) => (
          <>
            <p className="muted" role="status">
              Найдено вакансий: {number(res.total)}
            </p>
            <JobList jobs={res.data} />
            <div className="pager">
              <button disabled={page === 1} onClick={() => setPage(page - 1)}>
                Назад
              </button>
              <span>Страница {page}</span>
              <button
                disabled={page * res.limit >= res.total}
                onClick={() => setPage(page + 1)}
              >
                Далее
              </button>
            </div>
          </>
        )}
      </Resource>
    </section>
  );
}
