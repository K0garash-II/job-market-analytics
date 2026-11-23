import React from "react";
import { Resource, Ranking } from "./components";
import { useMarket } from "./api";

const categories = {
  language: "Языки",
  framework: "Фреймворки",
  database: "Базы данных",
  tool: "Инструменты",
  cloud: "Облака",
};
export default function Rankings({ skills }) {
  const result = useMarket("analytics/" + (skills ? "skills" : "companies"));
  const overview = useMarket(skills ? null : "analytics/overview");
  return (
    <Resource request={result}>
      {({ data }) => {
        const grouped = new Map();
        for (const row of data) {
          const category = row.category || "Прочее";
          if (!grouped.has(category)) grouped.set(category, []);
          grouped.get(category).push(row);
        }
        return (
          <div className={skills ? "charts" : ""}>
            <section className="panel">
              <h2>{skills ? "Частота упоминаний" : "Вакансии и доля рынка"}</h2>
              <Ranking
                rows={data}
                total={skills ? undefined : overview.data?.total}
              />
            </section>
            {skills && (
              <section className="panel">
                <h2>Карта навыков</h2>
                {[...grouped].map(([category, rows]) => (
                  <div className="category" key={category}>
                    <h3>{categories[category] || category}</h3>
                    <div>
                      {rows.map((row) => (
                        <span key={row.name}>
                          {row.name} <b>{row.count}</b>
                        </span>
                      ))}
                    </div>
                  </div>
                ))}
                {!data.length && (
                  <p className="muted">
                    Категории появятся после обработки вакансий.
                  </p>
                )}
              </section>
            )}
          </div>
        );
      }}
    </Resource>
  );
}
