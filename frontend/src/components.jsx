import React from "react";
import { number } from "./api";

const paths = {
  overview: "M3 3h7v7H3z M14 3h7v7h-7z M3 14h7v7H3z M14 14h7v7h-7z",
  skills: "m8 7-5 5 5 5 m8-10 5 5-5 5 m-3-13-2 20",
  companies: "M3 21h18 M5 21V4h10v17 M15 10h4v11 M8 8h4 M8 12h4 M8 16h4",
  vacancies: "M3 8h18v13H3z M8 8V4h8v4 M3 13q9 5 18 0",
  arrow: "M5 12h14 m-6-6 6 6-6 6",
  external: "M6 18 18 6 M6 6h12v12",
};
export function Icon({ name }) {
  return (
    <svg
      aria-hidden="true"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="1.5"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <path d={paths[name]} />
    </svg>
  );
}
export function Resource({ request, children }) {
  if (request.error)
    return (
      <div className="notice" role="alert">
        <strong>Не удалось загрузить данные</strong>
        <p>Проверьте соединение с API и повторите запрос.</p>
        <button onClick={request.retry}>Повторить запрос</button>
      </div>
    );
  if (request.loading)
    return (
      <div className="notice" role="status">
        Загружаем данные…
      </div>
    );
  return children(request.data);
}
export function Ranking({ rows, format = (name) => name, total }) {
  if (!rows.length)
    return <p className="muted">Пока нет данных для рейтинга.</p>;
  const max = Math.max(...rows.map((row) => row.count), 1);
  return (
    <ol className="ranking">
      {rows.map((row) => (
        <li key={row.name}>
          <div>
            <span>{format(row.name)}</span>
            <strong>
              {number(row.count)}{" "}
              {total > 0 && (
                <small> · {Math.round((row.count / total) * 100)}%</small>
              )}
            </strong>
          </div>
          <div aria-hidden="true" className="track">
            <i style={{ width: (row.count / (total || max)) * 100 + "%" }} />
          </div>
        </li>
      ))}
    </ol>
  );
}
