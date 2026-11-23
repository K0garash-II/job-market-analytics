import { defineConfig } from "vite";

const proxy = { "/api": process.env.API_TARGET || "http://127.0.0.1:8080" };
export default defineConfig({ server: { proxy }, preview: { proxy } });
