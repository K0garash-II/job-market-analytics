import React, { useState, useEffect } from "react";
import { createRoot } from "react-dom/client";
import Overview from "./Overview";
import Vacancies from "./Vacancies";
import Rankings from "./Rankings";
import { Icon } from "./components";
import "./styles.css";

const pages = {
  overview: ["Обзор", "Обзор рынка"],
  skills: ["Навыки", "Спрос на навыки"],
  companies: ["Компании", "Работодатели"],
  vacancies: ["Вакансии", "Вакансии"],
};
function route() {
  const hash = location.hash.slice(1);
  return Object.hasOwn(pages, hash) ? hash : "overview";
}
function App() {
  const [page, setPage] = useState(route);
  useEffect(() => {
    const change = () => {
      setPage(route());
      window.scrollTo(0, 0);
    };
    window.addEventListener("hashchange", change);
    return () => window.removeEventListener("hashchange", change);
  }, []);
  useEffect(() => {
    document.title = pages[page][1] + " · Job Market Analytics";
  }, [page]);
  return (
    <>
      <a
        className="skip"
        href="#content"
        onClick={(e) => {
          e.preventDefault();
          document.getElementById("content").focus();
        }}
      >
        К содержанию
      </a>
      <header>
        <div className="wrap mast">
          <a href="#overview" className="brand">
            <img src="/mark.svg" alt="" width="34" height="34" />
            Job Market Analytics
          </a>
          <nav aria-label="Разделы">
            {Object.entries(pages).map(([id, [label]]) => (
              <a
                href={"#" + id}
                key={id}
                aria-current={page === id ? "page" : undefined}
              >
                <Icon name={id} />
                {label}
              </a>
            ))}
          </nav>
        </div>
      </header>
      <main className="wrap" id="content" tabIndex="-1">
        <h1>{pages[page][1]}</h1>
        {page === "overview" ? (
          <Overview />
        ) : page === "vacancies" ? (
          <Vacancies />
        ) : (
          <Rankings key={page} skills={page === "skills"} />
        )}
      </main>
    </>
  );
}
createRoot(document.getElementById("root")).render(<App />);
